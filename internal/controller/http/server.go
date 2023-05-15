package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/roblesoft/topics/internal/usecase"
)

type Server struct {
	Port    string
	router  *gin.Engine
	Service *usecase.Service
}

func NewServer(port string, service usecase.Service) *Server {
	server := &Server{Port: port, Service: &service}
	server.setupRouter()
	return server
}

func (server *Server) Router() *gin.Engine {
	return server.router
}

func (server *Server) setupRouter() {
	var (
		router = gin.Default()
		api    = router.Group("/api")
		v1     = api.Group("/v1")
		users  = v1.Group("/users")
	)

	users.POST("/register", server.Register)
	users.POST("/login", server.Login)

	v1.GET("/healthcheck/", server.HealthCheck)

	server.router = router
}

func (r *Server) Start() {
	r.router.Run(r.Port)
}
