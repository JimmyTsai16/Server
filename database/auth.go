package database

import (
	"github.com/jimmy/server/model"
	"strings"
)

func (d *GormDatabase) GetUserAuthByBasic(userName, password string) *model.UserAuth {
	ua := new(model.UserAuth)
	d.DB.Where("user_name = ? && password = ?", userName, password).Find(&ua)
	/***** 不區分大小寫比較字串 *****/
	if strings.EqualFold(ua.UserName, userName) && strings.EqualFold(ua.Password, password) {
		return ua
	}
	return nil
}

func (d *GormDatabase) GetUserAuthByToken(token string) *model.UserAuth {
	ua := new(model.UserAuth)
	d.DB.Where("token = ?", token).Find(&ua)
	/***** 不區分大小寫比較字串 *****/
	if strings.EqualFold(ua.Token, token) {
		return ua
	}
	return nil
}
