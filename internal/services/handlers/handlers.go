package handlers

import (
	"context"
	"gw-currency-wallet/internal/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserService interface {
	CreateUser(ctx context.Context, user *models.User) (string, error)
	GetUser(ctx context.Context, user *models.Login) (string, error)
}

func UserRegistration(c *gin.Context, s UserService) {
	var user models.User
	var builder strings.Builder

	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	users, err := s.CreateUser(context.TODO(), &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user with this name or email already exists"})
		return
	}

	builder.WriteString("Bearer ")
	builder.WriteString(users)
	authHeader := builder.String()

	c.Header("Authorization", authHeader)
	c.JSON(http.StatusOK, gin.H{
		"message": "Registration successful",
	})
}

func UserAuthenticate(c *gin.Context, s UserService) {
	var login models.Login
	var builder strings.Builder

	if err := c.ShouldBind(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	users, err := s.GetUser(context.TODO(), &login)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user with this name or email already exists"})
		return
	}

	builder.WriteString("Bearer ")
	builder.WriteString(users)
	authHeader := builder.String()

	c.Header("Authorization", authHeader)
	c.JSON(http.StatusOK, gin.H{
		"message": "Successful",
	})
}
