package models

import (
	"gorm.io/gorm"
)

type MessageBasic struct {
	gorm.Model
	Identity     string `gorm:"coulmn:identity;type:varchar(64);" json:"identity"`           // 留言的唯一标识
	UserIdentity string `gorm:"coulmn:user_identity;type:varchar(64);" json:"user_identity"` // 用户唯一标识
	Username     string `gorm:"coulmn:username;type:varchar(100);" json:"username"`          // 用户名
	Usericon     string `gorm:"coulmn:usericon;type:varchar(255)" json:"usericon"`           // 头像
	Date         string `gorm:"coulmn:date;type:text" json:"date"`                           // 内容
	CreateTime   string `gorm:"coulmn:create_time;type:varchar(100);" json:"create_time"`    // 发布时间
}

/**
 * @Author jiang
 * @Description 信息参数
 * @Date 13:00 2023/1/13
 **/
type MessageStruct struct {
	Type         int    `json:"type"`          // 发送类型。-1：心跳检测。0：消息发送
	UserIdentity string `json:"user_identity"` // 用户唯一标识
	Username     string `json:"username"`      // 用户名
	Usericon     string `json:"usericon"`      // 头像
	Date         string `json:"date"`          // 内容
	CreateTime   string `json:"create_time"`   // 发布时间
}

func (table *MessageBasic) TableName() string {
	return "message_basic"
}
