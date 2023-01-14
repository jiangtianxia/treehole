package models

import "gorm.io/gorm"

type UserNote struct {
	gorm.Model
	AuthorIdentity string     `gorm:"coulmn:author_identity;type:varchar(64);" json:"author_identity"` // 用户唯一标识
	NoteIdentity   string     `gorm:"coulmn:note_identity;type:varchar(64);" json:"note_identity"`     // 帖子唯一标识
	NoteBasic      *NoteBasic `gorm:"foreignKey:identity;references:note_identity;" json:"note_basic"` // 关联帖子的基础信息表
}
