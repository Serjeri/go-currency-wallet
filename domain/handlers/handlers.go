package handlers

import (
	"context"
	"gw-currency-wallet/domain/models"
	"gw-currency-wallet/domain/services/auth"

	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserService interface {
	CreateUser(ctx context.Context, user *models.User) (string, error)
	GetUser(ctx context.Context, user *models.Login) (string, error)
	GetBalanceUser(ctx context.Context, id int) (*models.Balance, error)
	UpdateBalanceUser(ctx context.Context, id int, updateBalance *models.UpdateBalance) (*models.Balance, error)
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
		"token":   users,
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
		"token":   users,
	})
}

func GetUserBalance(c *gin.Context, s UserService) {
	token := c.GetHeader("Authorization")

	userID, err := auth.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user with this name or email already exists"})
		return
	}

	balance, err := s.GetBalanceUser(context.Background(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get balance"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"balance": gin.H{
			"EUR": float64(balance.EUR)/100,
			"RUB": float64(balance.RUB)/100,
			"USD": float64(balance.USD)/100,
		},
	})
}

func UpdateUserBalance(c *gin.Context, s UserService) {
	var update models.UpdateBalance

	if err := c.ShouldBind(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	token := c.GetHeader("Authorization")

	userID, err := auth.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user with this name or email already exists"})
		return
	}

	updatedBalance, err := s.UpdateBalanceUser(context.Background(), userID, &update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get balance"})
		return
	}

    c.JSON(http.StatusOK, gin.H{
        "balance": gin.H{
            "EUR": float64(updatedBalance.EUR)/100,
			"RUB": float64(updatedBalance.RUB)/100,
			"USD": float64(updatedBalance.USD)/100,
        },
    })
}
