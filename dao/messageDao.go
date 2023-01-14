package dao

import (
	"treehole/models"
	"treehole/utils"

	"gorm.io/gorm"
)

func CreateMessageBasic(messageBasic models.MessageBasic) error {
	return utils.DB.Create(&messageBasic).Error
}

// 查询数据库信息
func FindMessageBasic() *gorm.DB {
	tx := utils.DB.Model(new(models.MessageBasic))

	return tx.Order("create_time desc")
}
