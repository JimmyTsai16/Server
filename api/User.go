package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jimmy/server/model"
	"net/http"
)


type UserDatabase interface {
	GetUserProfile(UserId string) *model.UserProfile
}

type UserAPI struct {
	DB UserDatabase
}

func (a *UserAPI) GetUserProfile(ctx *gin.Context) {
	if ctx.GetString("UserId") == ctx.Param("id") {
		up := a.DB.GetUserProfile(ctx.Param("id"))
		if up.UserId != 0 {
			ctx.JSON(http.StatusOK, up)
		} else {
			ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "UnAuthorization"})
		}
	}else{
		ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "User and Token are Not Match."})
	}
}

func (a *UserAPI) UserInit(ctx *gin.Context) {
	up := a.DB.GetUserProfile(ctx.GetString("UserId"))
	ctx.JSON(http.StatusOK, up)
}