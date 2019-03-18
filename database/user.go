package database

import (
	"crypto/sha256"
	"fmt"
	"github.com/jimmy/server/model"
	"strconv"
	"time"
)

func (d *GormDatabase) GetUserProfile(UserId string) *model.UserProfile {
	up := new(model.UserProfile)
	d.DB.Where("id = ?", UserId).Find(&up)
	if strconv.Itoa(int(up.ID)) == UserId {
		return up
	}
	return nil
}

func (d *GormDatabase) CreateUser(ua *model.UserAuth, up *model.UserProfile) {
	t := ua.UserName + time.Now().String()
	ua.Token = fmt.Sprintf("%X", sha256.Sum256([]byte(t)))
	d.DB.Create(&ua)
	d.DB.Create(&up)
	d.DB.Create(&model.UserApplications{UserName: ua.UserName})
}