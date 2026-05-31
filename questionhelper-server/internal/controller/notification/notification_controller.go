package notification

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/service/notification"
	"questionhelper-server/pkg/response"
)

type NotificationController struct{}

func NewNotificationController() *NotificationController {
	return &NotificationController{}
}

// ListNotifications 通知列表
func (ctrl *NotificationController) ListNotifications(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.NotificationListRequest
	c.ShouldBindQuery(&req)

	list, total, err := notification.ListNotifications(userID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// GetUnreadCount 获取未读数量
func (ctrl *NotificationController) GetUnreadCount(c *gin.Context) {
	userID := c.GetUint("user_id")

	count, err := notification.GetUnreadCount(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, gin.H{"count": count})
}

// MarkAsRead 标记已读
func (ctrl *NotificationController) MarkAsRead(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的通知ID")
		return
	}

	if err := notification.MarkAsRead(uint(id), userID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "已标记已读", nil)
}

// MarkAllAsRead 全部标记已读
func (ctrl *NotificationController) MarkAllAsRead(c *gin.Context) {
	userID := c.GetUint("user_id")

	if err := notification.MarkAllAsRead(userID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "全部已读", nil)
}

// DeleteNotification 删除通知
func (ctrl *NotificationController) DeleteNotification(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的通知ID")
		return
	}

	if err := notification.DeleteNotification(uint(id), userID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "删除成功", nil)
}

// RecallNotification 撤回通知
func (ctrl *NotificationController) RecallNotification(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的通知ID")
		return
	}

	if err := notification.RecallNotification(uint(id), userID); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.SuccessWithMessage(c, "撤回成功", nil)
}

// BatchSend 群发通知
func (ctrl *NotificationController) BatchSend(c *gin.Context) {
	senderID := c.GetUint("user_id")

	var req dto.BatchSendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	result, err := notification.BatchSend(senderID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, result)
}

// BatchMarkAsRead 批量标记已读
func (ctrl *NotificationController) BatchMarkAsRead(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.NotificationBatchReadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := notification.BatchMarkAsRead(userID, req.IDs); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "批量标记已读成功", nil)
}

// BatchDeleteNotifications 批量删除
func (ctrl *NotificationController) BatchDeleteNotifications(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.NotificationBatchDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := notification.BatchDeleteNotifications(userID, req.IDs); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "批量删除成功", nil)
}

// CreateScheduled 创建定时通知
func (ctrl *NotificationController) CreateScheduled(c *gin.Context) {
	senderID := c.GetUint("user_id")

	var req dto.CreateScheduledRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	result, err := notification.CreateScheduled(senderID, &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, result)
}

// ListScheduled 定时通知列表
func (ctrl *NotificationController) ListScheduled(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.ScheduledListRequest
	c.ShouldBindQuery(&req)

	list, total, err := notification.ListScheduled(userID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// DeleteScheduled 取消定时通知
func (ctrl *NotificationController) DeleteScheduled(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的定时通知ID")
		return
	}

	if err := notification.DeleteScheduled(uint(id), userID); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.SuccessWithMessage(c, "取消成功", nil)
}

// GetNotificationSettings 获取通知设置
func (ctrl *NotificationController) GetNotificationSettings(c *gin.Context) {
	userID := c.GetUint("user_id")

	settings, err := notification.GetNotificationSettings(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, settings)
}

// UpdateNotificationSettings 更新通知设置
func (ctrl *NotificationController) UpdateNotificationSettings(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.UpdateNotificationSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := notification.UpdateNotificationSettings(userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "更新成功", nil)
}

// GetNotificationStats 获取通知统计
func (ctrl *NotificationController) GetNotificationStats(c *gin.Context) {
	userID := c.GetUint("user_id")

	stats, err := notification.GetNotificationStats(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, stats)
}
