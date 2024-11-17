package models

import "time"

// Bin 表用于存储被删除的文件和文件夹的基本信息
type Bin struct {
	ID          uint      `gorm:"primaryKey;autoIncrement;index;not null"`
	Type        int       `gorm:"not null;"` // file 类型或者 folder
	DeletedTime time.Time `gorm:"not null"`  // 移入回收站的时间
	UserID      uint      `gorm:"not null"`  // 用户ID，表示该回收站记录属于哪个用户

	// 外键约束，UserID -> UserData.ID
	User UserData `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
}

const (
	FOLDER = 1
	FILE   = 2
)

func (Bin) TableName() string {
	return "bins"
}

func init() {
	DataBase.AutoMigrate(&Bin{})
}
