package log

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/dto"
	logService "questionhelper-server/internal/service/log"
	"questionhelper-server/pkg/response"
)

type LogController struct{}

func NewLogController() *LogController {
	return &LogController{}
}

// ListOperationLogs 操作日志列表
func (ctrl *LogController) ListOperationLogs(c *gin.Context) {
	var req dto.OperationLogListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	list, total, err := logService.ListOperationLogs(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// GetOperationLog 获取操作日志详情
func (ctrl *LogController) GetOperationLog(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的日志ID")
		return
	}

	info, err := logService.GetOperationLog(uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, info)
}

// ExportOperationLogs 导出操作日志
func (ctrl *LogController) ExportOperationLogs(c *gin.Context) {
	var req dto.OperationLogListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	buf, err := logService.ExportOperationLogs(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=operation_logs_%s.csv",
		req.StartTime.Format("20060102")))
	c.Data(http.StatusOK, "text/csv; charset=utf-8", buf.Bytes())
}

// ListLoginLogs 登录日志列表
func (ctrl *LogController) ListLoginLogs(c *gin.Context) {
	var req dto.LoginLogListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	list, total, err := logService.ListLoginLogs(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// ExportLoginLogs 导出登录日志
func (ctrl *LogController) ExportLoginLogs(c *gin.Context) {
	var req dto.LoginLogListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	buf, err := logService.ExportLoginLogs(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=login_logs_%s.csv",
		req.StartTime.Format("20060102")))
	c.Data(http.StatusOK, "text/csv; charset=utf-8", buf.Bytes())
}

// CleanOperationLogs 清理操作日志
func (ctrl *LogController) CleanOperationLogs(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "90"))

	count, err := logService.CleanOperationLogs(days)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, fmt.Sprintf("清理 %d 条日志", count), nil)
}

// CleanLoginLogs 清理登录日志
func (ctrl *LogController) CleanLoginLogs(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "90"))

	count, err := logService.CleanLoginLogs(days)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, fmt.Sprintf("清理 %d 条日志", count), nil)
}
