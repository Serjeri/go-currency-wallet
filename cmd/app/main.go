package main

import (
	"gw-currency-wallet/domain/repository"
	"gw-currency-wallet/domain/repository/query"
	"gw-currency-wallet/domain/services"
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
	userService := services.NewUserService(repo)

	rest.Routers(router, userService)

	router.Run(":8080")
}
