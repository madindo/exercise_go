package config

import (
	"../models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {

	var err error
	dsn := "madindo:madindo123@tcp(127.0.0.1:3306)/test_learningin?charset=utf8mb4&parseTime=True&loc=Local"
	DB , err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	DB.AutoMigrate(&models.User{}, &models.Article{})
	/* DB.Model(&models.Article{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	DB.Model(&models.user{}).Related(&models.Article{}) */
}