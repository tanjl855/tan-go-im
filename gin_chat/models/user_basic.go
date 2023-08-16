package models

import (
	"gorm.io/gorm"
	"time"
)

type UserBasic struct {
	gorm.Model
	Name          string
	Password      string
	Phone         string
	Email         string
	Identity      string
	ClientIp      string
	ClientPort    string
	LoginTime     time.Time
	HeartbeatTime time.Time
	LogOutTime    time.Time `gorm:"column:login_out_time" json:"login_out_time"`
	IsLogout      int       // -1 out, 0 no out
	DeviceInfo    string
	IsAdmin       int
}

func (u *UserBasic) TableName() string {
	return "user_basic"
}
