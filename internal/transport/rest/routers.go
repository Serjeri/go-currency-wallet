package rest

import (
	"github.com/gin-gonic/gin"
	"gw-currency-wallet/internal/services"
)

func Routers(r *gin.Engine, client *services.Client) {
	users := r.Group("/api/v1")
	// wallet := r.Group("/api/v2")
	// exchange := r.Group("/api/v3")

	{
		users.POST("/register", client.UserRegistr)
		// users.POST("/login",)
	}
	// {
	// 	wallet.GET("/balance")
	// 	wallet.POST("wallet/deposit")
	// 	wallet.POST("wallet/withdraw")
	// }

	// {
	// 	exchange.GET("/exchange/rates")
	// 	exchange.POST("/exchange")
	// }
}
