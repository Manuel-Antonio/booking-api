package models

import "gorm.io/gorm"

type Reservation struct {
	gorm.Model
	UserID    uint   `json:"user_id"`
	Date      string `json:"date"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}
