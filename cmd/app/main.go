package main

import (
	//"context"
	"gw-currency-wallet/internal/database"
	"gw-currency-wallet/internal/database/query"
	"gw-currency-wallet/internal/services"
	"gw-currency-wallet/internal/transport/rest"
	"log"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	conn, err := database.Connect()
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	// defer conn.Close(context.Background())

	repo := query.NewRepository(conn)
	handl := services.NewClient(repo)

	rest.Routers(router, handl)

	router.Run(":8080")
}
