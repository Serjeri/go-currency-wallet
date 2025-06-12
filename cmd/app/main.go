package main

import (
	"gw-currency-wallet/domain/config"
	"gw-currency-wallet/domain/repository"
	"gw-currency-wallet/domain/repository/query"
	"gw-currency-wallet/domain/services"
	"gw-currency-wallet/domain/transport/gprc"
	"gw-currency-wallet/domain/transport/rest"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.MustLoad()

	conn, err := repository.Connect(cfg.Dburl)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	router := gin.Default()
	repo := query.NewRepository(conn)

	client, closer := gprc.New(cfg.Addressgrpc)
	defer closer()

	userService := services.NewUserService(repo, client)

	rest.Routers(router, userService, userService)
	router.Run(cfg.Address)
}
