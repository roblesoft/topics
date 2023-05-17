package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/roblesoft/topics/internal/usecase"
	"github.com/roblesoft/topics/pkg/token"
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
		router   = gin.Default()
		api      = router.Group("/api")
		v1       = api.Group("/v1")
		users    = v1.Group("/users")
		messages = users.Group("/:user/messages")
	)
	api.Use(gin.Logger())

	v1.GET("/healthcheck/", server.HealthCheck)
	users.POST("/register", server.Register)
	users.POST("/login", server.Login)

	messages.Use(JwtAuthMiddleware())
	messages.GET("/", server.ChatnetHandler)

	server.router = router
}

func (r *Server) Start() {
	r.router.Run(r.Port)
}

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := token.TokenValid(c)
		if err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
		c.Next()
	}
}
