package model

import "github.com/jinzhu/gorm"

type chatContent struct {
	gorm.Model
	UserId int
	RoomId string
	Content string
}