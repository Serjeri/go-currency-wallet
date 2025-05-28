package rest

import (
	"github.com/gin-gonic/gin"
	//"gw-currency-wallet/internal/services/auth"
	"gw-currency-wallet/internal/services/handlers"
)

func Routers(r *gin.Engine, client *handlers.Client) {
	api := r.Group("/api/v1")

	publicApi := api.Group("/")
	{
		publicApi.POST("/register", client.UserRegistr)
		// publicApi.POST("/login",)
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
