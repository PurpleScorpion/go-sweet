package venRouter

import (
	"service-auth/src/main/controller"

	"github.com/gin-gonic/gin"
)

func Init(router *gin.Engine) {
	{
		auth := router.Group("/auth")
		auth.POST("/login", controller.Login)
		auth.POST("/re-password", controller.RePassword)

	}
	{
		sys := router.Group("/sys")
		sys.POST("/user/page-data", controller.UserPageData)
		sys.POST("/menu/page-data", controller.MenuPageData)
		sys.POST("/role/page-data", controller.RolePageData)

		sys.POST("/user", controller.UserUpdate)
		sys.PUT("/user", controller.UserInsert)
		sys.GET("/user/:id", controller.GetUserById)
		sys.DELETE("/user/:id", controller.DeleteUserById)
		sys.POST("/user/status/change", controller.ChangeUserStatus)

		sys.GET("/menu", controller.AllParentMenu)
		sys.DELETE("/user/:id", controller.DeleteMenuById)
		sys.GET("/menu/:id", controller.GetMenuById)
		sys.PUT("/menu", controller.MenuInsert)
		sys.POST("/menu", controller.MenuUpdate)

		sys.GET("/role", controller.AllRole)
		sys.PUT("/role", controller.RoleInsert)
		sys.POST("/role", controller.RoleUpdate)
		sys.GET("/role/:id", controller.GetRoleById)
		sys.DELETE("/role/:id", controller.DeleteRoleById)
	}
}
