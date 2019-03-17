package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jimmy/server/api"
	"github.com/jimmy/server/auth"
	"github.com/jimmy/server/database"
)

func Create(db *database.GormDatabase) (router *gin.Engine){

	userAuthorization := auth.Auth{DB: db}
	loginHandler := api.LoginAPI{DB: db}
	chatHandler := api.ChatAPI{DB: db}
	userHandler := api.UserAPI{DB: db}

	router = gin.Default()

	corsProxy := "/proxy"

	//router.Use(userAuthorization.RequireAuth())
	router.POST(corsProxy+"/login", loginHandler.Login)

	user := router.Group(corsProxy+"/user")
	{
		user.Use(userAuthorization.RequireAuth())
		user.GET("", userHandler.UserInit)
		user.GET("/:id", userHandler.GetUserProfile)
	}

	chat := router.Group(corsProxy+"/chat")
	{
		chat.Use(userAuthorization.RequireAuth())
		chat.GET("/getrooms",chatHandler.GetRooms)
		chat.GET("/chatws/:roomid/:token", chatHandler.ChatWS)
	}

	return router
}