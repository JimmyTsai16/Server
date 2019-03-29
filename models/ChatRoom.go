package models

import "github.com/jinzhu/gorm"

type ChatRoom struct {
	gorm.Model
	RoomId string `gorm:"type:varchar(127);unique"`
	RoomName string
	Users string `gorm:"type:longtext"`
	Icon string
}