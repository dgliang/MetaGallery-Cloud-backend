package services

import (
	"MetaGallery-Cloud-backend/models"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var FileDirPath string

func init() {
	godotenv.Load()
	FileDirPath = os.Getenv("FILE_DIR_PATH")
}

func GenerateRootFolder(userID uint) error {
	rootFolderName := fmt.Sprintf("%d", userID)
	rootFolderPath := "/" + rootFolderName

	_, err := models.CreateRootFolder(userID, rootFolderName, rootFolderPath)
	if err != nil {
		return err
	}

	err = checkAndCreateFolder(rootFolderPath)
	if err != nil {
		return err
	}
	return nil
}

func IsValidFolderName(s string) bool {
	// 检查长度
	if len(s) == 0 || len(s) > 255 {
		return false
	}

	// 检查是否包含非法字符
	illegalChars := regexp.MustCompile(`[<>:"/\\|?*]`)
	if illegalChars.MatchString(s) {
		return false
	}

	// 检查是否以空格或句点结尾
	if strings.HasSuffix(s, " ") || strings.HasSuffix(s, ".") {
		return false
	}

	// 检查是否为保留名称（不区分大小写）
	reservedNames := map[string]bool{
		"CON": true, "PRN": true, "AUX": true, "NUL": true,
		"COM1": true, "COM2": true, "COM3": true, "COM4": true, "COM5": true, "COM6": true, "COM7": true, "COM8": true, "COM9": true,
		"LPT1": true, "LPT2": true, "LPT3": true, "LPT4": true, "LPT5": true, "LPT6": true, "LPT7": true, "LPT8": true, "LPT9": true,
	}

	upperName := strings.ToUpper(s)
	return !reservedNames[upperName]
}

func GenerateFolderPath(userID, parentID uint, folderName string) (string, error) {
	var folderPath string

	// 如果是 0 根目录，直接创建文件夹（注意要确认会不会重名）
	if parentID == 0 {

		// 首先检查是否有 /{userID} 这个根目录，没有的话创建一个
		rootFolder, err := models.GetRootFolderData(userID)
		if err != nil {
			return "", err
		}
		if rootFolder.ID == 0 {
			_ = GenerateRootFolder(userID)
		}

		folderPath = fmt.Sprintf("/%d/%s", userID, folderName)

		err = checkAndCreateFolder(folderPath)
		if err != nil {
			return "", err
		}
		return folderPath, nil
	}

	// 如果不是根目录，先获取父文件夹的目录，再加上文件夹名称（同样要验证会不会重名）
	parentPath, err := models.GetParentFolderPath(userID, parentID)
	if err != nil {
		return "", err
	}

	folderPath = path.Join(parentPath, folderName)
	err = checkAndCreateFolder(folderPath)
	if err != nil {
		return "", err
	}
	return folderPath, nil
}

func IsExist(userId, parentId uint, folderName string) (bool, error) {
	isExist, err := models.GetFolderId(userId, parentId, folderName)

	if err != nil {
		return false, err
	}

	if isExist == 0 {
		return false, nil
	}
	return true, nil
}

func checkAndCreateFolder(folderPath string) error {
	fullPath := path.Join(FileDirPath, folderPath)

	// 创建完整路径的所有父目录（如果不存在）
	fatherPath := path.Dir(fullPath)
	err := os.MkdirAll(fatherPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("checkAndCreateFolder: %w", err)
	}

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {

		// 如果文件夹不存在，创建文件夹
		err := os.Mkdir(fullPath, os.ModePerm)
		if err != nil {
			return fmt.Errorf("checkAndCreateFolder: %w", err)
		}
		log.Println("Folder created successfully: ", fullPath)
	}
	log.Println("Folder exist: ", fullPath)
	return nil
}

func RenameFolderAndUpdatePath(userId, folderId uint, newFolderName string) error {
	return models.DataBase.Transaction(func(tx *gorm.DB) error {

		// 获取当前文件夹信息
		var folder models.FolderData
		if err := tx.First(&folder, "id = ? AND belong_to = ?", folderId, userId).
			Error; err != nil {
			return fmt.Errorf("RenameFolderAndUpdatePath: %w", err)
		}

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
			return fmt.Errorf("RenameFolderAndUpdatePath: %w", err)
		}
		return nil
	})
}

// updateChildFolderPaths 递归更新子文件夹路径
func updateChildFolderPaths(tx *gorm.DB, userId uint, oldPath, newPath string) error {
	// 获取所有直接子文件夹
	oldPath = strings.ReplaceAll(strings.TrimSpace(oldPath), "\\", "/")
	newPath = strings.ReplaceAll(strings.TrimSpace(newPath), "\\", "/")

	var subFolders []models.FolderData
	if err := tx.Where("path LIKE ? AND belong_to = ?", oldPath+"/%", userId).
		Find(&subFolders).Error; err != nil {
		return fmt.Errorf("updateChildFolderPaths: %w", err)
	}

	log.Println(oldPath)
	log.Println(len(subFolders))

	// 遍历子文件夹并更新路径
	for _, folder := range subFolders {
		// 替换路径中的 oldPath 前缀为 newPath
		folder.Path = strings.Replace(folder.Path, oldPath, newPath, 1)

		// 保存更新后的路径
		if err := tx.Save(&folder).Error; err != nil {
			return fmt.Errorf("updateChildFolderPaths: %w", err)
		}
	}

	return nil
}

func updateFolderPath(oldPath, newPath string) error {
	oldFullPath := path.Join(FileDirPath, oldPath)
	newFullPath := path.Join(FileDirPath, newPath)

	log.Println("服务器旧路径：" + oldFullPath)
	log.Println("服务器新地址：" + newFullPath)

	err := os.Rename(oldFullPath, newFullPath)
	return err
}

func SetFolderFavorite(userID, folderID uint, isFavorite bool) error {
	return models.DataBase.Transaction(func(tx *gorm.DB) error {

		// 获取当前文件夹信息
		var folder models.FolderData
		if err := tx.First(&folder, "id = ? AND belong_to = ?", folderID, userID).Error; err != nil {
			return fmt.Errorf("SetFolderFavorite: %w", err)
		}

		// 获取当前文件夹的 Path
		rootPath := folder.Path

		// 更新当前文件夹的 Favorite
		folder.Favorite = isFavorite
		if err := tx.Save(&folder).Error; err != nil {
			return fmt.Errorf("SetFolderFavorite: %w", err)
		}

		// 更新所有子文件夹的 Favorite
		if err := setChildFolderFavorite(tx, userID, rootPath, isFavorite); err != nil {
			return fmt.Errorf("SetFolderFavorite: %w", err)
		}

		// 更新所有子文件的 Favorite
		// todo
		return nil
	})
}

func setChildFolderFavorite(tx *gorm.DB, userID uint, rootPath string, isFavorite bool) error {
	rootPath = strings.ReplaceAll(strings.TrimSpace(rootPath), "\\", "/")

	var subFolders []models.FolderData
	if err := tx.Where("path LIKE ? AND belong_to = ?", rootPath+"/%", userID).
		Find(&subFolders).Error; err != nil {
		return fmt.Errorf("setChildFolderFavorite: %w", err)
	}

	for _, folder := range subFolders {
		folder.Favorite = isFavorite
		if err := tx.Save(&folder).Error; err != nil {
			return fmt.Errorf("setChildFolderFavorite: %w", err)
		}
	}
	return nil
}
