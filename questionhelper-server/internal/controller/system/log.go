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
func (ctrl *SystemController) GetLogStats(c *gin.Context) {
	stats, err := system.GetLogStats()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, stats)
}
