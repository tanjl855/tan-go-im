package chatdb

import (
	"fmt"
	"gin_chat/models"
	"gin_chat/utils"
)

func GetUserList() []*models.UserBasic {
	data := make([]*models.UserBasic, 10)
	utils.DB.Find(&data)
	for _, v := range data {
		fmt.Println(v)
	}
	return data
}
