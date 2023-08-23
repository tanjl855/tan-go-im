package chatdb

import (
	"errors"
	"fmt"
	"gin_chat/models"
	"gin_chat/utils"
	"gorm.io/gorm"
	"time"
)

func CreateUser(user *models.UserBasic) (*gorm.DB, error) {
	if user == nil {
		fmt.Printf("user is nil")
		return nil, errors.New("nil user")
	}
	return utils.DB.Create(user), nil
}

func FindUserByNameAndPwd(name, password string) models.UserBasic {
	user := models.UserBasic{}
	utils.DB.Where("Name = ? and password = ?", name, password).First(&user)

	// token加密
	str := fmt.Sprintf("%d", time.Now().Unix())
	temp := utils.Md5Encode(str)

	utils.DB.Model(&user).Where("id = ?", user.ID).Update("identity", temp)
	return user
}

func FindUserByName(name string) models.UserBasic {
	user := models.UserBasic{}
	utils.DB.Where("Name = ?", name).First(&user)
	return user
}

func FindUserByPhone(phone string) *gorm.DB {
	user := models.UserBasic{}
	return utils.DB.Where("Phone = ?", phone).First(&user)
}

func FindUserByEmail(email string) *gorm.DB {
	user := models.UserBasic{}
	return utils.DB.Where("Email = ?", email).First(&user)
}

func GetUserList() []*models.UserBasic {
	data := make([]*models.UserBasic, 10)
	utils.DB.Find(&data)
	for _, v := range data {
		fmt.Println(v)
	}
	return data
}

func DeleteUser(user *models.UserBasic) *gorm.DB {
	return utils.DB.Delete(user)
}

func UpdateUser(user *models.UserBasic) *gorm.DB {
	return utils.DB.Model(&user).Updates(models.UserBasic{Name: user.Name, Password: user.Password, Phone: user.Phone, Email: user.Email})
}
