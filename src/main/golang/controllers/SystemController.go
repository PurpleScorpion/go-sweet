package controllers

import (
	"encoding/json"
	"sweet-common/vo"
	"sweet-src/main/golang/models"
	"sweet-src/main/golang/service"
)

type SystemController struct {
	BaseController
}

var systemService service.SystemService

func (that *SystemController) UserPageData() {
	var user vo.UserPageVO
	data := that.Ctx.Input.RequestBody
	json.Unmarshal(data, &user)
	r := systemService.UserPageData(user)
	that.Result(r)
}

func (that *SystemController) RolePageData() {
	var role vo.RolePageVO
	data := that.Ctx.Input.RequestBody
	json.Unmarshal(data, &role)
	r := systemService.RolePageData(role)
	that.Result(r)
}

func (that *SystemController) MenuPageData() {
	r := systemService.MenuPageData()
	that.Result(r)
}

func (that *SystemController) DeleteMenuById() {
	id, _ := that.GetInt32(":id")
	r := systemService.DeleteMenuById(id)
	that.Result(r)
}

func (that *SystemController) GetMenuById() {
	id, _ := that.GetInt32(":id")
	r := systemService.GetMenuById(id)
	that.Result(r)
}

func (that *SystemController) AllParentMenu() {
	r := systemService.AllParentMenu()
	that.Result(r)
}

func (that *SystemController) AllRole() {
	r := systemService.AllRole()
	that.Result(r)
}

func (that *SystemController) MenuInsert() {
	data := that.Ctx.Input.RequestBody
	var menu models.SysMenu
	json.Unmarshal(data, &menu)

	r := systemService.MenuInsert(menu)
	that.Result(r)
}

func (that *SystemController) MenuUpdate() {
	var menu models.SysMenu
	data := that.Ctx.Input.RequestBody
	json.Unmarshal(data, &menu)
	r := systemService.MenuUpdate(menu)
	that.Result(r)
}

func (that *SystemController) RoleInsert() {
	data := that.Ctx.Input.RequestBody
	var role vo.RolePageVO
	json.Unmarshal(data, &role)

	id := getUserId(that.Ctx)
	role.UserId = id

	r := systemService.RoleInsert(role)
	that.Result(r)
}

func (that *SystemController) GetRoleById() {
	id, _ := that.GetInt32(":id")
	r := systemService.GetRoleById(id)
	that.Result(r)
}

func (that *SystemController) RoleUpdate() {
	data := that.Ctx.Input.RequestBody
	var role vo.RolePageVO
	json.Unmarshal(data, &role)

	id := getUserId(that.Ctx)
	role.UserId = id

	r := systemService.RoleUpdate(role)
	that.Result(r)
}

func (that *SystemController) DeleteRoleById() {
	Id, _ := that.GetInt32(":id")
	r := systemService.DeleteRoleById(Id)
	that.Result(r)
}

func (that *SystemController) UserInsert() {
	var userVO vo.UserVO
	data := that.Ctx.Input.RequestBody
	json.Unmarshal(data, &userVO)

	id := getUserId(that.Ctx)
	userVO.CreatedBy = id
	userVO.LastModifiedBy = id

	r := systemService.UserInsert(userVO)
	that.Result(r)
}

func (that *SystemController) UserUpdate() {
	var userVO vo.UserVO
	data := that.Ctx.Input.RequestBody
	json.Unmarshal(data, &userVO)

	id := getUserId(that.Ctx)
	userVO.LastModifiedBy = id

	r := systemService.UserUpdate(userVO)
	that.Result(r)
}

func (that *SystemController) GetUserById() {
	id, _ := that.GetInt32(":id")
	r := systemService.GetUserById(id)
	that.Result(r)
}

func (that *SystemController) ChangeUserStatus() {
	var userVO vo.UserVO
	data := that.Ctx.Input.RequestBody
	json.Unmarshal(data, &userVO)

	id := getUserId(that.Ctx)
	userVO.LastModifiedBy = id

	r := systemService.ChangeUserStatus(userVO)
	that.Result(r)
}

func (that *SystemController) DeleteUserById() {
	id, _ := that.GetInt32(":id")

	userId := getUserId(that.Ctx)
	r := systemService.DeleteUserById(id, userId)
	that.Result(r)
}
