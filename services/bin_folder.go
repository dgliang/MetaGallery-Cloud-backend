package services

import (
	"MetaGallery-Cloud-backend/config"
	"MetaGallery-Cloud-backend/models"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"gorm.io/gorm"
)

func RemoveFolder(userId, folderID uint) error {
	return models.DataBase.Transaction(func(tx *gorm.DB) error {
		var folder models.FolderData
		if err := tx.Where("id = ? AND belong_to = ?", folderID, userId).
			First(&folder).Error; err != nil {
			return err
		}

		delTime := time.Now()

		// 重命名文件夹，修改相应的路径。根据时间戳设置重命名
		newFolderName := GenerateBinTimestamp(folder.FolderName, delTime)

		// 检查同一文件夹下是否存在同名的子文件夹
		var count int64
		if err := tx.Model(&models.FolderData{}).
			Where("parent_folder = ? AND folder_name = ? AND belong_to = ?", folder.ParentFolder,
				newFolderName, userId).Count(&count).Error; err != nil {
			return fmt.Errorf("RenameFolderAndUpdatePath: %w", err)
		}
		if count > 0 {
			return errors.New("RenameFolderAndUpdatePath: 重命名的文件夹已存在")
		}

		// 更新文件夹的 FolderName
		folder.FolderName = newFolderName

		// 创建新的 Path
		oldPath := folder.Path
		newPath := path.Join(path.Dir(oldPath), newFolderName)
		log.Println("oldPath: ", oldPath)
		log.Println("newPath: ", newPath)

		// 在服务器更新 Path
		if err := updateFolderPath(oldPath, newPath); err != nil {
			return fmt.Errorf("RenameFolderAndUpdatePath: %w", err)
		}

		// 数据库中更新当前文件夹的 Path
		folder.Path = newPath
		if err := tx.Save(&folder).Error; err != nil {
			return fmt.Errorf("RenameFolderAndUpdatePath: %w", err)
		}

		// 更新所有子文件夹的路径
		if err := updateChildFolderPaths(tx, userId, oldPath, newPath); err != nil {
			return fmt.Errorf("RenameFolderAndUpdatePath: %w", err)
		}

		// 更新所有子文件的路径
		if err := updateSubFilesPaths(tx, userId, oldPath, newPath); err != nil {
			return fmt.Errorf("UpdateSubFilesState: %w", err)
		}

		// 将文件夹数据插入到回收站表 Bin
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

		parentPath := folder.Path + "/"
		fmt.Println("parentPath" + parentPath)
		// 先软删除文件
		if err := removeSubFiles(tx, userId, newPath); err != nil {
			return err
		}

		// 再从原文件夹表中删除子文件夹（软删除）
		if err := removeSubfolder(tx, userId, parentPath); err != nil {
			return err
		}
		// 最后从原文件夹表中删除（软删除）
		if err := tx.Delete(&folder).Error; err != nil {
			return err
		}

		return nil
	})
}

func removeSubfolder(tx *gorm.DB, userId uint, parentPath string) error {
	fullParentPath := strings.ReplaceAll(strings.TrimSpace(parentPath), "\\", "/")

	var subFolders []models.FolderData
	if err := tx.Where("path LIKE ? AND belong_to = ?", fullParentPath+"%", userId).
		Find(&subFolders).Error; err != nil {
		return err
	}

	log.Println(subFolders)

	// 遍历子文件夹并进行软删除
	for _, subFolder := range subFolders {
		//// 将文件夹数据插入到回收站表 Bin
		//subBin := models.Bin{
		//	Type:        models.FOLDER,
		//	DeletedTime: deleteTime,
		//	UserID:      userId,
		//}
		//if err := tx.Create(&subBin).Error; err != nil {
		//	return err
		//}
		//
		//// 在 FolderBin 表中记录文件夹与回收站的关联
		//subFolderBin := models.FolderBin{
		//	FolderID: subFolder.ID,
		//	BinID:    subBin.ID,
		//}
		//if err := tx.Create(&subFolderBin).Error; err != nil {
		//	return err
		//}

		// 从原文件夹表中删除（软删除）
		if err := tx.Delete(&subFolder).Error; err != nil {
			return err
		}
	}
	return nil
}

