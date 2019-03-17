package main

import (
	"github.com/jimmy/server/database"
	"github.com/jimmy/server/router"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

func main() {
	db := database.New("mysql", "root", "root", "192.168.100.21", "chatServer", "3306")
	defer db.Close()

	r := router.Create(db)
	err := r.Run()
	if err != nil {
		log.Println("RouterError: ", err)
	}
}
