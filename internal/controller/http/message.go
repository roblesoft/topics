package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/gorilla/websocket"
	amqp "github.com/rabbitmq/amqp091-go"
	entity "github.com/roblesoft/topics/internal/entity"
)

var upgrader websocket.Upgrader

func (server *Server) MessageIndex(ctx *gin.Context) {
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

	if username == currentusr.Username {
		consumeMessage(username, currentusr.Username, server.Rabbitmq, conn)
	}

	closeCh := onDisconnect(conn)

loop:
	for {
		select {
		case <-closeCh:
			break loop
		default:
			onUserMessage(conn, ctx, server.RedisClient, username, currentusr.Username, server.Rabbitmq)
		}
	}
}

func onDisconnect(conn *websocket.Conn) chan struct{} {
	closeCh := make(chan struct{})

	conn.SetCloseHandler(func(code int, text string) error {
		close(closeCh)
		return nil
	})

	return closeCh
}

func consumeMessage(username string, currentusr string, connection *amqp.Connection, conn *websocket.Conn) {
	channel, err := connection.Channel()
	if err != nil {
		handleWSError(err, conn)
		return
	}

	defer channel.Close()

	msgs, err := channel.Consume(
		currentusr, // queue
		"",         // consumer
		true,       // auto ack
		false,      // exclusive
		false,      // no local
		false,      // no wait
		nil,        // args
	)

	if err != nil {
		handleWSError(err, conn)
		return
	}

	listenMessages := make(chan bool)
	go func() {
		for m := range msgs {
			var msg entity.Message

			if err := json.Unmarshal(m.Body, &msg); err != nil {
				handleWSError(err, conn)
				return
			}

			if err := conn.WriteJSON(msg); err != nil {
				handleWSError(err, conn)
			}
		}
	}()
	<-listenMessages
}

func publishMessage(
	username string,
	msg entity.Message,
	connection *amqp.Connection,
	conn *websocket.Conn) {

	channel, err := connection.Channel()
	if err != nil {
		handleWSError(err, conn)
	}
	defer channel.Close()

	queue, err := channel.QueueDeclare(
		username, // name
		false,    // durable
		false,    // auto delete
		false,    // exclusive
		false,    // no wait
		nil,      // args
	)

	if err != nil {
		handleWSError(err, conn)
		return
	}

	body, err := json.Marshal(msg)
	if err != nil {
		handleWSError(err, conn)
		return
	}

	err = channel.Publish(
		"",         // exchange
		queue.Name, // key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		},
	)
	if err != nil {
		handleWSError(err, conn)
		return
	}
	fmt.Println("Queue status:", queue)
	fmt.Println("Successfully published message")
}

func onUserMessage(
	conn *websocket.Conn,
	ctx *gin.Context,
	rdb *redis.Client,
	username string,
	currentusr string,
	connection *amqp.Connection) {

	var msg entity.Message

	msg.Username = currentusr
	if err := conn.ReadJSON(&msg); err != nil {
		handleWSError(err, conn)
		return
	}

	if currentusr != username {
		publishMessage(username, msg, connection, conn)
		return
	}
}
