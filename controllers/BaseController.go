package controllers

import (
	"github.com/PurpleScorpion/go-sweet-json/jsonutil"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"go-sweet/common/utils"
)

type BaseController struct {
	beego.Controller
}

type ReturnMsg struct {
	Code int32       `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (r *BaseController) Ok(msg string, data interface{}) {

	if msg == "" {
		msg = "success"
	}
	res := ReturnMsg{
		200, msg, data,
	}
	r.Data["json"] = res
	r.ServeJSON()
	r.StopRun()
}

func (r *BaseController) Error(code int32, msg string, data interface{}) {
	res := ReturnMsg{
		code, msg, data,
	}
	r.Data["json"] = res
	r.ServeJSON()
	r.StopRun()
}

func GetRequestHeader(headerName string, c *beego.Controller) string {
	headers := c.Ctx.Request.Header
	values := headers[headerName]

	if len(values) > 0 {
		return values[0]
	}

	return ""
}

func (r *BaseController) Result(rs utils.R) {
	res := ReturnMsg{
		rs.Code, rs.Msg, rs.Data,
	}
	r.Data["json"] = res
	r.ServeJSON()
	r.StopRun()
}

func getUserId(ctx *context.Context) int32 {
	headers := ctx.Request.Header
	values := headers["Token"]
	token := values[0]
	data, _ := utils.Decrypt(token)
	js := jsonutil.NewJSONObject()
	js.ParseObject(data)
	id := js.GetFloat64("id")
	return int32(id)
}
