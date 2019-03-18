package api

import (
	"crypto/sha256"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jimmy/server/model"
	"log"
	"net/http"
)


type UserDatabase interface {
	GetUserProfile(UserId string) *model.UserProfile
	CreateUser(ua *model.UserAuth, up *model.UserProfile)
}

type UserAPI struct {
	DB UserDatabase
}

func (a *UserAPI) CreateUser(ctx *gin.Context) {
	type reqInfo struct {
		UserName string
		Password string
		FirstName string
		LastName string
		Email string
		Age int
	}
	var r reqInfo

	if err := ctx.BindJSON(&r); err != nil {
		log.Println(err)
	}
	fmt.Println(r)

	if r.UserName != "" && r.Password != "" {

		r.Password = fmt.Sprintf("%x", sha256.Sum256([]byte(r.Password)))

		ua := &model.UserAuth{
			UserName: r.UserName,
			Password: r.Password,
		}

		up := &model.UserProfile{
			UserName:  r.UserName,
			FirstName: r.FirstName,
			LastName:  r.LastName,
			Email:     r.Email,
			Age:       r.Age,
		}
		a.DB.CreateUser(ua, up)
		//a.DB.CreateUser(model.UserAuth{UserName:"us", Password:"8C6976E5B5410415BDE908BD4DEE15DFB167A9C873FC4BB8A81F6F2AB448A918"}, model.UserProfile{UserName:"us"})
		fmt.Println("CreateUser: ", ua.UserName)
	}
}

func (a *UserAPI) GetUserProfile(ctx *gin.Context) {
	if ctx.GetString("UserId") == ctx.Param("id") {
		up := a.DB.GetUserProfile(ctx.Param("id"))
		if up.ID != 0 {
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