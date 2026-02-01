package controller

import (
	"service-auth/src/main/service"
	this "service-common/src/main/controller"
	"service-common/src/main/models"
	"shared/vo"

	"github.com/gin-gonic/gin"
)

type SystemController struct {
	this.BaseController
}

var systemService service.SystemService

func UserPageData(c *gin.Context) {
	user := this.BindJSON[vo.UserPageVO](c)
	r := systemService.UserPageData(user)
	this.Result(c, r)
}

func RolePageData(c *gin.Context) {
	role := this.BindJSON[vo.RolePageVO](c)
	r := systemService.RolePageData(role)
	this.Result(c, r)
}

func MenuPageData(c *gin.Context) {
	r := systemService.MenuPageData()
	this.Result(c, r)
}

func DeleteMenuById(c *gin.Context) {
	id := this.GetInt(c, "id")
	r := systemService.DeleteMenuById(id)
	this.Result(c, r)
}

func GetMenuById(c *gin.Context) {
	id := this.GetInt(c, "id")
	r := systemService.GetMenuById(id)
	this.Result(c, r)
}

func AllParentMenu(c *gin.Context) {
	r := systemService.AllParentMenu()
	this.Result(c, r)
}

func AllRole(c *gin.Context) {
	r := systemService.AllRole()
	this.Result(c, r)
}

func MenuInsert(c *gin.Context) {
	menu := this.BindJSON[models.SysMenu](c)
	r := systemService.MenuInsert(menu)
	this.Result(c, r)
}

func MenuUpdate(c *gin.Context) {
	menu := this.BindJSON[models.SysMenu](c)
	r := systemService.MenuUpdate(menu)
	this.Result(c, r)
}

func RoleInsert(c *gin.Context) {
	role := this.BindJSON[vo.RolePageVO](c)
	r := systemService.RoleInsert(role)
	this.Result(c, r)
}

func GetRoleById(c *gin.Context) {
	id := this.GetInt(c, "id")
	r := systemService.GetRoleById(id)
	this.Result(c, r)
}

func RoleUpdate(c *gin.Context) {
	role := this.BindJSON[vo.RolePageVO](c)
	r := systemService.RoleUpdate(role)
	this.Result(c, r)
}

func DeleteRoleById(c *gin.Context) {
	id := this.GetInt(c, "id")
	r := systemService.DeleteRoleById(id)
	this.Result(c, r)
}

func UserInsert(c *gin.Context) {
	userVO := this.BindJSON[vo.UserVO](c)
	r := systemService.UserInsert(userVO)
	this.Result(c, r)
}

func UserUpdate(c *gin.Context) {
	userVO := this.BindJSON[vo.UserVO](c)
	r := systemService.UserUpdate(userVO)
	this.Result(c, r)
}

func GetUserById(c *gin.Context) {
	id := this.GetInt(c, "id")
	r := systemService.GetUserById(id)
	this.Result(c, r)
}

func ChangeUserStatus(c *gin.Context) {
	userVO := this.BindJSON[vo.UserVO](c)
	r := systemService.ChangeUserStatus(userVO)
	this.Result(c, r)
}

func DeleteUserById(c *gin.Context) {
	id := this.GetInt(c, "id")
	r := systemService.DeleteUserById(id)
	this.Result(c, r)
}
