package models

import (
	"MetaGallery-Cloud-backend/dao"
	"time"
)

var DataBase = dao.DataBase

type User_Data struct {
	Account       string    `gorm:"type:varchar(20); not null; primarykey" json:"account" binding:"required"`
	Password      string    `gorm:"type:varchar(255); not null;" json:"password" binding:"required"`
	Brief_Intro   string    `gorm:"type:text;" json:"Brief_Intro"`
	Profile_Photo string    `gorm:"type:text;" json:"Profile_Photo"`
	CreatedAt     time.Time // 创建时间（由GORM自动管理）
	UpdatedAt     time.Time // 最后一次更新时间（由GORM自动管理）
}

func CreateAccount(Account string, Password string) {

	User := User_Data{Account: Account, Password: Password}

	DataBase.Create(&User)

}

func CreateUserData(UserData User_Data) {

	DataBase.Create(&UserData)

}

func UpdatePassword(Account string, OriPassword string, NewPassword string) {

	User := User_Data{Account: Account, Password: OriPassword}

	DataBase.Model(&User).Where("account = ?", User.Account).Updates(User_Data{Account: Account, Password: NewPassword})

}

func GetPassword(Account string) string {

	var UserData User_Data

	DataBase.Model(&UserData).Where("account= ? ", Account).Find(&UserData)

	return UserData.Password
}

func GetUserData(Account string) User_Data {

	var UserData User_Data

	DataBase.Where("account= ? ", Account).Find(&UserData)

	return UserData
}

func UpdateUserData(Account string, NewUserData User_Data) {

	DataBase.Model(&NewUserData).Where("account = ?", Account).Updates(User_Data{Account: Account, Brief_Intro: NewUserData.Brief_Intro, Profile_Photo: NewUserData.Profile_Photo})

}

// 实时更新数据库表结构
func init() {
	DataBase.AutoMigrate(&User_Data{})
}
