package models

import (
	"MetaGallery-Cloud-backend/dao"
	"MetaGallery-Cloud-backend/services"
	"fmt"
	"time"
)

var DataBase = dao.DataBase

type User_Data struct {
	Account       string    `gorm:"type:varchar(20); not null; primarykey" json:"account" binding:"required"`
	UserName      string    `gorm:"type:varchar(20); not null;" json:"username" binding:"required"`
	Password      string    `gorm:"type:varchar(255); not null;" json:"password" binding:"required"`
	Brief_Intro   string    `gorm:"type:text;" json:"Brief_Intro"`
	Profile_Photo string    `gorm:"type:text;" json:"Profile_Photo"`
	CreatedAt     time.Time // 创建时间（由GORM自动管理）
	UpdatedAt     time.Time // 最后一次更新时间（由GORM自动管理）
}

// 由账号密码创建账号信息
func CreateAccount(Account string, Password string) {
	ProfilePhotoURL, Err := services.GetAvatarUrl(Account)
	var User User_Data
	if Err != nil {
		fmt.Println(Err)
		User = User_Data{Account: Account, Password: Password, Brief_Intro: "这个人很懒，什么都没有写"}
	} else {
		User = User_Data{Account: Account, Password: Password, Brief_Intro: "这个人很懒，什么都没有写", Profile_Photo: ProfilePhotoURL}
	}

	DataBase.Create(&User)

}

// 由结构体User_Data创建账号信息
func CreateUserData(UserData User_Data) {

	DataBase.Create(&UserData)

}

// 更新账号的密码
func UpdatePassword(Account string, OriPassword string, NewPassword string) {

	User := User_Data{Account: Account, Password: OriPassword}

	DataBase.Model(&User).Where("account = ?", User.Account).Updates(User_Data{Account: Account, Password: NewPassword})

}

// 获取账号的密码
func GetPassword(Account string) string {

	var UserData User_Data

	DataBase.Model(&UserData).Where("account= ? ", Account).Find(&UserData)

	return UserData.Password
}

// 获取账号相关信息
func GetUserData(Account string) User_Data {

	var UserData User_Data

	DataBase.Where("account= ? ", Account).Find(&UserData)

	return UserData
}

// 更新账号相关信息，不更新密码
func UpdateUserData(Account string, NewUserData User_Data) {
	UserData := User_Data{Account: Account}

	if NewUserData.Brief_Intro == "" {
		NewUserData.Brief_Intro = "这个人很懒，什么都没有写"
	}

	DataBase.Model(&UserData).Where("account = ?", Account).Updates(User_Data{
		Account:       Account,
		UserName:      NewUserData.UserName,
		Brief_Intro:   NewUserData.Brief_Intro,
		Profile_Photo: NewUserData.Profile_Photo})

}

// 实时更新数据库表结构
func init() {
	DataBase.AutoMigrate(&User_Data{})
}
