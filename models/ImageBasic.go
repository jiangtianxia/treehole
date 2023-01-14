package models

import "gorm.io/gorm"

type ImageBasic struct {
	gorm.Model
	Identity string `gorm:"coulmn:identity;type:varchar(64);" json:"identity"` // 用户唯一标识
	Url      string `gorm:"coulmn:url;type:varchar(255);" json:"url"`          // 图片url
	Type     int    `gorm:"coulmn:type;type:tinyint(1);" json:"type"`          // 图片类型   1：头像图片
}

func (table *ImageBasic) TableName() string {
	return "image_basic"
}
