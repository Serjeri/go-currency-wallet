package rest

import (
	"github.com/gin-gonic/gin"
	"gw-currency-wallet/domain/handlers"
	"gw-currency-wallet/domain/services/auth"
)

func Routers(r *gin.Engine, client handlers.UserService, service handlers.GrpcService) {
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
			wallet.POST("/update", func(c *gin.Context) {
				handlers.UpdateUserBalance(c, client)
			})
		}

		exchange := privateApi.Group("/exchange")
		{
			exchange.GET("/rates", func(c *gin.Context) {
				handlers.GetExchangeRates(c, service)
			})
			exchange.POST("/", func(c *gin.Context) {
				handlers.PerformExchange(c, service)
			})
		}
	}

}
