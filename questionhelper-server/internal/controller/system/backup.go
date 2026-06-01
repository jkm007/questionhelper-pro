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
func (ctrl *SystemController) GetBackupConfigs(c *gin.Context) {
	configs, err := system.GetBackupConfigs()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, configs)
}

// UpdateBackupConfig 更新备份配置
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
