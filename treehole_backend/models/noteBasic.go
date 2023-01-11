package models

import (
	"gorm.io/gorm"
)

type NoteBasic struct {
	gorm.Model
	Identity   string `gorm:"coulmn:identity;type:varchar(64);" json:"identity"`        // 帖子唯一标识
	Title      string `gorm:"coulmn:title;type:varchar(255);" json:"title"`             // 帖子标题
	Content    string `gorm:"coulmn:content;type:text" json:"content"`                  // 内容
	Urls       string `gorm:"coulmn:urls;type:text" json:"urls"`                        // 图片url，用逗号隔开
	Score      int    `gorm:"coulmn:score;type:int;" json:"score"`                      // 分数
	Approve    int    `gorm:"coulmn:approve;type:int;" json:"approve"`                  // 赞票数
	Against    int    `gorm:"coulmn:against;type:int;" json:"against"`                  // 反对票数
	Visit      int    `gorm:"coulmn:visit;type:int;" json:"visit"`                      // 访问人数
	CreateTime string `gorm:"coulmn:create_time;type:varchar(100);" json:"create_time"` // 发布时间
}

func (table *NoteBasic) TableName() string {
	return "note_basic"
}
