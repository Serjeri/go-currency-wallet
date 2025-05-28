package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Password hashing failed"})
		return
	}

	user.Password = string(hashedPassword)

	exists, err := client.repository.RegistrUser(context.TODO(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Registration failed"})
		return
	}

	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	token, err := auth.CreateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}

	c.SetCookie(
		"auth_token",
		token,
		3600,
		"/",
		"localhost",
		false,
		true,
	)

	builder.WriteString("Bearer ")
	builder.WriteString(token)
	authHeader := builder.String()

	c.Header("Authorization", authHeader)
	c.JSON(http.StatusOK, gin.H{
		"message": "Registration successful",
	})
}
