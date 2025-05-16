package main

import (
	"booking-api/config"
	"booking-api/controllers"
	"booking-api/database"
	"booking-api/middlewares"
	"booking-api/repositories"
	"booking-api/services"
	"booking-api/utils"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, usando variables de entorno del sistema")
	}

	secret := os.Getenv("JWT_SECRET")
	utils.SetJWTSecret(secret)

	cfg := config.LoadConfig()

	if err := database.ConnectDatabase(cfg); err != nil {
		log.Fatalf("❌ Error al conectar la base de datos: %v", err)
	}

	userRepo := repositories.NewUserRepository(database.DB)
	reservationRepo := repositories.NewReservationRepository(database.DB)

	authService := services.NewAuthService(userRepo)
	reservationService := services.NewReservationService(reservationRepo)

	authController := controllers.NewAuthController(authService)
	reservationController := controllers.NewReservationController(reservationService)

	app := fiber.New()

	app.Post("/register", authController.Register)
	app.Post("/login", authController.Login)

	app.Use("/reservations", jwtware.New(jwtware.Config{
		SigningKey:   []byte(cfg.JWTSecret),
		ContextKey:   "user",
		ErrorHandler: middlewares.JWTError,
	}))

	protected := app.Group("/reservations", middlewares.Protected())

	protected.Post("/", middlewares.Protected(), reservationController.CreateReservation)
	protected.Get("/", middlewares.Protected(), reservationController.GetReservationsByDate)

	port := cfg.Port
	if port == "" {
		port = "3000"
	}

	fmt.Printf("Starting server on port %s...\n", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
