package services

import (
	"MetaGallery-Cloud-backend/config"
	"MetaGallery-Cloud-backend/models"
	"fmt"
	"log"
	"mime/multipart"
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

func FileType(fileName string) string {
	ext := filepath.Ext(fileName)
	return ext
}

func SaveFileByName(userID, uintPID uint, fileName string, file multipart.File) error {

	//生成路径
	path, err := models.GenerateFilePath(userID, uintPID, fileName)
	if err != nil {
		fmt.Println("文件路径生成失败:", err)

		return fmt.Errorf("文件路径生成失败")
	}

	//在本地创建文件
	out, err := os.Create("resources/files" + path)
	log.Println("resources/files" + path)
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

func SaveFileByID(userID, uintPID uint, fileID uint, file multipart.File) error {

	//生成路径
	path, err := models.GenerateFilePath(userID, uintPID, strconv.FormatUint(uint64(fileID), 10))
	if err != nil {
		fmt.Println("文件路径生成失败:", err)

		return fmt.Errorf("文件路径生成失败")
	}

	//在本地创建文件
	out, err := os.Create("resources/files" + path)
	log.Println("resources/files" + path)
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

// 获取一个文件的详细信息
func GetFileDetail(fileID uint) (models.FileData, error) {
	fileData, err := models.GetFileData(fileID)
	if err != nil {
		return models.FileData{}, err
	}
	// 将文件路径转为正确形式
	path := fileData.Path
	dirParts := strings.Split(path, "/")
	dirParts[len(dirParts)-1] = fileData.FileName
	newPath := strings.Join(dirParts[2:], "/")

	fileData.Path = newPath

	return fileData, err
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

	if err := models.RenameFileWithFileID(fileID, newFileName); err != nil {
		return err
	}
	// models.UpdateFileType(fileID, newFileName)
	return nil
}

func IsFileBelongto(userID, fileID uint) bool {

	fileData, err := models.UnscopedGetFileData(fileID)
	if err != nil {
		return false
	}

	if fileData.BelongTo != userID {
		return false
	}

	return true
}

// 获取文件夹中的子文件
func GetSubFiles(parentFolderID uint) ([]models.FileBrief, error) {
	subFiles, err := models.GetSubFiles(parentFolderID)
	if err != nil {
		return nil, err
	}

	var fileBriefs []models.FileBrief

	for _, source := range subFiles {
		destination := models.FileBrief{
			ID:       source.ID,
			FileName: source.FileName,
			FileType: source.FileType,
			Favorite: source.Favorite,
			Share:    source.Share,
			InBin:    source.DeletedAt.Time,
		}
		fileBriefs = append(fileBriefs, destination)
	}

	return fileBriefs, nil
}

func GetAllFavorFiles(userID uint) ([]models.FileBrief, error) {
	favorFileData, err := models.SearchAllFavorFile(userID)
	if err != nil {
		return nil, err
	}

	var fileBriefs []models.FileBrief

	for _, source := range favorFileData {
		destination := models.FileBrief{
			ID:       source.ID,
			FileName: source.FileName,
			FileType: source.FileType,
			Favorite: source.Favorite,
			Share:    source.Share,
			InBin:    source.DeletedAt.Time,
		}
		fileBriefs = append(fileBriefs, destination)
	}
	return fileBriefs, nil
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

// 将被移入回收站内的文件夹的子文件标记软删除
func removeSubFiles(tx *gorm.DB, userId uint, parentPath string) error {
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

		// 对文件表中的记录软删除
		if err := tx.Delete(&file).Error; err != nil {
			return err
		}
	}

	return nil
}

// 将被移出回收站内的文件夹的子文件软删除标记去除
// func recoverSubFiles(tx *gorm.DB, userId uint, parentPath string) error {
// 	parentPath = strings.ReplaceAll(strings.TrimSpace(parentPath), "\\", "/")
//
// 	var subFiles []models.FileData
// 	if err := tx.Model(models.FileData{}).Where("path LIKE ? AND belong_to = ?", parentPath+"%", userId).
// 		Find(&subFiles).Error; err != nil {
// 		return fmt.Errorf("RecoverSubFiles Error: %w", err)
// 	}
//
// 	log.Println("路径：" + parentPath)
// 	log.Println("影响文件:" + strconv.Itoa(len(subFiles)))
//
// 	// 遍历子文件进行恢复
// 	for _, file := range subFiles {
// 		if err := models.RecoverFile(file.ID); err != nil {
// 			return fmt.Errorf("RecoverSubFiles error: %w", err)
// 		}
// 	}
//
// 	return nil
// }

// 将文件移入回收站
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

// 将文件移出回收站
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
	// count := 1
	fileName := fileData.FileName
	for {
		if FileExist(userID, fileData.ParentFolderID, fileName) {
			return fmt.Errorf("RecoverFile error: 原文件夹下已有重名文件")
			// models.UnscopedRenameFile2(fileID, fileData.FileName+" ("+strconv.Itoa(count)+")")
			// fileName = fileData.FileName + " (" + strconv.Itoa(count) + ")"
			// count += 1
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
	if err := models.DeleteBinItem(fileBinItem.BinID); err != nil {
		return fmt.Errorf("RecoverFile error: %w", err)
	}

	return nil
}

// 获取用户回收站内所有文件的简讯
func GetBinFiles(userID uint) ([]models.FileBrief, error) {
	//获取这个用户在回收站内所有文件的fileBinItem的ID
	_, binItemIDs := models.SearchBinItems(userID)

	//遍历所有的这些fileBinItem，获取文件简讯
	var fileBriefs []models.FileBrief
	for _, binItemID := range binItemIDs {

		fileID := models.GetFileIDInBin(binItemID)

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

// 彻底粉碎文件
func ActuallyDeleteFile(fileID uint) error {
	// 检查回收站内是否有该文件
	if !models.FileBinItemExist(fileID) {
		return fmt.Errorf("RecoverFile error: fileRecycleItem do not exist")
	}

	fileData, err := models.GetDeletedFileData(fileID)
	if err != nil {
		return err
	}
	//文件系统中删除文件
	filePath := path.Join(config.FILE_RES_PATH, fileData.Path)

	err2 := os.Remove(filePath)
	if err2 != nil {
		fmt.Println("Error deleting file:", err)
	} else {
		fmt.Println("File successfully deleted")
	}

	// 删除对应文件回收项
	fileBinItem, err := models.DeleteFileBinItem(fileID)
	if err != nil {
		return err
	}
	// 删除对应回收站项
	if err := models.DeleteBinItem(fileBinItem.BinID); err != nil {
		return fmt.Errorf("RecoverFile error: %w", err)
	}

	_, err3 := models.UnscopedDeleteFileData(fileID)
	if err3 != nil {
		return fmt.Errorf("RecoverFile error: %w", err)
	}

	return nil
}

// 下载文件
func DownloadFile(c *gin.Context, userID uint, fileID uint) (multipart.File, error) {
	// 获取文件信息（路径），同时检查文件是否存在
	fileData, err := models.GetFileData(fileID)
	if err != nil {
		return nil, err
	}
	if fileData.ID == 0 {
		return nil, fmt.Errorf("文件不存在")
	}

	//返回文件
	filePath := path.Join(config.FILE_RES_PATH, fileData.Path)
	c.FileAttachment(filePath, fileData.FileName)

	return nil, nil
}

// 在非回收站内查找名称中包含某个子字符串的文件
func SearchFile(userID uint, pattern string) ([]models.FileData, error) {
	fileDatas, err := models.SearchFile(userID, pattern)
	if err != nil {
		return nil, err
	}
	return fileDatas, nil
}

// 查找已标记为收藏的名称中包含某个子字符串的文件
func SearchFavorFile(userID uint, pattern string) ([]models.FileData, error) {
	fileDatas, err := models.SearchFavorFile(userID, pattern)
	if err != nil {
		return nil, err
	}
	return fileDatas, nil
}

// 在回收站内查找名称中包含某个子字符串的文件
func SearchBinFile(userID uint, pattern string) ([]models.FileBrief, error) {
	_, binItemIDs := models.SearchBinItems(userID)

	var fileBriefs []models.FileBrief
	for _, binItemID := range binItemIDs {
		fileID := models.GetFileIDInBin(binItemID)

		fileData, err := models.GetDeletedFileData(fileID)
		if err != nil {
			return nil, err
		}
		// 字符串查找
		if strings.Contains(fileData.FileName, pattern) {
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

	}

	return fileBriefs, nil
}

func GetBinFileData(fileID uint) (models.FileData, error) {
	fileData, err := models.UnscopedGetFileData(fileID)
	if err != nil {
		return models.FileData{}, err
	}
	path := fileData.Path
	dirParts := strings.Split(path, "/")
	dirParts[len(dirParts)-1] = fileData.FileName
	newPath := strings.Join(dirParts[2:], "/")

	fileData.Path = newPath

	return fileData, nil
}
