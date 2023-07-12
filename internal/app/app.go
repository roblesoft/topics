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
	var (
		redis_addr      = fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
		rdb             = redis.NewClient(&redis.Options{Addr: redis_addr, Password: os.Getenv("REDIS_PASSWORD")})
		port            = os.Getenv("PORT")
		dbUrl           = os.Getenv("DB_URL")
		db              = db.Init(dbUrl)
		connection, err = amqp.Dial(os.Getenv("RABBITMQ_ADDRESS"))
		userRepo        = &repo.UserRepository{Db: db}
		service         = usecase.NewService(userRepo)
		server          = http.NewServer(port, *service, *rdb, *connection)
	)

	fmt.Println(redis_addr)
	fmt.Println(rdb)
	fmt.Println(port)
	fmt.Println(dbUrl)
	fmt.Println(db)
	fmt.Println(connection)
	if err != nil {
		panic(err)
	}

	defer connection.Close()

	server.Start()
}
