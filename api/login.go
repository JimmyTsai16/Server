package api

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/JimmyTsai16/server/header"
	"github.com/JimmyTsai16/server/jwt"
	"github.com/JimmyTsai16/server/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

type LoginDatabase interface {
	GetUserAuthByBasic(userName, password string) *models.UserAuth
}

type LoginAPI struct {
	DB LoginDatabase
}

func (u *LoginAPI) Login(c *gin.Context) {
	//c.JSON(http.StatusOK, map[string]string{ "status": "Logged in." })
	type reqAuth struct {
		Username string
		Password string
	}
	header.HeaderWrite(c)

	var req reqAuth
	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		log.Println("Auth json decode failed: ", err)
	}
	fmt.Println(req.Username, req.Password, "login.")
	req.Password = fmt.Sprintf("%X", sha256.Sum256([]byte(req.Password)))
	ua := u.DB.GetUserAuthByBasic(req.Username, req.Password)
	fmt.Println(ua)
	/****** if query result just have one match *********/
	if ua != nil {
		j := jwt.UserJwt{
			Token: ua.Token,
			Exp: time.Now(),
		}
		if jwtString, ok := j.JwtSignedString(); ok {
			c.JSON(http.StatusOK, map[string]string{"token": jwtString})
		}
	}else{
		/****** Fail Auth, return null and ok status *********/
		log.Println("Auth failed: ", req)
		c.JSON(http.StatusUnauthorized, map[string]string{"error": "Wrong Username or Password."})
	}
}