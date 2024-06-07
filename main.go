package main

import (
	"fmt"
	"github.com/PurpleScorpion/go-sweet-json/jsonutil"
	"go-sweet/common/constants"
	"go-sweet/common/logger"
	"go-sweet/common/utils"
	sweetyml "go-sweet/common/yaml"
	_ "go-sweet/routers"
	"go-sweet/service"
	"net/http"
	"strings"
	"time"

	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"github.com/beego/beego/v2/server/web/filter/cors"
)

var fullExcludeUrls []string
var prefixExcludeUrls []string

func init() {
	sweetyml.ReadYml()
	conf := sweetyml.GetYmlConf()
	web.BConfig.Listen.HTTPPort = conf.Server.Port
	web.BConfig.AppName = conf.Server.Name
	web.BConfig.CopyRequestBody = true

	web.BConfig.WebConfig.StaticDir[conf.Sweet.Img.MappingUrl] = conf.Sweet.Img.Path

	logger.InitLogger()
	service.InitService()

	fullExcludeUrls = conf.Sweet.ExcUrl.Full
	prefixExcludeUrls = conf.Sweet.ExcUrl.Prefix
}

func main() {

	corsMiddleware := cors.Allow(&cors.Options{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "X-Requested-With", "Content-Type", "Accept", "User-Id", "Token"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	})

	// 开启恐慌函数处理
	web.BConfig.RecoverPanic = true
	web.BConfig.RecoverFunc = RecoverPanic

	// 注册中间件
	web.InsertFilter("*", web.BeforeRouter, corsMiddleware)

	web.InsertFilter("/*", web.BeforeRouter, FilterUser)

	web.Run()
}

func RecoverPanic(ctx *context.Context, cfg *web.Config) {
	if err := recover(); err != nil {
		// 在这里处理 panic
		ctx.Output.SetStatus(http.StatusInternalServerError)
		ctx.Output.JSON(getErrorToken(fmt.Sprintf("%v", err)), false, false)
	}
}

var FilterUser = func(ctx *context.Context) {
	if excludeUrl(ctx.Input.URL()) {
		return
	}

	headers := ctx.Request.Header
	values := headers["Token"]

	if len(values) == 0 {
		ctx.Output.JSON(getErrorToken("Token does not exist"), false, false)
		return
	}

	token := values[0]

	if utils.IsEmpty(token) {
		ctx.Output.JSON(getErrorToken("Token does not exist"), false, false)
		return
	}

	data, err := utils.Decrypt(token)
	if err != nil {
		ctx.Output.JSON(getErrorToken("Token error"), false, false)
		return
	}

	flag := hasExpire(data)

	if flag {
		ctx.Output.JSON(getErrorToken("Token has expired"), false, false)
		return
	}
}

/*
true: 过期
false: 未过期
*/
func hasExpire(data string) bool {
	js := jsonutil.NewJSONObject()
	js.ParseObject(data)
	if !js.HasKey("expirationTime") {
		return true
	}
	expirationTime := js.GetString("expirationTime")

	flag, err := compareDate(expirationTime, data)
	if err != nil {
		return true
	}
	return flag
}

func getErrorToken(msg string) interface{} {
	js := jsonutil.NewJSONObject()
	js.FluentPut("code", constants.TOKEN_ERROR)
	js.FluentPut("msg", msg)
	js.FluentPut("data", "")
	return js.GetData()
}

func excludeUrl(url string) bool {
	if len(prefixExcludeUrls) > 0 {
		for _, v := range prefixExcludeUrls {
			if strings.HasPrefix(url, v) {
				return true
			}
		}
	}

	if len(fullExcludeUrls) > 0 {
		for _, v := range fullExcludeUrls {
			if url == v {
				return true
			}
			tmp := "/" + v
			if url == tmp {
				return true
			}
		}
	}
	return false
}

/*
true: 过期
false: 未过期
*/
func compareDate(utcStr string, data string) (bool, error) {
	t, err := time.Parse(constants.UTC_LAYOUT, utcStr)
	if err != nil {
		return false, fmt.Errorf("Error parsing the date: %w", err)
	}
	now := time.Now().UTC()
	if now.After(t) {
		return true, nil
	}
	js := jsonutil.NewJSONObject()
	js.ParseObject(data)
	id := js.GetFloat64("id")

	expire := service.GetCache(constants.GetUserExpireTimeKey(int32(id)))
	if utils.IsEmpty(expire) {
		return true, nil
	}

	utcDate, err := utils.ParseUTC(expire)
	if err != nil {
		return true, nil
	}

	if now.After(utcDate) {
		return true, nil
	}

	return false, nil
}
