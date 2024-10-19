package dao

import (
	"MetaGallery-Cloud-backend/config"
	"MetaGallery-Cloud-backend/models"
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
	DataBase.AutoMigrate(&models.User_Data{})

	return DataBase
}

func GetDB() *gorm.DB {
	return DataBase
}
