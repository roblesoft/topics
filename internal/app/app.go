package app

import (
	http "github.com/roblesoft/topics/internal/controller/http"
	"github.com/spf13/viper"
)

func Run() {
	viper.SetConfigFile("./pkg/envs/.env")
	viper.ReadInConfig()

	port := viper.Get("PORT").(string)

	server := http.NewServer(port)
	server.Start()
}
