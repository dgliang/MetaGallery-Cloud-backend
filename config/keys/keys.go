package config

import (
	"os"

	"github.com/joho/godotenv"
)

var SSL_CRT_PATH string
var SSL_KEY_PATH string

func init() {
	godotenv.Load()
	SSL_CRT_PATH = os.Getenv("SSL_CRT_PATH")
	SSL_KEY_PATH = os.Getenv("SSL_KEY_PATH")
}
