package controllers

import (
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

	// 数据库创建用户账号
	_, err := services.HashPassword(password)
	if err != nil {
		log.Println(err)
		ReturnServerError(c, "服务器加密密码失败")
		return
	}
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

	// 获取用户密码，并验证密码

	log.Printf("from %s 登录 %s %s\n", c.Request.Host, account, password)
	ReturnSuccess(c, "SUCCESS", "", "")
}

func (u UerController) GetUserInfo(c *gin.Context) {

}

func (u UerController) UpdateUserPassword(c *gin.Context) {

}

func (u UerController) UpdateUserInfo(c *gin.Context) {

}
