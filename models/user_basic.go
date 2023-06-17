package models

import (
	"fmt"
	"ginchat/utils"
	"gorm.io/gorm"
	"time"
)

type UserBasic struct {
	gorm.Model
	Name          string `gorm:"column:name" json:"name"`
	Password      string `gorm:"column:password" json:"password"`
	Phone         string `valid:"matches(^1[3-9]{1}\\d{9}$)" gorm:"column:phone" json:"phone"`
	Email         string `valid:"email" gorm:"column:email" json:"email"`
	Identity      string `gorm:"column:identity" json:"identity"`
	ClientIp      string `gorm:"column:client_ip" json:"client_ip"`
	ClientPort    string `gorm:"column:client_port" json:"client_port"`
	Salt          string `gorm:"column:salt" json:"salt"`
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

func FindUserByName(name string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name = ?", name).First(&user)
	return user
}

func FindUserByNameAndPwd(name, password string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name = ? and password = ?", name, password).First(&user)

	//token加密
	str := fmt.Sprintf("%d", time.Now().Unix())
	temp := utils.Md5Encode(str)
	utils.DB.Model(&user).Where("id = ?", user.ID).Update("identity", temp)
	return user
}

func FindUserByPhone(phone string) *gorm.DB {
	user := UserBasic{}
	return utils.DB.Where("phone = ?", phone).First(user)
}

func FindUserByEmail(email string) *gorm.DB {
	user := UserBasic{}
	return utils.DB.Where("email = ?", email).First(user)
}

func CreateUser(user *UserBasic) *gorm.DB {
	return utils.DB.Create(user)
}

func DeleteUser(user *UserBasic) *gorm.DB {
	return utils.DB.Delete(user)
}

func UpdateUser(user *UserBasic) *gorm.DB {
	return utils.DB.Model(user).Updates(UserBasic{Name: user.Name, Password: user.Password, Phone: user.Phone, Email: user.Email})
}

// FindById 查找某个用户
func FindById(id uint) UserBasic {
	user := UserBasic{}
	utils.DB.Where("id = ?", id).First(&user)
	return user
}
