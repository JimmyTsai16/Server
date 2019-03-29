package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	)

type UserAuth struct {
	gorm.Model
	UserName string `gorm:"type:varchar(100);unique"`
	Password string
	Token string
}