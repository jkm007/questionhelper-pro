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
// @Summary      模板列表
// @Description  分页获取通知模板列表
// @Tags         通知模板
// @Accept       json
// @Produce      json
// @Param        page       query     int     false  "页码"
// @Param        page_size  query     int     false  "每页数量"
// @Param        keyword    query     string  false  "搜索关键词"
// @Success      200  {object}  response.Response{data=object{list=[]dto.TemplateInfo,total=int64,page=int,page_size=int}}  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/notifications/templates [get]
// @Security     BearerAuth
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
// @Summary      创建模板
// @Description  创建一条新的通知模板
// @Tags         通知模板
// @Accept       json
// @Produce      json
// @Param        req  body      dto.CreateTemplateRequest  true  "模板信息"
// @Success      200  {object}  response.Response{data=dto.TemplateInfo}  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/notifications/templates [post]
// @Security     BearerAuth
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
// @Summary      更新模板
// @Description  更新指定通知模板
// @Tags         通知模板
// @Accept       json
// @Produce      json
// @Param        id   path      uint                       true  "模板ID"
// @Param        req  body      dto.UpdateTemplateRequest   true  "模板信息"
// @Success      200  {object}  response.Response{data=dto.TemplateInfo}  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/notifications/templates/{id} [put]
// @Security     BearerAuth
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
// @Summary      删除模板
// @Description  删除指定通知模板
// @Tags         通知模板
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "模板ID"
// @Success      200  {object}  response.Response  "删除成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/notifications/templates/{id} [delete]
// @Security     BearerAuth
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
// @Summary      渠道列表
// @Description  分页获取通知渠道列表
// @Tags         通知渠道
// @Accept       json
// @Produce      json
// @Param        page       query     int  false  "页码"
// @Param        page_size  query     int  false  "每页数量"
// @Success      200  {object}  response.Response{data=object{list=[]dto.ChannelInfo,total=int64,page=int,page_size=int}}  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/notifications/channels [get]
// @Security     BearerAuth
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
// @Summary      更新渠道配置
// @Description  更新指定通知渠道配置
// @Tags         通知渠道
// @Accept       json
// @Produce      json
// @Param        id   path      uint                     true  "渠道ID"
// @Param        req  body      dto.UpdateChannelRequest  true  "渠道配置"
// @Success      200  {object}  response.Response{data=dto.ChannelInfo}  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/notifications/channels/{id} [put]
// @Security     BearerAuth
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
