package controllers

import (
	"MetaGallery-Cloud-backend/models"
	"MetaGallery-Cloud-backend/services"
	"fmt"
	"strconv"

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

func (receiver FolderController) GetFolderInfo(c *gin.Context) {
	account := c.Query("account")
	folderId := c.Query("folder_id")

	if account == "" || folderId == "" {
		ReturnError(c, "FAILED", "提供的 account 和 folder_id 信息不全")
		return
	}

	userID, err := models.GetUserID(account)
	if err != nil {
		ReturnServerError(c, "GetUserID"+err.Error())
		return
	}
	if userID == 0 {
		ReturnError(c, "FAILED", "提供的 account 用户不存在")
		return
	}

	folderIdUint64, _ := strconv.ParseUint(folderId, 10, 64)
	folderData, err1 := models.GetFolderDataByID(uint(folderIdUint64))
	if err1 != nil {
		ReturnServerError(c, "GetFolderDataByID"+err1.Error())
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
		IPFSHash:   folderData.IPFSInformation,
	}
	ReturnSuccess(c, "SUCCESS", "", folderRes)
}

type createFolderRequest struct {
	Account    string `json:"account" binding:"required"`
	ParentID   uint   `json:"parent_id" binding:"required"`
	FolderName string `json:"folder_name" binding:"required"`
}

func (receiver FolderController) CreateFolder(c *gin.Context) {
	var req createFolderRequest

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

func (receiver FolderController) GetChildFolders(c *gin.Context) {
	// 获取查询参数
	account := c.Query("account")       // 获取 account 参数
	folderIDStr := c.Query("folder_id") // 获取 folder_id 参数

	// 校验参数是否齐全
	if account == "" || folderIDStr == "" {
		ReturnError(c, "FAILED", "提供查看子文件夹的信息不全")
		return
	}

	// 转换 folder_id 为 uint
	folderID, err := strconv.Atoi(folderIDStr)
	if err != nil {
		ReturnError(c, "FAILED", "folder_id 无效")
		return
	}

	// 获取用户 ID
	userID, err := models.GetUserID(account)
	if err != nil {
		ReturnServerError(c, "GetUserID: "+err.Error())
		return
	}
	if userID == 0 {
		ReturnError(c, "FAILED", "用户不存在")
		return
	}

	// 获取子文件夹数据
	foldersData, err := models.ListChildFolders(uint(userID), uint(folderID))
	if err != nil {
		ReturnServerError(c, "ListChildFolders: "+err.Error())
		return
	}

	// 检查 folder_id 是否有效
	if foldersData == nil || len(foldersData) == 0 {
		// 判断当前 folder_id 是否是有效文件夹
		var folder models.FolderData
		err = models.DataBase.Where("belong_to = ? AND id = ?", uint(userID), uint(folderID)).First(&folder).Error
		if err != nil {
			// 数据库中找不到该文件夹，返回错误
			ReturnError(c, "FAILED", "folder_id 对应的文件夹不存在")
			return
		}

		// 文件夹有效，但为空，返回空列表
		ReturnSuccess(c, "SUCCESS", "文件夹加载成功", []models.FolderData{})
		return
	}

	// 返回成功响应，包含子文件夹数据
	folderRes := matchFolderResJson(foldersData)
	ReturnSuccess(c, "SUCCESS", "文件夹加载成功", folderRes)
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

type favoriteFolderRequest struct {
	Account    string `json:"account" binding:"required"`
	FolderId   uint   `json:"folder_id" binding:"required"`
	IsFavorite int    `json:"is_favorite" binding:"required"`
}

func (receiver FolderController) FavoriteFolder(c *gin.Context) {
	var req favoriteFolderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, "FAILED", "提供的信息不全"+err.Error())
		return
	}

	// 验证 IsFavorite 的取值是否为 1 或者 2
	var favoriteStatus bool
	if req.IsFavorite == 1 {
		favoriteStatus = false
	} else if req.IsFavorite == 2 {
		favoriteStatus = true
	} else {
		ReturnError(c, "FAILED", "is_favorite 的取值只能是 1 或者 2")
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

	folderData, err1 := models.GetFolderDataByID(req.FolderId)
	if err1 != nil {
		ReturnServerError(c, "GetFolderDataByID: "+err1.Error())
		return
	}
	if folderData.ID == 0 {
		ReturnError(c, "FAILED", "文件夹不存在")
		return
	}
	//if folderData.Favorite == favoriteStatus {
	//	ReturnSuccess(c, "SUCCESS", fmt.Sprintf("修改文件夹收藏状态为 %t", favoriteStatus))
	//	return
	//}

	// 更新文件夹的收藏状态
	err = services.SetFolderFavorite(userID, req.FolderId, favoriteStatus)
	if err != nil {
		ReturnServerError(c, "SetFolderFavorite: "+err.Error())
		return
	}

	ReturnSuccess(c, "SUCCESS", fmt.Sprintf("成功将 %s 的 %d 文件收藏状态改为 %t",
		req.Account, req.FolderId, favoriteStatus))
}
