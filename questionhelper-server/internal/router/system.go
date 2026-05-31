package router

import (
	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/controller/system"
)

func SetupAdminSystemRoutes(r *gin.RouterGroup, ctrl *system.SystemController) {
	// 系统设置 (设计文档路径: /settings)
	settings := r.Group("/settings")
	{
		settings.GET("", ctrl.GetSettings)
		settings.PUT("", ctrl.UpdateSettings)
	}

	// 操作日志 (设计文档路径: /logs/operation, /logs/login)
	logs := r.Group("/logs")
	{
		logs.GET("/operation", ctrl.ListOperationLogs)
		logs.GET("/login", ctrl.ListLoginLogs)
	}
}
