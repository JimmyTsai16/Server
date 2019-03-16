package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	)

type UserAuth struct {
	gorm.Model
	UserId int `gorm:"type:varchar(127);AUTO_INCREMENT;unique"`
	UserName string `gorm:"type:varchar(100);unique"`
	Password string
	Token string
}