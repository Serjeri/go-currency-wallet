package services

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"gw-currency-wallet/internal/database/query"
	"gw-currency-wallet/internal/models"
	"net/http"
)

type Client struct {
	repository *query.Repository
}

func NewClient(repository *query.Repository) *Client {
	return &Client{repository: repository}
}

func (client *Client) UserRegistr(c *gin.Context) {
	var user models.User

	err := c.ShouldBind(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reg, _ := client.repository.RegistrUser(context.TODO(), user)
	fmt.Println(reg)
	c.String(200, "Success")
}
