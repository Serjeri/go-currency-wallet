package tests

import (
	"gw-currency-wallet/internal/services/auth"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

var secretKey = []byte("your-secret-key")
func TestCreateToken_Success(t *testing.T) {
	userID := 123
	tokenString, err := auth.CreateToken(userID)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	assert.NoError(t, err)
	assert.True(t, token.Valid)

	claims, ok := token.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	assert.Equal(t, float64(userID), claims["sub"])
	assert.Equal(t, "app", claims["iss"])

	exp := time.Unix(int64(claims["exp"].(float64)), 0)
	assert.WithinDuration(t, time.Now().Add(time.Minute), exp, time.Second)
}

func TestCreateToken_Expiration(t *testing.T) {
	userID := 789
	tokenString, _ := auth.CreateToken(userID)

	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	}, jwt.WithTimeFunc(func() time.Time {
		return time.Now().Add(2 * time.Minute)
	}))

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token is expired")
}

func TestCreateToken_InvalidSecret(t *testing.T) {
	userID := 100
	tokenString, _ := auth.CreateToken(userID)

	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("wrong-secret-key"), nil
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "signature is invalid")
}

func TestCreateToken_EdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		userID  int
		wantErr bool
	}{
		{"Zero ID", 0, false},
		{"Negative ID", -1, false},
		{"Max Int", 1<<31 - 1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenString, err := auth.CreateToken(tt.userID)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, tokenString)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, tokenString)
			}
		})
	}
}
