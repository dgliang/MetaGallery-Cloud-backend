package models

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type FolderData struct {
	ID              uint   `gorm:"primaryKey;index;not null;autoIncrement"`
	BelongTo        uint   `gorm:"default:NULL"`
	FolderName      string `gorm:"type:varchar(255); not null;"`
	ParentFolder    uint   `gorm:"default:NULL"`
	Path            string
	Favorite        bool `gorm:"index; default:false"`
	Share           bool `gorm:"default:false"`
	IPFSInformation string
	InBin           bool      `gorm:"default:false"`
	BinDate         time.Time `gorm:"default:NULL"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`

	//外键约束
	User       UserData     `gorm:"foreignKey:BelongTo;"`
	SubFolders []FolderData `gorm:"foreignKey:ParentFolder;references:ID;constraint:OnDelete:CASCADE;"`
	Files      []FileData   `gorm:"foreignKey:ParentFolderID;constraint:OnDelete:CASCADE;"`
}

func init() {
	DataBase.AutoMigrate(&FolderData{})
}

func CreateFolder(userID, parentID uint, folderName, path string) (FolderData, error) {
	newFolder := FolderData{
		BelongTo:     userID,
		FolderName:   folderName,
		ParentFolder: parentID,
		Path:         path,
	}

	if err := DataBase.Create(&newFolder).Error; err != nil {
		return FolderData{}, err
	}
	return newFolder, nil
}

func CreateRootFolder(userId uint, folderName, path string) (FolderData, error) {
	newFolder := FolderData{
		BelongTo:   userId,
		FolderName: folderName,
		Path:       path,
	}

	if err := DataBase.Create(&newFolder).Error; err != nil {
		return FolderData{}, err
	}
	return newFolder, nil
}

func GetRootFolderData(userID uint) (FolderData, error) {
	var folder FolderData

	folderName := fmt.Sprintf("%d", userID)
	if err := DataBase.Where("belong_to = ? AND folder_name = ?", userID, folderName).
		First(&folder).Error; err != nil {
		return FolderData{}, err
	}
	return folder, nil
}

func GetParentFolderPath(userID, parentID uint) (string, error) {
	var folderData FolderData

	if err := DataBase.Where("belong_to = ? AND id = ?", userID, parentID).
		First(&folderData).Error; err != nil {
		return "", err
	}
	return folderData.Path, nil
}
