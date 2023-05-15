package app

import (
	"os"

	http "github.com/roblesoft/topics/internal/controller/http"
	entity "github.com/roblesoft/topics/internal/entity"
	"github.com/roblesoft/topics/internal/usecase"
	repo "github.com/roblesoft/topics/internal/usecase/repo"
	"github.com/roblesoft/topics/pkg/db"
)

func Run() {
	var (
		port     = os.Getenv("PORT")
		dbUrl    = os.Getenv("DB_URL")
		db       = db.Init(dbUrl)
		userRepo = &repo.UserRepository{Db: db}
		service  = usecase.NewService(userRepo)
		server   = http.NewServer(port, *service)
	)

	db.AutoMigrate(&entity.User{})
	server.Start()
}
