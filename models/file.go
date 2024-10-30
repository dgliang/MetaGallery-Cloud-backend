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
	BelongTo       string
	FileName       string `gorm:"type:varchar(255); not null;"`
	FileType       string `gorm:"type:varchar(63); not null;"`
	ParentFolderID uint
	Path           string
	Favorate       bool `gorm:"index"`
	Share          bool
	IPFSInfomation string
	InBin          bool
	BinDate        time.Time

	//外键约束
	User         UserData   `gorm:"foreignKey:BelongTo"`
	ParentFolder FolderData `gorm:"foreignKey:ParentFolderID"`
}

// type User struct {
// 	gorm.Model
// 	CreditCards []CreditCard `gorm:"foreignKey:UserRefer"`
// }

// type CreditCard struct {
// 	gorm.Model
// 	Number    string
// 	UserRefer uint
// }

func init() {
	DataBase.AutoMigrate(&FileData{})
}
