package config

import (
	"os"

	"github.com/joho/godotenv"
)

var JWT_SECRET_KEY string

func init() {
	godotenv.Load()
	JWT_SECRET_KEY = os.Getenv("JWT_SECRET_KEY")
}
