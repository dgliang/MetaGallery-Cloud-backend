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

func (s FolderShareController) GetFolderInfo(c *gin.Context) {
	ownerAccount := c.Query("owner_account")
	folderName := c.Query("folder_name")
	ipfsHash := c.Query("ipfs_hash")

	if ownerAccount == "" || folderName == "" || ipfsHash == "" {
		ReturnError(c, "FAILED", "owner_account, folder_name, ipfs_hash 字段不能为空")
		return
	}

	userId, err := models.GetUserID(ownerAccount)
	if err != nil || userId == 0 {
		ReturnError(c, "FAILED", "获取用户ID失败，用户不存在")
		return
	}

	sharedFolder, err := services.GetSharedFolderByOwnerAndName(userId, folderName)
	if err != nil || sharedFolder.ID == 0 {
		ReturnError(c, "FAILED", "要删除的共享文件夹不存在")
		return
	}

	res, err := services.GetSharedFolderInfoFromIPFS(ownerAccount, ipfsHash)
	if err != nil {
		ReturnServerError(c, "获取共享文件夹信息失败"+err.Error())
		return
	}

	ReturnSuccess(c, "SUCCESS", "", res)
}

func (s FolderShareController) DownloadSharedFile(c *gin.Context) {
	account := c.Query("account")
	ipfsHash := c.Query("ipfs_hash")

	if account == "" || ipfsHash == "" {
		ReturnError(c, "FAILED", "account, ipfs_hash 字段不能为空")
		return
	}

	userId, err := models.GetUserID(account)
	if err != nil || userId == 0 {
		ReturnError(c, "FAILED", "获取用户ID失败，用户不存在")
		return
	}

	sharedFolder, err := services.GetSharedFolderByIPFSHash(userId, ipfsHash)
	if err == nil && sharedFolder.ID != 0 {
		ReturnError(c, "FAILED", "要下载的是共享文件夹而不是共享文件，不支持下载共享文件夹")
		return
	}

	// 从缓存文件中获取共享文件数据并传送给前端
	err = services.DownloadSharedFile(c, ipfsHash)
	if err != nil {
		ReturnServerError(c, "下载共享文件失败"+err.Error())
		return
	}
}
