package models

import "gorm.io/gorm"

type VotedNote struct {
	gorm.Model
	UserIdentity string `gorm:"coulmn:user_identity;type:varchar(64);" json:"user_identity"` // 用户唯一标识
	NoteIdentity string `gorm:"coulmn:note_identity;type:varchar(64);" json:"note_identity"` // 帖子唯一标识
	Isvoted      string `gorm:"coulmn:isvoted;type:varchar(5);" json:"isvoted"`              // 是否投票，0：未点赞，也未踩，1：点赞，但未踩，-1：踩
}

func (table *VotedNote) TableName() string {
	return "voted_note"
}
