package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jimmy/server/api"
	"github.com/jimmy/server/auth"
	"github.com/jimmy/server/database"
	"net/http"
	_ "net/http/pprof"
)

func Create(db *database.GormDatabase, sysInfoDb *database.GormDatabase) (router *gin.Engine){

	// System Information Handle API.
	sysInfoHandle := api.NewSysInfoApi(sysInfoDb)

	// Member Information Handle API.
	userAuthorization := auth.Auth{DB: db}
	loginHandler := api.LoginAPI{DB: db}
	chatHandler := api.NewChatAPI(db)
	userHandler := api.UserAPI{DB: db}

	router = gin.Default()

	corsProxy := "/proxy"

	//router.Use(userAuthorization.MiddleWare())
	router.POST(corsProxy+"/login", loginHandler.Login)

	// Handle static files
	router.Static(corsProxy + "/static", "./static")

	sysInfo := router.Group(corsProxy + "/sysinfo")
	{
		sysInfo.GET("/multiple/:info/:startDate/:endDate", sysInfoHandle.GetSysInfo)
		sysInfo.GET("/host", sysInfoHandle.GetHostInfo)
		sysInfo.GET("/ws", sysInfoHandle.GetSysInfoWS)
		// sysInfo.GET("/cpuinfo/:startDate/:endDate", sysInfoHandle.GetCpuInfoBetween)
		// sysInfo.GET("/cputemp/:startDate/:endDate", sysInfoHandle.GetTemp)
	}

	user := router.Group(corsProxy+"/user")
	{
		user.Use(userAuthorization.MiddleWare())
		user.GET("", userHandler.UserInfo)
		user.POST("", userHandler.CreateUser)
		user.GET("/:id", userHandler.GetUserProfile)
	}

	chat := router.Group(corsProxy+"/chat")
	{
		chat.Use(userAuthorization.MiddleWare())
		chat.POST("/createroom", chatHandler.CreateRoom)
		chat.GET("/getrooms",chatHandler.GetRooms)
		chat.GET("/chatws/:roomid/:token", chatHandler.ChatWS)
	}

	//rt := make(chan string, 10)

	//go func() {
	//	gr := ""
	//	for {
	//		ngr := fmt.Sprintln("Number of go routine: ", runtime.NumGoroutine())
	//		if gr != ngr {
	//			gr = ngr
	//			rt <- gr
	//		}
	//
	//		time.Sleep(time.Millisecond * 600)
	//	}
	//}()

	go func() {
		http.ListenAndServe("0.0.0.0:8081", nil)
	}()

	//go func() {
	//	for {
	//		var gcs debug.GCStats
	//		fmt.Println("Number of go routine: ", runtime.NumGoroutine())
	//		fmt.Println("Garbage collection stat: ")
	//		debug.ReadGCStats(&gcs)
	//		fmt.Println(gcs)
	//
	//		fmt.Println("MemStat: ")
	//		var mems runtime.MemStats
	//		runtime.ReadMemStats(&mems)
	//		fmt.Printf("%+v\n", mems)
	//
	//
	//		time.Sleep(time.Millisecond * 2000)
	//	}
	//}()

	return router
}