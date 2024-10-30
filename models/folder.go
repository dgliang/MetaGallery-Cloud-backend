package models

import (
	"time"

	"gorm.io/gorm"
)

type FolderData struct {
	ID             uint `gorm:"primaryKey;index;not null;"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
	BelongTo       string
	FolderName     string `gorm:"type:varchar(255); not null;"`
	ParentFolder   uint
	Path           string
	Favorate       bool `gorm:"index"`
	Share          bool
	IPFSInfomation string
	InBin          bool
	BinDate        time.Time

	//外键约束
	User     UserData     `gorm:"foreignKey:BelongTo;references:Account"`
	SubFiles []FolderData `gorm:"foreignKey:ParentFolder;references:ID"`
	Files    []FileData   `gorm:"foreignKey:ParentFolderID"`
}

func init() {
	DataBase.AutoMigrate(&FolderData{})
}
