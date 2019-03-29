package models

import "github.com/jinzhu/gorm"

type ChatContent struct {
	gorm.Model
	UserId int
	RoomId string
	Content string
}