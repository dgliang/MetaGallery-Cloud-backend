package config

import (
	"os"

	"github.com/joho/godotenv"
)

var PinataJWT string
var PinataHostUrl string
var PinataGatewayKey string

func init() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	PinataJWT = os.Getenv("PINATA_JWT")
	PinataHostUrl = os.Getenv("PINATA_HOST_URL")
	PinataGatewayKey = os.Getenv("PINATA_GATEWAY_KEY")
}
