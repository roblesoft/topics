package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
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

	var (
		queue, channel = connectChannel(username, currentusr.Username, server.Rabbitmq, conn)
		closeCh        = onDisconnect(conn, channel)
	)

	if username == currentusr.Username {
		consumeMessage(channel, queue, conn)
	}

loop:
	for {
		select {
		case <-closeCh:
			break loop
		default:
			onUserMessage(currentusr.Username, channel, queue, conn)
		}
	}
}

func onDisconnect(conn *websocket.Conn, channel *amqp.Channel) chan struct{} {
	closeCh := make(chan struct{})

	conn.SetCloseHandler(func(code int, text string) error {
		close(closeCh)
		channel.Close()
		return nil
	})

	return closeCh
}

func connectChannel(
	username string,
	currentusr string,
	connection *amqp.Connection,
	conn *websocket.Conn) (*amqp.Queue, *amqp.Channel) {

	channel, err := connection.Channel()

	if err != nil {
		handleWSError(err, conn)
		return nil, nil
	}

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
		return nil, nil
	}

	return &queue, channel
}

func consumeMessage(channel *amqp.Channel, queue *amqp.Queue, conn *websocket.Conn) {
	defer channel.Close()

	msgs, err := channel.Consume(
		queue.Name, // queue
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
	msg entity.Message,
	channel *amqp.Channel,
	queue *amqp.Queue,
	conn *websocket.Conn) {

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
	currentusr string,
	channel *amqp.Channel,
	queue *amqp.Queue,
	conn *websocket.Conn) {

	var msg entity.Message

	if err := conn.ReadJSON(&msg); err != nil {
		handleWSError(err, conn)
		return
	}

	msg.Username = currentusr

	publishMessage(msg, channel, queue, conn)
}
