package dao

import (
	"treehole/models"
	"treehole/utils"
)

// 根据帖子id和用户id查询投票表数据
func FindVotedNodeByIdentity(user_identity, note_identity string) (voted models.VotedNote) {
	utils.DB.Where("user_identity = ? and note_identity = ?", user_identity, note_identity).First(&voted)
	return
}

// 插入投票表数据
func CreateVotedNote(voted models.VotedNote) error {
	err := utils.DB.Create(&voted).Error
	return err
}

// 修改投票表数据
func ModifyVotedNote(voted models.VotedNote) error {
	return utils.DB.Where("note_identity = ? and user_identity = ?", voted.NoteIdentity, voted.UserIdentity).Updates(voted).Error
}