func DeleteFolder(userId uint, binId uint, folderID uint) error {
	return models.DataBase.Transaction(func(tx *gorm.DB) error {
		var folder models.FolderData

		if err := tx.Unscoped().First(&folder, "id = ? AND belong_to = ?", folderID, userId).
			Error; err != nil {
			return err
		}

		// 删除 bins 表中的记录
		var bin models.Bin
		if err := tx.Where("id = ?", binId).First(&bin).Error; err != nil {
			return err
		}

		if err := tx.Delete(&bin).Error; err != nil {
			return err
		}

		// 删除 folder_data 中的文件夹及其子文件夹，file_data 中的文件，folder_bins 中的记录
		if err := deleteOSFolder(folder.Path); err != nil {
			return err
		}
		if err := tx.Unscoped().Delete(&folder).Error; err != nil {
			return err
		}
		return nil
	})
}

func deleteOSFolder(folderPath string) error {
	folderPath = strings.ReplaceAll(folderPath, "\\", "/")
	folderPath = path.Join(config.FileResPath, folderPath)
	err := os.RemoveAll(folderPath)
	return err
}

type FolderBinInfo struct {
	models.FolderData
	BinId   uint
	DelTime time.Time
}

func ListBinFolders(userId uint) ([]FolderBinInfo, error) {
	// 获取 Bin 回收站中文件夹
	var binRecord []models.Bin
	if err := models.DataBase.Where("user_id = ? AND type = ?", userId, models.FOLDER).
		Find(&binRecord).Error; err != nil {
		return nil, err
	}

	var folderBinInfo []FolderBinInfo
	for _, bin := range binRecord {
		// 获取 folderBin 中对应的记录
		folderBin, err := getFolderBinDataByBinId(bin.ID)
		if err != nil {
			return nil, err
		}

		// 获取 folderData 中对应的记录
		var folder models.FolderData
		if err := models.DataBase.Unscoped().First(&folder, "id = ?", folderBin.FolderID).Error; err != nil {
			return nil, err
		}

		// 将 folderName 和 path 处理，去除时间戳
		fullFolderName, _ := SplitBinTimestamp(folder.FolderName)
		fullFolderPath := strings.ReplaceAll(folder.Path, folder.FolderName, fullFolderName)
		folderBinInfo = append(folderBinInfo, FolderBinInfo{
			FolderData: models.FolderData{
				ID:              folder.ID,
				BelongTo:        folder.BelongTo,
				FolderName:      fullFolderName,
				ParentFolder:    folder.ParentFolder,
				Path:            fullFolderPath,
				Favorite:        folder.Favorite,
				Share:           folder.Share,
				IPFSInformation: folder.IPFSInformation,
				CreatedAt:       folder.CreatedAt,
				UpdatedAt:       folder.UpdatedAt,
				DeletedAt:       folder.DeletedAt,
			},
			BinId:   bin.ID,
			DelTime: bin.DeletedTime,
		})
	}

	return folderBinInfo, nil
}

func getFolderBinDataByBinId(binId uint) (models.FolderBin, error) {
	var folderBin models.FolderBin
	if err := models.DataBase.Where("bin_id = ?", binId).First(&folderBin).Error; err != nil {
		return folderBin, err
	}
	return folderBin, nil
}

func IsFolderInBin(userId, binId uint) bool {
	return models.DataBase.Where("user_id = ? AND id = ?", userId, binId).First(&models.Bin{}).Error == nil
}

// getFolderDataByBinId 根据 userId，binId 获取已经软删除的 folderData
func getFolderDataByBinId(userId, binId uint) (models.FolderData, error) {
	folderBin, err := getFolderBinDataByBinId(binId)
	if err != nil {
		return models.FolderData{}, err
	}

	var folder models.FolderData
	if err := models.DataBase.Unscoped().First(&folder, "id = ? AND belong_to = ?", folderBin.FolderID,
		userId).Error; err != nil {
		return models.FolderData{}, err
	}
	return folder, nil
}

