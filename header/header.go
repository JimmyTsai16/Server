package header

import (
	"github.com/gin-gonic/gin"
)

func HeaderWrite(c *gin.Context){
	c.Header("content-type", "application/x-www-form-urlencoded")
	c.Header("accept", "application/json, text/plain, */*")

	/*** Allow CORS(Cross-Origin Resource Sharing ) Header ***/
	//c.Header().Add("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Origin", "*") //"http://127.0.0.1:3000")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Access-Control-Allow-Methods", "GET POST PUT DELETE OPTIONS")
	c.Header("Access-Control-Allow-Headers", "authorization access-control-allow-origin X-Requested-With Content-Type Accept")
	/*********************************************************/
	//c.WriteHeader(http.StatusOK)
}