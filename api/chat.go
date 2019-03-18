package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jimmy/server/model"
	"github.com/jimmy/server/ws"
	"net/http"
)

type ChatDatabase interface {
	GetUserProfile(UserId string) *model.UserProfile
	GetUserApplicationByUserId(UserId string) *model.UserApplications
	GetRoomsByRoomIds(RoomIds[]string) []model.ChatRoom
	SaveChatContent(cc *model.ChatContent)
	GetChatContent(RoomId string) []model.ChatContent
}

func NewChatAPI(db ChatDatabase) ChatAPI {
	return ChatAPI{DB: db, rooms: make(map[string]*ws.ChatWS)}
}

type ChatAPI struct {
	DB ChatDatabase
	rooms map[string]*ws.ChatWS
}

func (c *ChatAPI) GetRooms(ctx *gin.Context) {
	uApp := c.DB.GetUserApplicationByUserId(ctx.GetString("UserId"))
	if uApp != nil {
		var roomIds []string
		json.Unmarshal([]byte(uApp.RoomIds), &roomIds)

		crs := c.DB.GetRoomsByRoomIds(roomIds)
		ctx.JSON(http.StatusOK, crs)
	}
}

func (c *ChatAPI) ChatWS(ctx *gin.Context) {
	roomId := ctx.GetString("RoomId")
	if c.rooms[roomId] == nil {
		c.rooms[roomId] = ws.NewRoom(ctx.GetString("RoomId"), c.DB)
		fmt.Println(ctx.GetString("RoomId"))
	}
	c.rooms[roomId].AddClient(ctx, ctx.GetString("UserId"))
}