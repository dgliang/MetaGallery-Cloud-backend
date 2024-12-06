package models

import (
	"MetaGallery-Cloud-backend/dao"
	"errors"
	"time"
)

var DataBase = dao.DataBase

type UserData struct {
	ID           uint      `gorm:"primarykey;not null;autoIncrement"`
	Account      string    `gorm:"type:varchar(20); not null; primarykey" json:"account" binding:"required"`
	UserName     string    `gorm:"type:varchar(20); not null;" json:"username" binding:"required"`
	Password     string    `gorm:"type:varchar(255); not null;" json:"password" binding:"required"`
	BriefIntro   string    `gorm:"type:text;" json:"BriefIntro"`
	ProfilePhoto string    `gorm:"type:text;" json:"ProfilePhoto"`
	CreatedAt    time.Time // 创建时间（由GORM自动管理）
	UpdatedAt    time.Time // 最后一次更新时间（由GORM自动管理）

	//外键约束
	FileData   []FileData   `gorm:"foreignKey:BelongTo"`
	FolderData []FolderData `gorm:"foreignKey:BelongTo"`
}

// 由账号密码创建账号信息
func CreateAccount(Account, Password, avatar, userName string) (uint, error) {
	// var userData UserData
	userData := UserData{
		Account:      Account,
		Password:     Password,
		BriefIntro:   "这个人很懒，什么都没有写",
		ProfilePhoto: avatar,
		UserName:     userName,
	}

	if err := DataBase.Create(&userData).Error; err != nil {
		return 0, err
	}
	return userData.ID, nil
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

func GetUserID(account string) (uint, error) {
	if account == "" {
		return 0, errors.New("数据库查询 UserID 时账号不能为空")
	}

	var userData UserData
	DataBase.Where("account= ? ", account).Find(&userData)
	return userData.ID, nil
}

func GetUserAccountById(id uint) (string, error) {
	if id == 0 {
		return "", errors.New("数据库查询 Account 时 ID 不能为空")
	}

	var userData UserData
	DataBase.Where("id= ? ", id).Find(&userData)
	return userData.Account, nil
}
