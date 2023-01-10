package dao

import (
	"treehole/models"
	"treehole/utils"
)

// 插入数据
func CreateUserNode(userNote models.UserNote) error {
	return utils.DB.Create(&userNote).Error
}
