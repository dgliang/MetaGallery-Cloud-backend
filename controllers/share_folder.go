package controllers

import (
	"MetaGallery-Cloud-backend/config"
	"MetaGallery-Cloud-backend/models"
	"MetaGallery-Cloud-backend/services"
	"fmt"
	"log"
	"path"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type FolderShareController struct{}

func (s FolderShareController) SetFolderShared(c *gin.Context) {
	account := c.DefaultPostForm("account", "")
	folderIdStr := c.DefaultPostForm("folder_id", "0")
	sharedName := c.DefaultPostForm("shared_name", "")
	intro := c.DefaultPostForm("intro", "")
	folderId, _ := strconv.ParseUint(folderIdStr, 10, 64)

	if account == "" || sharedName == "" || intro == "" || folderId == 0 {
		ReturnError(c, "FAILED", "提供的信息不全")
		return
	}

	userId, err := models.GetUserID(account)
	if err != nil {
		ReturnServerError(c, "获取 GetUserID: "+err.Error())
		return
	}
	if userId == 0 {
		ReturnError(c, "FAILED", "用户不存在")
		return
	}

	folderData, err := models.GetFolderDataByID(uint(folderId))
	if err != nil || folderData.ID == 0 {
		ReturnError(c, "FAILED", "文件夹不存在")
		return
	}

	// 判断共享文件夹是否已经存在
	sharedFolder, _ := services.GetSharedFolderByOwnerAndName(userId, sharedName)
	if sharedFolder.ID != 0 {
		ReturnError(c, "FAILED", "共享文件夹已经存在")
		return
	}

	defaultUrl := config.HostURL + "/resources/img/cover_img/default.png"
	file, err := c.FormFile("cover_img")
	var savedUrl string

	if err != nil {
		savedUrl = defaultUrl
		log.Print(savedUrl)
	} else {
		timestamp := time.Now().Unix()
		fileName := fmt.Sprintf("%d_%s", timestamp, path.Base(file.Filename))
		filePath := path.Join("./resources/img/cover_img", fileName)

		if err := c.SaveUploadedFile(file, filePath); err != nil {
			ReturnServerError(c, "保存封面图片文件失败")
			return
		}

		savedUrl = config.HostURL + "/resources/img/cover_img/" + fileName
		log.Print(savedUrl)
	}

	err = services.SetFolderShareState(userId, folderData.ID, true, sharedName, intro, savedUrl)
	if err != nil {
		ReturnServerError(c, "SetFolderShareState: "+err.Error())
		return
	}

	ReturnSuccess(c, "SUCCESS", fmt.Sprintf("设置文件夹分享状态成功为%v", true))
}

func (s FolderShareController) SetFolderUnShared(c *gin.Context) {
	account := c.DefaultPostForm("account", "")
	sharedName := c.DefaultPostForm("shared_name", "")

	if account == "" || sharedName == "" {
		ReturnError(c, "FAILED", "提供的信息不全")
		return
	}

	userId, err := models.GetUserID(account)
	if err != nil || userId == 0 {
		ReturnError(c, "FAILED", "用户不存在")
		return
	}

	// 判断收藏文件夹是否存在
	sharedFolder, err := services.GetSharedFolderByOwnerAndName(userId, sharedName)
	if err != nil || sharedFolder.ID == 0 {
		ReturnError(c, "FAILED", "收藏文件夹不存在")
		return
	}

	err = services.SetFolderShareState(userId, sharedFolder.FolderID, false)
	if err != nil {
		ReturnServerError(c, "SetFolderShareState: "+err.Error())
		return
	}

	ReturnSuccess(c, "SUCCESS", fmt.Sprintf("设置文件夹分享状态成功为%v", false))
}
