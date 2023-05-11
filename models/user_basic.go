package models

import (
	"fmt"
	"ginchat/utils"
	"gorm.io/gorm"
)

type UserBasic struct {
	gorm.Model
	Name          string `gorm:"column:name" json:"name"`
	Password      string `gorm:"column:password" json:"password"`
	Phone         string `gorm:"column:phone" json:"phone"`
	Email         string `gorm:"column:email" json:"email"`
	Identity      string `gorm:"column:identity" json:"identity"`
	ClientIp      string `gorm:"column:client_ip" json:"client_ip"`
	ClientPort    string `gorm:"column:client_port" json:"client_port"`
	LoginTime     uint64 `gorm:"column:login_time" json:"login_time"`
	HeartbeatTime uint64 `gorm:"column:heartbeat_time" json:"heartbeat_time"`
	LoginOutTime  uint64 `gorm:"column:login_out_time" json:"login_out_time"`
	IsLogout      bool   `gorm:"column:is_logout" json:"is_logout"`
	DeviceInfo    string `gorm:"column:device_info" json:"device_info"`
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}

func GetUserList() []*UserBasic {
	data := make([]*UserBasic, 10)
	utils.DB.Find(&data)
	for _, v := range data {
		fmt.Println(v)
	}
	return data
}

func CreateUser(user UserBasic) *gorm.DB {
	return utils.DB.Create(user)
}
