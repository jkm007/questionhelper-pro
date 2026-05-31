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
	}
}
