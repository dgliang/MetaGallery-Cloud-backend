package models

type SharedFolder struct {
	ID       uint   `gorm:"primaryKey;index;not null;auto_increment"`
	OwnerID  uint   `gorm:"not null;index"`
	FolderID uint   `gorm:"not null;index"`
	Intro    string `gorm:"not null;column:intro"`
	IPFSHash string `gorm:"type:varchar(255);not null;column:ipfs_hash"`

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
