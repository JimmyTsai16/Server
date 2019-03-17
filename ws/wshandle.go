package ws

import (
	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

type RoomChan struct {
	RoomId string
	Content chan string
	DB *gorm.DB
	WS *websocket.Conn
}

func (a *RoomChan) Handler() {

	ch := make(chan []byte, 10)
	go func() {
		for {
			select {
			case d := <- ch:
				a.WS.WriteMessage(websocket.TextMessage, d)
			}
		}
	}()

	for {
		_, buf, err := a.WS.ReadMessage()
		if err != nil {
			log.Printf("WS connection error: ", err)
			a.WS.Close()
			return
		}
		ch <- buf
	}
}