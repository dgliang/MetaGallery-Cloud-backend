package services

import (
	"MetaGallery-Cloud-backend/models"
	"errors"
)

const (
	PAGE_SIZE = 10 // 每页条数
)

type sharedFolderResponse struct {
	OwnerAccount string `json:"owner_account"`
	FolderName   string `json:"folder_name"`
	IPFSHash     string `json:"ipfs_hash"`
	Intro        string `json:"intro"`
	PinDate      string `json:"pin_date"`
	TotalPage    int    `json:"total_page"`
}

// 获取用户自己共享的所有文件夹列表
func ListUserSharedFolders(ownerId uint, pageNum int) ([]sharedFolderResponse, error) {
	if pageNum <= 0 {
		return nil, errors.New("page number should be greater than 0")
	}

	// 计算偏移量
	offset := (pageNum - 1) * PAGE_SIZE
	var sharedFolders []models.SharedFolder
	if err := models.DataBase.Model(&models.SharedFolder{}).
		Where("owner_id = ?", ownerId).
		// Order("created_at DESC"). // 按创建时间降序排序，新的在前
		Limit(PAGE_SIZE).
		Offset(offset).
		Find(&sharedFolders).Error; err != nil {
		return nil, err
	}

	// 计算总页数
	var totalRecords int64
	if err := models.DataBase.Model(&models.SharedFolder{}).
		Where("owner_id = ?", ownerId).
		Count(&totalRecords).Error; err != nil {
		return nil, err
	}

	totalPage := int((totalRecords + int64(PAGE_SIZE) - 1) / int64(PAGE_SIZE))

	res := matchSharedFolderModelToResponse(sharedFolders, totalPage)
	return res, nil
}

// 获取所有用户的共享文件夹列表
func ListAllSharedFolders(pageNum int) ([]sharedFolderResponse, error) {
	if pageNum <= 0 {
		return nil, errors.New("page number should be greater than 0")
	}

	// 计算偏移量
	offset := (pageNum - 1) * PAGE_SIZE
	var sharedFolders []models.SharedFolder
	if err := models.DataBase.Model(&models.SharedFolder{}).
		Order("created_at DESC"). // 按创建时间降序排序，新的在前
		Limit(PAGE_SIZE).
		Offset(offset).
		Find(&sharedFolders).Error; err != nil {
		return nil, err
	}

	// 计算总页数
	var totalRecords int64
	if err := models.DataBase.Model(&models.SharedFolder{}).
		Count(&totalRecords).Error; err != nil {
		return nil, err
	}

	totalPage := int((totalRecords + int64(PAGE_SIZE) - 1) / int64(PAGE_SIZE))

	res := matchSharedFolderModelToResponse(sharedFolders, totalPage)
	return res, nil
}

func matchSharedFolderModelToResponse(folders []models.SharedFolder, totalPage int) []sharedFolderResponse {
	var res []sharedFolderResponse
	for _, folder := range folders {
		ownerAccount, _ := models.GetUserAccountById(folder.OwnerID)

		res = append(res, sharedFolderResponse{
			OwnerAccount: ownerAccount,
			FolderName:   folder.SharedName,
			IPFSHash:     folder.IPFSHash,
			Intro:        folder.Intro,
			PinDate:      folder.CreatedAt.Format("2006-01-02 15:04:05"),
			TotalPage:    totalPage,
		})
	}
	return res
}

type sharedFolderInfoResponse struct {
	OwnerAccount string `json:"owner_account"`
	FolderInfo   folder `json:"folder_info"`
}

func GetSharedFolderInfo(ownerAccount, cid string) (sharedFolderInfoResponse, error) {
	var sharedFolderInfo sharedFolderInfoResponse

	folderData, err := GetFolderJsonFromIPFS(cid)
	if err != nil {
		return sharedFolderInfo, err
	}

	sharedFolderInfo = sharedFolderInfoResponse{
		OwnerAccount: ownerAccount,
		FolderInfo:   folderData,
	}
	return sharedFolderInfo, nil
}
