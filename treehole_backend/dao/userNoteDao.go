package dao

import (
	"treehole/models"
	"treehole/utils"

	"gorm.io/gorm"
)

// 插入数据
func CreateUserNode(userNote models.UserNote) error {
	return utils.DB.Create(&userNote).Error
}

// 根据帖子id查询，作者id
func FindUserNoteByNoteIdentity(note_identity string) (models.UserNote, error) {
	author := models.UserNote{}
	err := utils.DB.Where("note_identity = ?", note_identity).First(&author).Error
	return author, err
}

// 根据作者id查询，帖子id
func FindUserNoteByAuthorIdentity(author string) *gorm.DB {
	tx := utils.DB.Model(new(models.UserNote)).Preload("NoteBasic").Where("author_identity = ?", author)
	return tx
}

// 根据帖子id查询帖子的详细信息
func FindUserNoteByNoteIdentityFind(note_identity string) (user models.UserNote, err error) {
	err = utils.DB.Model(new(models.UserNote)).Preload("NoteBasic").
		Where("note_identity = ?", note_identity).First(&user).Error
	return
}

// 根据identity删除帖子数据
func DeleteUserNote(identity string) error {
	return utils.DB.Where("note_identity = ?", identity).Delete(new(models.UserNote)).Error
}
