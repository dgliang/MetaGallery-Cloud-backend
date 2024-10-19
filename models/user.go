package models

import (
	"time"
)

type User_Data struct {
	Account       string    `gorm:"type:varchar(20); not null; primarykey" json:"account" binding:"required"`
	Password      string    `gorm:"type:varchar(20); not null;" json:"password" binding:"required"`
	Brief_Intro   *string   `gorm:"type:text;" json:"Brief_Intro"`
	Profile_Photo *string   `gorm:"type:text;" json:"Profile_Photo"`
	CreatedAt     time.Time // 创建时间（由GORM自动管理）
	UpdatedAt     time.Time // 最后一次更新时间（由GORM自动管理）
}
