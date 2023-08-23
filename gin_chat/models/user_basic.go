package models

import (
	"gorm.io/gorm"
)

var UserStatus = struct {
	Normal, ParamsError, IntervalError int
}{
	Normal:        0,
	ParamsError:   2,
	IntervalError: 3,
}

type UserBasic struct {
	gorm.Model
	Name          string
	Password      string
	Phone         string `valid:"matches(^1[3-9]{1}\\d{9}$)"`
	Email         string `valid:"email"`
	Identity      string
	ClientIp      string
	ClientPort    string
	Salt          string
	LoginTime     int64
	HeartbeatTime int64
	LogOutTime    int64 `gorm:"column:login_out_time" json:"login_out_time"`
	IsLogout      int   // -1 out, 0 no out
	DeviceInfo    string
	IsAdmin       int
}

func (u *UserBasic) TableName() string {
	return "user_basic"
}
