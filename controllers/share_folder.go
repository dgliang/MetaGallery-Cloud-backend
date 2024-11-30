package controllers

import (
	"MetaGallery-Cloud-backend/models"
	"MetaGallery-Cloud-backend/services"
	"fmt"

	"github.com/gin-gonic/gin"
)

type FolderShareController struct{}

type shareFolderRequest struct {
	Account  string `json:"account" binding:"required"`
	FolderId uint   `json:"folder_id" binding:"required"`
	IsShared int    `json:"is_shared" binding:"required"` // 1: not shared, 2: shared
	Intro    string `json:"intro"`
}

func (s FolderShareController) SetFolderShared(c *gin.Context) {
	var req shareFolderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, "FAILED", "提供的信息不全。解析 JSON Request："+err.Error())
		return
	}

	userId, err := models.GetUserID(req.Account)
	if err != nil {
		ReturnServerError(c, "获取 GetUserID: "+err.Error())
		return
	}
	if userId == 0 {
		ReturnError(c, "FAILED", "用户不存在")
		return
	}

	folderData, err := models.GetFolderDataByID(req.FolderId)
	if err != nil || folderData.ID == 0 {
		ReturnError(c, "FAILED", "文件夹不存在")
		return
	}

	var shareStatus bool
	if req.IsShared == 1 {
		shareStatus = false
	} else if req.IsShared == 2 {
		shareStatus = true

		// 判断此时 intro 是否为空
		if req.Intro == "" {
			ReturnError(c, "FAILED", "分享文件夹时，必须提供简介")
			return
		}
	} else {
		ReturnError(c, "FAILED", "is_shared 参数不正确，只能取值 1 或 2")
		return
	}

	err = services.SetFolderShareState(userId, folderData.ID, shareStatus, req.Intro)
	if err != nil {
		ReturnServerError(c, "SetFolderShareState: "+err.Error())
		return
	}

	ReturnSuccess(c, "SUCCESS", fmt.Sprintf("设置文件夹分享状态成功为%v", shareStatus))
}
