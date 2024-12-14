package config

import (
	"os"

	"github.com/joho/godotenv"
)

var FileResPath string
var CacheResPath string

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	FileResPath = os.Getenv("FILE_DIR_PATH")
	CacheResPath = os.Getenv("LOCAL_CACHE_PATH")
}
