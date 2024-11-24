package services

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

var hostUrl = ""

func init() {
	godotenv.Load(".env")
	hostUrl = os.Getenv("HOST_URL")
}

func IsValidAccount(s string) bool {
	re := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9]{4,9}$`)
	return re.MatchString(s)
}

func IsValidPassword(s string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9_@#$%.?]{6,}$`)
	return re.MatchString(s)
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
	if len(account) < 3 {
		return "", errors.New("account is too short, less than 3 characters")
	}

	prefix := account[:3]
	charSet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	rand.Seed(time.Now().UnixNano())
	remainLen := 10 - len("MGC") - len(prefix)
	randomSuffix := make([]byte, remainLen)
	for i := range randomSuffix {
		randomSuffix[i] = charSet[rand.Intn(len(charSet))]
	}

	userName := "MGC" + prefix + string(randomSuffix)
	return userName, nil
}
