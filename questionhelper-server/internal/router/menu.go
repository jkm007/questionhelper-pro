package router

import (
	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/controller/menu"
)

func SetupMenuRoutes(r *gin.RouterGroup, ctrl *menu.MenuController) {
	m := r.Group("/menu")
	{
		m.GET("/user", ctrl.GetUserMenus)
		m.GET("/buttons", ctrl.GetUserButtons)
	}

	// 兼容前端 /menus/routes 接口
	menus := r.Group("/menus")
	{
		menus.GET("/routes", ctrl.GetUserRoutes)
	}
}

func SetupAdminMenuRoutes(r *gin.RouterGroup, ctrl *menu.MenuController) {
	m := r.Group("/menus")
	{
		m.GET("", ctrl.ListMenus)
		m.GET("/tree", ctrl.GetMenuTree)
		m.GET("/:id", ctrl.GetMenu)
		m.POST("", ctrl.CreateMenu)
		m.PUT("/:id", ctrl.UpdateMenu)
		m.DELETE("/:id", ctrl.DeleteMenu)
	}
}
