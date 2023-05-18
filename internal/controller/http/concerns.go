package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	entity "github.com/roblesoft/topics/internal/entity"
)

func handleWSError(err error, conn *websocket.Conn) {
	_ = conn.WriteJSON(entity.Message{Err: err.Error()})
}

func (server *Server) HealthCheck(ctx *gin.Context) {
	ctx.Status(200)
}
