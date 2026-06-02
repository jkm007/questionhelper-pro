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
// @Summary      获取操作日志列表
// @Description  获取操作日志列表，支持分页和条件筛选
// @Tags         操作日志
// @Accept       json
// @Produce      json
// @Param        page       query     int  false  "页码"
// @Param        page_size  query     int  false  "每页数量"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/operation-logs [get]
// @Security     BearerAuth
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
// @Summary      获取操作日志详情
// @Description  根据ID获取操作日志详情
// @Tags         操作日志
// @Accept       json
// @Produce      json
// @Param        id  path      uint  true  "日志ID"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "无效的日志ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/operation-logs/{id} [get]
// @Security     BearerAuth
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
// @Summary      导出操作日志
// @Description  根据条件导出操作日志为CSV文件
// @Tags         操作日志
// @Accept       json
// @Produce      text/csv
// @Param        start_time  query     string  false  "开始时间"
// @Param        end_time    query     string  false  "结束时间"
// @Success      200  {string}  string  "CSV文件内容"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/operation-logs/export [get]
// @Security     BearerAuth
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
// @Summary      获取登录日志列表
// @Description  获取登录日志列表，支持分页和条件筛选
// @Tags         登录日志
// @Accept       json
// @Produce      json
// @Param        page       query     int  false  "页码"
// @Param        page_size  query     int  false  "每页数量"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/login-logs [get]
// @Security     BearerAuth
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
// @Summary      导出登录日志
// @Description  根据条件导出登录日志为CSV文件
// @Tags         登录日志
// @Accept       json
// @Produce      text/csv
// @Param        start_time  query     string  false  "开始时间"
// @Param        end_time    query     string  false  "结束时间"
// @Success      200  {string}  string  "CSV文件内容"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/login-logs/export [get]
// @Security     BearerAuth
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
// @Summary      清理操作日志
// @Description  清理指定天数之前的旧操作日志
// @Tags         操作日志
// @Accept       json
// @Produce      json
// @Param        days  query     int  false  "保留天数"  default(90)
// @Success      200  {object}  response.Response  "清理成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/operation-logs/clean [post]
// @Security     BearerAuth
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
// @Summary      清理登录日志
// @Description  清理指定天数之前的旧登录日志
// @Tags         登录日志
// @Accept       json
// @Produce      json
// @Param        days  query     int  false  "保留天数"  default(90)
// @Success      200  {object}  response.Response  "清理成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/login-logs/clean [post]
// @Security     BearerAuth
func (ctrl *LogController) CleanLoginLogs(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "90"))

	count, err := logService.CleanLoginLogs(days)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, fmt.Sprintf("清理 %d 条日志", count), nil)
}
