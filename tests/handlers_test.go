package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"gw-currency-wallet/internal/models"
	"gw-currency-wallet/internal/services/handlers"
	"gw-currency-wallet/internal/transport/rest"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type testRepo struct {
	result int
	err    error
}

func (t testRepo) GetUser(ctx context.Context, user models.Login) (int, error) {
	return t.result, nil
}

func (t testRepo) RegistrUser(ctx context.Context, user models.User) (int, error) {
	return t.result, nil
}

func TestRegisterUser_Validation(t *testing.T) {
	repo := testRepo{}
	client := handlers.NewClient(repo)
	router := gin.Default()
	rest.Routers(router, client)

	testCases := []struct {
		name         string
		payload      models.User
		expectedCode int
		expectedBody string
	}{
		{
			name: "missing name",
			payload: models.User{
				Password: "test123",
				Email:    "test@test.com",
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error": "Invalid request data"}`,
		},
		{
			name: "missing password",
			payload: models.User{
				Name:  "Test User",
				Email: "test@test.com",
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error": "Invalid request data"}`,
		},
		{
			name: "invalid email",
			payload: models.User{
				Name:     "Test User",
				Password: "test123",
				Email:    "not-an-email",
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error": "Invalid request data"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			payload, _ := json.Marshal(tc.payload)
			req, _ := http.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(payload))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedCode, w.Code)
			assert.JSONEq(t, tc.expectedBody, w.Body.String())
		})
	}
}

func TestAuthenticateUser_Validation(t *testing.T) {
	repo := testRepo{}
	client := handlers.NewClient(repo)
	router := gin.Default()
	rest.Routers(router, client)

	testCases := []struct {
		name         string
		payload      models.Login
		expectedCode int
		expectedBody string
	}{
		{
			name: "missing name",
			payload: models.Login{
				Password: "test123",
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error": "Invalid request data"}`,
		},
		{
			name: "missing password",
			payload: models.Login{
				Name: "Test User",
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error": "Invalid request data"}`,
		},
		{
			name: "number name",
			payload: models.Login{
				Name:     "gkkhffrdr",
				Password: "test123",
			},
			expectedCode: http.StatusOK,
			expectedBody: `{"message": "Successful"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			payload, _ := json.Marshal(tc.payload)
			req, _ := http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(payload))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedCode, w.Code)
			assert.JSONEq(t, tc.expectedBody, w.Body.String())
		})
	}
}
