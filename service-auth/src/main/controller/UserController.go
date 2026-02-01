package controller

import (
	"service-auth/src/main/service"
	this "service-common/src/main/controller"
	"service-common/src/main/models"
	"shared/vo"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	this.BaseController
}

var userService service.UserService

func Login(c *gin.Context) {
	user := this.BindJSON[models.User](c)
	r := userService.Login(user)
	this.Result(c, r)
}

func RePassword(c *gin.Context) {
	user := this.BindJSON[vo.UserVO](c)
	r := userService.RePassword(user)
	this.Result(c, r)
}
