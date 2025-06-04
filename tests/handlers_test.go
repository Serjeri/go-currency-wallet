package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gw-currency-wallet/domain/models"
	"gw-currency-wallet/domain/services"
	"gw-currency-wallet/domain/transport/rest"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

type testRepo struct {
	result int
	err    error
}

func (t testRepo) Create(ctx context.Context, user *models.User) (int, error) {
	return t.result, nil
}

func (t testRepo) Get(ctx context.Context, user *models.Login) (int, error) {
	return t.result, nil
}

func (t testRepo) GetBalance(ctx context.Context, id int) (*models.Balance, error) {
	return nil, nil
}

func (t testRepo) UpdateBalance(ctx context.Context, id int, updateBalance *models.UpdateBalance, newAmount int) error {
	return nil
}

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetBalanceUser(ctx context.Context, userID int) (*models.Balance, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*models.Balance), args.Error(1)
}

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) ParseToken(tokenString string) (int, error) {
	args := m.Called(tokenString)
	return args.Int(0), args.Error(1)
}

func (m *MockUserService) UpdateBalanceUser(ctx context.Context, userID string, update *models.UpdateBalance) (*models.Balance, error) {
	args := m.Called(ctx, userID, update)
	return args.Get(0).(*models.Balance), args.Error(1)
}

func TestRegisterUser_Validation(t *testing.T) {
	repo := testRepo{}
	client := services.NewUserService(repo)
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
	client := services.NewUserService(repo)
	router := gin.Default()
	rest.Routers(router, client)

	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
		checkHeader  bool
	}{
		{
			name: "missing name/password",
			payload: map[string]interface{}{
				"password": "test123",
			},
			expectedCode: http.StatusBadRequest,
			checkHeader:  false,
		},
		{
			name: "missing password",
			payload: map[string]interface{}{
				"name": "test",
			},
			expectedCode: http.StatusBadRequest,
			checkHeader:  false,
		},
		{
			name: "valid credentials",
			payload: map[string]interface{}{
				"name":     "test",
				"password": "test123",
			},
			expectedCode: http.StatusOK,
			checkHeader:  true,
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

			if tc.checkHeader {
				authHeader := w.Header().Get("Authorization")
				assert.True(t, strings.HasPrefix(authHeader, "Bearer "))
				assert.NotEmpty(t, authHeader)
			}
		})
	}
}

