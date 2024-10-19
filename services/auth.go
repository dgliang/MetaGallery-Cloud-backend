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
	tokenString, err := token.SignedString([]byte(secretKey))

	// token 加上 Bearer 前缀
	tokenString = "Bearer " + tokenString
	return tokenString, err
}
