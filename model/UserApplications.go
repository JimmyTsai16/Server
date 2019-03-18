package model

import "github.com/jinzhu/gorm"

type UserApplications struct {
	gorm.Model
	UserName string `gorm:"type:varchar(100);unique"`
	RoomIds string `gorm:"type:longtext"` // store RoomId with type []string
}