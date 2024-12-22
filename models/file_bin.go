package models

// FileBin 表用于存储被删除的文件夹的基本信息
type FileBin struct {
	ID     uint `gorm:"primaryKey;autoIncrement;not null;index;"`
	FileID uint `gorm:"not null;"` // 被删除的文件 ID
	BinID  uint `gorm:"not null;"` // 对应的回收站记录 ID

	// 外键约束，指定 FileID 和 BinID 的引用关系
	File FileData `gorm:"foreignKey:FileID;references:ID;"`                           // FileID -> FolderData.ID
	Bin  Bin      `gorm:"foreignKey:BinID;references:ID;constraint:OnDelete:CASCADE"` // BinID -> Bin.ID
}

func (FileBin) TableName() string {
	return "file_bins"
}

func init() {
	DataBase.AutoMigrate(&FileBin{})
}

func InsertFileBinItem(fileBinItem FileBin) error {
	if err := DataBase.Create(&fileBinItem).Error; err != nil {
		return err
	}
	return nil
}

func DeleteFileBinItem(fileID uint) (FileBin, error) {
	var deletedItem FileBin
	if err := DataBase.Model(&FileBin{}).Where("file_id = ?", fileID).Find(&deletedItem).Error; err != nil {
		return FileBin{}, err
	}
	if err := DataBase.Model(&FileBin{}).Where("file_id = ?", fileID).Delete(&deletedItem).Error; err != nil {
		return FileBin{}, err
	}
	return deletedItem, nil
}

func FileBinItemExist(fileID uint) bool {
	var Item FileBin
	if err := DataBase.Model(&FileBin{}).Where("file_id = ?", fileID).Find(&Item).Error; err != nil {
		return false
	}

	if Item.BinID != 0 {
		return true
	}
	return false
}

func GetFileIDInBin(binItemID uint) uint {
	var Item FileBin
	if err := DataBase.Model(&FileBin{}).Where("bin_id = ?", binItemID).Find(&Item).Error; err != nil {
		return 0
	}
	return Item.FileID
}
