package rest

import (
	"gw-currency-wallet/domain/handlers"
	"gw-currency-wallet/domain/services/auth"

	"github.com/gin-gonic/gin"
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

	privateApi := api.Group("/")
	privateApi.Use(auth.AuthenticateMiddleware)
	{
	    wallet := privateApi.Group("/wallet")
	    {
	        wallet.GET("/balance", func(c *gin.Context) {
				handlers.GetUserBalance(c, client)
			})
	        wallet.POST("/deposit", func(c *gin.Context) {
				handlers.UpdateUserBalance(c, client)
			})
	        // wallet.POST("/withdraw", client.Withdraw)
	    }

	//     exchange := privateApi.Group("/exchange")
	//     {
	//         exchange.GET("/rates", client.GetExchangeRates)
	//         exchange.POST("/", client.ExchangeCurrency)
	//     }
	}
}
