package router

import (
	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/controller/permission"
)

func SetupAdminButtonPermissionRoutes(r *gin.RouterGroup, ctrl *permission.PermissionController) {
	p := r.Group("/button-permissions")
	{
		p.GET("", ctrl.ListButtonPermissions)
		p.GET("/:id", ctrl.GetButtonPermission)
		p.POST("", ctrl.CreateButtonPermission)
		p.PUT("/:id", ctrl.UpdateButtonPermission)
		p.DELETE("/:id", ctrl.DeleteButtonPermission)
		p.GET("/menu/:menu_id", ctrl.ListMenuButtonPermissions)
	}
}
