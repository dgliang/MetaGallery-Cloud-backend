package controllers

import (
	"MetaGallery-Cloud-backend/models"
	"MetaGallery-Cloud-backend/services"
	"github.com/gin-gonic/gin"
)

type FolderController struct{}

type FolderJson struct {
	ID         uint   `json:"id"`
	User       uint   `json:"user"`
	FolderName string `json:"folder_name"`
	ParentID   uint   `json:"parent_id"`
	Path       string `json:"path"`
	IsFavorite bool   `json:"is_favorite"`
	IsShare    bool   `json:"is_share"`
	IPFSHash   string `json:"ipfs_hash"`
}

func (receiver FolderController) GetRootFolder(c *gin.Context) {
	account := c.Query("account")

	if account == "" {
		ReturnError(c, "FAILED", "未提供账号")
		return
	}

	userID, err := models.GetUserID(account)
	if err != nil {
		ReturnServerError(c, "GetUserID"+err.Error())
		return
	}
	if userID == 0 {
		ReturnError(c, "FAILED", "用户不存在")
		return
	}

	folderData, err := models.GetRootFolderData(userID)
	if err != nil {
		ReturnServerError(c, "GetRootFolderData"+err.Error())
		return
	}

	folderRes := FolderJson{
		ID:         folderData.ID,
		User:       folderData.BelongTo,
		FolderName: folderData.FolderName,
		Path:       folderData.Path,
		IsFavorite: folderData.Favorite,
		IsShare:    folderData.Share,
	}
	ReturnSuccess(c, "SUCCESS", "", folderRes)
}

type folderRequest struct {
	Account    string `json:"account" binding:"required"`
	ParentID   uint   `json:"parent_id" binding:"required"`
	FolderName string `json:"folder_name" binding:"required"`
}

func (receiver FolderController) CreateFolder(c *gin.Context) {
	var req folderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, "FAILED", "解析 JSON Request："+err.Error())
		return
	}

	userID, err := models.GetUserID(req.Account)
	if err != nil {
		ReturnServerError(c, "获取 GetUserID: "+err.Error())
		return
	}

	if userID == 0 {
		ReturnError(c, "Failed", "用户不存在")
		return
	}

	// 判断文件夹是否已经被创建过了
	isExist, err := services.IsExist(userID, req.ParentID, req.FolderName)
	if err != nil {
		ReturnServerError(c, "IsExist"+err.Error())
		return
	}
	if isExist {
		ReturnError(c, "FAILED", "文件夹已经被创建了，无法再次创建")
		return
	}

	path, err := services.GenerateFolderPath(userID, req.ParentID, req.FolderName)
	if err != nil {
		ReturnServerError(c, "GenerateFolderPath: "+err.Error())
		return
	}

	folderData, err := models.CreateFolder(userID, req.ParentID, req.FolderName, path)
	if err != nil {
		ReturnServerError(c, err.Error())
		return
	}

	folderRes := FolderJson{
		ID:         folderData.ID,
		User:       folderData.BelongTo,
		FolderName: folderData.FolderName,
		ParentID:   folderData.ParentFolder,
		Path:       folderData.Path,
		IsFavorite: folderData.Favorite,
		IsShare:    folderData.Share,
	}
	ReturnSuccess(c, "SUCCESS", "", folderRes)
}

type ChildFolderReq struct {
	Account  string `json:"account" binding:"required"`
	FolderId uint   `json:"folder_id" binding:"required"`
}

func (receiver FolderController) GetChildFolders(c *gin.Context) {
	var req ChildFolderReq

	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, "FAILED", "提供查看子文件夹的信息不全")
		return
	}

	userID, err := models.GetUserID(req.Account)
	if err != nil {
		ReturnServerError(c, "GetUserID"+err.Error())
		return
	}
	if userID == 0 {
		ReturnError(c, "FAILED", "用户不存在")
		return
	}

	foldersData, err := models.ListChildFolders(userID, req.FolderId)
	if err != nil {
		ReturnServerError(c, "ListChildFolders: "+err.Error())
		return
	}
	if foldersData == nil {
		ReturnError(c, "FAILED", "folder_id 对应的文件夹不存在")
		return
	}

	folderRes := matchFolderResJson(foldersData)
	ReturnSuccess(c, "SUCCESS", "", folderRes)
}

func matchFolderResJson(foldersData []models.FolderData) []FolderJson {
	if len(foldersData) == 0 {
		return nil
	}

	var folderJson []FolderJson
	for _, folderData := range foldersData {
		folderJson = append(folderJson, FolderJson{
			ID:         folderData.ID,
			User:       folderData.BelongTo,
			FolderName: folderData.FolderName,
			ParentID:   folderData.ParentFolder,
			Path:       folderData.Path,
			IsFavorite: folderData.Favorite,
			IsShare:    folderData.Share,
			IPFSHash:   folderData.IPFSInformation,
		})
	}
	return folderJson
}

type renameFolderRequest struct {
	Account       string `json:"account" binding:"required"`
	FolderID      uint   `json:"folder_id" binding:"required"`
	NewFolderName string `json:"new_folder_name" binding:"required"`
}

func (receiver FolderController) RenameFolder(c *gin.Context) {
	var req renameFolderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, "FAILED", "提供的 JSON 数据出错"+err.Error())
		return
	}

	userID, err := models.GetUserID(req.Account)
	if err != nil {
		ReturnServerError(c, "GetUserID"+err.Error())
		return
	}
	if userID == 0 {
		ReturnError(c, "FAILED", "提供的用户不存在")
		return
	}

	folderData, err1 := models.GetFolderDataByID(req.FolderID)
	if folderData.ID == 0 {
		ReturnError(c, "FAILED", "文件夹不存在")
		return
	}

	// 更改前后的文件夹名称相同
	if folderData.FolderName == req.NewFolderName {
		folderRes := FolderJson{
			ID:         folderData.ID,
			User:       folderData.BelongTo,
			ParentID:   folderData.ParentFolder,
			FolderName: folderData.FolderName,
			Path:       folderData.Path,
			IsFavorite: folderData.Favorite,
			IsShare:    folderData.Share,
		}
		ReturnSuccess(c, "SUCCESS", "", folderRes)
		return
	}

	err = services.RenameFolderAndUpdatePath(userID, req.FolderID, req.NewFolderName)
	if err != nil {
		ReturnServerError(c, "RenameFolderAndUpdatePath: "+err.Error())
		return
	}

	folderData, err1 = models.GetFolderDataByID(req.FolderID)
	if err1 != nil {
		ReturnServerError(c, "GetFolderDataByID: "+err1.Error())
	}

	folderRes := FolderJson{
		ID:         folderData.ID,
		User:       folderData.BelongTo,
		FolderName: folderData.FolderName,
		Path:       folderData.Path,
		IsFavorite: folderData.Favorite,
		IsShare:    folderData.Share,
		ParentID:   folderData.ParentFolder,
	}
	ReturnSuccess(c, "SUCCESS", "", folderRes)
}
