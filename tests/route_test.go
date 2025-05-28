package tests

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gw-currency-wallet/internal/models"
	"gw-currency-wallet/internal/services/handlers"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPingRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}

func TestRegisterRoute(t *testing.T) {
	// 1. Настройка тестового окружения
	r := gin.Default()
	mockRepo := &mocks.UserRepository{} // Мок репозитория
	userHandler := &handlers.Client{Repository: mockRepo}
	Routers(r, userHandler)

	// 2. Тестовые данные
	testUser := models.User{
		Name:     "testuser",
		Password: "validPassword123",
		Email:    "test@example.com",
	}

	// 3. Настройка ожиданий для мока
	mockRepo.On("RegistrUser", mock.Anything, mock.AnythingOfType("models.User")).
		Return(1, nil). // Ожидаем успешную регистрацию с ID=1
		Once()

	// 4. Подготовка запроса
	jsonData, _ := json.Marshal(testUser)
	req, _ := http.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// 5. Выполнение запроса
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 6. Проверки
	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, 1, response.ID)
	assert.Equal(t, testUser.Name, response.Name)
	assert.Equal(t, testUser.Email, response.Email)

	// 7. Проверка вызова мока
	mockRepo.AssertExpectations(t)
}
