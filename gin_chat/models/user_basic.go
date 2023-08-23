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
	Name          string `json:"name"`
	Password      string `json:"password"`
	Phone         string `valid:"matches(^1[3-9]{1}\\d{9}$)" json:"phone"`
	Email         string `valid:"email" json:"email"`
	Identity      string `json:"identity"`
	ClientIp      string `json:"client_ip"`
	ClientPort    string `json:"client_port"`
	Salt          string `json:"salt"`
	LoginTime     int64  `json:"login_time"`
	HeartbeatTime int64  `json:"heartbeat_time"`
	LogOutTime    int64  `gorm:"column:login_out_time" json:"login_out_time"`
	IsLogout      int    `json:"is_logout"` // -1 out, 0 no out
	DeviceInfo    string `json:"device_info"`
	IsAdmin       int    `json:"is_admin"`
}

func (u *UserBasic) TableName() string {
	return "user_basic"
}
