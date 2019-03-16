package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

type UserProfile struct {
	gorm.Model
	UserId int `gorm:"type:varchar(127);AUTO_INCREMENT;unique"`
	UserName string `gorm:"type:varchar(100);unique"`
	FirstName string
	LastName string
	Age int
	BirthDay time.Time
	Email string
}