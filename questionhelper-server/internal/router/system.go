package router

import (
	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/controller/system"
)

func SetupAdminSystemRoutes(r *gin.RouterGroup, ctrl *system.SystemController) {
	// 系统设置
	configs := r.Group("/configs")
	{
		configs.GET("", ctrl.GetSettings)
		configs.PUT("", ctrl.UpdateSettings)
	}
}
