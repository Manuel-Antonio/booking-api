package tests

import (
	"booking-api/controllers"
	"booking-api/models"
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockAuthService) Login(email, password string) (string, error) {
	args := m.Called(email, password)
	return args.String(0), args.Error(1)
}

func setupFiber() *fiber.App {
	return fiber.New()
}

func TestAuthController_Login_Success(t *testing.T) {
	app := setupFiber()
	mockService := new(MockAuthService)
	controller := controllers.NewAuthController(mockService)

	app.Post("/login", controller.Login)

	loginInput := map[string]string{
		"email":    "test@example.com",
		"password": "secret",
	}
	body, _ := json.Marshal(loginInput)

	mockService.On("Login", "test@example.com", "secret").Return("mocked-token", nil)

	req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)

	defer resp.Body.Close()
	bodyBytes, _ := io.ReadAll(resp.Body)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	assert.Contains(t, string(bodyBytes), "mocked-token")
	mockService.AssertExpectations(t)
}

func TestAuthController_Register_Success(t *testing.T) {
	app := setupFiber()
	mockService := new(MockAuthService)
	controller := controllers.NewAuthController(mockService)

	app.Post("/register", controller.Register)

	registerInput := map[string]string{
		"name":     "John Doe",
		"email":    "john@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(registerInput)

	expectedUser := &models.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	mockService.On("Register", mock.MatchedBy(func(u *models.User) bool {
		return u.Name == expectedUser.Name &&
			u.Email == expectedUser.Email &&
			u.Password == expectedUser.Password
	})).Return(nil)

	req := httptest.NewRequest("POST", "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)

	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	assert.Contains(t, string(respBody), "token")
	mockService.AssertExpectations(t)
}
