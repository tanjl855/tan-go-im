package main

import (
	"fmt"
	"gin_chat/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(mysql.Open("root:123@tcp(127.0.0.1:3307)/ginchat?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		fmt.Printf("[E] connect mysql error:%v", err)
		return
	}
	db.AutoMigrate(&models.UserBasic{})

	user := &models.UserBasic{}
	user.Name = "tanjl"

	db.Create(user)

	fmt.Println(db.First(user, 1))

	db.Model(user).Update("PassWord", "123")

	fmt.Println(db.First(user, 1))
}
