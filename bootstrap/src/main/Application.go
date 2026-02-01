package bootstrap

import (
	"bootstrap/src/main/filter"
	"fmt"
	"io"
	"runtime/debug"
	"service-common/src/main/service"
	demoRouter "service-demo/src/main/router"
	"shared/logger"
	"shared/utils"

	"github.com/PurpleScorpion/go-sweet-keqing/keqing"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitApp(router *gin.Engine) {
	utils.Init()
	// 优先设置全局panic处理
	globalPanicRecover(router)
	// 然后注册全局Filter
	filterInit(router)
	// 最后在注册router
	routerInit(router)
	serviceInit()
	queueInit()
	runApp(router)
}

func routerInit(router *gin.Engine) {
	demoRouter.Init(router)
}

func serviceInit() {
	service.Init()
}

func queueInit() {
}

func filterInit(router *gin.Engine) {
	filter.Init()
	router.Use(filter.AuthFilter())
}

func middlewareInit(router *gin.Engine) {
	router.Use(corsMiddleware())
}

func globalPanicRecover(router *gin.Engine) {
	router.Use(gin.CustomRecoveryWithWriter(io.Discard, func(c *gin.Context, err any) {
		uuid := keqing.UUID()
		logger.Error("ErrorId:[{}]panic recovered start============================================================", uuid)
		logger.Error("ErrorId:[{}]error info: {}", uuid, err)
		logger.Error("ErrorId:[{}]stack info: {}", uuid, string(debug.Stack()))
		logger.Error("ErrorId:[{}]panic recovered end==============================================================", uuid)
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "server error",
		})
	}))
}

func corsMiddleware() gin.HandlerFunc {

	headers := make([]string, 0)
	headers = append(headers, "*", "Origin", "X-Requested-With", "Content-Type", "Accept", "User-Id", "Token")
	headers = append(headers, "X-Tenant-Id", "Authorization", "client-id", "Location", "session-state")
	headers = append(headers, "X-Current-User-Id", "X-Current-User-Name", "X-Current-User-Email", "X-Current-User-First-Name", "X-Current-User-Last-Name")
	headers = append(headers, "X-Current-User-Roles", "X-Current-User-Cluster", "X-Current-System-Type")

	return cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     headers,
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	})
}

func runApp(router *gin.Engine) {
	middlewareInit(router)

	port := keqing.ValueInt("${server.port}")
	err := router.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		logger.Error("Failed to start server: {}", err)
		return
	}
}
