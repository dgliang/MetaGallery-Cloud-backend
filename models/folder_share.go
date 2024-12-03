package models

import "time"

type SharedFolder struct {
	ID         uint   `gorm:"primaryKey;index;not null;auto_increment"`
	OwnerID    uint   `gorm:"not null;index"`
	FolderID   uint   `gorm:"not null;index"`
	SharedName string `gorm:"type:varchar(255);not null;column:shared_name"`
	Intro      string `gorm:"not null;column:intro"`
	IPFSHash   string `gorm:"type:varchar(255);not null;column:ipfs_hash"`
	CreatedAt  time.Time

	// Relationships
	Owner  UserData   `gorm:"foreignKey:OwnerID;references:ID;constraint:OnDelete:CASCADE;"`
	Folder FolderData `gorm:"foreignKey:FolderID;references:ID;constraint:OnDelete:CASCADE;"`
}

func init() {
	DataBase.AutoMigrate(&SharedFolder{})
}

func (SharedFolder) TableName() string {
	return "shared_folders"
}
