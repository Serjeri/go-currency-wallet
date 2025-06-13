package auth

import (
	"gw-currency-wallet/domain/lib/jwttoken"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthenticateMiddleware(c *gin.Context) {

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error":    "Authentication required",
			"redirect": "/login",
		})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Token is missing",
		})
		return
	}

	token, err := jwttoken.VerifyToken(tokenString)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error":    "Invalid or expired token",
			"redirect": "/login",
		})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		c.Set("sub", claims)
	}

	c.Next()
}
