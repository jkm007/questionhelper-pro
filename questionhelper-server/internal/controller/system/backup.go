package system

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/service/system"
	"questionhelper-server/pkg/response"
)

// CreateBackup 创建备份
// @Summary      创建备份
// @Description  创建系统备份
// @Tags         备份管理
// @Accept       json
// @Produce      json
// @Param        req  body      dto.CreateBackupRequest  true  "备份参数"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/backup/create [post]
// @Security     BearerAuth
func (ctrl *SystemController) CreateBackup(c *gin.Context) {
	var req dto.CreateBackupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	record, err := system.CreateBackup(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, record)
}

// ListBackupRecords 备份列表
// @Summary      获取备份列表
// @Description  获取备份记录列表，支持分页
// @Tags         备份管理
// @Accept       json
// @Produce      json
// @Param        page       query     int  false  "页码"
// @Param        page_size  query     int  false  "每页数量"
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/backup/list [get]
// @Security     BearerAuth
func (ctrl *SystemController) ListBackupRecords(c *gin.Context) {
	var query dto.BackupListQuery
	c.ShouldBindQuery(&query)
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 || query.PageSize > 100 {
		query.PageSize = 10
	}

	records, total, err := system.ListBackupRecords(&query)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, records, total, query.Page, query.PageSize)
}

// RestoreBackup 恢复备份
// @Summary      恢复备份
// @Description  根据备份ID恢复系统备份
// @Tags         备份管理
// @Accept       json
// @Produce      json
// @Param        id  path      uint  true  "备份ID"
// @Success      200  {object}  response.Response  "恢复任务已提交"
// @Failure      400  {object}  response.Response  "无效的备份ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/backup/{id}/restore [post]
// @Security     BearerAuth
func (ctrl *SystemController) RestoreBackup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的备份ID")
		return
	}

	if err := system.RestoreBackup(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "恢复任务已提交", nil)
}

// DeleteBackup 删除备份
// @Summary      删除备份
// @Description  根据备份ID删除备份记录
// @Tags         备份管理
// @Accept       json
// @Produce      json
// @Param        id  path      uint  true  "备份ID"
// @Success      200  {object}  response.Response  "删除成功"
// @Failure      400  {object}  response.Response  "无效的备份ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/backup/{id} [delete]
// @Security     BearerAuth
func (ctrl *SystemController) DeleteBackup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的备份ID")
		return
	}

	if err := system.DeleteBackup(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "删除成功", nil)
}

// GetBackupConfigs 获取备份配置
// @Summary      获取备份配置
// @Description  获取系统备份配置
// @Tags         备份管理
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/backup/config [get]
// @Security     BearerAuth
func (ctrl *SystemController) GetBackupConfigs(c *gin.Context) {
	configs, err := system.GetBackupConfigs()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, configs)
}

// UpdateBackupConfig 更新备份配置
// @Summary      更新备份配置
// @Description  更新系统备份配置
// @Tags         备份管理
// @Accept       json
// @Produce      json
// @Param        req  body      dto.BackupConfigRequest  true  "配置数据"
// @Success      200  {object}  response.Response  "保存成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/backup/config [put]
// @Security     BearerAuth
func (ctrl *SystemController) UpdateBackupConfig(c *gin.Context) {
	var req dto.BackupConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := system.CreateOrUpdateBackupConfig(0, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "保存成功", nil)
}
