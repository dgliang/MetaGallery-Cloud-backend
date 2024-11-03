package services

import (
	"MetaGallery-Cloud-backend/models"
	"fmt"
)

func GenerateRootFolder(userID uint) error {
	_, err := models.CreateRootFolder(userID, fmt.Sprintf("%d", userID),
		fmt.Sprintf("/%d", userID))

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
		return path, nil
	}

	// 如果不是根目录，先获取父文件夹的目录，再加上文件夹名称（同样要验证会不会重名）
	parentPath, err := models.GetParentFolderPath(userID, parentID)
	if err != nil {
		return "", err
	}

	path = fmt.Sprintf("%s/%s", parentPath, folderName)
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
