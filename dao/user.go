package dao

import (
	"MetaGallery-Cloud-backend/models"
)

func CreateAccount(Account string, Password string) {

	User := models.User_Data{Account: Account, Password: Password}

	DataBase.Create(&User)

}

func CreateUserData(UserData models.User_Data) {

	DataBase.Create(&UserData)

}

func UpdatePassword(Account string, OriPassword string, NewPassword string) {

	User := models.User_Data{Account: Account, Password: OriPassword}

	DataBase.Model(&User).Where("account = ?", User.Account).Updates(models.User_Data{Account: Account, Password: NewPassword})

}

func GetPassword(Account string) string {

	var UserData models.User_Data

	DataBase.Model(&UserData).Where("account= ? ", Account).Find(&UserData) 

	return UserData.Password
}

func GetUserData(Account string) models.User_Data {

	var UserData models.User_Data

	DataBase.Where("account= ? ", Account).Find(&UserData)

	return UserData
}

func UpdateUserData(Account string, NewUserData models.User_Data) {

	DataBase.Model(&NewUserData).Where("account = ?", Account).Updates(models.User_Data{Account: Account, Brief_Intro: NewUserData.Brief_Intro, Profile_Photo: NewUserData.Profile_Photo})

}
