package models

import (
	"gorm.io/gorm"
)

type UserBasic struct {
	gorm.Model
	Identity string `gorm:"coulmn:identity;type:varchar(64);" json:"identity"`         // 用户唯一标识
	Username string `gorm:"coulmn:username;type:varchar(100);" json:"username"`        // 用户名
	Usericon string `gorm:"coulmn:usericon;type:varchar(255)" json:"usericon"`         // 头像
	Password string `gorm:"coulmn:password;type:varchar(32);" json:"password"`         // 密码
	Email    string `gorm:"coulmn:email;type:varchar(50);" json:"email" valid:"email"` // 邮箱
	Age      int    `gorm:"coulmn:age;type:int;" json:"age"`                           // 年龄
	Sex      int    `gorm:"coulmn:sex;type:tinyint(1);" json:"sex"`                    // 性别  0：无性别 1：男 2：女
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}
