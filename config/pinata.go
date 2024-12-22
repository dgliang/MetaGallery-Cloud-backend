package config

import (
	"os"

	"github.com/joho/godotenv"
)

var PINATA_JWT string
var PINATA_HOST_URL string
var PINATA_GATEWAY_KEY string

func init() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	PINATA_JWT = os.Getenv("PINATA_JWT")
	PINATA_HOST_URL = os.Getenv("PINATA_HOST_URL")
	PINATA_GATEWAY_KEY = os.Getenv("PINATA_GATEWAY_KEY")
}
