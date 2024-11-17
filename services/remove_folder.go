package services

import (
	"MetaGallery-Cloud-backend/models"
	"gorm.io/gorm"
	"log"
	"strings"
	"time"
)

func RemoveFolder(userId, folderID uint) error {
	return models.DataBase.Transaction(func(tx *gorm.DB) error {
		var folder models.FolderData
		if err := tx.Where("id = ? AND belong_to = ?", folderID, userId).
			First(&folder).Error; err != nil {
			return err
		}

		// 将文件夹数据插入到回收站表 Bin
		delTime := time.Now()
		bin := models.Bin{
			Type:        models.FOLDER,
			DeletedTime: delTime,
			UserID:      userId,
		}
		if err := tx.Create(&bin).Error; err != nil {
			return err
		}

		// 在 FolderBin 表中记录文件夹与回收站的关联
		folderBin := models.FolderBin{
			FolderID: folder.ID,
			BinID:    bin.ID,
		}
		if err := tx.Create(&folderBin).Error; err != nil {
			return err
		}

		// 先软删除文件

		// 再从原文件夹表中删除子文件夹（软删除）
		parentPath := folder.Path + "/"
		if err := removeSubfolder(tx, userId, parentPath, delTime); err != nil {
			return err
		}
		// 最后从原文件夹表中删除（软删除）
		if err := tx.Delete(&folder).Error; err != nil {
			return err
		}

		return nil
	})
}

func removeSubfolder(tx *gorm.DB, userId uint, parentPath string, deleteTime time.Time) error {
	fullParentPath := strings.ReplaceAll(strings.TrimSpace(parentPath), "\\", "/")

	var subFolders []models.FolderData
	if err := tx.Where("path LIKE ? AND belong_to = ?", fullParentPath+"%", userId).
		Find(&subFolders).Error; err != nil {
		return err
	}

	log.Println(subFolders)

	// 遍历子文件夹并进行软删除
	for _, subFolder := range subFolders {
		// 将文件夹数据插入到回收站表 Bin
		subBin := models.Bin{
			Type:        models.FOLDER,
			DeletedTime: deleteTime,
			UserID:      userId,
		}
		if err := tx.Create(&subBin).Error; err != nil {
			return err
		}

		// 在 FolderBin 表中记录文件夹与回收站的关联
		subFolderBin := models.FolderBin{
			FolderID: subFolder.ID,
			BinID:    subBin.ID,
		}
		if err := tx.Create(&subFolderBin).Error; err != nil {
			return err
		}

		// 从原文件夹表中删除（软删除）
		if err := tx.Delete(&subFolder).Error; err != nil {
			return err
		}
	}
	return nil
}
