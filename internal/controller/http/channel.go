package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/gorilla/websocket"
	entity "github.com/roblesoft/topics/internal/entity"
	"github.com/roblesoft/topics/internal/usecase"
)

var connectedUsers = make(map[string]*usecase.Service)

const (
	commandSubscribe = iota
	commandUnsubscribe
	commandChat
)

func (server *Server) GetChannel(ctx *gin.Context) {
	var (
		currentusr = server.CurrentUser(ctx)
		username   = ctx.Param("user")
		_, err     = server.Service.GetUserByUsername(username)
	)

	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		handleWSError(err, conn)
		return
	}

	err = onConnect(ctx, conn, server.RedisClient, username)
	if err != nil {
		handleWSError(err, conn)
		return
	}

	closeCh := onUserDisconnect(ctx, conn, server.RedisClient, username)

	onChannelMessage(conn, ctx, *currentusr, username)

loop:
	for {
		select {
		case <-closeCh:
			break loop
		default:
			onUsersMessage(conn, ctx, server.RedisClient, username)
		}
	}
}

func onConnect(ctx *gin.Context, conn *websocket.Conn, rdb *redis.Client, username string) error {
	u, err := usecase.Connect(rdb, username)
	if err != nil {
		return err
	}
	connectedUsers[username] = u
	return nil
}

func onUserDisconnect(ctx *gin.Context, conn *websocket.Conn, rdb *redis.Client, username string) chan struct{} {
	closeCh := make(chan struct{})

	conn.SetCloseHandler(func(code int, text string) error {
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

func onUsersMessage(conn *websocket.Conn, ctx *gin.Context, rdb *redis.Client, username string) {
	var msg entity.Message

	if err := conn.ReadJSON(&msg); err != nil {
		handleWSError(err, conn)
		return
	}
	u := connectedUsers[username]

	switch msg.Command {
	case commandSubscribe:
		if err := u.Subscribe(rdb, username); err != nil {
			handleWSError(err, conn)
		}
	case commandUnsubscribe:
		if err := u.Unsubscribe(rdb, username); err != nil {
			handleWSError(err, conn)
		}
	case commandChat:
		if err := usecase.Chat(rdb, username, msg.Content); err != nil {
			handleWSError(err, conn)
		}
	}
}

func onChannelMessage(conn *websocket.Conn, ctx *gin.Context, currentusr entity.User, username string) {
	u := connectedUsers[username]
	go func() {
		for m := range u.MessageChan {

			msg := entity.Message{
				Content: m.Payload,
				Channel: m.Channel,
			}

			if err := conn.WriteJSON(msg); err != nil {
				handleWSError(err, conn)
			}
		}
	}()
}