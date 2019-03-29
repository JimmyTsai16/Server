package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jimmy/server/models"
	"log"
)

type Person struct {
	UserName     string
	Password  string
}

func main() {
	//route := gin.Default()
	//route.POST("/testing", startPage)
	//route.Run(":8085")
	models.InitDatabase()
}

func startPage(c *gin.Context) {
	var person Person
	// If `GET`, only `Form` binding engine (`query`) used.
	// If `POST`, first checks the `content-type` for `JSON` or `XML`, then uses `Form` (`form-data`).
	// See more at https://github.com/gin-gonic/gin/blob/master/binding/binding.go#L48
	if c.ShouldBind(&person) == nil {
		log.Println(person.UserName)
		log.Println(person.Password)
	}

	c.String(200, "Success")
}