package models

import "gorm.io/gorm"

// Message 消息
type Message struct {
	gorm.Model
	FromId   uint   //发送者
	TargetId uint   //接受者
	Type     string //消息类型 1.群聊 2.私聊 3.广播
	Media    int    //消息类型 1.文字 2.图片 3.音频
	Content  string //消息内容
	Pic      string
	Url      string
	Desc     string
	Amount   int //其他数字统计
}

func (table *Message) TableName() string {
	return "message"
}
