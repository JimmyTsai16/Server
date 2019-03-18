package database

import (
	"github.com/jimmy/server/model"
	"strconv"
)

func (d *GormDatabase) GetUserApplicationByUserId(UserId string) *model.UserApplications {
	uApp := new(model.UserApplications)
	d.DB.Where("user_id = ?", UserId).Find(&uApp)
	if strconv.Itoa(uApp.UserId) == UserId {
		return uApp
	}
	return nil
}

func (d *GormDatabase) GetRoomsByRoomIds(RoomIds[]string) []model.ChatRoom {
	var cr []model.ChatRoom
	d.DB.Where("room_id IN (?)", RoomIds).Find(&cr)
	if len(cr) > 0 {
		return cr
	}
	return nil
}

func (d *GormDatabase) SaveChatContent(cc *model.ChatContent) {
	d.DB.Create(&cc)
}

func (d *GormDatabase) GetChatContent(RoomId string) []model.ChatContent {
	var cc []model.ChatContent
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