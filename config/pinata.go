package config

import (
	"os"

	"github.com/joho/godotenv"
)

var PinataJWT string

func init() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	PinataJWT = os.Getenv("PINATA_JWT")
}
