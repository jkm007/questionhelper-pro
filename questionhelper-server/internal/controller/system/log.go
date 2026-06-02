package system

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/service/system"
	"questionhelper-server/pkg/response"
)

// ListOperationLogs 操作日志列表
// @Summary      获取操作日志列表
// @Description  获取系统操作日志列表，支持分页
// @Tags         系统日志
// @Accept       json
// @Produce      json
// @Param        page       query     int  false  "页码"
// @Param        page_size  query     int  false  "每页数量"
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/logs/operation [get]
// @Security     BearerAuth
func (ctrl *SystemController) ListOperationLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	logs, total, err := system.ListOperationLogs(page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, logs, total, page, pageSize)
}

// ListLoginLogs 登录日志列表
// @Summary      获取登录日志列表
// @Description  获取系统登录日志列表，支持分页
// @Tags         系统日志
// @Accept       json
// @Produce      json
// @Param        page       query     int  false  "页码"
// @Param        page_size  query     int  false  "每页数量"
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/logs/login [get]
// @Security     BearerAuth
func (ctrl *SystemController) ListLoginLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	logs, total, err := system.ListLoginLogs(page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, logs, total, page, pageSize)
}

// ListSystemLogs 系统日志列表
// @Summary      获取系统日志列表
// @Description  获取系统日志列表，支持分页和条件筛选
// @Tags         系统日志
// @Accept       json
// @Produce      json
// @Param        page       query     int  false  "页码"
// @Param        page_size  query     int  false  "每页数量"
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/logs/system [get]
// @Security     BearerAuth
func (ctrl *SystemController) ListSystemLogs(c *gin.Context) {
	var query dto.SystemLogQuery
	c.ShouldBindQuery(&query)
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 || query.PageSize > 100 {
		query.PageSize = 10
	}

	logs, total, err := system.ListSystemLogs(&query)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, logs, total, query.Page, query.PageSize)
}

// ListErrorLogs 错误日志列表
// @Summary      获取错误日志列表
// @Description  获取系统错误日志列表，支持分页和条件筛选
// @Tags         系统日志
// @Accept       json
// @Produce      json
// @Param        page       query     int  false  "页码"
// @Param        page_size  query     int  false  "每页数量"
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/logs/error [get]
// @Security     BearerAuth
func (ctrl *SystemController) ListErrorLogs(c *gin.Context) {
	var query dto.ErrorLogQuery
	c.ShouldBindQuery(&query)
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 || query.PageSize > 100 {
		query.PageSize = 10
	}

	logs, total, err := system.ListErrorLogs(&query)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, logs, total, query.Page, query.PageSize)
}

// SearchLogs 日志搜索
// @Summary      搜索日志
// @Description  根据条件搜索系统日志
// @Tags         系统日志
// @Accept       json
// @Produce      json
// @Param        page       query     int     false  "页码"
// @Param        page_size  query     int     false  "每页数量"
// @Param        keyword    query     string  false  "搜索关键词"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/logs/search [get]
// @Security     BearerAuth
func (ctrl *SystemController) SearchLogs(c *gin.Context) {
	var query dto.LogSearchRequest
	c.ShouldBindQuery(&query)
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 || query.PageSize > 100 {
		query.PageSize = 10
	}

	logs, total, err := system.SearchLogs(&query)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Page(c, logs, total, query.Page, query.PageSize)
}

// ArchiveLogs 日志归档
// @Summary      归档日志
// @Description  根据条件归档系统日志
// @Tags         系统日志
// @Accept       json
// @Produce      json
// @Param        req  body      dto.LogArchiveRequest  true  "归档参数"
// @Success      200  {object}  response.Response  "归档成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/logs/archive [post]
// @Security     BearerAuth
func (ctrl *SystemController) ArchiveLogs(c *gin.Context) {
	var req dto.LogArchiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	count, err := system.ArchiveLogs(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "归档成功", gin.H{"archived_count": count})
}

// GetLogStats 日志统计
// @Summary      获取日志统计
// @Description  获取系统日志统计数据
// @Tags         系统日志
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/logs/stats [get]
// @Security     BearerAuth
func (ctrl *SystemController) GetLogStats(c *gin.Context) {
	stats, err := system.GetLogStats()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, stats)
}
