package controllers

import (
	"MetaGallery-Cloud-backend/models"
	"MetaGallery-Cloud-backend/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SearchController struct{}

func (s SearchController) SearchFilesAndFolders(c *gin.Context) {
	account := c.Query("account")
	parentFolder := c.Query("parent_folder")
	keyword := c.Query("keyword")

	if account == "" || parentFolder == "" || keyword == "" {
		ReturnError(c, "FAILED", "提供的参数不完整")
		return
	}

	userId, err := models.GetUserID(account)
	if err != nil || userId == 0 {
		ReturnError(c, "FAILED", "获取用户 ID 失败，用户不存在")
		return
	}

	parentFolderId, _ := strconv.ParseUint(parentFolder, 10, 64)
	folderData, err := models.GetFolderDataByID(uint(parentFolderId))
	if err != nil || folderData.ID == 0 {
		ReturnError(c, "FAILED", "获取文件夹信息失败，父文件夹不存在")
		return
	}

	// 根据参数查询文件和文件夹
	res, err := services.SearchFilesAndFolders(userId, folderData.Path, keyword)
	if err != nil {
		ReturnServerError(c, "查询失败"+err.Error())
		return
	}

	ReturnSuccess(c, "SUCCESS", "", res)
}

func (s SearchController) SearchBinFilesAndFolders(c *gin.Context) {
	account := c.Query("account")
	keyword := c.Query("keyword")

	if account == "" || keyword == "" {
		ReturnError(c, "FAILED", "提供的参数不完整")
		return
	}

	userId, err := models.GetUserID(account)
	if err != nil || userId == 0 {
		ReturnError(c, "FAILED", "获取用户 ID 失败，用户不存在")
		return
	}

	// 根据参数查询回收站的文件和文件夹
	res, err := services.SearchBinFilesAndFolders(userId, keyword)
	if err != nil {
		ReturnServerError(c, "查询失败"+err.Error())
		return
	}
	ReturnSuccess(c, "SUCCESS", "", res)
}
