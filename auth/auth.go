package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jimmy/server/header"
	"github.com/jimmy/server/jwt"
	"github.com/jimmy/server/model"
	"net/http"
	"strconv"
)

type AuthDatabase interface {
	GetUserAuthByToken(token string) *model.UserAuth
	GetUserProfile(UserId string) *model.UserProfile
}

type Auth struct {
	TokenHeader string
	DB AuthDatabase
}

const (
	JwtHeader = "X-AccessToken"
)

func (a *Auth) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		header.HeaderWrite(c)
		JwtString := c.Request.Header.Get(JwtHeader)
		/***** Handle WebSocket *****/
		if c.Request.Header.Get("Upgrade") == "websocket" {
			type RoomIdToken struct {
				RoomId string
				AccessToken string
			}
			var rt RoomIdToken
			fmt.Println("Upgrade WebSocket")
			rt.RoomId = c.Param("roomid")
			rt.AccessToken = c.Param("token")

			JwtString = rt.AccessToken
			c.Set("RoomId", rt.RoomId)
			//c.AbortWithStatus(http.StatusUnauthorized)
			//return
		}

		if JwtString != "" {
			if status, up := a.JwtAuth(JwtString); status {
				fmt.Println("JWTAuth Pass.")
				//c.JSON(http.StatusOK, userId)
				c.Set("UserId", fmt.Sprintf("%d", up.UserId))
			}else {
				c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"Error": "Please login."})
			}
		}else {
			//c.Redirect(http.StatusMovedPermanently, "/login")
			c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"Error": "Please login."})
		}
		c.Next()
	}
}

//func (a *Auth) BasicAuth(UserName, Password string) (status bool, jwtString string, userId int) {
//	var ua model.UserAuth
//	a.DB.Where("user_name = ?, password = ?", UserName, Password).Find(&ua)
//	if ua.UserName == UserName {
//		return true, "success", ua.UserId
//	}
//	return false, "", 0
//}

func (a *Auth) JwtAuth(jwtString string) (status bool, userProfile *model.UserProfile) {
	j := jwt.UserJwt{}
	j.JwtParse(jwtString)

	ua := a.DB.GetUserAuthByToken(j.Token)
	if ua.UserId != 0 {
		userProfile = a.DB.GetUserProfile(strconv.Itoa(ua.UserId))
		return true, userProfile
	}else{
		return false, nil
	}
}