package controllers

import (
	"MetaGallery-Cloud-backend/models"
	"MetaGallery-Cloud-backend/services"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type removeFolderRequest struct {
	Account  string `json:"account" binding:"required"`
	FolderId uint   `json:"folder_id" binding:"required"`
}

func (b BinController) RemoveFolder(c *gin.Context) {
	// var req removeFolderRequest
	// if err := c.ShouldBindJSON(&req); err != nil {
	// 	ReturnError(c, "FAILED", "将文件夹移至回收站时信息不全"+err.Error())
	// 	return
	// }
	req, _ := c.Get("jsondata")
	jsondata, ok := req.(map[string]interface{})
	if !ok {
		c.JSON(400, gin.H{
			"status": "ERROR",
			"msg":    "服务端获取请求体错误",
		})
		c.Abort()
		return
	}

	folderID := jsondata["folder_id"].(float64)
	uintFolderID := uint(folderID)
	// fmt.Println("Type of jsondata[\"folder_id\"]:", folderID, reflect.TypeOf(folderID))

	account := jsondata["account"].(string)
	// fmt.Println("Type of jsondata[\"account\"]:", account, reflect.TypeOf(account))

	userId, err := models.GetUserID(account)
	if err != nil {
		ReturnServerError(c, "GetUserID: "+err.Error())
		return
	}
	if userId == 0 {
		ReturnError(c, "FAILED", fmt.Sprintf("%v 不存在", account))
		return
	}

	folderData, err1 := models.GetFolderDataByID(uintFolderID)
	if err1 != nil || folderData.ID == 0 {
		ReturnError(c, "FAILED", "要删除的文件夹不存在")
		return
	}

	if userId != folderData.BelongTo {
		ReturnError(c, "FAILED", "仅允许删除当前用户的文件夹")
		return
	}

	err = services.RemoveFolder(userId, uintFolderID)
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
	// var req deleteFolderRequest
	// if err := c.ShouldBindJSON(&req); err != nil {
	// 	ReturnError(c, "FAILED", "将文件夹彻底删除时信息不全"+err.Error())
	// 	return
	// }
	req, _ := c.Get("jsondata")
	jsondata, ok := req.(map[string]interface{})
	if !ok {
		c.JSON(400, gin.H{
			"status": "ERROR",
			"msg":    "服务端获取请求体错误",
		})
		c.Abort()
		return
	}

	folderID := jsondata["folder_id"].(float64)
	uintFolderID := uint(folderID)
	// fmt.Println("Type of jsondata[\"folder_id\"]:", folderID, reflect.TypeOf(folderID))

	account := jsondata["account"].(string)
	// fmt.Println("Type of jsondata[\"account\"]:", account, reflect.TypeOf(account))

	binID := jsondata["bin_id"].(float64)
	uintBinID := uint(binID)
	// fmt.Println("Type of jsondata[\"bin_id\"]:", binID, reflect.TypeOf(binID))

	userId, err := models.GetUserID(account)
	if err != nil {
		ReturnServerError(c, "GetUserID: "+err.Error())
		return
	}
	if userId == 0 {
		ReturnError(c, "FAILED", fmt.Sprintf("%v 不存在", account))
		return
	}

	// 检查 bins 中是否有该记录
	if !services.IsFolderInBin(userId, uintBinID) {
		ReturnError(c, "FAILED", "要删除的 bin record 不在回收站中")
		return
	}

	// 检查 folder 文件夹是否存在
	folderData, err1 := models.GetBinFolderDataByID(uintFolderID)
	if err1 != nil || folderData.ID == 0 {
		ReturnError(c, "FAILED", "要删除的文件夹不存在或已经被删除")
		return
	}

	// 检查 folder、 bins 和 user 的记录是否能对应上
	if !services.CheckFolderBinAndUserRel(userId, uintFolderID, uintBinID) {
		ReturnError(c, "FAILED", "要删除的回收站记录，文件夹和用户对应不上")
		return
	}

	err = services.DeleteFolder(userId, uintBinID, uintFolderID)
	if err != nil {
		ReturnServerError(c, "DeleteFolder: "+err.Error())
		return
	}
	ReturnSuccess(c, "SUCCESS", "删除成功")
}

type FolderBinJson struct {
	FolderJson
	BinId   uint   `json:"bin_id"`
	DelTime string `json:"del_time"`
}

func (b BinController) ListBinFolder(c *gin.Context) {
	account := c.Query("account")
	if account == "" {
		ReturnError(c, "FAILED", "查看回收站中文件夹信息时 account 不全")
		return
	}

	userId, err := models.GetUserID(account)
	if err != nil {
		ReturnServerError(c, "GetUserID: "+err.Error())
		return
	}
	if userId == 0 {
		ReturnError(c, "FAILED", fmt.Sprintf("%v 不存在", account))
		return
	}

	folderData, err := services.ListBinFolders(userId)
	if err != nil {
		ReturnServerError(c, "ListBinFolders: "+err.Error())
		return
	}

	folderBinRes := matchFolderBinResJson(folderData)
	ReturnSuccess(c, "SUCCESS", "", folderBinRes)
}

func (b BinController) GetBinFolderInfo(c *gin.Context) {
	account := c.Query("account")
	binIdStr := c.Query("bin_id")
	folderIdStr := c.Query("folder_id")

	if account == "" || binIdStr == "" || folderIdStr == "" {
		ReturnError(c, "FAILED", "查看回收站中文件夹信息时 account、bin_id 或 folder_id 不全")
		return
	}

	userId, err := models.GetUserID(account)
	if err != nil || userId == 0 {
		ReturnError(c, "FAILED", fmt.Sprintf("%v 用户不存在", account))
		return
	}

	binId, _ := strconv.Atoi(binIdStr)
	folderId, _ := strconv.Atoi(folderIdStr)

	if !services.IsFolderInBin(userId, uint(binId)) {
		ReturnError(c, "FAILED", "要查看的回收站记录不存在")
		return
	}

	if !services.CheckFolderBinAndUserRel(userId, uint(folderId), uint(binId)) {
		ReturnError(c, "FAILED", "要查看的回收站记录，文件夹和用户对应不上")
		return
	}

	var binFolderInfos []services.FolderBinInfo
	binFolderInfo, err := services.GetBinFolderInfoByID(userId, uint(binId), uint(folderId))
	if err != nil {
		ReturnServerError(c, "GetBinFolderInfoByID: "+err.Error())
		return
	}
	binFolderInfos = append(binFolderInfos, binFolderInfo)

	ReturnSuccess(c, "SUCCESS", "", matchFolderBinResJson(binFolderInfos)[0])
}

func matchFolderBinResJson(folderData []services.FolderBinInfo) []FolderBinJson {
	if len(folderData) == 0 {
		return nil
	}

	var folderBinJson []FolderBinJson
	for _, folder := range folderData {
		folderBinJson = append(folderBinJson, FolderBinJson{
			FolderJson{
				ID:         folder.FolderData.ID,
				User:       folder.FolderData.BelongTo,
				FolderName: folder.FolderData.FolderName,
				ParentID:   folder.FolderData.ParentFolder,
				Path:       TrimPathPrefix(folder.FolderData.Path),
				IsFavorite: folder.FolderData.Favorite,
				IsShare:    folder.FolderData.Share,
			},
			folder.BinId,
			folder.DelTime.Format("2006-01-02 15:04:05"),
		})
	}
	return folderBinJson
}

type recoverFolderRequest struct {
	Account string `json:"account" binding:"required"`
	BinId   uint   `json:"bin_id" binding:"required"`
}

func (b BinController) RecoverBinFolder(c *gin.Context) {
	// var req recoverFolderRequest
	// if err := c.ShouldBindJSON(&req); err != nil {
	// 	ReturnError(c, "FAILED", "RecoverBinFolder: "+err.Error())
	// 	return
	// }
	req, _ := c.Get("jsondata")
	jsondata, ok := req.(map[string]interface{})
	if !ok {
		c.JSON(400, gin.H{
			"status": "ERROR",
			"msg":    "服务端获取请求体错误",
		})
		c.Abort()
		return
	}

	binID := jsondata["bin_id"].(float64)
	uintBinID := uint(binID)
	// fmt.Println("Type of jsondata[\"folder_id\"]:", folderID, reflect.TypeOf(folderID))

	account := jsondata["account"].(string)
	// fmt.Println("Type of jsondata[\"account\"]:", account, reflect.TypeOf(account))

	userId, err := models.GetUserID(account)
	if err != nil {
		ReturnServerError(c, "GetUserID: "+err.Error())
		return
	}
	if userId == 0 {
		ReturnError(c, "FAILED", fmt.Sprintf("%v 不存在", account))
		return
	}

	// 先检查文件夹是否在回收站中
	if !services.IsFolderInBin(userId, uintBinID) {
		ReturnError(c, "FAILED", fmt.Sprintf("%v 不在回收站中", uintBinID))
		return
	}

	// 检查恢复的文件夹会不会与现有文件夹产生冲突
	if services.CheckBinFolderAndFolder(userId, uintBinID) {
		ReturnError(c, "FAILED", fmt.Sprintf("恢复文件夹 %v 会导致冲突", uintBinID))
		return
	}

	// 恢复文件夹
	if err := services.RecoverBinFolder(userId, uintBinID); err != nil {
		ReturnServerError(c, "RecoverBinFolder: "+err.Error())
		return
	}
	ReturnSuccess(c, "SUCCESS", fmt.Sprintf("恢复文件夹 %v 成功", uintBinID))
}
