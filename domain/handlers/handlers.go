package handlers

import (
	"context"
	"database/sql"
	"errors"
	pb "github.com/Serjeri/proto-exchange/exchange"
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

type GrpcService interface {
	GetRates(ctx context.Context) (*pb.ExchangeRatesResponse, error)
	Exchange(ctx context.Context, user *models.Exchange, id int) (*models.Balance, error)
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
		"message": "User registered successfully",
		"token":   users,
	})
}

func UserAuthenticate(c *gin.Context, s UserService) {
	var login models.Login
	var builder strings.Builder

	if err := c.ShouldBind(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "Invalid request data format",
		})
		return
	}

	token, err := s.GetUser(c.Request.Context(), &login)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "invalid_credentials",
				"message": "Invalid email or password",
			})
		case strings.Contains(err.Error(), "token generation"):
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "token_error",
				"message": "Failed to generate authentication token",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "authentication_error",
				"message": "Failed to authenticate user",
			})
		}
		return
	}

	builder.WriteString("Bearer ")
	builder.WriteString(token)
	authHeader := builder.String()

	c.Header("Authorization", authHeader)
	c.JSON(http.StatusOK, gin.H{
		"message": "Successful",
		"token":   token,
	})
}

func GetUserBalance(c *gin.Context, s UserService) {
	token := c.GetHeader("Authorization")

	userID, err := auth.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	balance, err := s.GetBalanceUser(context.Background(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get balance"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"balance": gin.H{
			"EUR": float64(balance.EUR) / 100,
			"RUB": float64(balance.RUB) / 100,
			"USD": float64(balance.USD) / 100,
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
			"EUR": float64(updatedBalance.EUR) / 100,
			"RUB": float64(updatedBalance.RUB) / 100,
			"USD": float64(updatedBalance.USD) / 100,
		},
	})
}

func GetExchangeRates(c *gin.Context, g GrpcService) {
	rates, err := g.GetRates(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get balance"})
	}
	c.JSON(http.StatusOK, gin.H{
		"rates": rates.Rates,
	})
}

func PerformExchange(c *gin.Context, g GrpcService) {
	var exchange models.Exchange

	if err := c.ShouldBind(&exchange); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request data"})
	}
	token := c.GetHeader("Authorization")

	userID, _ := auth.ParseToken(token)

	rates, err := g.Exchange(context.Background(), &exchange, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to cheange money"})
	}

	c.JSON(http.StatusOK, gin.H{
		"EUR": float64(rates.EUR) / 100,
		"RUB": float64(rates.RUB) / 100,
		"USD": float64(rates.USD) / 100,
	})
}
