package system

import (
	"github.com/PurpleScorpion/go-sweet-keqing/keqing"
	"github.com/gin-gonic/gin"
)

func webServerInit(router *gin.Engine) {
	port := keqing.ValueInt("${server.port}")
	if port == 0 {
		port = 8080
	}

	if port < 0 || port > 65535 {
		panic("web server port is invalid")
	}

	name := keqing.ValueString("${server.name}")
	if keqing.IsEmpty(name) {
		name = "go-sweet"
	}

	mappingUrl := keqing.ValueString("${sweet.img.mappingUrl}")
	path := keqing.ValueString("${sweet.img.path}")

	router.Use(func(c *gin.Context) {
		c.Set("app_name", name)
		c.Next()
	})

	router.Static(mappingUrl, path)
}
