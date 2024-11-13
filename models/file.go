package models

import (
	"time"

	"gorm.io/gorm"
)

type FileData struct {
	ID             uint `gorm:"primaryKey;index;not null;"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
	BelongTo       uint
	FileName       string `gorm:"type:varchar(255); not null;"`
	FileType       string `gorm:"type:varchar(63); not null;"`
	ParentFolderID uint
	Path           string
	Favorite       bool `gorm:"index"`
	Share          bool
	IPFSInfomation string
	InBin          bool
	BinDate        time.Time

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
	InBin    bool
}

func init() {
	DataBase.AutoMigrate(&FileData{})
}

func GetFilePath(userID uint, fileName string, parentFolderID uint) (string, error) {
	parentFloderPath, err := GetParentFolderPath(userID, parentFolderID)
	if err != nil {
		return "", err
	}
	filepath := parentFloderPath + "/" + fileName

	return filepath, nil
}

func CreateFileData(userID uint, fileName string, parentFolderID uint) (FileData, error) {
	filePath, err := GetFilePath(userID, fileName, parentFolderID)
	if err != nil {
		return FileData{}, err
	}

	newFile := FileData{
		BelongTo:       userID,
		FileName:       fileName,
		ParentFolderID: parentFolderID,
		Path:           filePath,
	}

	if err := DataBase.Create(&newFile).Error; err != nil {
		return FileData{}, err
	}
	return newFile, nil
}

func DeleteWithFileData(fileID uint, userID uint, fileName string, parentFolderID uint) (FileData, error) {
	fileToDelete := FileData{
		ID:             fileID,
		BelongTo:       userID,
		FileName:       fileName,
		ParentFolderID: parentFolderID,
	}

	if err := DataBase.Delete(&fileToDelete).Error; err != nil {
		return FileData{}, err
	}

	return fileToDelete, nil
}

func DeleteWithFileID(fileID uint) (FileData, error) {
	fileToDelete := FileData{
		ID: fileID,
	}

	if err := DataBase.Delete(&fileToDelete).Error; err != nil {
		return FileData{}, err
	}

	return fileToDelete, nil
}

func DeleteWithFileName(userID uint, fileName string, parentFolderID uint) (FileData, error) {
	fileToDelete := FileData{
		BelongTo:       userID,
		FileName:       fileName,
		ParentFolderID: parentFolderID,
	}

	if err := DataBase.Delete(&fileToDelete).Error; err != nil {
		return FileData{}, err
	}

	return fileToDelete, nil
}

func RenameFile(userID uint, oldfilename string, newFileName string, parentFolderID uint) error {
	var originFileData FileData
	DataBase.Model(&FileData{}).Where("belong_to = ? AND parent_folder_id = ? AND file_name = ?", userID, parentFolderID, oldfilename).First(&originFileData)

	newFilePath, err := GetFilePath(originFileData.BelongTo, newFileName, originFileData.ParentFolderID)
	if err != nil {
		return err
	}

	DataBase.Model(&FileData{}).Where("belong_to = ? AND parent_folder_id = ? AND file_name = ?", userID, parentFolderID, oldfilename).Updates(FileData{FileName: newFileName, Path: newFilePath})
	return nil
}

func RenameFileWithFileID(fileID uint, newFileName string) error {
	File := FileData{
		ID: fileID,
	}
	var originFileData FileData
	DataBase.Model(&File).First(&originFileData)

	newFilePath, err := GetFilePath(originFileData.BelongTo, newFileName, originFileData.ParentFolderID)
	if err != nil {
		return err
	}

	DataBase.Model(&File).Where("ID = ?", fileID).Updates(FileData{FileName: newFileName, Path: newFilePath})
	return nil
}

func GetFileData(fileID uint) (FileData, error) {
	var fileData FileData
	// 预加载
	DataBase.Preload("User").Preload("ParentFolder").Model(&FileData{ID: fileID}).Find(&fileData)

	return fileData, nil
}

func GetSubFiles(parentFolderID uint) ([]FileBrief, error) {
	var subFiles []FileData

	DataBase.Set("gorm:auto_preload", false).Model(&FileData{}).Where("parent_folder_id = ?", parentFolderID).Find(&subFiles)

	var fileBriefs []FileBrief

	for _, source := range subFiles {
		destination := FileBrief{
			ID:       source.ID,
			FileName: source.FileName,
			FileType: source.FileType,
			Favorite: source.Favorite,
			Share:    source.Share,
			InBin:    source.InBin,
		}
		fileBriefs = append(fileBriefs, destination)
	}

	return fileBriefs, nil
}
