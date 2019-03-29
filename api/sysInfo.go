package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jimmy/server/models"
	"github.com/jimmy/server/ws"
	"log"
	"net/http"
	"strconv"
	"time"
)

type sysInfoDatabase interface {
	GetCpuInfoBetween(start time.Time, end time.Time) []models.CPUInfo
	GetSysInfoBetween(start time.Time, end time.Time, t interface{}) interface{}
}

type SysInfoAPI struct {
	DB sysInfoDatabase
}

func (d *SysInfoAPI) GetCpuInfoBetween(ctx *gin.Context) {
	//start, _ := time.Parse(time.RFC3339, "2019-03-22T02:00:00+08:00")
	//end, _ := time.Parse(time.RFC3339, "2019-03-22T03:00:00+08:00")
	//fmt.Println(ctx.Param("startDate"))
	//fmt.Println(ctx.Param("endDate"))

	start, end := parseTime(ctx)

	x := d.DB.GetCpuInfoBetween(start, end)
	ctx.JSON(http.StatusOK, &x)
}

//func (d *SysInfoAPI) GetTemp(ctx *gin.Context) {
//	start, end := parseTime(ctx)
//
//	t := d.DB.GetSysInfoBetween(start, end, &models.CPUTemp{})
//	cpuTemps, err := t.([]models.CPUTemp)
//	if !err {
//		log.Fatal("Type assertion if fail: ", err)
//	}
//
//}

func(d *SysInfoAPI) GetSysInfo(ctx *gin.Context) {
	var info string
	start, end := parseTime(ctx)
	if d, ok := ctx.Params.Get("info"); ok {
		info = d
	}

	switch info {
	case "cpuTemp":
		t := d.DB.GetSysInfoBetween(start, end, &models.CPUTemp{})
		cpuTemps, err := t.([]models.CPUTemp)
		if !err {
			log.Fatal("Type assertion if fail: ", err)
		}
		ctx.JSON(http.StatusOK, &cpuTemps)

	case "cpuInfo":
		t := d.DB.GetSysInfoBetween(start, end, &models.CPUInfo{})
		cpuInfos, err := t.([]models.CPUInfo)
		if !err {
			log.Fatal("Type assertion if fail: ", err)
		}
		ctx.JSON(http.StatusOK, &cpuInfos)

	}

}

func (d *SysInfoAPI) GetSysInfoWS(ctx *gin.Context) {
	conn := ws.NewWS(ctx)
	//_, p, _ := conn.ReadMessage()
	type Info struct {
		Info string `json:"info"`
	}
	var i Info
	type x struct {
		Test string `json:"test"`
	}
	var xx x

	conn.ReadJSON(&i)
	fmt.Println(i.Info)

	conn.ReadJSON(&xx)
	fmt.Println(xx)
	conn.WriteJSON(&Info{
		Info: "cpuInfo",
	})
	conn.Close()
}

func parseTime(ctx *gin.Context) (start time.Time, end time.Time) {
	// var start, end time.Time

	if s, ok := ctx.Params.Get("startDate");ok {
		t, _ := strconv.Atoi(s)
		t = t / 1000
		start = time.Unix(int64(t), 0)
		// fmt.Println(start)
	}else{
		log.Println("Get startDate Failed.")
	}

	if e, ok := ctx.Params.Get("endDate");ok {
		t, _ := strconv.Atoi(e)
		t = t / 1000
		end = time.Unix(int64(t), 0)
		// fmt.Println(end)
	}else{
		log.Println("Get endDate Failed.")
	}

	return
}