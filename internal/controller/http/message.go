package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/gorilla/websocket"
	"github.com/roblesoft/topics/internal/usecase"
)

func (server *Server) HealthCheck(ctx *gin.Context) {
	ctx.Status(200)
}

var upgrader websocket.Upgrader

var connectedUsers = make(map[string]*usecase.Service)

func H(rdb *redis.Client, fn func(http.ResponseWriter, *http.Request, *redis.Client)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, rdb)
	}
}

type msg struct {
	Content string `json:"content,omitempty"`
	Channel string `json:"channel,omitempty"`
	Command int    `json:"command,omitempty"`
	Err     string `json:"err,omitempty"`
}

const (
	commandSubscribe = iota
	commandUnsubscribe
	commandChat
)

func (server *Server) ChatnetHandler(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)

	if err != nil {
		handleWSError(err, conn)
		return
	}

	err = onConnect(ctx.Request, conn, server.RedisClient)
	if err != nil {
		handleWSError(err, conn)
		return
	}

	closeCh := onDisconnect(ctx.Request, conn, server.RedisClient)

	onChannelMessage(conn, ctx.Request)

loop:
	for {
		select {
		case <-closeCh:
			break loop
		default:
			onUserMessage(conn, ctx.Request, server.RedisClient)
		}
	}
}

func onConnect(r *http.Request, conn *websocket.Conn, rdb *redis.Client) error {
	username := r.URL.Query()["username"][0]
	fmt.Println("connected from:", conn.RemoteAddr(), "user:", username)

	fmt.Println(rdb)
	u, err := usecase.Connect(rdb, username)
	if err != nil {
		return err
	}
	connectedUsers[username] = u
	return nil
}

func onDisconnect(r *http.Request, conn *websocket.Conn, rdb *redis.Client) chan struct{} {

	closeCh := make(chan struct{})

	username := r.URL.Query()["username"][0]

	conn.SetCloseHandler(func(code int, text string) error {
		fmt.Println("connection closed for user", username)

		u := connectedUsers[username]
		if err := u.Disconnect(); err != nil {
			return err
		}
		delete(connectedUsers, username)
		close(closeCh)
		return nil
	})

	return closeCh
}

func onUserMessage(conn *websocket.Conn, r *http.Request, rdb *redis.Client) {

	var msg msg

	if err := conn.ReadJSON(&msg); err != nil {
		handleWSError(err, conn)
		return
	}

	username := r.URL.Query()["username"][0]
	u := connectedUsers[username]

	switch msg.Command {
	case commandSubscribe:
		if err := u.Subscribe(rdb, msg.Channel); err != nil {
			handleWSError(err, conn)
		}
	case commandUnsubscribe:
		if err := u.Unsubscribe(rdb, msg.Channel); err != nil {
			handleWSError(err, conn)
		}
	case commandChat:
		if err := usecase.Chat(rdb, msg.Channel, msg.Content); err != nil {
			handleWSError(err, conn)
		}
	}
}

func onChannelMessage(conn *websocket.Conn, r *http.Request) {

	username := r.URL.Query()["username"][0]
	u := connectedUsers[username]

	go func() {
		for m := range u.MessageChan {

			msg := msg{
				Content: m.Payload,
				Channel: m.Channel,
			}

			if err := conn.WriteJSON(msg); err != nil {
				fmt.Println(err)
			}
		}

	}()
}

func handleWSError(err error, conn *websocket.Conn) {
	_ = conn.WriteJSON(msg{Err: err.Error()})
}
