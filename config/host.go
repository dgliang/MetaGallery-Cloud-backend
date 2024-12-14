package config

import (
	"os"

	"github.com/joho/godotenv"
)

var HostURL string

func init() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	HostURL = os.Getenv("HOST_URL")
}
