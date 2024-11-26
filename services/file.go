package services

import (
	"MetaGallery-Cloud-backend/models"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func FileExist(userID, parentId uint, fileName string) bool {
	isExist, err := models.GetFileID(userID, parentId, fileName)
	if err != nil {
		return false
	}
	return isExist != 0
}

func DetectFileType(file multipart.File) (string, error) {
	// 创建一个缓冲区，读取文件的前512字节进行检测
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		return "", err
	}

	// 使用http.DetectContentType来检测文件的MIME类型
	fileType := http.DetectContentType(buffer)
	file.Seek(0, 0)
	return fileType, nil
}

func detectFileTypeByExtension(fileName string) string {
	ext := filepath.Ext(fileName)
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".pdf":
		return "application/pdf"
	default:
		return "unknown"
	}
}

func FileType(fileName string) string {
	ext := filepath.Ext(fileName)
	return ext
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

func SaveFile2(userID, uintPID uint, fileID uint, file multipart.File) error {

	//生成路径
	path, err := models.GenerateFilePath(userID, uintPID, strconv.FormatUint(uint64(fileID), 10))
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
	oldFullPath := path.Join(FileDirPath, oldPath)
	newFullPath := path.Join(FileDirPath, newPath)
	os.Rename(oldFullPath, newFullPath)
	//修改数据库相关内容
	models.RenameFileWithFileID(fileID, newFileName)
	return nil
}

func RenameFile(userID, fileID uint, newFileName string) error {

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

	models.RenameFileWithFileID2(fileID, newFileName)
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
			return fmt.Errorf("updateChildFilesPaths: %w", err)
		}
	}

	return nil
}

func removeSubFiles(tx *gorm.DB, userId uint, parentPath string, deleteTime time.Time) error {
	// 获取所有直接子文件夹
	parentPath = strings.ReplaceAll(strings.TrimSpace(parentPath), "\\", "/")

	var subFiles []models.FileData
	if err := tx.Model(models.FileData{}).Where("path LIKE ? AND belong_to = ?", parentPath+"%", userId).
		Find(&subFiles).Error; err != nil {
		return fmt.Errorf("removeSubFiles Error: %w", err)
	}

	log.Println("路径：" + parentPath)
	log.Println("影响文件:" + strconv.Itoa(len(subFiles)))

	// 遍历子文件进行删除
	for _, file := range subFiles {

		//创建回收站记录
		// binItem, err := models.InsertBinItem(models.Bin{
		// 	Type:        models.FILE,
		// 	DeletedTime: deleteTime,
		// 	UserID:      file.BelongTo,
		// })
		// if err != nil {
		// 	return fmt.Errorf("RemoveSubFiles error: %w", err)
		// }

		//创建文件回收项
		// fileBinItem := models.FileBin{
		// 	FileID: file.ID,
		// 	BinID:  binItem.ID,
		// }
		// if err := models.InsertFileBinItem(fileBinItem); err != nil {
		// 	return err
		// }

		// 对文件表中的记录软删除
		if err := models.RemoveFile(file.ID); err != nil {
			return fmt.Errorf("RemoveSubFiles error: %w", err)
		}
	}

	return nil
}

func recoverSubFiles(tx *gorm.DB, userId uint, parentPath string, deleteTime time.Time) error {
	parentPath = strings.ReplaceAll(strings.TrimSpace(parentPath), "\\", "/")

	var subFiles []models.FileData
	if err := tx.Model(models.FileData{}).Where("path LIKE ? AND belong_to = ?", parentPath+"%", userId).
		Find(&subFiles).Error; err != nil {
		return fmt.Errorf("RecoverSubFiles Error: %w", err)
	}

	log.Println("路径：" + parentPath)
	log.Println("影响文件:" + strconv.Itoa(len(subFiles)))

	// 遍历子文件进行恢复
	for _, file := range subFiles {
		if err := models.RecoverFile(file.ID); err != nil {
			return fmt.Errorf("RecoverSubFiles error: %w", err)
		}
	}

	return nil
}

func RemoveFile(fileID uint) error {

	fileData, err := models.GetFileData(fileID)
	if err != nil {
		return fmt.Errorf("RemoveFile error: %w", err)
	}
	//创建回收站项
	binItem, err := models.InsertBinItem(models.Bin{
		Type:        models.FILE,
		DeletedTime: time.Now(),
		UserID:      fileData.BelongTo,
	})
	if err != nil {
		return fmt.Errorf("RemoveFile error: %w", err)
	}
	//创建文件回收项
	fileBinItem := models.FileBin{
		FileID: fileData.ID,
		BinID:  binItem.ID,
	}
	if err := models.InsertFileBinItem(fileBinItem); err != nil {
		return err
	}
	// 对文件表中的记录软删除
	if err := models.RemoveFile(fileID); err != nil {
		return fmt.Errorf("RemoveFile error: %w", err)
	}

	return nil
}

func RecoverFile(userID uint, fileID uint) error {
	//回收站内是否有该文件
	if !models.FileBinItemExist(fileID) {
		return fmt.Errorf("RecoverFile error: fileRecycleItem do not exist")
	}

	fileData, err := models.GetDeletedFileData(fileID)
	if err != nil {
		return err
	}

	//查询原本名字是否被占用
	count := 1
	fileName := fileData.FileName
	for {
		if FileExist(userID, fileData.ParentFolderID, fileName) {
			models.UnscopedRenameFile2(fileID, fileData.FileName+" ("+strconv.Itoa(count)+")")
			fileName = fileData.FileName + " (" + strconv.Itoa(count) + ")"
			count += 1
		} else {
			break
		}
	}
	// 将文件的delete at重新置为空
	if err := models.RecoverFile(fileID); err != nil {
		return fmt.Errorf("RecoverFile error: %w", err)
	}
	// 删除对应文件回收项
	fileBinItem, err := models.DeleteFileBinItem(fileID)
	if err != nil {
		return err
	}
	// 删除对应回收站项
	if err := models.DeleteBinItem(fileBinItem.ID); err != nil {
		return fmt.Errorf("RecoverFile error: %w", err)
	}

	return nil
}

func GetBinFiles(userID uint) ([]models.FileBrief, error) {
	_, binItemIDs := models.SearchBinItems(userID)

	var fileBriefs []models.FileBrief
	for _, binItemID := range binItemIDs {

		fileID := models.GetFileIDInBIN(binItemID)

		fileData, err := models.GetDeletedFileData(fileID)
		if err != nil {
			return nil, err
		}

		fileBrief := models.FileBrief{
			ID:       fileData.ID,
			FileName: fileData.FileName,
			FileType: fileData.FileType,
			Favorite: fileData.Favorite,
			Share:    fileData.Share,
			InBin:    fileData.DeletedAt.Time,
		}
		fileBriefs = append(fileBriefs, fileBrief)
	}
	return fileBriefs, nil
}

func DownloadFile(c *gin.Context, userID uint, fileID uint) (multipart.File, error) {

	fileData, err := models.GetFileData(fileID)
	if err != nil {
		return nil, err
	}
	if fileData.ID == 0 {
		return nil, fmt.Errorf("文件不存在")
	}

	filePath := path.Join(FileDirPath, fileData.Path)
	c.FileAttachment(filePath, fileData.FileName)

	return nil, nil
}
