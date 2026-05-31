package notification

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/service/notification"
	"questionhelper-server/pkg/response"
)

type NotificationAdminController struct{}

func NewNotificationAdminController() *NotificationAdminController {
	return &NotificationAdminController{}
}

// ==================== 通知模板管理 ====================

// ListTemplates 模板列表
func (ctrl *NotificationAdminController) ListTemplates(c *gin.Context) {
	var req dto.NotificationTemplateListRequest
	c.ShouldBindQuery(&req)

	list, total, err := notification.ListTemplates(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// CreateTemplate 创建模板
func (ctrl *NotificationAdminController) CreateTemplate(c *gin.Context) {
	var req dto.CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	result, err := notification.CreateTemplate(&req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, result)
}

// UpdateTemplate 更新模板
func (ctrl *NotificationAdminController) UpdateTemplate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的模板ID")
		return
	}

	var req dto.UpdateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	result, err := notification.UpdateTemplate(uint(id), &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, result)
}

// DeleteTemplate 删除模板
func (ctrl *NotificationAdminController) DeleteTemplate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的模板ID")
		return
	}

	if err := notification.DeleteTemplate(uint(id)); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.SuccessWithMessage(c, "删除成功", nil)
}

// ==================== 通知渠道管理 ====================

// ListChannels 渠道列表
func (ctrl *NotificationAdminController) ListChannels(c *gin.Context) {
	var req dto.ChannelListRequest
	c.ShouldBindQuery(&req)

	list, total, err := notification.ListChannels(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// UpdateChannel 更新渠道配置
func (ctrl *NotificationAdminController) UpdateChannel(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的渠道ID")
		return
	}

	var req dto.UpdateChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	result, err := notification.UpdateChannel(uint(id), &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, result)
}
