package controllers

import (
	"github.com/gin-gonic/gin"
)

func (server *Server) HealthCheck(ctx *gin.Context) {
	ctx.Status(200)
}
