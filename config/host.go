package config

import (
	"os"

	"github.com/joho/godotenv"
)

var HOST_URL string

func init() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	HOST_URL = os.Getenv("HOST_URL")
}
