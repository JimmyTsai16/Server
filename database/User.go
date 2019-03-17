package database

import (
	"github.com/jimmy/server/model"
	"strconv"
)

func (d *GormDatabase) GetUserProfile(UserId string) *model.UserProfile {
	up := new(model.UserProfile)
	d.DB.Where("id = ?", UserId).Find(&up)
	if strconv.Itoa(up.UserId) == UserId {
		return up
	}
	return nil
}