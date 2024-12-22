package services

import (
	"MetaGallery-Cloud-backend/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(payload interface{}) (string, error) {
	claims := jwt.MapClaims{
		"payload": payload,
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.JWT_SECRET_KEY))

	// token 加上 Bearer 前缀
	tokenString = "Bearer " + tokenString
	return tokenString, err
}
