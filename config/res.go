package config

import (
	"os"

	"github.com/joho/godotenv"
)

var FILE_RES_PATH string
var CACHE_RES_PATH string

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	FILE_RES_PATH = os.Getenv("FILE_DIR_PATH")
	CACHE_RES_PATH = os.Getenv("LOCAL_CACHE_PATH")
}
