/**
* @Author: Lanhai Bai
* @Date: 2021/8/25 10:21
* @Description:
 */
package router

import (
	"github.com/gin-gonic/gin"
	"proxy-pool/handlers"
	"proxy-pool/utils"
)
var routerEngine *gin.Engine
func init() {
	utils.Logger().Infof("router init...")
	routerEngine = gin.New()
	v1 := routerEngine.Group("v1")
	{
		v1.GET("/ips", handlers.IndexHandler)
	}
}

func Run(addr string) error {
	return routerEngine.Run(addr)
}