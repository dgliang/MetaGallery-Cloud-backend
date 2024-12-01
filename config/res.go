package config

import (
	"os"

	"github.com/joho/godotenv"
)

var FileResPath string

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	FileResPath = os.Getenv("FILE_DIR_PATH")
}
