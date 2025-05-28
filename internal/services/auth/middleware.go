package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthenticateMiddleware(c *gin.Context) {

	tokenString, err := c.Cookie("auth_token")
	if err != nil {
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
            "error": "Authentication required",
            "redirect": "/login",
        })
        return
    }

	token, err := verifyToken(tokenString)
    if err != nil {
    	c.SetCookie("auth_token", "", -1, "/", "", false, true)
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
            "error": "Invalid or expired token",
            "redirect": "/login",
        })
        return
        }

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
            c.Set("sub", claims)
        }
	c.Next()
}
