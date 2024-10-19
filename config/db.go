package config

import (
	"os"

	"github.com/joho/godotenv"
)

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
