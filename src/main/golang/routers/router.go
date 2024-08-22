package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"sweet-src/main/golang/controllers"
)

func Init() {
	beego.Router("/login", &controllers.UserController{}, "post:Login")
	beego.Router("/rePassword", &controllers.UserController{}, "post:RePassword")
	beego.Router("/sys/user/healthCheck", &controllers.UserController{}, "get:HealthCheck")
	// system

	beego.Router("/sys/user/pageData", &controllers.SystemController{}, "post:UserPageData")
	beego.Router("/sys/menu/pageData", &controllers.SystemController{}, "post:MenuPageData")
	beego.Router("/sys/role/pageData", &controllers.SystemController{}, "post:RolePageData")

	beego.Router("/sys/user/userUpdate", &controllers.SystemController{}, "post:UserUpdate")
	beego.Router("/sys/user/userInsert", &controllers.SystemController{}, "post:UserInsert")
	beego.Router("/sys/user/getUserById/:id", &controllers.SystemController{}, "get:GetUserById")
	beego.Router("/sys/user/changeUserStatus", &controllers.SystemController{}, "post:ChangeUserStatus")
	beego.Router("/sys/user/deleteUserById/:id", &controllers.SystemController{}, "get:DeleteUserById")

	beego.Router("/sys/menu/allParentMenu", &controllers.SystemController{}, "get:AllParentMenu")
	beego.Router("/sys/menu/deleteMenuById/:id", &controllers.SystemController{}, "get:DeleteMenuById")
	beego.Router("/sys/menu/getMenuById/:id", &controllers.SystemController{}, "get:GetMenuById")
	beego.Router("/sys/menu/menuInsert", &controllers.SystemController{}, "post:MenuInsert")
	beego.Router("/sys/menu/menuUpdate", &controllers.SystemController{}, "post:MenuUpdate")
	beego.Router("/sys/menu/menuUpdate", &controllers.SystemController{}, "post:MenuUpdate")

	beego.Router("/sys/role/roleInsert", &controllers.SystemController{}, "post:RoleInsert")
	beego.Router("/sys/role/roleUpdate", &controllers.SystemController{}, "post:RoleUpdate")
	beego.Router("/sys/role/getRoleById/:id", &controllers.SystemController{}, "get:GetRoleById")
	beego.Router("/sys/role/deleteRoleById/:id", &controllers.SystemController{}, "get:DeleteRoleById")
	beego.Router("/sys/role/getAllRole", &controllers.SystemController{}, "get:AllRole")

	beego.Router("/upload/uploadImg", &controllers.FileController{}, "post:UploadImg")
	beego.Router("/upload/getBaseImg", &controllers.FileController{}, "get:GetBaseImg")

	//
}
