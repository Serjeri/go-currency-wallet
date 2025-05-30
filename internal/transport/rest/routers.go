package rest

import (
	"github.com/gin-gonic/gin"
	//"gw-currency-wallet/internal/services/auth"
	"gw-currency-wallet/internal/services/handlers"
)

func Routers(r *gin.Engine, client handlers.UserService) {
	api := r.Group("/api/v1")

	publicApi := api.Group("/")
	{
		publicApi.POST("/register", func(c *gin.Context) {
			handlers.UserRegistration(c, client)
		})
		publicApi.POST("/login", func(c *gin.Context) {
			handlers.UserAuthenticate(c, client)
		})
	}

	// privateApi := api.Group("/")
	// privateApi.Use(auth.AuthenticateMiddleware)
	// {
	//     wallet := privateApi.Group("/wallet")
	//     {
	//         wallet.GET("/balance", client.GetBalance)
	//         wallet.POST("/deposit", client.Deposit)
	//         wallet.POST("/withdraw", client.Withdraw)
	//     }

	//     exchange := privateApi.Group("/exchange")
	//     {
	//         exchange.GET("/rates", client.GetExchangeRates)
	//         exchange.POST("/", client.ExchangeCurrency)
	//     }
	// }
}
