package services

import (
	"MetaGallery-Cloud-backend/models"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
)

var FileDirPath string

func init() {
	godotenv.Load()
	FileDirPath = os.Getenv("FILE_DIR_PATH")
}

func GenerateRootFolder(userID uint) error {
	_, err := models.CreateRootFolder(userID, fmt.Sprintf("%d", userID),
		fmt.Sprintf("/%d", userID))

	if err != nil {
		return err
	}

	err = checkAndCreateFolder(fmt.Sprintf("/%d", userID))
	if err != nil {
		return err
	}
	return nil
}

func GenerateFolderPath(userID, parentID uint, folderName string) (string, error) {
	var path string

	// 如果是 0 根目录，直接创建文件夹（注意要确认会不会重名）
	if parentID == 0 {

		// 首先检查是否有 /{userID} 这个根目录，没有的话创建一个
		rootFolder, err := models.GetRootFolderData(userID)
		if err != nil {
			return "", err
		}
		if rootFolder.ID == 0 {
			_, err := models.CreateRootFolder(userID, fmt.Sprintf("%d", userID),
				fmt.Sprintf("/%d", parentID))
			if err != nil {
				return "", err
			}
		}

		path = fmt.Sprintf("/%d/%s", userID, folderName)

		err = checkAndCreateFolder(path)
		if err != nil {
			return "", err
		}
		return path, nil
	}

	// 如果不是根目录，先获取父文件夹的目录，再加上文件夹名称（同样要验证会不会重名）
	parentPath, err := models.GetParentFolderPath(userID, parentID)
	if err != nil {
		return "", err
	}

	path = fmt.Sprintf("%s/%s", parentPath, folderName)
	err = checkAndCreateFolder(path)
	if err != nil {
		return "", err
	}
	return path, nil
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
	fullPath := filepath.Join(FileDirPath, folderPath)

	// 创建完整路径的所有父目录（如果不存在）
	fatherPath := filepath.Dir(fullPath)
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
