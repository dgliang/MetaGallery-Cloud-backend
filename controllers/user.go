package controllers

import "github.com/gin-gonic/gin"

type UerController struct{}

func (u UerController) Register(c *gin.Context) {
	account := c.DefaultPostForm("account", "")
	password := c.DefaultPostForm("password", "")
	confirmPassword := c.DefaultPostForm("confirm_password", "")

	if account == "" || password == "" || confirmPassword == "" {
		ReturnError(c, "FAILED", "提供的账号、密码、确认密码不全")
	}

	if password != confirmPassword {
		ReturnError(c, "FAILED", "提供的密码与确认密码不相同")
	}

	ReturnSuccess(c, "SUCCESS", "账号注册成功")
}

func (u UerController) Login(c *gin.Context) {
	account := c.DefaultPostForm("account", "")
	password := c.DefaultPostForm("password", "")

	if account == "" || password == "" {

	}

}
