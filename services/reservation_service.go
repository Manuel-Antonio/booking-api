package services

import (
	"booking-api/models"
	"booking-api/repositories"
	"errors"
	"fmt"
	"time"
)

type ReservationService interface {
	CreateReservation(reservation *models.Reservation) error
	GetReservationsByDate(date string) ([]models.Reservation, error)
}

type reservationService struct {
	repo repositories.ReservationRepository
}

func NewReservationService(repo repositories.ReservationRepository) ReservationService {
	return &reservationService{repo: repo}
}

func (s *reservationService) CreateReservation(res *models.Reservation) error {

	startTime, err := time.Parse("15:04", res.StartTime)
	if err != nil {
		return errors.New("invalid start_time format (expected HH:MM)")
	}
	endTime, err := time.Parse("15:04", res.EndTime)
	if err != nil {
		return errors.New("invalid end_time format (expected HH:MM)")
	}

	opening, _ := time.Parse("15:04", "09:00")
	closing, _ := time.Parse("15:04", "18:00")

	if startTime.Before(opening) || endTime.After(closing) {
		return errors.New("reservation must be between 09:00 and 18:00")
	}

	if !endTime.After(startTime) {
		return errors.New("end_time must be after start_time")
	}

	duration := endTime.Sub(startTime)
	if duration < time.Hour {
		return errors.New("minimum reservation duration is 1 hour")
	}

	overlapping, err := s.repo.FindOverlapping(res.Date, res.StartTime, res.EndTime)
	if err != nil {
		return fmt.Errorf("failed to check overlaps: %v", err)
	}
	if len(overlapping) > 0 {
		return errors.New("reservation time overlaps with an existing reservation")
	}

	return s.repo.Create(res)
}

func (s *reservationService) GetReservationsByDate(date string) ([]models.Reservation, error) {
	return s.repo.FindByDate(date)
}
