package services

import (
	"MetaGallery-Cloud-backend/models"
	"strings"

	"gorm.io/gorm"
)

func SetFolderShareState(userId, folderId uint, shareState bool) error {
	// 1. 更新数据库的 share 字段
	// if err := setShareFolderState(userId, folderId, shareState); err != nil {
	// 	return fmt.Errorf("SetFolderShareState: database: %w", err)
	// }

	// workingDir, err := os.Getwd()
	// if err != nil {
	// 	return fmt.Errorf("SetFolderShareState: get working dir: %w", err)
	// }
	// workingDir = strings.ReplaceAll(workingDir, "\\", "/")
	// log.Println(workingDir)

	// absolutePath := path.Join(workingDir, "resources", "./img/B.png")
	// log.Println(absolutePath)

	// ipfsHash, err := UploadFileToPinata(absolutePath)
	// if err != nil {
	// 	return fmt.Errorf("SetFolderShareState: upload file: %w", err)
	// }
	// log.Println(ipfsHash)
	CreateGroup("")

	return nil
}

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

// func getFilesInFolder()
