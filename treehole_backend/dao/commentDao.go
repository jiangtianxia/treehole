package dao

import (
	"treehole/models"
	"treehole/utils"
)

// 插入数据
func CreateComment(comment models.CommentNote) error {
	return utils.DB.Create(&comment).Error
}

// 查询数据
func FindCommentByNoteIdentity(identity string) (list []models.CommentNote, err error) {
	err = utils.DB.Where("note_identity = ?", identity).Preload("UserBasic").Order("create_time desc").Find(&list).Error
	return
}

// 查询数据
func FindCommentByUserIdentity(user_identity string) (list []models.CommentNote, err error) {
	err = utils.DB.Where("user_identity = ?", user_identity).Preload("NoteBasic").Order("create_time desc").Find(&list).Error
	return
}

// 删除数据
func DeleteNoteComment(identity string, note_identity string) error {
	return utils.DB.Where("identity = ? and note_identity = ?", identity, note_identity).Delete(new(models.CommentNote)).Error
}
