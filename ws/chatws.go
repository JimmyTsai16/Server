package ws

import (
	"crypto/sha256"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jimmy/server/model"
	"log"
	"sync"
	"time"
)

type ChatContentDatabase interface {
	GetUserProfile(UserId string) *model.UserProfile
	SaveChatContent(cc *model.ChatContent)
	GetChatContent(RoomId string) []model.ChatContent
}

/*
Init the new room connection.
Launch the write routine to send message to each connection.
 */
func NewRoom(Id string, db ChatContentDatabase) *ChatWS {
	c := &ChatWS{
		Clients: make(map[string]*Client),
		msg:    make(chan *model.ChatContent, 5),
		RoomId: Id,
		DB:     db,
	}
	go c.writeRoutine()
	return c
}

type ChatWS struct {
	Clients map[string]*Client
	lock    sync.RWMutex
	msg     chan *model.ChatContent
	DB      ChatContentDatabase
	RoomId  string
}

func (c *ChatWS) Close(connHash string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.Clients, connHash)
}

func (c *ChatWS) AddClient(ctx *gin.Context, UserId string) {
	ws := NewWS(ctx)
	hash := fmt.Sprintf("%X", sha256.Sum256([]byte(time.Now().String())))
	c.Clients[hash] = &Client{Conn: ws, UP: c.DB.GetUserProfile(UserId)}
	// go c.writeRoutine()
	c.writeChatContent(ws)
	go c.readRoutine(ws, hash)
}

func (c *ChatWS) writeChatContent(ws *websocket.Conn) {
	cc := c.DB.GetChatContent(c.RoomId)
	ws.WriteJSON(cc)
}

func (c *ChatWS) readRoutine(ws *websocket.Conn, connHash string) {
	var cc *model.ChatContent
	for {
		err := ws.ReadJSON(&cc)
		if err != nil {
			log.Println("WebSocket read message error: ", err)
			ws.Close()
			c.Close(connHash)
			return
		}
		cc.RoomId = c.RoomId
		cc.ID = 0
		c.DB.SaveChatContent(cc)
		c.msg <- cc
	}
}

func (c *ChatWS) writeRoutine() {
	for {
		select {
		case buf := <-c.msg:
			for _, b := range c.Clients {
				if err := b.Conn.WriteJSON(buf); err != nil {
					log.Println("WebSocket write message error: ", err)
				}
			}
		}
	}
}