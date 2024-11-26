package services

import (
	"MetaGallery-Cloud-backend/config"
	"MetaGallery-Cloud-backend/models"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

func FileExist(userID, parentId uint, fileName string) bool {
	isExist, err := models.GetFileID(userID, parentId, fileName)
	if err != nil {
		return false
	}
	return isExist != 0
}

func SaveFile(userID, uintPID uint, fileName string, file multipart.File) error {

	//生成路径
	path, err := models.GenerateFilePath(userID, uintPID, fileName)
	if err != nil {
		fmt.Println("文件路径生成失败:", err)

		return fmt.Errorf("文件路径生成失败")
	}

	//在本地创建文件
	out, err := os.Create("resources/files" + path)
	log.Printf("resources/files" + path)
	if err != nil {

		fmt.Println("文件创建失败:", err)
		return fmt.Errorf("服务器创建文件失败")
	}
	defer out.Close()

	// 将上传的文件内容写入到本地文件
	_, err = out.ReadFrom(file)
	if err != nil {
		log.Printf("写入文件失败")

		return fmt.Errorf("服务器写入保存文件失败")
	}
	return nil
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
	oldFullPath := path.Join(config.FileResPath, oldPath)
	newFullPath := path.Join(config.FileResPath, newPath)
	os.Rename(oldFullPath, newFullPath)
	//修改数据库相关内容
	models.RenameFileWithFileID(fileID, newFileName)
	return nil
}

// 批量更新文件路径
func updateSubFilesPaths(tx *gorm.DB, userId uint, oldPath, newPath string) error {
	// 获取所有直接子文件夹
	oldPath = strings.ReplaceAll(strings.TrimSpace(oldPath), "\\", "/")
	newPath = strings.ReplaceAll(strings.TrimSpace(newPath), "\\", "/")

	var subFiles []models.FileData
	if err := tx.Model(models.FileData{}).Where("path LIKE ? AND belong_to = ?", oldPath+"/%", userId).
		Find(&subFiles).Error; err != nil {
		return fmt.Errorf("updateChildFolderPaths: %w", err)
	}

	log.Println("原路径：" + oldPath)
	log.Println("新路径：" + newPath)
	log.Println("影响文件:" + strconv.Itoa(len(subFiles)))

	// 遍历子文件并更新路径
	for _, file := range subFiles {
		// 替换路径中的 oldPath 前缀为 newPath
		file.Path = strings.Replace(file.Path, oldPath, newPath, 1)

		// 保存更新后的路径
		if err := tx.Save(&file).Error; err != nil {
			return fmt.Errorf("updateChildFolderPaths: %w", err)
		}
	}

	return nil
}
