package models

import (
	"time"

	"gorm.io/gorm"
)

type User_data struct {
	Account       string    `gorm:"type:varchar(20); not null; primarykey" json:"account" binding:"required"`
	Password      string    `gorm:"type:varchar(20); not null;" json:"password" binding:"required"`
	Brief_Intro   *string   `gorm:"type:text;" json:"Brief_Intro"`
	Profile_Photo *string   `gorm:"type:text;" json:"Profile_Photo"`
	CreatedAt     time.Time // 创建时间（由GORM自动管理）
	UpdatedAt     time.Time // 最后一次更新时间（由GORM自动管理）
}

func CreateUser(Account string, Password string, db *gorm.DB) {

	user := User_data{Account: Account, Password: Password}

	db.Create(&user)

}

func CreateUser_data(user_data User_data, db *gorm.DB) {

	db.Create(&user_data)

}

func UPdatePassword(Account string, OriPassword string, NewPassword string, db *gorm.DB) {

	user := User_data{Account: Account, Password: OriPassword}

	db.Model(&user).Where("account = ?", user.Account).Updates(User_data{Account: Account, Password: NewPassword})

}

func GetPassword(Account string, db *gorm.DB) string {

	var userdata User_data

	db.Where("account= ? ", Account).Find(&userdata) // 查询数据库

	return userdata.Password
}

func GetUser_data(Account string, db *gorm.DB) User_data {

	var userdata User_data

	db.Where("account= ? ", Account).Find(&userdata)

	return userdata
}

func UpdateUser_data(Account string, new_user_data User_data, db *gorm.DB) {

	db.Model(&new_user_data).Where("account = ?", Account).Updates(User_data{Account: Account, Brief_Intro: new_user_data.Brief_Intro, Profile_Photo: new_user_data.Profile_Photo})

}
