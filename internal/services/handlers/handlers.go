package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"gw-currency-wallet/internal/database/query"
	"gw-currency-wallet/internal/models"
	"gw-currency-wallet/internal/services/auth"
	"net/http"
	"strings"
)

type Client struct {
	repository *query.Repository
}

func NewClient(repository *query.Repository) *Client {
	return &Client{repository: repository}
}

func (client *Client) UserRegistr(c *gin.Context) {
	var user models.User
	var builder strings.Builder

	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	hashedPassword := HashedPassword(user.Password, c)

	user.Password = hashedPassword

	id, err := client.repository.RegistrUser(context.TODO(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Registration failed"})
		return
	}

	if id == 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	token, err := auth.CreateToken(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}

	builder.WriteString("Bearer ")
	builder.WriteString(token)
	authHeader := builder.String()

	c.Header("Authorization", authHeader)
	c.JSON(http.StatusOK, gin.H{
		"message": "Registration successful",
	})
}

func (client *Client) AuthenticateUser(c *gin.Context) {
	var login models.Login
	var builder strings.Builder

	if err := c.ShouldBind(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}
	hashedPassword := HashedPassword(login.Password, c)

	login.Password = hashedPassword

	id, err := client.repository.GetUser(context.TODO(), login )
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}

	token, err := auth.CreateToken(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}

	builder.WriteString("Bearer ")
	builder.WriteString(token)
	authHeader := builder.String()

	c.Header("Authorization", authHeader)
	c.JSON(http.StatusOK, gin.H{
		"message": "Successful",
	})
}
