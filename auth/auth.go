package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jimmy/server/header"
	"github.com/jimmy/server/jwt"
	"github.com/jimmy/server/models"
	"net/http"
	"strconv"
)

type AuthDatabase interface {
	GetUserAuthByToken(token string) *models.UserAuth
	GetUserProfile(UserId string) *models.UserProfile
}

type Auth struct {
	TokenHeader string
	DB AuthDatabase
}

const (
	JwtHeader = "X-AccessToken"
)

func (a *Auth) MiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		// First MiddleWare
		a.RequireAuth(c)
		// Add more middleware. ...
	}
}

func (a *Auth) RequireAuth(ctx *gin.Context) {
	header.HeaderWrite(ctx)
	JwtString := ctx.Request.Header.Get(JwtHeader)

	/***** Handle WebSocket *****/
	if ctx.Request.Header.Get("Upgrade") == "websocket" {
		type RoomIdToken struct {
			RoomId string
			AccessToken string
		}
		var rt RoomIdToken
		fmt.Println("Upgrade WebSocket")
		rt.RoomId = ctx.Param("roomid")
		rt.AccessToken = ctx.Param("token")

		JwtString = rt.AccessToken
		ctx.Set("RoomId", rt.RoomId)
		//c.AbortWithStatus(http.StatusUnauthorized)
		//return
	}

	if JwtString != "" {
		if status, up := a.JwtAuth(JwtString); status {
			fmt.Println("JWTAuth Pass.")
			//c.JSON(http.StatusOK, userId)
			ctx.Set("UserId", fmt.Sprintf("%d", up.ID))
			ctx.Next()
		}else {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"Error": "Please login."})
		}
	}else {
		//c.Redirect(http.StatusMovedPermanently, "/login")
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"Error": "Please login."})
	}
}

//func (a *Auth) BasicAuth(UserName, Password string) (status bool, jwtString string, userId int) {
//	var ua models.UserAuth
//	a.DB.Where("user_name = ?, password = ?", UserName, Password).Find(&ua)
//	if ua.UserName == UserName {
//		return true, "success", ua.UserId
//	}
//	return false, "", 0
//}

func (a *Auth) JwtAuth(jwtString string) (status bool, userProfile *models.UserProfile) {
	j := jwt.UserJwt{}
	j.JwtParse(jwtString)

	ua := a.DB.GetUserAuthByToken(j.Token)
	if ua.ID != 0 {
		userProfile = a.DB.GetUserProfile(strconv.Itoa(int(ua.ID)))
		return true, userProfile
	}else{
		return false, nil
	}
}