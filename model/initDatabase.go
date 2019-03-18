package model

import (
	"crypto/sha256"
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

func InitDatabase()  {
	db, err := gorm.Open("mysql",
		"root:root@tcp(192.168.100.21)/chatServer?charset=utf8&parseTime=True&loc=Local")
	if err!= nil {
		log.Println(err)
		panic("failed connect database!")
	}
	defer db.Close()

	db.AutoMigrate(&UserAuth{})
	db.AutoMigrate(&UserProfile{})
	db.AutoMigrate(&UserApplications{})
	db.AutoMigrate(&ChatRoom{})
	db.AutoMigrate(&ChatContent{})
	fmt.Println("Database Init Finish.")

	db.Create(&UserAuth{
		UserName: "admin",
		Password: sha256String("admin"),
		Token: sha256String("admin" + time.Now().String()),
	})
	db.Create(&UserProfile{
		UserName: "admin",
		FirstName: "admin",
		LastName: "admin",
		Email: "admin@admin.com",
	})
	db.Create(&UserApplications{
		UserName: "admin",
	})

	db.Create(&UserAuth{
		UserName: "user1",
		Password: sha256String("user1"),
		Token: sha256String("admin" + time.Now().String()),
	})
	db.Create(&UserProfile{
		UserName: "user1",
		FirstName: "admin",
		LastName: "admin",
		Email: "admin@admin.com",
	})
	db.Create(&UserApplications{
		UserName: "user1",
	})
}

func sha256String(s string) string {
	return fmt.Sprintf("%X", sha256.Sum256([]byte(s)))
}