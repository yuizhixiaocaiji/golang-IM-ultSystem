package main

import (
	"ginchat/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(mysql.Open("root:343527@tcp(101.42.54.168:3306)/ginchat?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 迁移 schema
	//db.AutoMigrate(&models.UserBasic{})
	db.AutoMigrate(&models.Community{})
	// Create
	//user := &models.UserBasic{}
	//user.Name = "一只小菜狗"
	//db.Create(user)

	// Read
	//fmt.Println(db.First(user, 1)) // 根据整型主键查找
	//db.First(user, "code = ?", "D42") // 查找 code 字段值为 D42 的记录

	// Update - 将 product 的 price 更新为 200
	//db.Model(user).Update("Password", "343527")
	// Update - 更新多个字段
	//db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // 仅更新非零值字段
	//db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// Delete - 删除 product
	//db.Delete(&product, 1)
}
