package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtKey []byte

func SetJWTSecret(secret string) {
	jwtKey = []byte(secret)
}

func GenerateJWT(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
