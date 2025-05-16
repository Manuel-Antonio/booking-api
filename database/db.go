package database

import (
	"fmt"
	"strings"

	"booking-api/config"
	"booking-api/models"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase(cfg config.Config) error {
	var err error

	dsn := cfg.DBUrl
	if strings.HasPrefix(dsn, "postgres://") || strings.HasPrefix(dsn, "postgresql://") {
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
	} else {
		DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
	}

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	fmt.Println("📦 Conectado a la base de datos correctamente")

	if err := DB.AutoMigrate(&models.User{}, &models.Reservation{}); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	fmt.Println("✅ Migraciones completadas")
	return nil
}
