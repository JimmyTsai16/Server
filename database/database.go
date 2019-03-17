package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

func New(dialect, dbUsername, dbPassword, dbAddress, dbName, dbPort string) *GormDatabase {

	//db, err := gorm.Open("mysql","root:root@tcp(192.168.100.21)/chatServer?charset=utf8&parseTime=True&loc=Local")
	connStr := fmt.Sprint(dbUsername, ":", dbPassword, "@tcp(", dbAddress, ":", dbPort, ")/", dbName, "?charset=utf8&parseTime=True&loc=Local")
	db, err := gorm.Open(dialect, connStr)
	if err!= nil {
		log.Println(err)
		panic("failed connect database!")
	}
	return &GormDatabase{DB: db}
}

type GormDatabase struct {
	DB *gorm.DB
}

func (d *GormDatabase) Close() {
	err := d.DB.Close()
	if err != nil {
		log.Println("GormDatabase Close Failed: ", err)
	}
}