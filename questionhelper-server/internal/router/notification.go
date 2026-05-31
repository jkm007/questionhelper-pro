package router

import (
	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/controller/notification"
)

func SetupNotificationRoutes(r *gin.RouterGroup, ctrl *notification.NotificationController) {
	// 兼容前端 /notifications 接口
	n := r.Group("/notifications")
	{
		n.GET("", ctrl.ListNotifications)
		n.GET("/unread-count", ctrl.GetUnreadCount)
		n.PUT("/:id/read", ctrl.MarkAsRead)
		n.PUT("/read-all", ctrl.MarkAllAsRead)
		n.DELETE("/:id", ctrl.DeleteNotification)

		// 撤回通知
		n.PUT("/:id/recall", ctrl.RecallNotification)

		// 群发通知
		n.POST("/batch-send", ctrl.BatchSend)

		// 批量操作
		n.PUT("/batch-read", ctrl.BatchMarkAsRead)
		n.DELETE("/batch-delete", ctrl.BatchDeleteNotifications)

		// 定时通知
		n.POST("/scheduled", ctrl.CreateScheduled)
		n.GET("/scheduled", ctrl.ListScheduled)
		n.DELETE("/scheduled/:id", ctrl.DeleteScheduled)

		// 通知设置
		n.GET("/settings", ctrl.GetNotificationSettings)
		n.PUT("/settings", ctrl.UpdateNotificationSettings)

		// 通知统计
		n.GET("/stats", ctrl.GetNotificationStats)
	}
}

func SetupAdminNotificationRoutes(r *gin.RouterGroup, ctrl *notification.NotificationAdminController) {
	// 通知模板管理
	templates := r.Group("/notifications/templates")
	{
		templates.GET("", ctrl.ListTemplates)
		templates.POST("", ctrl.CreateTemplate)
		templates.PUT("/:id", ctrl.UpdateTemplate)
		templates.DELETE("/:id", ctrl.DeleteTemplate)
	}

	// 通知渠道管理
	channels := r.Group("/notifications/channels")
	{
		channels.GET("", ctrl.ListChannels)
		channels.PUT("/:id", ctrl.UpdateChannel)
	}
}
