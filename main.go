package main

import (
	bootstrap "bootstrap/src/main"
	"bootstrap/src/main/system"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// 初始化yml之类的文件
	system.Init(router)

	// 初始化服务
	bootstrap.InitApp(router)
}
