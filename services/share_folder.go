package services

import (
	"MetaGallery-Cloud-backend/models"
	"fmt"
	"log"
	"strings"

	"gorm.io/gorm"
)

func SetFolderShareState(userId, folderId uint, shareState bool, intro ...string) error {
	folderData, err := models.GetFolderDataByID(folderId)
	if err != nil {
		return fmt.Errorf("SetFolderShareState: get folder data: %w", err)
	}

	// 如果两次的 Share 状态相同，则直接返回
	if folderData.Share == shareState {
		return nil
	}

	if shareState {
		var introStr string
		if len(intro) > 0 {
			introStr = intro[0]
		} else {
			return fmt.Errorf("SetFolderShareState: intro is empty")
		}

		// 如果 Share 状态为 true，则将文件夹上传到 IPFS
		// 1. 更新数据库的 share 字段
		if err := setShareFolderState(userId, folderId, shareState); err != nil {
			return fmt.Errorf("SetFolderShareState: database: %w", err)
		}

		// 2. 上传文件夹到 IPFS
		var folderCID string
		folderCID, err = uploadFolderToIPFS(folderData)
		if err != nil {
			return fmt.Errorf("SetFolderShareState: upload folder: %w", err)
		}

		// 3. 更新数据库，在 shared_folders 表中插入一条记录
		log.Println(folderCID)
		if _, err := CreateSharedFolder(userId, folderId, introStr, folderCID); err != nil {
			return fmt.Errorf("SetFolderShareState: create shared folder: %w", err)
		}
	} else {
		// 如果 Share 状态为 false，则从 IPFS 删除文件夹
		// 1. 更新数据库的 share 字段
		if err := setShareFolderState(userId, folderId, shareState); err != nil {
			return fmt.Errorf("SetFolderShareState: database: %w", err)
		}

		// 2. 从 IPFS 删除文件夹

		// 3. 更新数据库，在 shared_folders 表中删除对应记录
		if err := DeleteSharedFolder(userId, folderId); err != nil {
			return fmt.Errorf("SetFolderShareState: delete shared folder: %w", err)
		}
	}

	return nil
}

// 数据库中更新 folder_data 表和 file_data 表的 Share 字段
func setShareFolderState(userId, folderId uint, shareState bool) error {
	return models.DataBase.Transaction(func(tx *gorm.DB) error {
		var folder models.FolderData
		if err := tx.First(&folder, "id = ? AND belong_to = ?", folderId, userId).Error; err != nil {
			return err
		}

		if folder.Share == shareState {
			return nil
		}

		folder.Share = shareState
		tmpParentPath := folder.Path

		// 更新父文件夹的 Share 字段
		if err := tx.Save(&folder).Error; err != nil {
			return err
		}

		// 更新所有子文件和子文件夹的 Share 字段
		parentPath := strings.ReplaceAll(strings.TrimSpace(tmpParentPath), "\\", "/")
		var subFolders []models.FolderData
		if err := tx.Where("path LIKE ? AND belong_to = ?", parentPath+"/%", userId).
			Find(&subFolders).Error; err != nil {
			return err
		}

		for _, subFolder := range subFolders {
			subFolder.Share = shareState
			if err := tx.Save(&subFolder).Error; err != nil {
				return err
			}
		}

		// TODO: 更新所有子文件的 Share 字段

		return nil
	})
}

/*
uploadFolderToIPFS
 1. 使用 ListChildFolders 检查是否有子文件夹
 2. 对每个子文件夹，继续递归，直到所有子文件夹和文件都被处理。
 3. 每一层的文件夹都会上传到 IPFS，并为每个文件夹生成一个唯一的 CID（内容标识符）。
 4. 子文件夹中的文件也需要上传，并生成对应的 CID。在文件夹上传之前，必须先处理该文件夹中的所有内容（文件和子文件夹）。可以将文件夹和文件信息（包括子文件夹的 CID）组织成一个 JSON 数据结构，然后上传该结构到 IPFS，获得文件夹的 CID。
*/
func uploadFolderToIPFS(folderData models.FolderData) (string, error) {
	// 1. 上传当前文件夹中的文件，使用 UploadFileToIPFS 接口
	// TODO: 上传当前文件夹中的文件，同时获取所有的 file 构成 filesMap
	/*
		"filesMap": [
			{
				"file_name": "file1.txt",
				"cid": "Qm...cid_of_file1",
			},
			{
				"file_name": "file2.txt",
				"cid": "Qm...cid_of_file2",
			},
			...
		]
	*/
	var filesMap []map[string]interface{}
	filesMap = append(filesMap, map[string]interface{}{
		"file_name": "file1.txt",
		"cid":       "Qm...cid_of_file1",
	})

	// 2. 递归上传子文件夹，使用 ListChildFolders 获取所有的子文件夹
	var subFoldersMap []map[string]interface{}
	subFolders, err := models.ListChildFolders(folderData.BelongTo, folderData.ID)
	if err != nil {
		return "", err
	}

	for _, subFolder := range subFolders {
		subFolderCID, err := uploadFolderToIPFS(subFolder) // 递归上传子文件夹
		if err != nil {
			return "", err
		}
		// 将 subFolderCID 记录起来
		subFoldersMap = append(subFoldersMap, map[string]interface{}{
			"folder_name": subFolder.FolderName,
			"cid":         subFolderCID,
		})
	}

	// 3. 构建文件夹的元数据
	folderMetadata := generateMetaInFolder(folderData.FolderName, filesMap, subFoldersMap)

	// 4. 将文件夹结构上传到 IPFS
	folderCID, err := UploadJsonToIPFS(folderMetadata)
	if err != nil {
		return "", err
	}

	// 5. 返回文件夹的 CID
	return folderCID, nil
}

// 获取文件夹内部子文件夹和子文件的结构，并创建 JSON 格式文件
func generateMetaInFolder(folderName string, files, subFolders []map[string]interface{}) map[string]interface{} {
	folderMeta := map[string]interface{}{
		"folderName": folderName,
		"files":      files,
		"subFolders": subFolders,
	}
	return folderMeta
}

// 数据库中 shared_folders 表中创建新的共享文件夹的记录
func CreateSharedFolder(userId, folderId uint, intro, ipfsHash string) (uint, error) {
	sharedFolder := models.SharedFolder{
		OwnerID:  userId,
		FolderID: folderId,
		Intro:    intro,
		IPFSHash: ipfsHash,
	}

	return sharedFolder.ID, models.DataBase.Create(&sharedFolder).Error
}

// 数据库中 shared_folders 表中删除的共享文件夹的记录
func DeleteSharedFolder(userId, folderId uint) error {
	return models.DataBase.Where("owner_id = ? AND folder_id = ?", userId, folderId).Delete(&models.SharedFolder{}).Error
}