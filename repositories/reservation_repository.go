package repositories

import (
	"booking-api/models"

	"gorm.io/gorm"
)

type ReservationRepository interface {
	Create(reservation *models.Reservation) error
	FindOverlapping(date, start, end string) ([]models.Reservation, error)
	FindByDate(date string) ([]models.Reservation, error)
}

type reservationRepository struct {
	db *gorm.DB
}

func NewReservationRepository(db *gorm.DB) ReservationRepository {
	return &reservationRepository{db: db}
}

func (r *reservationRepository) Create(reservation *models.Reservation) error {
	return r.db.Create(reservation).Error
}

func (r *reservationRepository) FindOverlapping(date, start, end string) ([]models.Reservation, error) {
	var reservations []models.Reservation
	err := r.db.Where("date = ? AND ((start_time < ? AND end_time > ?) OR (start_time >= ? AND start_time < ?))",
		date, end, start, start, end).Find(&reservations).Error
	return reservations, err
}

func (r *reservationRepository) FindByDate(date string) ([]models.Reservation, error) {
	var reservations []models.Reservation
	err := r.db.Where("date = ?", date).Find(&reservations).Error
	return reservations, err
}
