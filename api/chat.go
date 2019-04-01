package api

import (
	"encoding/json"
	"fmt"
	"github.com/JimmyTsai16/server/models"
	"github.com/JimmyTsai16/server/ws"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type ChatDatabase interface {
	GetUserProfile(UserId string) *models.UserProfile
	GetUserApplicationByUserId(UserId string) *models.UserApplications
	GetRoomsByRoomIds(RoomIds[]string) []models.ChatRoom
	SaveChatContent(cc *models.ChatContent)
	GetChatContent(RoomId string) []models.ChatContent
	CreateRoom(cr *models.ChatRoom)
}

func NewChatAPI(db ChatDatabase) ChatAPI {
	return ChatAPI{DB: db, rooms: make(map[string]*ws.ChatWS)}
}

type ChatAPI struct {
	DB ChatDatabase
	rooms map[string]*ws.ChatWS
}

func (c *ChatAPI) CreateRoom(ctx *gin.Context) {
	type reqInfo struct {
		RoomName string `json:"RoomName"`
		Users []string	`json:"Users"`
	}
	var r reqInfo
	if err := ctx.BindJSON(&r); err != nil {
		log.Println(err)
	}

	usersStr, _ := json.Marshal(r.Users)

	cr := &models.ChatRoom{
		RoomName: r.RoomName,
		Users: string(usersStr),
	}

	c.DB.CreateRoom(cr)

	fmt.Println(r.RoomName, r.Users)
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