package models

import (
	"fmt"
	"ginchat/utils"
	"gorm.io/gorm"
)

type Community struct {
	gorm.Model
	Name    string
	OwnerId uint
	Img     string
	Desc    string
}

func CreateCommunity(community Community) (int, string) {
	if len(community.Name) == 0 {
		return -1, "群名称不能为空"
	}
	if community.OwnerId == 0 {
		return -1, "请先登录"
	}
	if err := utils.DB.Create(&community).Error; err != nil {
		fmt.Println(err)
		return -1, "建群失败"
	}
	contact := Contact{}
	contact.OwnerId = community.OwnerId
	contact.TargetId = community.ID
	contact.Type = 2 //群关系
	if err := utils.DB.Create(&contact).Error; err != nil {
		return -1, "添加群关系失败"
	}
	return 0, "建群成功"
}

func LoadCommunity(ownerId uint) ([]*Community, int, string) {
	contacts := make([]Contact, 0)
	objIds := make([]uint, 0)
	utils.DB.Where("owner_id = ? and type = 2", ownerId).Find(&contacts)
	for _, v := range contacts {
		objIds = append(objIds, v.TargetId)
	}

	data := make([]*Community, 10)
	utils.DB.Where("id in ?", objIds).Find(&data)
	for _, v := range data {
		fmt.Println(v)
	}
	return data, 0, "查询成功"
}

func JoinGroup(userId uint, comInfo string) (int, string) {
	contact := Contact{}
	contact.OwnerId = userId
	contact.Type = 2
	community := Community{}

	utils.DB.Where("id = ? or name = ?", comInfo, comInfo).Find(&community)
	if community.Name == "" {
		return -1, "没有找到群"
	}
	utils.DB.Where("owner_id = ? and target_id = ? and type = 2", userId, comInfo).Find(&contact)
	if !contact.CreatedAt.IsZero() {
		return -1, "已加入过此群"
	} else {
		contact.TargetId = community.ID
		utils.DB.Create(&contact)
		return 0, "加群成功"
	}
}
