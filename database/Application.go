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
