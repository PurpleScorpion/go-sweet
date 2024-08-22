package sweetyml

import (
	"github.com/PurpleScorpion/go-sweet-keqing/keqing"
	"github.com/beego/beego/v2/server/web"
)

func beegoInit() {
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

	web.BConfig.Listen.HTTPPort = port
	web.BConfig.AppName = name
	web.BConfig.CopyRequestBody = true
	if keqing.IsNotEmpty(path) && keqing.IsNotEmpty(mappingUrl) {
		web.BConfig.WebConfig.StaticDir[mappingUrl] = path
	}
}
