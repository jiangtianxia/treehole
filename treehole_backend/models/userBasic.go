package models

import "gorm.io/gorm"

type UserBasic struct {
	gorm.Model
	Identity string `gorm:"coulmn:identity;type:varchar(64);" json:"identity"`         // 用户唯一标识
	Username string `gorm:"coulmn:username;type:varchar(100);" json:"username"`        // 用户名
	Password string `gorm:"coulmn:password;type:varchar(32);" json:"password"`         // 密码
	Email    string `gorm:"coulmn:email;type:varchar(50);" json:"email" valid:"email"` // 邮箱
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}
