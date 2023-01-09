package dao

import (
	"treehole/models"
	"treehole/utils"
)

// 创建图片
func CreateImage(imagebasic models.ImageBasic) error {
	return utils.DB.Create(&imagebasic).Error
}

// 查询图片
func FindImageByIdentity(identity, url string) ([]*models.ImageBasic, error) {
	ImageList := make([]*models.ImageBasic, 0)
	err := utils.DB.Where("identity = ? AND type = 1 AND url <> ?", identity, url).Find(&ImageList).Error
	return ImageList, err
}

// 删除数据库图片信息
func DeleteImages(imageList []*models.ImageBasic) error {
	return utils.DB.Delete(imageList).Error
}
