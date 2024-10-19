package dao

import (
	"MetaGallery-Cloud-backend/config"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DataBase *gorm.DB

func DataBaseInit() *gorm.DB {
	DBHost, DBPort, DBUser, DBPassword, DBName, Err := config.GetDBEnv()
	if Err != nil {
		log.Fatalf("Error loading .env file")
	}
	DSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		DBUser, DBPassword, DBHost, DBPort, DBName)
	var Err2 error
	DataBase, Err2 = gorm.Open(mysql.Open(DSN), &gorm.Config{})
	if Err2 != nil {
		log.Fatalf("Error connecting database")
	}

	return DataBase
}

// init 函数在创建包的时候执行
func init() {
	DataBaseInit()
}
