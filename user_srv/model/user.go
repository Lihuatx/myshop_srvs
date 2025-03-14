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

type User struct {
	BaseModel
	Mobile   string     `gorm:"index:idx_mobile;unique;not null;type:varchar(11)"`
	Password string     `gorm:"type:varchar(100);not null"`
	NickName string     `gorm:"type:varchar(20)"`
	Birthday *time.Time `gorm:"type:datetime"`
	Gender   string     `gorm:"column:gender;default:male;type:varchar(6) comment 'female表示女, male表示男'"`
	Role     int        `gorm:"column:role;default:1;type:int comment '1表示普通用户，2表示管理员'"`
}
