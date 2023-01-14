package models

import (
	"gorm.io/gorm"
)

type CommentNote struct {
	gorm.Model
	Identity     string     `gorm:"coulmn:identity;type:varchar(64);" json:"identity"`               // 留言的唯一标识
	UserIdentity string     `gorm:"coulmn:user_identity;type:varchar(64);" json:"user_identity"`     // 用户唯一标识
	UserBasic    *UserBasic `gorm:"foreignKey:identity;references:user_identity;" json:"user_basic"` // 关联用户的基础信息表
	NoteIdentity string     `gorm:"coulmn:note_identity;type:varchar(64);" json:"note_identity"`     // 帖子唯一标识
	NoteBasic    *NoteBasic `gorm:"foreignKey:identity;references:note_identity;" json:"note_basic"` // 关联帖子的基础信息表
	Conetent     string     `gorm:"coulmn:content;type:text" json:"content"`                         // 内容
	CreateTime   string     `gorm:"coulmn:create_time;type:varchar(100);" json:"create_time"`        // 发布时间
}

func (table *CommentNote) TableName() string {
	return "comment_note"
}
