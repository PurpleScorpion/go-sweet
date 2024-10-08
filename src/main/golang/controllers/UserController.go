package controllers

import (
	"encoding/json"
	"sweet-common/vo"
	"sweet-src/main/golang/models"
	"sweet-src/main/golang/service"
)

type UserController struct {
	BaseController
}

var userService service.UserService

func (that *UserController) Login() {
	var user models.User
	data := that.Ctx.Input.RequestBody
	json.Unmarshal(data, &user)
	r := userService.Login(user)
	that.Result(r)
}

func (that *UserController) RePassword() {
	var user vo.UserVO
	data := that.Ctx.Input.RequestBody
	json.Unmarshal(data, &user)
	r := userService.RePassword(user)
	that.Result(r)
}

func (that *UserController) HealthCheck() {
	id := getUserId(that.Ctx)
	r := userService.HealthCheck(id)
	that.Result(r)
}
