package app

import (
	"fmt"
	"os"

	"github.com/go-redis/redis/v7"
	http "github.com/roblesoft/topics/internal/controller/http"
	entity "github.com/roblesoft/topics/internal/entity"
	"github.com/roblesoft/topics/internal/usecase"
	repo "github.com/roblesoft/topics/internal/usecase/repo"
	"github.com/roblesoft/topics/pkg/db"
)

func Run() {
	var (
		redis_addr = fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
		rdb        = redis.NewClient(&redis.Options{Addr: redis_addr, Password: os.Getenv("REDIS_PASSWORD")})
		port       = os.Getenv("PORT")
		dbUrl      = os.Getenv("DB_URL")
		db         = db.Init(dbUrl)
		userRepo   = &repo.UserRepository{Db: db}
		service    = usecase.NewService(userRepo)
		server     = http.NewServer(port, *service, *rdb)
	)

	db.AutoMigrate(&entity.User{})
	server.Start()
}
