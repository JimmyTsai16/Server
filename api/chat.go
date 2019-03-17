package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jimmy/server/model"
	"log"
	"net/http"
)

type ChatDatabase interface {
	GetUserProfile(UserId string) *model.UserProfile
	GetUserApplicationByUserId(UserId string) *model.UserApplications
	GetRoomsByRoomIds(RoomIds[]string) []model.ChatRoom
}

type ChatAPI struct {
	DB ChatDatabase
}

func (c *ChatAPI) GetRooms(ctx *gin.Context) {
	uApp := c.DB.GetUserApplicationByUserId(ctx.GetString("UserId"))
	if uApp != nil {
		var roomIds []string
		json.Unmarshal([]byte(uApp.RoomIds), &roomIds)

		crs := c.DB.GetRoomsByRoomIds(roomIds)
		fmt.Println(crs)
		ctx.JSON(http.StatusOK, crs)
	}
}

func (c *ChatAPI) ChatWS(ctx *gin.Context) {

	var upGrader = websocket.Upgrader{
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
		EnableCompression: false,
		/*** CORS ***/
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	/******* Upgrade the connection to WebSocket **********/
	ws, err := upGrader.Upgrade(ctx.Writer, ctx.Request , nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer ws.Close()

}