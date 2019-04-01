package database

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/JimmyTsai16/server/models"
	"strconv"
	"time"
)

func (d *GormDatabase) GetUserApplicationByUserId(UserId string) *models.UserApplications {
	uApp := new(models.UserApplications)
	d.DB.Where("ID = ?", UserId).Find(&uApp)
	if strconv.Itoa(int(uApp.ID)) == UserId {
		return uApp
	}
	return nil
}

func (d *GormDatabase) GetRoomsByRoomIds(RoomIds[]string) []models.ChatRoom {
	var cr []models.ChatRoom
	d.DB.Where("room_id IN (?)", RoomIds).Find(&cr)
	if len(cr) > 0 {
		return cr
	}
	return nil
}

func (d *GormDatabase) CreateRoom(cr *models.ChatRoom) {
	cr.RoomId = fmt.Sprintf("%x", sha256.Sum256([]byte(time.Now().String())))
	d.DB.Create(&cr)

	var users []string
	json.Unmarshal([]byte(cr.Users), &users)
	var ca []models.UserApplications
	d.DB.Where("user_name IN (?)", users).Find(&ca)
	fmt.Println(ca)
	for _, db := range ca {
		var roomsId []string
		json.Unmarshal([]byte(db.RoomIds), &roomsId)
		fmt.Println(roomsId)
		roomsId = append(roomsId, cr.RoomId)
		roomsIdStr, _ := json.Marshal(roomsId)
		db.RoomIds = string(roomsIdStr)
		d.DB.Save(&db)
	}
}

func (d *GormDatabase) SaveChatContent(cc *models.ChatContent) {
	d.DB.Create(&cc)
}

func (d *GormDatabase) GetChatContent(RoomId string) []models.ChatContent {
	var cc []models.ChatContent
	d.DB.Where("room_id = ?", RoomId).Order("id desc").Limit(15).Find(&cc)
	//q := "SELECT *	from chat_contents	WHERE id IN (SELECT * from (select id FROM chat_contents ORDER BY id desc	LIMIT 15) as t)	ORDER BY id"
	//sub1 := d.DB.Table("chat_contents").Select("id").Where("room_id = ?", RoomId).Order("id desc").Limit(15).QueryExpr()
	//d.DB.Raw("select id from chat_contents where id = ? order_by id desc limit 15", 1).Find(&cc)
	//sub2 := d.DB.Raw("SELECT * FROM = ?", sub1).QueryExpr()
	//d.DB.Where("id IN ?", sub1).Find(&cc)
	for i, j := 0, len(cc)-1; i < j; i, j = i+1, j-1 {
		cc[i], cc[j] = cc[j], cc[i]
	}
	if len(cc) > 0 {
		return cc
	}
	return nil
}