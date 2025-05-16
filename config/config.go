package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	JWTSecret string
	DBUrl     string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No se pudo cargar .env, usando variables del sistema...")
	}

	return Config{
		Port:      getEnv("PORT", "3000"),
		JWTSecret: getEnv("JWT_SECRET", "defaultsecret"),
		DBUrl:     getEnv("DB_URL", "booking.db"),
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
