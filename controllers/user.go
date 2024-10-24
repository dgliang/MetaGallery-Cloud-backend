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
		Name:    userData.UserName,
		Intro:   userData.BriefIntro,
		Avatar:  userData.ProfilePhoto,
	}
	log.Printf("from %s 登录 %s %s %v\n", c.Request.Host, account, password, userInfo)
	ReturnSuccess(c, "SUCCESS", "", struct {
		UserInfo UserInfo `json:"userInfo"`
		Token    string   `json:"token"`
	}{userInfo, token})
}

func (u UerController) GetUserInfo(c *gin.Context) {
	account := c.Query("account")

	if account == "" {
		log.Printf("from %s 登录提供的账号不全\n", c.Request.Host)
		ReturnError(c, "FAILED", "提供的账号不全")
		return
	}

	userData := models.GetUserData(account)
	if userData.Account == "" {
		log.Printf("from %s 提供的账号 %s 不存在\n", c.Request.Host, account)
		ReturnError(c, "FAILED", "账号不存在")
		return
	}

	userInfo := UserInfo{
		Account: userData.Account,
		Name:    userData.UserName,
		Intro:   userData.BriefIntro,
		Avatar:  userData.ProfilePhoto,
	}
	ReturnSuccess(c, "SUCCESS", "", userInfo)
}

func (u UerController) UpdateUserPassword(c *gin.Context) {
	account := c.DefaultPostForm("account", "")
	oldPassword := c.DefaultPostForm("old_password", "")
	newPassword := c.DefaultPostForm("new_password", "")
	confirmPassword := c.DefaultPostForm("confirm_password", "")

	if account == "" || oldPassword == "" || newPassword == "" || confirmPassword == "" {
		log.Printf("from %s 修改密码提供的信息不全\n", c.Request.Host)
		ReturnError(c, "FAILED", "修改密码提供的信息不全")
		return
	}

	if newPassword != confirmPassword {
		log.Printf("from %s 提供的密码与确认密码不相同\n", c.Request.Host)
		ReturnError(c, "FAILED", "提供的密码与确认密码不相同")
		return
	}

	if oldPassword == newPassword {
		ReturnError(c, "FAILED", "新密码与旧密码相同，无需修改")
		return
	}

	userPd := models.GetPassword(account)
	if userPd == "" {
		log.Printf("from %s %s 用户不存在\n", c.Request.Host, account)
		ReturnError(c, "NOT EXIST", "用户不存在")
		return
	}

	if invalid := services.VerifyPassword(userPd, oldPassword); invalid == false {
		log.Printf("from %s %s 密码错误\n", c.Request.Host, account)
		ReturnError(c, "FAILED", "修改密码失败，原密码错误")
		return
	}

	hashedPd, err := services.HashPassword(newPassword)
	if err != nil {
		log.Println(err)
		ReturnServerError(c, "服务器加密密码失败")
		return
	}

	models.UpdatePassword(account, userPd, hashedPd)
	log.Printf("from %s 用户 %s 修改密码成功 \n", c.Request.Host, account)
	ReturnSuccess(c, "SUCCESS", "密码修改成功")
}

func (u UerController) UpdateUserInfo(c *gin.Context) {
	account := c.DefaultPostForm("account", "")
	newUsername := c.DefaultPostForm("name", "")
	newIntro := c.DefaultPostForm("info", "")

	if newUsername == "" && newIntro == "" {
		log.Printf("from %s 修改用户信息未提供新用户名\n", c.Request.Host)
		ReturnError(c, "FAILED", "提供的账号、原密码不全")
		return
	}

	if newUsername == "" {
		models.UpdateUserData(account, models.UserData{BriefIntro: newIntro})
	} else {
		NewURL, err := services.GetAvatarUrl(newUsername)

		if err != nil {
			log.Printf("from %s 提供的新用户名 %s 不合法\n", c.Request.Host, newUsername)
			ReturnError(c, "FAILED", "用户名不合法")
			return
		}

		models.UpdateUserData(account, models.UserData{
			UserName:     newUsername,
			BriefIntro:   newIntro,
			ProfilePhoto: NewURL,
		})
	}
	log.Printf("from %s 用户 %s 修改资料成功 \n", c.Request.Host, account)
	ReturnSuccess(c, "SUCCESS", "修改成功")
}