// 检查 folder，bin 和 user 的所属对应关系
func CheckFolderBinAndUserRel(userId, folderId, binId uint) bool {
	folderData, err := models.GetBinFolderDataByID(folderId)
	if err != nil {
		return false
	}

	var binRecord models.Bin
	err = models.DataBase.Where("id = ?", binId).First(&binRecord).Error
	if err != nil {
		return false
	}

	var folderBin models.FolderBin
	err = models.DataBase.Where("folder_id = ? AND bin_id = ?", folderId, binId).First(&folderBin).Error
	if err != nil {
		return false
	}

	return folderData.BelongTo == userId && binRecord.UserID == userId &&
		folderBin.FolderID == folderId && folderBin.BinID == binId
}

// 检查恢复的文件夹会不会与现有文件夹产生冲突
func CheckBinFolderAndFolder(userId, binId uint) bool {
	binFolderData, err := getFolderDataByBinId(userId, binId)
	if err != nil {
		return false
	}

	binFolderOriginName, _ := SplitBinTimestamp(binFolderData.FolderName)
	var folderData models.FolderData
	if err := models.DataBase.Where("belong_to = ? AND folder_name = ? AND parent_folder = ?", userId,
		binFolderOriginName, binFolderData.ParentFolder).First(&folderData).Error; err != nil {
		return false
	}

	return true
}

func RecoverBinFolder(userId, binId uint) error {
	return models.DataBase.Transaction(func(tx *gorm.DB) error {
		folder, err := getFolderDataByBinId(userId, binId)
		if err != nil {
			return err
		}

		if err := tx.Unscoped().Model(&folder).Update("deleted_at", nil).Error; err != nil {
			return err
		}

		// 恢复文件夹的子文件夹
		var subFolders []models.FolderData
		if err := tx.Unscoped().Where("path LIKE ? AND belong_to = ?", strings.ReplaceAll(
			strings.TrimSpace(folder.Path), "\\", "/")+"/%", userId).Find(&subFolders).Error; err != nil {
			return err
		}
		for _, subFolder := range subFolders {
			if err := tx.Unscoped().Model(&subFolder).Update("deleted_at", nil).Error; err != nil {
				return err
			}
		}

		// 重命名文件夹，修改相应的路径。根据时间戳设置重命名
		newFolderName, _ := SplitBinTimestamp(folder.FolderName)

		// 检查同一文件夹下是否存在同名的子文件夹
		var count int64
		if err := tx.Model(&models.FolderData{}).
			Where("parent_folder = ? AND folder_name = ? AND belong_to = ?", folder.ParentFolder,
				newFolderName, userId).Count(&count).Error; err != nil {
			return fmt.Errorf("RenameFolderAndUpdatePath: %w", err)
		}
		if count > 0 {
			return errors.New("RenameFolderAndUpdatePath: 重命名的文件夹已存在")
		}

		// 更新文件夹的 FolderName
		folder.FolderName = newFolderName

		// 创建新的 Path
		oldPath := folder.Path
		newPath := path.Join(path.Dir(oldPath), newFolderName)

		// 在服务器更新 Path
		if err := updateFolderPath(oldPath, newPath); err != nil {
			return fmt.Errorf("RenameFolderAndUpdatePath: %w", err)
		}

		// 数据库中更新当前文件夹的 Path
		folder.Path = newPath
		if err := tx.Save(&folder).Error; err != nil {
			return fmt.Errorf("RenameFolderAndUpdatePath: %w", err)
		}

		// 更新所有子文件夹的路径
		if err := updateChildFolderPaths(tx, userId, oldPath, newPath); err != nil {
			return fmt.Errorf("RenameFolderAndUpdatePath: %w", err)
		}

		// 更新所有子文件的路径
		if err := updateSubFilesPaths(tx, userId, oldPath, newPath); err != nil {
			return fmt.Errorf("UpdateSubFilesState: %w", err)
		}

		// 从 bins 和 folder_bins 表中删除相应记录
		// 删除 bins 表中的记录
		var bin models.Bin
		if err := tx.Where("id = ?", binId).First(&bin).Error; err != nil {
			return err
		}

		if err := tx.Delete(&bin).Error; err != nil {
			return err
		}

		return nil
	})
}
