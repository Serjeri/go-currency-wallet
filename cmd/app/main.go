package main

import (
	"gw-currency-wallet/domain/repository"
	"gw-currency-wallet/domain/repository/query"
	"gw-currency-wallet/domain/services"
	"gw-currency-wallet/domain/transport/gprc"
	"gw-currency-wallet/domain/transport/rest"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	conn, err := repository.Connect()
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	repo := query.NewRepository(conn)

	client, closer := gprc.New("localhost", "50051")
	defer closer()

	userService := services.NewUserService(repo, client)

	rest.Routers(router, userService, userService)
	router.Run(":8080")
}
