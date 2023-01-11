package dao

import (
	"treehole/models"
	"treehole/utils"

	"gorm.io/gorm"
)

// 创建帖子数据
func CreateNote(note models.NoteBasic) error {
	return utils.DB.Create(&note).Error
}

// 根据关键词查询帖子信息
func SearchNotes(keyword string) *gorm.DB {
	tx := utils.DB.Model(new(models.NoteBasic)).Where("title like ? or content like ?", "%"+keyword+"%", "%"+keyword+"%")

	return tx.Order("score desc")
}

// 查询全部帖子信息
func GetNoteList(identity string) (models.NoteBasic, error) {
	res := models.NoteBasic{}

	err := utils.DB.Model(new(models.NoteBasic)).Where("identity = ?", identity).First(&res).Error
	return res, err
}

// 根据identity删除帖子数据
func DeleteNote(identity string) error {
	return utils.DB.Where("identity = ?", identity).Delete(new(models.NoteBasic)).Error
}

// 修改帖子访问量
func UpdateNote(note models.NoteBasic) error {
	return utils.DB.Where("identity = ?", note.Identity).Updates(note).Error
}
