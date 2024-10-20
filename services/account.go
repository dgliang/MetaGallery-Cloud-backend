package services

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

var hostUrl = ""

func init() {
	godotenv.Load(".env")
	hostUrl = os.Getenv("HOST_URL")
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func GetAvatarUrl(account string) (string, error) {
	if account == "" {
		return "", errors.New("account is empty")
	}

	firstLetter := strings.ToUpper(string(account[0]))
	avatarUrl := fmt.Sprintf("%s/resources/img/%s.png", hostUrl, firstLetter)
	return avatarUrl, nil
}

func RandomUsername(account string) (string, error) {
	if account == "" {
		return "", errors.New("account is empty")
	}

	firstLetter := strings.ToUpper(string(account[0]))
	avatarUrl := fmt.Sprintf("%s/resources/img/%s.png", hostUrl, firstLetter)
	return avatarUrl, nil
}
