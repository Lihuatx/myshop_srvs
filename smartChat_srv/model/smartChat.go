package model

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID       int32     `gorm:"primarykey"`
	CreateAt time.Time `gorm:"column:add_time;autoCreateTime"`
	UpdateAt time.Time `gorm:"column:update_time;autoUpdateTime"`
	DeleteAt gorm.DeletedAt
	IsDelete bool
}

type ChatLog struct {
	BaseModel
	UserID   int32  `gorm:"type:int;index"`
	Question string `gorm:"type:text"`
	Answer   string `gorm:"type:text"`
	Intent   string `gorm:"type:varchar(20)"`
}
