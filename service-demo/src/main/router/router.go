package venRouter

import (
	"service-demo/src/main/controller"

	"github.com/gin-gonic/gin"
)

func Init(router *gin.Engine) {
	{
		ven := router.Group("/demo")
		ven.GET("/test", controller.Demo1)
		ven.GET("/test2/:bas-code", controller.Demo2)
		ven.POST("/test3", controller.Demo3)
		ven.POST("/upload", controller.UploadFile)
		ven.POST("/download", controller.DownloadFile)

	}
}
