package dao

import (
	"treehole/models"
	"treehole/utils"
)

// 创建帖子数据
func CreateNote(note models.NoteBasic) error {
	return utils.DB.Create(&note).Error
}
