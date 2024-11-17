package models

// FolderBin 表用于存储被删除的文件夹的基本信息
type FolderBin struct {
	ID       uint `gorm:"primaryKey;autoIncrement;not null;index;"`
	FolderID uint `gorm:"not null;"` // 被删除的文件夹 ID
	BinID    uint `gorm:"not null;"` // 对应的回收站记录 ID

	// 外键约束，指定 FolderID 和 BinID 的引用关系
	Folder FolderData `gorm:"foreignKey:FolderID;references:ID;constraint:OnDelete:CASCADE"` // 外键约束：FolderID -> FolderData.ID
	Bin    Bin        `gorm:"foreignKey:BinID;references:ID;constraint:OnDelete:CASCADE"`    // 外键约束：BinID -> Bin.ID
}

func (FolderBin) TableName() string {
	return "folder_bins"
}

func init() {
	DataBase.AutoMigrate(&FolderBin{})
}
