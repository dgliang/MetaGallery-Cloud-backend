package models

import (
	"path"
	"time"

	"gorm.io/gorm"
)

type FileData struct {
	ID             uint `gorm:"primaryKey;index;not null;"`
	BelongTo       uint
	FileName       string `gorm:"type:varchar(256); not null;"`
	FileType       string `gorm:"type:varchar(64); not null;"`
	ParentFolderID uint
	Path           string
	Favorite       bool `gorm:"index"`
	Share          bool
	IPFSInfomation string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
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
	InBin    time.Time
}

func init() {
	DataBase.AutoMigrate(&FileData{})
}

// func GetFilePath(userID uint, fileName string, parentFolderID uint) (string, error) {
// 	parentFloderPath, err := GetParentFolderPath(userID, parentFolderID)
// 	if err != nil {
// 		return "", err
// 	}
// 	filepath := parentFloderPath + "/" + fileName

// 	return filepath, nil
// }

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

func CreateFileData(userID uint, fileName string, parentFolderID uint) (FileData, error) {
	filePath, err := GenerateFilePath(userID, parentFolderID, fileName)
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

	newFilePath, err := GenerateFilePath(originFileData.BelongTo, originFileData.ParentFolderID, newFileName)
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
	DataBase.Model(&FileData{}).Where("id = ?", fileID).First(&originFileData)

	newFilePath, err := GenerateFilePath(originFileData.BelongTo, originFileData.ParentFolderID, newFileName)
	if err != nil {
		return err
	}

	DataBase.Model(&File).Where("ID = ?", fileID).Updates(FileData{FileName: newFileName, Path: newFilePath})
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
			InBin:    source.DeletedAt.Time,
		}
		fileBriefs = append(fileBriefs, destination)
	}

	return fileBriefs, nil
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
