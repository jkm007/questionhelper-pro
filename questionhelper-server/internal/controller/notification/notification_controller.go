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
// @Summary      通知列表
// @Description  分页获取当前用户的通知列表
// @Tags         通知管理
// @Accept       json
// @Produce      json
// @Param        page       query     int     false  "页码"
// @Param        page_size  query     int     false  "每页数量"
// @Param        type       query     string  false  "通知类型"
// @Param        is_read    query     bool    false  "是否已读"
// @Success      200  {object}  response.Response{data=object{list=[]dto.NotificationInfo,total=int64,page=int,page_size=int}}  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /notifications [get]
// @Security     BearerAuth
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
// @Summary      获取未读数量
// @Description  获取当前用户的未读通知数量
// @Tags         通知管理
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response{data=object{count=int64}}  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /notifications/unread-count [get]
// @Security     BearerAuth
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
// @Summary      标记已读
// @Description  将指定通知标记为已读
// @Tags         通知管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "通知ID"
// @Success      200  {object}  response.Response  "已标记已读"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /notifications/{id}/read [put]
// @Security     BearerAuth
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
// @Summary      全部标记已读
// @Description  将当前用户所有通知标记为已读
// @Tags         通知管理
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response  "全部已读"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /notifications/read-all [put]
// @Security     BearerAuth
func (ctrl *NotificationController) MarkAllAsRead(c *gin.Context) {
	userID := c.GetUint("user_id")

	if err := notification.MarkAllAsRead(userID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "全部已读", nil)
}

// DeleteNotification 删除通知
// @Summary      删除通知
// @Description  删除指定通知
// @Tags         通知管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "通知ID"
// @Success      200  {object}  response.Response  "删除成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /notifications/{id} [delete]
// @Security     BearerAuth
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
// @Summary      撤回通知
// @Description  撤回已发送的通知
// @Tags         通知管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "通知ID"
// @Success      200  {object}  response.Response  "撤回成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /notifications/{id}/recall [put]
// @Security     BearerAuth
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
// @Summary      群发通知
// @Description  批量发送通知给多个用户
// @Tags         通知管理
// @Accept       json
// @Produce      json
// @Param        req  body      dto.BatchSendRequest  true  "群发信息"
// @Success      200  {object}  response.Response{data=dto.BatchSendResponse}  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /notifications/batch-send [post]
// @Security     BearerAuth
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
// @Summary      批量标记已读
// @Description  批量将指定通知标记为已读
// @Tags         通知管理
// @Accept       json
// @Produce      json
// @Param        req  body      dto.NotificationBatchReadRequest  true  "通知ID列表"
// @Success      200  {object}  response.Response  "批量标记已读成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /notifications/batch-read [put]
// @Security     BearerAuth
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
// @Summary      批量删除通知
// @Description  批量删除指定通知
// @Tags         通知管理
// @Accept       json
// @Produce      json
// @Param        req  body      dto.NotificationBatchDeleteRequest  true  "通知ID列表"
// @Success      200  {object}  response.Response  "批量删除成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /notifications/batch-delete [delete]
// @Security     BearerAuth
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
// @Summary      创建定时通知
// @Description  创建一条定时发送的通知
// @Tags         通知管理
// @Accept       json
// @Produce      json
// @Param        req  body      dto.CreateScheduledRequest  true  "定时通知信息"
// @Success      200  {object}  response.Response{data=dto.ScheduledNotificationInfo}  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /notifications/scheduled [post]
// @Security     BearerAuth
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
// @Summary      定时通知列表
// @Description  分页获取当前用户的定时通知列表
// @Tags         通知管理
// @Accept       json
// @Produce      json
// @Param        page       query     int  false  "页码"
// @Param        page_size  query     int  false  "每页数量"
// @Success      200  {object}  response.Response{data=object{list=[]dto.ScheduledNotificationInfo,total=int64,page=int,page_size=int}}  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /notifications/scheduled [get]
// @Security     BearerAuth
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
// @Summary      取消定时通知
// @Description  取消指定的定时通知
// @Tags         通知管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "定时通知ID"
// @Success      200  {object}  response.Response  "取消成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /notifications/scheduled/{id} [delete]
// @Security     BearerAuth
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
// @Summary      获取通知设置
// @Description  获取当前用户的通知设置
// @Tags         通知管理
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response{data=dto.NotificationSettingInfo}  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /notifications/settings [get]
// @Security     BearerAuth
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
// @Summary      更新通知设置
// @Description  更新当前用户的通知设置
// @Tags         通知管理
// @Accept       json
// @Produce      json
// @Param        req  body      dto.UpdateNotificationSettingRequest  true  "通知设置"
// @Success      200  {object}  response.Response  "更新成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /notifications/settings [put]
// @Security     BearerAuth
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
// @Summary      获取通知统计
// @Description  获取当前用户的通知统计数据
// @Tags         通知管理
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response{data=dto.NotificationStats}  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /notifications/stats [get]
// @Security     BearerAuth
func (ctrl *NotificationController) GetNotificationStats(c *gin.Context) {
	userID := c.GetUint("user_id")

	stats, err := notification.GetNotificationStats(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, stats)
}
