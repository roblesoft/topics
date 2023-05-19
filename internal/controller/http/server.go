package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/roblesoft/topics/internal/usecase"
	"github.com/roblesoft/topics/pkg/token"
)

type Server struct {
	Port        string
	Router      *gin.Engine
	Service     *usecase.Service
	RedisClient *redis.Client
	Rabbitmq    *amqp.Connection
}

func NewServer(
	port string,
	service usecase.Service,
	redisClient redis.Client,
	rabbitConnection amqp.Connection) *Server {

	server := &Server{
		Port:        port,
		Service:     &service,
		RedisClient: &redisClient,
		Rabbitmq:    &rabbitConnection,
	}
	server.setupRouter()

	return server
}

func (server *Server) setupRouter() {
	var (
		router   = gin.Default()
		api      = router.Group("/api")
		v1       = api.Group("/v1")
		users    = v1.Group("/users")
		messages = users.Group("/:user/messages")
		channels = v1.Group("/channels")
	)
	api.Use(gin.Logger())

	v1.GET("/healthcheck/", server.HealthCheck)
	users.POST("/register", server.Register)
	users.POST("/login", server.Login)

	messages.Use(JwtAuthMiddleware())
	messages.GET("/", server.MessageIndex)

	channels.Use(JwtAuthMiddleware())
	channels.GET("/:channel", server.GetChannel)

	server.Router = router
}

func (server *Server) Start() {
	server.Router.Run(server.Port)
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
