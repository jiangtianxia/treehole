package main

import (
	"treehole/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(mysql.Open("root:124578@tcp(127.0.0.1:3306)/treehole?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 迁移 schema
	// db.AutoMigrate(&models.UserBasic{})
	// db.AutoMigrate(&models.ImageBasic{})
	// db.AutoMigrate(&models.NoteBasic{})
	// db.AutoMigrate(&models.UserNote{})
	// db.AutoMigrate(&models.VotedNote{})
	// db.AutoMigrate(&models.CommentNote{})
	db.AutoMigrate(&models.MessageBasic{})
}
