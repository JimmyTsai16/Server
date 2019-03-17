package database

import "github.com/jimmy/server/model"

func (d *GormDatabase) GetUserAuthByBasic(userName, password string) *model.UserAuth {
	ua := new(model.UserAuth)
	d.DB.Where("user_name = ? && password = ?", userName, password).Find(&ua)
	if ua.UserName == userName && ua.Password == password {
		return ua
	}
	return nil
}

func (d *GormDatabase) GetUserAuthByToken(token string) *model.UserAuth {
	ua := new(model.UserAuth)
	d.DB.Where("token = ?", token).Find(&ua)
	if ua.Token == token {
		return ua
	}
	return nil
}
