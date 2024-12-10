package models

// import "time"

// type SharedFile struct {
// 	ID        uint   `gorm:"primaryKey;index;not null;auto_increment"`
// 	OwnerID   uint   `gorm:"not null;index"`
// 	FileID    uint   `gorm:"not null;index"`
// 	Intro     string `gorm:"not null;column:intro"`
// 	IPFSHash  string `gorm:"type:varchar(255);not null;column:ipfs_hash"`
// 	CreatedAt time.Time

// 	// Relationships
// 	Owner UserData   `gorm:"foreignKey:OwnerID;references:ID;constraint:OnDelete:CASCADE;"`
// 	File  FolderData `gorm:"foreignKey:FileID;references:ID;constraint:OnDelete:CASCADE;"`
// }

// func init() {
// 	DataBase.AutoMigrate(&SharedFile{})
// }

// func (SharedFile) TableName() string {
// 	return "shared_files"
// }
