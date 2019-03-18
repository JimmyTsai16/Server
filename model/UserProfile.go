package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

type UserProfile struct {
	gorm.Model
	UserName string `gorm:"type:varchar(100);unique"`
	FirstName string
	LastName string
	Friends string `gorm:"type:longtext"`
	FriendChecking string `gorm:"type:longtext"`
	Age int
	BirthDay time.Time
	Email string
}