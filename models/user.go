package models

import (
	"MetaGallery-Cloud-backend/dao"
	"MetaGallery-Cloud-backend/services"
	"log"
	"time"
)

var DataBase = dao.DataBase

type UserData struct {
	Account      string    `gorm:"type:varchar(20); not null; primarykey" json:"account" binding:"required"`
	UserName     string    `gorm:"type:varchar(20); not null;" json:"username" binding:"required"`
	Password     string    `gorm:"type:varchar(255); not null;" json:"password" binding:"required"`
	BriefIntro   string    `gorm:"type:text;" json:"BriefIntro"`
	ProfilePhoto string    `gorm:"type:text;" json:"ProfilePhoto"`
	CreatedAt    time.Time // 创建时间（由GORM自动管理）
	UpdatedAt    time.Time // 最后一次更新时间（由GORM自动管理）
}

// 由账号密码创建账号信息
func CreateAccount(Account string, Password string) {
	profilePhotoURL, err1 := services.GetAvatarUrl(Account)
	defaultUserName, err2 := services.RandomUsername(Account)

	var userData UserData
	if err1 != nil || err2 != nil {
		log.Printf("err1: %v, err2: %v", err1, err2)
		userData = UserData{
			Account:    Account,
			Password:   Password,
			BriefIntro: "这个人很懒，什么都没有写",
		}
	} else {
		userData = UserData{
			Account:      Account,
			Password:     Password,
			BriefIntro:   "这个人很懒，什么都没有写",
			ProfilePhoto: profilePhotoURL,
			UserName:     defaultUserName,
		}
	}

	DataBase.Create(&userData)

}

// 由结构体User_Data创建账号信息
func CreateUserData(userData UserData) {

	DataBase.Create(&userData)

}

// 更新账号的密码
func UpdatePassword(account string, oldPassword string, newPassword string) {

	userData := UserData{
		Account:  account,
		Password: oldPassword,
	}

	DataBase.Model(&userData).Where("account = ?", userData.Account).Updates(UserData{
		Account:  account,
		Password: newPassword,
	})

}

// 获取账号的密码
func GetPassword(account string) string {

	var userData UserData

	DataBase.Model(&userData).Where("account= ? ", account).Find(&userData)

	return userData.Password
}

// 获取账号相关信息
func GetUserData(account string) UserData {

	var userData UserData

	DataBase.Where("account= ? ", account).Find(&userData)

	return userData
}

// 更新账号相关信息，不更新密码
func UpdateUserData(account string, newUserData UserData) {
	userData := UserData{Account: account}

	if newUserData.BriefIntro == "" {
		newUserData.BriefIntro = "这个人很懒，什么都没有写"
	}

	DataBase.Model(&userData).Where("account = ?", account).Updates(UserData{
		Account:      account,
		UserName:     newUserData.UserName,
		BriefIntro:   newUserData.BriefIntro,
		ProfilePhoto: newUserData.ProfilePhoto,
	})

}

// 实时更新数据库表结构
func init() {
	DataBase.AutoMigrate(&UserData{})
}
