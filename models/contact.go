package models

import (
	"fmt"
	"ginchat/utils"
	"gorm.io/gorm"
)

// Contact 人员关系
type Contact struct {
	gorm.Model
	OwnerId  uint //谁的关系
	TargetId uint //对应的谁
	Type     int  //对应的类型 1好友 2群 3xx
	Desc     string
}

func (table *Contact) TableName() string {
	return "contract"
}

func SearchFriend(userId uint) []UserBasic {
	contacts := make([]Contact, 0)
	objIds := make([]uint64, 0)
	utils.DB.Where("owner_id = ? and type = 1", userId).Find(&contacts)

	for _, v := range contacts {
		fmt.Println(v)
		objIds = append(objIds, uint64(v.TargetId))
	}

	users := make([]UserBasic, 0)
	utils.DB.Where("id in ?", objIds).Find(&users)

	return users
}

// AddFriend 添加好友
func AddFriend(userId uint, targetId uint) (int, string) {
	user := UserBasic{}

	if targetId != 0 {
		user = FindById(targetId)
		if user.Salt != "" {
			if userId == user.ID {
				return -1, "不能加自己为好友"
			}

			contact0 := Contact{}
			utils.DB.Where("owner_id = ? and target_id = ? and type = 1", userId, targetId).Find(&contact0)
			if contact0.ID != 0 {
				return -1, "该用户已经是您的好友了"
			}

			utils.DB.Where("")

			tx := utils.DB.Begin()
			//事务一旦开始，不论什么异常都会Rollback
			defer func() {
				if err := recover(); err != nil {
					tx.Rollback()
				}
			}()
			contact := Contact{}
			contact.OwnerId = userId
			contact.TargetId = targetId
			contact.Type = 1
			if err := utils.DB.Create(&contact).Error; err != nil {
				tx.Rollback()
				return -1, "添加好友失败"
			}
			contact1 := Contact{}
			contact1.OwnerId = targetId
			contact1.TargetId = userId
			contact1.Type = 1
			if err := utils.DB.Create(&contact1).Error; err != nil {
				tx.Rollback()
				return -1, "添加好友失败"
			}
			tx.Commit()
			return 0, "添加好友成功"
		}
		return -1, "没有找到此用户"
	}
	return -1, "好友id不能为空"
}

func SearchUserByGroupId(comId uint) []uint {
	contacts := make([]Contact, 0)
	objIds := make([]uint, 0)

	utils.DB.Where("target_id = ? and type = 2", comId).Find(&contacts)
	for _, v := range contacts {
		objIds = append(objIds, v.OwnerId)
	}
	return objIds
}
