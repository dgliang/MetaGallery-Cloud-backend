package config

import (
	"os"

	"github.com/joho/godotenv"
)

var PinataJWT string
var FileResPath string

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	PinataJWT = os.Getenv("PINATA_JWT")
	FileResPath = os.Getenv("FILE_DIR_PATH")
}

func GetDBEnv() (string, string, string, string, string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", "", "", "", "", err
	}
	DBHost := os.Getenv("DB_HOST")
	DBPort := os.Getenv("DB_PORT")
	DBUser := os.Getenv("DB_USER")
	DBPassword := os.Getenv("DB_PASSWORD")
	DBName := os.Getenv("DB_NAME")

	return DBHost, DBPort, DBUser, DBPassword, DBName, nil
}
