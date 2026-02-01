package controller

import (
	"net/url"
	this "service-common/src/main/controller"
	"service-demo/src/main/service"
	"shared/vo"

	"github.com/gin-gonic/gin"
)

type DemoController struct {
	this.BaseController
}

var demoService = service.DemoService{}

func Demo1(c *gin.Context) {
	userInfo := this.GetUserInfo(c)
	r := demoService.Demo1(userInfo)
	this.Result(c, r)
}

func Demo2(c *gin.Context) {
	fullBaseCode := this.GetString(c, "bas-code")
	r := demoService.Demo2(fullBaseCode)
	this.Result(c, r)
}

func Demo3(c *gin.Context) {
	pageVO := this.BindJSON[vo.UserVO](c)
	r := demoService.Demo3(pageVO)
	this.Result(c, r)
}

func UploadFile(c *gin.Context) {
	file := this.GetFile(c)
	r := demoService.UploadFile(file.File)
	this.Result(c, r)
}

func DownloadFile(c *gin.Context) {
	r, fileName := demoService.DownloadFile()
	if r.Code != 200 {
		this.HttpFail(c, 500, r.Msg)
	} else {
		filePath := r.Data.(string)
		escaped := url.QueryEscape(fileName)
		c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Disposition", "attachment; filename*=UTF-8''"+escaped)
		c.File(filePath)
	}
}
