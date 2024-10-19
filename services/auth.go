package services

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"os"
	"time"
)

var secretKey = ""

func init() {
	godotenv.Load()
	secretKey = os.Getenv("JWT_SECRET_KEY")
}

func GenerateToken(payload interface{}) (string, error) {
	claims := jwt.MapClaims{
		"payload": payload,
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}
