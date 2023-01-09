package dao

import (
	"treehole/models"
	"treehole/utils"
)

// 查询数据库判断该identity是否已经存在
func FindUserByIdentityCount(identity string) (cnt int64, err error) {
	err = utils.DB.Where("identity = ?", identity).Model(new(models.UserBasic)).Count(&cnt).Error
	return
}

// 查询数据库判断该邮箱是否已经存在
func FindUserByEmailCount(email string) (cnt int64, err error) {
	err = utils.DB.Where("email = ?", email).Model(new(models.UserBasic)).Count(&cnt).Error
	return
}

// 查询数据库判断该用户名是否已经注册
func FindUserByNameCount(username string) (cnt int64, err error) {
	err = utils.DB.Where("username = ?", username).Model(new(models.UserBasic)).Count(&cnt).Error
	return
}

// 插入数据
func CreateUser(user models.UserBasic) error {
	err := utils.DB.Create(&user).Error
	return err
}

// 根据用户名查询信息
func FindUserByName(username string) (models.UserBasic, error) {
	user := models.UserBasic{}
	err := utils.DB.Where("username = ?", username).First(&user).Error
	return user, err
}

// 根据邮箱查询信息
func FindByEmail(email string) (user models.UserBasic, err error) {
	err = utils.DB.Where("email = ?", email).First(&user).Error
	return
}

// 根据identidy查询信息
func FindByIdentity(identity string) (user models.UserBasic, err error) {
	err = utils.DB.Where("identity = ?", identity).First(&user).Error
	return
}

// 修改用户信息
func ModifyUserInfo(user models.UserBasic) error {
	return utils.DB.Where("identity = ?", user.Identity).Updates(user).Error
}

// 判断该用户名是否已经被注册了
func FindUserByNameAndIdentity(username, identity string) (cnt int64, err error) {
	err = utils.DB.Model(new(models.UserBasic)).Where("username = ? and identity != ?", username, identity).Count(&cnt).Error
	return
}
