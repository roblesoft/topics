package app

import (
	"fmt"
	"os"

	"github.com/go-redis/redis/v7"
	amqp "github.com/rabbitmq/amqp091-go"
	http "github.com/roblesoft/topics/internal/controller/http"
	"github.com/roblesoft/topics/internal/usecase"
	repo "github.com/roblesoft/topics/internal/usecase/repo"
	"github.com/roblesoft/topics/pkg/db"
)

func Run() {
	var dbUrl = os.Getenv("DB_URL")
	var db = db.Init(dbUrl)
	fmt.Println(dbUrl)
	fmt.Println(db)

	var redis_addr = fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	var rdb = redis.NewClient(&redis.Options{Addr: redis_addr, Password: os.Getenv("REDIS_PASSWORD")})

	fmt.Println(redis_addr)
	fmt.Println(rdb)

	var port = os.Getenv("PORT")

	var connection, err = amqp.Dial(os.Getenv("RABBITMQ_ADDRESS"))
	var userRepo = &repo.UserRepository{Db: db}

	fmt.Println(port)
	fmt.Println(connection)

	var service = usecase.NewService(userRepo)
	var server = http.NewServer(port, *service, *rdb, *connection)

	if err != nil {
		panic(err)
	}

	defer connection.Close()

	server.Start()
}
