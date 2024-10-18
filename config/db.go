package config

import (
	"os"

	"github.com/joho/godotenv"
)

func Getdb() (string, string, string, string, string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", "", "", "", "", err
	}
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	return dbHost, dbPort, dbUser, dbPassword, dbName, nil
}
