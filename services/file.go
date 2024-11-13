package services

import (
	"MetaGallery-Cloud-backend/models"
	"fmt"
	"os"
)

func FileExist(userID, parentId uint, fileName string) bool {
	isExist, err := models.GetFileID(userID, parentId, fileName)
	if err != nil {
		return false
	}
	return isExist != 0
}

func RenameFileAndUpdatePath(userID, fileID uint, newFileName string) error {

	//获取文件的父文件夹id
	uintPID, err := models.GetParentFolderID(fileID)
	if err != nil {
		return err
	}
	//判断是否有同名文件
	isExist := FileExist(userID, uintPID, newFileName)
	if isExist {
		return fmt.Errorf("FileName already exists: %s", newFileName)
	}
	//获取原文件路径
	oldPath, err := models.GetFilePath(fileID)
	if err != nil {
		return err
	}
	//生成新文件路径
	newPath, err := models.GenerateFilePath(userID, uintPID, newFileName)
	if err != nil {
		return err
	}
	//修改本地文件名称
	os.Rename("resources/files"+oldPath, "resources/files"+newPath)
	//修改数据库相关内容
	models.RenameFileWithFileID(fileID, newFileName)
	return nil
}
