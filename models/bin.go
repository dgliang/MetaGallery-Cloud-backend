package models

import (
	"time"
)

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

func InsertBinItem(binItem Bin) (Bin, error) {
	if err := DataBase.Create(&binItem).Error; err != nil {
		return Bin{}, err
	}
	return binItem, nil
}

func DeleteBinItem(binItemID uint) error {
	if err := DataBase.Model(&Bin{}).Where("id = ?", binItemID).Delete(&Bin{ID: binItemID}).Error; err != nil {
		return err
	}
	return nil
}

func SearchBinItems(userID uint) ([]uint, []uint) {
	var binItems []Bin
	DataBase.Model(&Bin{}).Where("user_id = ?", userID).Find(&binItems)

	var folderItemID, fileItemID []uint
	for _, binItem := range binItems {
		// log.Printf("from 用户ID：%d 查询到回收项： %d %d %s %d\n", userID, binItem.ID, binItem.Type, binItem.DeletedTime, binItem.UserID)
		if binItem.Type == 1 {
			folderItemID = append(folderItemID, binItem.ID)
		} else {
			fileItemID = append(fileItemID, binItem.ID)
		}
	}

	return folderItemID, fileItemID
}
