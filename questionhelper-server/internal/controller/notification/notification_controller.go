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
