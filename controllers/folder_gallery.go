package controllers

import (
	"MetaGallery-Cloud-backend/models"
	"MetaGallery-Cloud-backend/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 获取用户自己共享的所有文件夹列表
func (s FolderShareController) GetUserSharedFolders(c *gin.Context) {
	account := c.Query("account")
	pageNumStr := c.Query("page_num")
	pageNum, _ := strconv.Atoi(pageNumStr)

	// 需要先判断 account 字段是否为空
	if account == "" {
		ReturnError(c, "FAILED", "account 字段不能为空")
		return
	}

	// 判断 page_num 是否大于 0
	if pageNum <= 0 {
		ReturnError(c, "FAILED", "page_num 必须大于 0")
		return
	}

	// 获取用户ID，判断用户是否存在
	userId, err := models.GetUserID(account)
	if err != nil || userId == 0 {
		ReturnError(c, "FAILED", "获取用户ID失败，用户不存在")
		return
	}

	// 获取用户共享的文件夹列表
	res, err := services.ListUserSharedFolders(userId, pageNum)
	if err != nil {
		ReturnServerError(c, "获取用户共享的文件夹列表失败"+err.Error())
		return
	}

	ReturnSuccess(c, "SUCCESS", "", res)
}

// 获取所有用户的共享文件夹列表
func (s FolderShareController) GetAllSharedFolders(c *gin.Context) {
	pageNumStr := c.Query("page_num")
	pageNum, _ := strconv.Atoi(pageNumStr)

	// 判断 page_num 是否大于 0
	if pageNum <= 0 {
		ReturnError(c, "FAILED", "page_num 必须大于 0")
		return
	}

	// 获取所有共享的文件夹列表
	res, err := services.ListAllSharedFolders(pageNum)
	if err != nil {
		ReturnServerError(c, "获取所有共享的文件夹列表失败"+err.Error())
		return
	}

	ReturnSuccess(c, "SUCCESS", "", res)
}