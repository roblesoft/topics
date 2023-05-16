package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/roblesoft/topics/internal/usecase"
)

type Server struct {
	Port        string
	router      *gin.Engine
	Service     *usecase.Service
	RedisClient *redis.Client
}

func NewServer(port string, service usecase.Service, redisClient redis.Client) *Server {
	server := &Server{Port: port, Service: &service, RedisClient: &redisClient}
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
		chat   = v1.Group("/chat")
	)

	users.POST("/register", server.Register)
	users.POST("/login", server.Login)
	chat.GET("/index", server.ChatnetHandler)

	v1.GET("/healthcheck/", server.HealthCheck)

	server.router = router
}

func (r *Server) Start() {
	r.router.Run(r.Port)
}
