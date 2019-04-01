package database

import (
	"github.com/JimmyTsai16/server/models"
	"strings"
)

func (d *GormDatabase) GetUserAuthByBasic(userName, password string) *models.UserAuth {
	ua := new(models.UserAuth)
	d.DB.Where("user_name = ? && password = ?", userName, password).Find(&ua)

	/***** 不區分大小寫比較字串 *****/
	if strings.EqualFold(ua.UserName, userName) && strings.EqualFold(ua.Password, password) {
		return ua
	}
	return nil
}

func (d *GormDatabase) GetUserAuthByToken(token string) *models.UserAuth {
	ua := new(models.UserAuth)
	d.DB.Where("token = ?", token).Find(&ua)
	/***** 不區分大小寫比較字串 *****/
	if strings.EqualFold(ua.Token, token) {
		return ua
	}
	return nil
}
