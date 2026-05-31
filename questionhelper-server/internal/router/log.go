package router

import (
	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/controller/log"
)

func SetupAdminLogRoutes(r *gin.RouterGroup, ctrl *log.LogController) {
	// 操作日志
	operationLog := r.Group("/operation-logs")
	{
		operationLog.GET("", ctrl.ListOperationLogs)
		operationLog.GET("/:id", ctrl.GetOperationLog)
		operationLog.GET("/export", ctrl.ExportOperationLogs)
		operationLog.POST("/clean", ctrl.CleanOperationLogs)
	}

	// 登录日志
	loginLog := r.Group("/login-logs")
	{
		loginLog.GET("", ctrl.ListLoginLogs)
		loginLog.GET("/export", ctrl.ExportLoginLogs)
		loginLog.POST("/clean", ctrl.CleanLoginLogs)
	}
}
