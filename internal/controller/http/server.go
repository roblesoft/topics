package controllers

import (
	"github.com/gin-gonic/gin"
)

type Server struct {
	Port   string
	router *gin.Engine
}

func NewServer(port string) *Server {
	server := &Server{Port: port}
	server.setupRouter()
	return server
}

func (server *Server) Router() *gin.Engine {
	return server.Router()
}

func (server *Server) setupRouter() {
	router := gin.Default()
	router.GET("/healthcheck/", server.HealthCheck)
	server.router = router
}

func (r *Server) Start() {
	r.router.Run(r.Port)
}
