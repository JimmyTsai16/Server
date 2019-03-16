package model

import "github.com/jinzhu/gorm"

type ChatRoom struct {
	gorm.Model
	RoomId string `gorm:"type:varchar(127);unique"`
	RoomName string
	UserId int
	Icon string
}