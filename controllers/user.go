package controllers

import (
	"MetaGallery-Cloud-backend/models"
	"MetaGallery-Cloud-backend/services"
	"log"

	"github.com/gin-gonic/gin"
)

type UerController struct{}

func (u UerController) Register(c *gin.Context) {
	account := c.DefaultPostForm("account", "")
	password := c.DefaultPostForm("password", "")
	confirmPassword := c.DefaultPostForm("confirm_password", "")

	if account == "" || password == "" || confirmPassword == "" {
		log.Printf("from %s 注册提供的账号、密码、确认密码不全\n", c.Request.Host)
		ReturnError(c, "FAILED", "提供的账号、密码、确认密码不全")
		return
	}

	if password != confirmPassword {
		log.Printf("from %s 提供的密码与确认密码不相同\n", c.Request.Host)
		ReturnError(c, "FAILED", "提供的密码与确认密码不相同")
		return
	}

	// 调用数据库接口
	// 判断用户是否已经创建
	pd := models.GetPassword(account)
	if pd != "" {
		log.Printf("from %s %s 用户已经创建过了\n", c.Request.Host, account)
		ReturnError(c, "FAILED", "用户已经存在")
		return
	}

	// 数据库创建用户账号
	hashedPd, err := services.HashPassword(password)
	if err != nil {
		log.Println(err)
		ReturnServerError(c, "服务器加密密码失败")
		return
	}

	models.CreateAccount(account, hashedPd)
	log.Printf("from %s 注册 %s %s %s\n", c.Request.Host, account, password, confirmPassword)
	ReturnSuccess(c, "SUCCESS", "账号注册成功")
}

type UserInfo struct {
	Account string `json:"account"`
	Name    string `json:"name"`
	Intro   string `json:"intro"`
	Avatar  string `json:"avatar"`
}

func (u UerController) Login(c *gin.Context) {
	account := c.DefaultPostForm("account", "")
	password := c.DefaultPostForm("password", "")

	if account == "" || password == "" {
		log.Printf("from %s 登录提供的账号、密码不全\n", c.Request.Host)
		ReturnError(c, "FAILED", "提供的账号、密码不全")
		return
	}

	// 调用数据库接口，用户是否存在
	userPd := models.GetPassword(account)
	if userPd == "" {
		log.Printf("from %s %s 用户不存在\n", c.Request.Host, account)
		ReturnError(c, "NOT EXIST", "用户不存在")
		return
	}

	// log.Println(userPd)
	// 验证密码
	if invalid := services.VerifyPassword(userPd, password); invalid == false {
		log.Printf("from %s %s 密码错误\n", c.Request.Host, account)
		ReturnError(c, "FAILED", "密码错误")
		return
	}

	userData := models.GetUserData(account)
	token, err := services.GenerateToken(userData)
	if err != nil {
		log.Println(err)
		ReturnServerError(c, "生成 jwt token 失败")
		return
	}

	userInfo := UserInfo{
		Account: userData.Account,
		Name:    "name",
		Intro:   userData.Brief_Intro,
		Avatar:  userData.Profile_Photo,
	}
	log.Printf("from %s 登录 %s %s %v\n", c.Request.Host, account, password, userInfo)
	ReturnSuccess(c, "SUCCESS", "", struct {
		UserInfo UserInfo `json:"userInfo"`
		Token    string   `json:"token"`
	}{userInfo, token})
}

func (u UerController) GetUserInfo(c *gin.Context) {

}

func (u UerController) UpdateUserPassword(c *gin.Context) {

}

func (u UerController) UpdateUserInfo(c *gin.Context) {

}