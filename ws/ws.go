package ws

import (
	"github.com/JimmyTsai16/server/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"net/http"
)

func newUpgrade() *websocket.Upgrader {
	return &websocket.Upgrader{
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
		EnableCompression: false,
		/*** CORS ***/
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
}

func NewWS(ctx *gin.Context) *websocket.Conn {
	conn, err := newUpgrade().Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println("Connection upgrade to websocket error: ", err)
	}
	return conn
}

type Client struct {
	Conn *websocket.Conn
	UP *models.UserProfile
}
