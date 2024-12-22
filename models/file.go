package models

import (
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type FileData struct {
	ID             uint   `gorm:"primaryKey;index;not null;"`
	BelongTo       uint   `gorm:"index"`
	FileName       string `gorm:"type:varchar(256); not null; index ;"`
	FileType       string `gorm:"type:varchar(64); not null;"`
	ParentFolderID uint
	Path           string
	Favorite       bool `gorm:"index"`
	Share          bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`

	//外键约束
	User         UserData   `gorm:"foreignKey:BelongTo"`
	ParentFolder FolderData `gorm:"foreignKey:ParentFolderID;constraint:OnDelete:CASCADE;"`
}

type FileBrief struct {
	ID       uint
	FileName string
	FileType string
	Favorite bool
	Share    bool
	InBin    time.Time
}

func init() {
	DataBase.AutoMigrate(&FileData{})
}

func GetFilePath(fileID uint) (string, error) {
	var fileData FileData
	if err := DataBase.Model(&FileData{}).Where(" id = ?", fileID).Find(&fileData).Error; err != nil {
		return "", err
	}
	filePath := fileData.Path

	return filePath, nil
}

func GetParentFolderID(fileID uint) (uint, error) {
	var fileData FileData
	DataBase.Model(&FileData{}).Where(" id = ?", fileID).Find(&fileData)

	PID := fileData.ParentFolderID

	return PID, nil
}

func GenerateFilePath(userID, parentFolderID uint, fileName string) (string, error) {
	var filePath string

	parentPath, err := GetParentFolderPath(userID, parentFolderID)
	if err != nil {
		return "", err
	}

	filePath = path.Join(parentPath, fileName)

	return filePath, nil
}

func CreateFileData(userID uint, fileName string, parentFolderID uint, fileType string) (FileData, error) {
	filePath, err := GenerateFilePath(userID, parentFolderID, fileName)
	if err != nil {
		return FileData{}, err
	}

	newFile := FileData{
		BelongTo:       userID,
		FileName:       fileName,
		FileType:       fileType,
		ParentFolderID: parentFolderID,
		Path:           filePath,
	}

	if err := DataBase.Create(&newFile).Error; err != nil {
		return FileData{}, err
	}
	newFile.Path = strings.Replace(newFile.Path, fileName, strconv.FormatUint(uint64(newFile.ID), 10), 1)
	if err := DataBase.Where("id = ?", newFile.ID).Updates(&newFile).Error; err != nil {
		return FileData{}, err
	}

	return newFile, nil
}

func UnscopedDeleteFileData(fileID uint) (FileData, error) {
	var deletedData FileData
	err := DataBase.Model(&FileData{}).Unscoped().Where("id = ?", fileID).Delete(&deletedData).Error
	return deletedData, err
}

func RenameFileWithFileID(fileID uint, newFileName string) error {
	File := FileData{
		ID: fileID,
	}
	var originFileData FileData
	DataBase.Model(&FileData{}).Where("id = ?", fileID).First(&originFileData)
	newType := filepath.Ext(newFileName)
	DataBase.Model(&File).Where("ID = ?", fileID).Updates(FileData{FileName: newFileName, FileType: newType})
	return nil
}

func UpdateFileType(fileID uint, newFileName string) error {
	return nil
}

func UnscopedRenameFile(fileID uint, newFileName string) error {
	File := FileData{
		ID: fileID,
	}
	var originFileData FileData
	DataBase.Model(&FileData{}).Unscoped().Where("id = ?", fileID).First(&originFileData)

	DataBase.Model(&File).Unscoped().Where("ID = ?", fileID).Updates(FileData{FileName: newFileName})
	return nil
}

func GetFileID(userId, parentId uint, fileName string) (uint, error) {
	var file FileData

	if err := DataBase.Where("belong_to = ? AND parent_folder_id = ? AND file_name = ?",
		userId, parentId, fileName).First(&file).Error; err != nil {
		return 0, nil
	}
	return file.ID, nil
}

func GetFileData(fileID uint) (FileData, error) {
	var fileData FileData
	// 预加载
	DataBase.Preload("ParentFolder").Model(&FileData{ID: fileID}).Where("id = ?", fileID).Find(&fileData)

	return fileData, nil
}

func UnscopedGetFileData(fileID uint) (FileData, error) {
	var fileData FileData
	// 预加载
	DataBase.Preload("ParentFolder").Model(&FileData{ID: fileID}).Unscoped().Where("id = ?", fileID).Find(&fileData)

	return fileData, nil
}

func GetDeletedFileData(fileID uint) (FileData, error) {
	var fileData FileData
	// 预加载
	DataBase.Preload("User").Preload("ParentFolder").Unscoped().Model(&FileData{ID: fileID}).Where("id = ?", fileID).Find(&fileData)

	return fileData, nil
}

func GetSubFiles(parentFolderID uint) ([]FileData, error) {
	var subFiles []FileData

	if err := DataBase.Set("gorm:auto_preload", false).Model(&FileData{}).Where("parent_folder_id = ?", parentFolderID).Find(&subFiles).Error; err != nil {
		return nil, err
	}
	return subFiles, nil
}

func GetSubFileDatas(parentFolderID uint) ([]FileData, error) {
	var subFiles []FileData

	DataBase.Set("gorm:auto_preload", false).Model(&FileData{}).Where("parent_folder_id = ?", parentFolderID).Find(&subFiles)

	return subFiles, nil
}

func SetFileFavorite(fileID uint) {
	DataBase.Model(&FileData{}).Where("id = ?", fileID).Updates(FileData{
		Favorite: true,
	})
}

func CancelFileFavorite(fileID uint) {
	var file FileData
	DataBase.Model(&FileData{}).Where("id = ?", fileID).Find(&file)
	file.Favorite = false
	DataBase.Save(&file)
}

func SearchAllFavorFile(userID uint) ([]FileData, error) {
	var favorFiles []FileData
	err := DataBase.Model(&FileData{}).Where("belong_to = ? and favorite = ?", userID, true).Find(&favorFiles).Error
	if err != nil {
		return nil, err
	}
	return favorFiles, nil
}

func RemoveFile(fileID uint) error {
	err := DataBase.Model(&FileData{}).Delete(&FileData{ID: fileID}).Error
	return err
}

func GetBinFiles(userID uint) ([]FileBrief, error) {
	var binFiles []FileData

	DataBase.Set("gorm:auto_preload", false).Model(&FileData{}).Unscoped().Where("belong_to = ? and deleted_at IS NOT NULL", userID).Find(&binFiles)

	var fileBriefs []FileBrief

	for _, source := range binFiles {
		destination := FileBrief{
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

func RecoverFile(fileID uint) error {
	err := DataBase.Model(&FileData{}).Unscoped().Where("id = ?", fileID).Update("deleted_at", nil).Error
	return err
}

func SearchFile(userID uint, pattern string) ([]FileData, error) {
	var fileDatas []FileData
	err := DataBase.Model(&FileData{}).Where("belong_to = ? and file_name LIKE ?", userID, "%"+pattern+"%").Find(&fileDatas)
	if err != nil {
		return nil, err.Error
	}
	return fileDatas, nil
}

func SearchFavorFile(userID uint, pattern string) ([]FileData, error) {
	var fileDatas []FileData
	err := DataBase.Model(&FileData{}).Where("belong_to = ? and favorite = 1 and file_name LIKE ?", userID, "%"+pattern+"%").Find(&fileDatas)
	if err != nil {
		return nil, err.Error
	}
	return fileDatas, nil
}
