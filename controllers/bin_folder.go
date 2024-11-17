package controllers

import (
	"MetaGallery-Cloud-backend/models"
	"MetaGallery-Cloud-backend/services"
	"fmt"
	"github.com/gin-gonic/gin"
)

type removeFolderRequest struct {
	Account  string `json:"account" binding:"required"`
	FolderId uint   `json:"folder_id" binding:"required"`
}

func (b BinController) RemoveFolder(c *gin.Context) {
	var req removeFolderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, "FAILED", "将文件夹移至回收站时信息不全"+err.Error())
		return
	}

	userId, err := models.GetUserID(req.Account)
	if err != nil {
		ReturnServerError(c, "GetUserID: "+err.Error())
		return
	}
	if userId == 0 {
		ReturnError(c, "FAILED", fmt.Sprintf("%v 不存在", req.Account))
		return
	}

	folderData, err1 := models.GetFolderDataByID(req.FolderId)
	if err1 != nil {
		ReturnServerError(c, "GetFolderDataByID: "+err1.Error())
		return
	}
	if folderData.ParentFolder == 0 {
		ReturnError(c, "FAILED", "要删除的文件不存在")
		return
	}
	if userId != folderData.BelongTo {
		ReturnError(c, "FAILED", "仅允许删除当前用户的文件夹")
		return
	}

	err = services.RemoveFolder(userId, req.FolderId)
	if err != nil {
		ReturnServerError(c, "RemoveFolder: "+err.Error())
		return
	}
	ReturnSuccess(c, "SUCCESS", "删除（成功移到回收站）")
}

type deleteFolderRequest struct {
	Account  string `json:"account" binding:"required"`
	BinId    uint   `json:"bin_id" binding:"required"`
	FolderId uint   `json:"folder_id" binding:"required"`
}

func (b BinController) DeleteFolder(c *gin.Context) {
	var req deleteFolderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, "FAILED", "将文件夹彻底删除时信息不全"+err.Error())
		return
	}

	userId, err := models.GetUserID(req.Account)
	if err != nil {
		ReturnServerError(c, "GetUserID: "+err.Error())
		return
	}
	if userId == 0 {
		ReturnError(c, "FAILED", fmt.Sprintf("%v 不存在", req.Account))
		return
	}

	folderData, err1 := models.GetBinFolderDataByID(req.FolderId)
	if err1 != nil {
		ReturnServerError(c, "GetFolderDataByID: "+err1.Error())
		return
	}
	if folderData.ParentFolder == 0 {
		ReturnError(c, "FAILED", "要删除的文件不存在")
		return
	}
	if userId != folderData.BelongTo {
		ReturnError(c, "FAILED", "仅允许删除当前用户的文件夹")
		return
	}

	err = services.DeleteFolder(userId, req.BinId, req.FolderId)
	if err != nil {
		ReturnServerError(c, "DeleteFolder: "+err.Error())
		return
	}
	ReturnSuccess(c, "SUCCESS", "删除成功")
}
