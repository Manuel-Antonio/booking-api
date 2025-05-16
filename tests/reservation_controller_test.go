package tests

import (
	"booking-api/controllers"
	"booking-api/models"
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockReservationService struct {
	mock.Mock
}

func (m *MockReservationService) CreateReservation(reservation *models.Reservation) error {
	return m.Called(reservation).Error(0)
}

func (m *MockReservationService) GetReservationsByDate(date string) ([]models.Reservation, error) {
	args := m.Called(date)
	return args.Get(0).([]models.Reservation), args.Error(1)
}

func setupTestFiber() *fiber.App {
	return fiber.New()
}

func createTestToken(userID uint) *jwt.Token {
	claims := jwt.MapClaims{
		"user_id": float64(userID),
		"exp":     time.Now().Add(time.Hour).Unix(),
	}
	return &jwt.Token{
		Claims: claims,
		Method: jwt.SigningMethodHS256,
		Valid:  true,
	}
}

func TestReservationController_CreateReservation_Success(t *testing.T) {
	app := setupTestFiber()
	mockService := new(MockReservationService)
	controller := controllers.NewReservationController(mockService)

	app.Post("/reservations", func(c *fiber.Ctx) error {
		c.Locals("user", createTestToken(123))
		return controller.CreateReservation(c)
	})

	inputDate := time.Now().Format("2006-01-02")
	input := map[string]string{
		"date":       inputDate,
		"start_time": "12:00",
		"end_time":   "13:00",
	}
	body, _ := json.Marshal(input)

	mockService.On("GetReservationsByDate", inputDate).
		Return([]models.Reservation{}, nil)

	mockService.
		On("CreateReservation", mock.MatchedBy(func(r *models.Reservation) bool {
			return r.UserID == 123 &&
				r.StartTime == input["start_time"] &&
				r.EndTime == input["end_time"]
		})).
		Return(nil)

	req := httptest.NewRequest("POST", "/reservations", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestReservationController_CreateReservation_InvalidTokenClaims(t *testing.T) {
	app := setupTestFiber()
	mockService := new(MockReservationService)
	controller := controllers.NewReservationController(mockService)

	app.Post("/reservations", func(c *fiber.Ctx) error {
		invalid := &jwt.Token{
			Claims: jwt.MapClaims{"foo": "bar"},
			Valid:  false,
		}
		c.Locals("user", invalid)
		return controller.CreateReservation(c)
	})

	req := httptest.NewRequest("POST", "/reservations", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestReservationController_CreateReservation_ServiceError(t *testing.T) {
	app := setupTestFiber()
	mockService := new(MockReservationService)
	controller := controllers.NewReservationController(mockService)

	app.Post("/reservations", func(c *fiber.Ctx) error {
		c.Locals("user", createTestToken(456))
		return controller.CreateReservation(c)
	})

	inputDate := time.Now().Format("2006-01-02")
	input := map[string]string{
		"date":       inputDate,
		"start_time": "14:00",
		"end_time":   "15:00",
	}
	body, _ := json.Marshal(input)

	mockService.On("GetReservationsByDate", inputDate).
		Return([]models.Reservation{}, nil)
	mockService.On("CreateReservation", mock.Anything).
		Return(errors.New("service error"))

	req := httptest.NewRequest("POST", "/reservations", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	mockService.AssertExpectations(t)
}
