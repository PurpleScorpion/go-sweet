package appMain

import (
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/filter/cors"
	"sweet-common/utils"
	sweetyml "sweet-common/yaml"
	"sweet-src/main/golang/filter"
	"sweet-src/main/golang/routers"
	"sweet-src/main/golang/service"
)

func Main() {
	initMain()
	// 跨域配置
	corsMiddleware := cors.Allow(&cors.Options{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "X-Requested-With", "Content-Type", "Accept", "User-Id", "Token"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	})

	// 开启恐慌函数处理
	web.BConfig.RecoverPanic = true
	web.BConfig.RecoverFunc = filter.RecoverPanic

	// 注册中间件
	web.InsertFilter("*", web.BeforeRouter, corsMiddleware)
	// 权限拦截
	web.InsertFilter("/*", web.BeforeRouter, filter.FilterUser)

	web.Run()
}

func initMain() {
	sweetyml.ReadYml()
	conf := sweetyml.GetYmlConf()
	web.BConfig.Listen.HTTPPort = conf.Server.Port
	web.BConfig.AppName = conf.Server.Name
	web.BConfig.CopyRequestBody = true
	if conf.Sweet.Img.Active {
		web.BConfig.WebConfig.StaticDir[conf.Sweet.Img.MappingUrl] = conf.Sweet.Img.Path
	}
	routers.InitRouters()
	service.ServiceInit()
	utils.InitUtils()
	filter.InitFilter()
}
