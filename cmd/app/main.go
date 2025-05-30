package main

import (
	//"context"
	"github.com/gin-gonic/gin"
	"gw-currency-wallet/internal/database"
	"gw-currency-wallet/internal/database/query"
	"gw-currency-wallet/internal/services"
	"gw-currency-wallet/internal/transport/rest"
	"log"
)

func main() {
	router := gin.Default()

	conn, err := database.Connect()
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	// defer conn.Close(context.Background())

	repo := query.NewRepository(conn)
	userService := services.NewUserService(repo)

	rest.Routers(router, userService)

	router.Run(":8080")
}