func TestGetUserBalance(t *testing.T) {
	router := gin.Default()

	tests := []struct {
		name           string
		setupAuth      func() string
		mockUserID     int
		mockAuthError  error
		mockBalance    *models.Balance
		mockBalanceErr error
		expectedCode   int
		expectedBody   string
	}{
		{
			name:         "successful request",
			mockUserID:   1,
			mockBalance:  &models.Balance{EUR: 1000, RUB: 5000, USD: 3000},
			expectedCode: http.StatusOK,
			expectedBody: `{"balance":{"EUR":10,"RUB":50,"USD":30}}`,
		},
		{
			name: "invalid token",
			setupAuth: func() string {
				return "Bearer invalid_token"
			},
			mockAuthError: errors.New("invalid token"),
			expectedCode:  http.StatusUnauthorized,
			expectedBody:  `{"error":"invalid token"}`,
		},
		{
			name:           "balance not found",
			mockUserID:     1,
			mockBalanceErr: errors.New("not found"),
			expectedCode:   http.StatusInternalServerError,
			expectedBody:   `{"error":"failed to get balance"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authService := new(MockAuthService)
			userService := new(MockUserService)

			authService.On("ParseToken", mock.Anything).Return(tt.mockUserID, tt.mockAuthError)

			if tt.mockAuthError == nil {
				userService.On("GetBalanceUser", mock.Anything, tt.mockUserID).
					Return(tt.mockBalance, tt.mockBalanceErr)
			}

			router.GET("/api/v1/wallet/balance", func(c *gin.Context) {
				token := c.GetHeader("Authorization")

				userID, err := authService.ParseToken(token)
				if err != nil {
					c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
					return
				}

				balance, err := userService.GetBalanceUser(c, userID)
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
			})

			req, _ := http.NewRequest("GET", "/api/v1/wallet/balance", nil)
			if tt.setupAuth != nil {
				req.Header.Set("Authorization", tt.setupAuth())
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())

			authService.AssertExpectations(t)
			if tt.mockAuthError == nil {
				userService.AssertExpectations(t)
			}
		})
	}
}

func TestUpdateUserBalance(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		setupAuth      func() string
		requestBody    string
		mockUserID     int
		mockAuthError  error
		mockBalance    *models.Balance
		mockError      error
		expectedCode   int
		expectedBody   string
		expectAuthCall bool
	}{
		{
			name: "successful request",
			setupAuth: func() string {
				return "Bearer valid_token"
			},
			requestBody: `{
                "amount": 3305,
                "currency": "EUR",
                "status": "withdrawal"
            }`,
			mockUserID:     1,
			mockBalance:    &models.Balance{EUR: 1000, RUB: 5000, USD: 3000},
			expectedCode:   http.StatusOK,
			expectedBody:   `{"balance":{"EUR":10,"RUB":50,"USD":30}}`,
			expectAuthCall: true,
		},
		{
			name: "invalid token",
			setupAuth: func() string {
				return "Bearer invalid_token"
			},
			requestBody:    `{"amount": 100, "currency": "EUR", "status": "withdrawal"}`,
			mockAuthError:  errors.New("invalid token"),
			expectedCode:   http.StatusUnauthorized,
			expectedBody:   `{"error":"invalid token"}`,
			expectAuthCall: true,
		},
		{
			name: "invalid request data",
			setupAuth: func() string {
				return "Bearer valid_token"
			},
			requestBody:    `invalid_json`,
			mockUserID:     1,
			expectedCode:   http.StatusBadRequest,
			expectedBody:   `{"error":"Invalid request data"}`,
			expectAuthCall: false,
		},
		{
			name: "update error",
			setupAuth: func() string {
				return "Bearer valid_token"
			},
			requestBody: `{
                "amount": 3305,
                "currency": "EUR",
                "status": "withdrawal"
            }`,
			mockUserID:     1,
			mockError:      errors.New("update failed"),
			expectedCode:   http.StatusInternalServerError,
			expectedBody:   `{"error":"failed to get balance"}`,
			expectAuthCall: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authService := new(MockAuthService)
			userService := new(MockUserService)

			// Настраиваем мок только если ожидаем вызов ParseToken
			if tt.expectAuthCall {
				token := tt.setupAuth()
				authService.On("ParseToken", strings.TrimPrefix(token, "Bearer ")).
					Return(tt.mockUserID, tt.mockAuthError)
			}

			if tt.mockAuthError == nil && tt.requestBody != "invalid_json" {
				var expectedUpdate models.UpdateBalance
				err := json.Unmarshal([]byte(tt.requestBody), &expectedUpdate)
				if err != nil {
					t.Fatalf("Failed to unmarshal test request body: %v", err)
				}

				userService.On("UpdateBalanceUser", mock.Anything, strconv.Itoa(tt.mockUserID), &expectedUpdate).
					Return(tt.mockBalance, tt.mockError)
			}

			router := gin.New()
			router.PUT("/api/v1/wallet/update", func(c *gin.Context) {
				var update models.UpdateBalance

				if err := c.ShouldBindJSON(&update); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
					return
				}

				token := c.GetHeader("Authorization")
				userID, err := authService.ParseToken(strings.TrimPrefix(token, "Bearer "))
				if err != nil {
					c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
					return
				}

				balance, err := userService.UpdateBalanceUser(c, strconv.Itoa(userID), &update)
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
			})

			req, _ := http.NewRequest(http.MethodPut, "/api/v1/wallet/update", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			if tt.setupAuth != nil {
				req.Header.Set("Authorization", tt.setupAuth())
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, w.Body.String())
			}

			// Проверяем ожидания только если они должны были быть
			if tt.expectAuthCall {
				authService.AssertExpectations(t)
			}
			if tt.mockAuthError == nil && tt.requestBody != "invalid_json" {
				userService.AssertExpectations(t)
			}
		})
	}
}
