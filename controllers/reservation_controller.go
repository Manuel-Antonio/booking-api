package controllers

import (
	"booking-api/models"
	"booking-api/services"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type ReservationController struct {
	Service services.ReservationService
}

func NewReservationController(service services.ReservationService) *ReservationController {
	return &ReservationController{Service: service}
}

func (rc *ReservationController) CreateReservation(c *fiber.Ctx) error {

	userToken := c.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userIDf, ok := claims["user_id"].(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
	}
	userID := uint(userIDf)

	var input struct {
		Date      string `json:"date"`
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if _, err := time.Parse("2006-01-02", input.Date); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid date format"})
	}

	start, err := time.Parse("15:04", input.StartTime)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid start_time format"})
	}
	end, err := time.Parse("15:04", input.EndTime)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid end_time format"})
	}

	if start.Equal(end) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf(
				"The reservation start time (%s) and end time (%s) cannot be the same—duration must be at least 1 hour.",
				input.StartTime, input.EndTime,
			),
		})
	}

	if !start.Before(end) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf(
				"The reservation start time (%s) must be earlier than the end time (%s). For example, 10:00–11:00 is valid.",
				input.StartTime, input.EndTime,
			),
		})
	}

	open, _ := time.Parse("15:04", "09:00")
	close, _ := time.Parse("15:04", "18:00")
	if start.Before(open) || end.After(close) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Reservation must be between 09:00 and 18:00",
		})
	}

	if end.Sub(start) < time.Hour {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Minimum reservation duration is 1 hour",
		})
	}

	existing, err := rc.Service.GetReservationsByDate(input.Date)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	for _, ex := range existing {
		exStart, _ := time.Parse("15:04", ex.StartTime)
		exEnd, _ := time.Parse("15:04", ex.EndTime)
		if start.Before(exEnd) && end.After(exStart) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Reservation time overlaps with existing reservation",
			})
		}
	}

	reservation := models.Reservation{
		UserID:    userID,
		Date:      input.Date,
		StartTime: start.Format("15:04"),
		EndTime:   end.Format("15:04"),
	}
	if err := rc.Service.CreateReservation(&reservation); err != nil {

		if strings.Contains(err.Error(), "foreign key") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "The user does not exist or is not authorized to create reservations.",
			})
		}

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Could not create the reservation",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(reservation)
}

func (rc *ReservationController) GetReservationsByDate(c *fiber.Ctx) error {
	date := c.Query("date")
	if date == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "date query param is required"})
	}

	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid date format"})
	}

	reservations, err := rc.Service.GetReservationsByDate(date)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(reservations)
}
