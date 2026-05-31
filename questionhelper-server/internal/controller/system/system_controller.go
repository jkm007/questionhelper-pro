package system

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/service/system"
	"questionhelper-server/pkg/response"
)

type SystemController struct{}

func NewSystemController() *SystemController {
	return &SystemController{}
}

// ==================== 系统设置(原有) ====================

// GetSettings 获取系统设置
func (ctrl *SystemController) GetSettings(c *gin.Context) {
	settings, err := system.GetSettings()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, settings)
}

// UpdateSettings 更新系统设置
func (ctrl *SystemController) UpdateSettings(c *gin.Context) {
	var settings map[string]string
	if err := c.ShouldBindJSON(&settings); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := system.UpdateSettings(settings); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "更新成功", nil)
}

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

// ==================== 分类设置 ====================

// GetClassSettings 获取班级设置
func (ctrl *SystemController) GetClassSettings(c *gin.Context) {
	settings, err := system.GetClassSettings()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, settings)
}

// UpdateClassSettings 更新班级设置
func (ctrl *SystemController) UpdateClassSettings(c *gin.Context) {
	var req dto.ClassSettings
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	if err := system.UpdateClassSettings(&req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "更新成功", nil)
}

// GetResourceSettings 获取资源设置
func (ctrl *SystemController) GetResourceSettings(c *gin.Context) {
	settings, err := system.GetResourceSettings()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, settings)
}

// UpdateResourceSettings 更新资源设置
func (ctrl *SystemController) UpdateResourceSettings(c *gin.Context) {
	var req dto.ResourceSettings
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	if err := system.UpdateResourceSettings(&req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "更新成功", nil)
}

// ==================== 系统日志 ====================

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

// ==================== 通知渠道 ====================

// ListNotificationChannels 渠道列表
func (ctrl *SystemController) ListNotificationChannels(c *gin.Context) {
	channels, err := system.ListNotificationChannels()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, channels)
}

// UpdateNotificationChannel 更新渠道
func (ctrl *SystemController) UpdateNotificationChannel(c *gin.Context) {
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

	if err := system.UpdateNotificationChannel(uint(id), &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "更新成功", nil)
}

// ==================== 数据备份 ====================

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

// ==================== 功能开关 ====================

// ListFeatureFlags 功能开关列表
func (ctrl *SystemController) ListFeatureFlags(c *gin.Context) {
	flags, err := system.ListFeatureFlags()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, flags)
}

// UpdateFeatureFlag 更新功能开关
func (ctrl *SystemController) UpdateFeatureFlag(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		response.Error(c, http.StatusBadRequest, "功能标识不能为空")
		return
	}

	var req dto.UpdateFeatureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := system.UpdateFeatureFlag(key, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "更新成功", nil)
}

// ==================== 安全配置 ====================

// GetSecurityConfigs 获取安全配置
func (ctrl *SystemController) GetSecurityConfigs(c *gin.Context) {
	configs, err := system.GetSecurityConfigs()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, configs)
}

// UpdateSecurityConfigs 更新安全配置
func (ctrl *SystemController) UpdateSecurityConfigs(c *gin.Context) {
	var req dto.SecurityConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := system.UpdateSecurityConfigs(&req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "更新成功", nil)
}

// ==================== 存储配置 ====================

// ListStorageConfigs 存储配置列表
func (ctrl *SystemController) ListStorageConfigs(c *gin.Context) {
	configs, err := system.ListStorageConfigs()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, configs)
}

// CreateStorageConfig 创建存储配置
func (ctrl *SystemController) CreateStorageConfig(c *gin.Context) {
	var req dto.StorageConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := system.CreateStorageConfig(&req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "创建成功", nil)
}

// UpdateStorageConfig 更新存储配置
func (ctrl *SystemController) UpdateStorageConfig(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的配置ID")
		return
	}

	var req dto.StorageConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := system.UpdateStorageConfig(uint(id), &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "更新成功", nil)
}

// ==================== 邮件配置 ====================

// GetEmailConfig 获取邮件配置
func (ctrl *SystemController) GetEmailConfig(c *gin.Context) {
	config, err := system.GetEmailConfig()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, config)
}

// UpdateEmailConfig 更新邮件配置
func (ctrl *SystemController) UpdateEmailConfig(c *gin.Context) {
	var req dto.EmailConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := system.UpdateEmailConfig(&req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "更新成功", nil)
}

// ListEmailTemplates 邮件模板列表
func (ctrl *SystemController) ListEmailTemplates(c *gin.Context) {
	templates, err := system.ListEmailTemplates()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, templates)
}

// CreateEmailTemplate 创建邮件模板
func (ctrl *SystemController) CreateEmailTemplate(c *gin.Context) {
	var req dto.EmailTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := system.CreateEmailTemplate(&req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "创建成功", nil)
}

// ==================== 短信配置 ====================

// GetSMSConfig 获取短信配置
func (ctrl *SystemController) GetSMSConfig(c *gin.Context) {
	config, err := system.GetSMSConfig()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, config)
}

// UpdateSMSConfig 更新短信配置
func (ctrl *SystemController) UpdateSMSConfig(c *gin.Context) {
	var req dto.SMSConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := system.UpdateSMSConfig(&req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "更新成功", nil)
}

// ListSMSTemplates 短信模板列表
func (ctrl *SystemController) ListSMSTemplates(c *gin.Context) {
	templates, err := system.ListSMSTemplates()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, templates)
}

// CreateSMSTemplate 创建短信模板
func (ctrl *SystemController) CreateSMSTemplate(c *gin.Context) {
	var req dto.SMSTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := system.CreateSMSTemplate(&req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "创建成功", nil)
}

// ==================== 缓存管理 ====================

// GetCacheStats 缓存统计
func (ctrl *SystemController) GetCacheStats(c *gin.Context) {
	stats, err := system.GetCacheStats()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, stats)
}

// ClearCache 清除缓存
func (ctrl *SystemController) ClearCache(c *gin.Context) {
	var req dto.ClearCacheRequest
	// 支持空body
	c.ShouldBindJSON(&req)

	if err := system.ClearCache(&req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "缓存清除成功", nil)
}

// ==================== 主题配置 ====================

// GetThemeConfig 获取主题配置
func (ctrl *SystemController) GetThemeConfig(c *gin.Context) {
	config, err := system.GetThemeConfig()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, config)
}

// UpdateThemeConfig 更新主题配置
func (ctrl *SystemController) UpdateThemeConfig(c *gin.Context) {
	var req dto.ThemeConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := system.UpdateThemeConfig(&req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "更新成功", nil)
}

// ==================== 告警管理 ====================

// ListAlertRules 告警规则列表
func (ctrl *SystemController) ListAlertRules(c *gin.Context) {
	var query dto.AlertRuleQuery
	c.ShouldBindQuery(&query)
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 || query.PageSize > 100 {
		query.PageSize = 10
	}

	rules, total, err := system.ListAlertRules(&query)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, rules, total, query.Page, query.PageSize)
}

// CreateAlertRule 创建告警规则
func (ctrl *SystemController) CreateAlertRule(c *gin.Context) {
	var req dto.AlertRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := system.CreateAlertRule(&req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "创建成功", nil)
}

// UpdateAlertRule 更新告警规则
func (ctrl *SystemController) UpdateAlertRule(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的规则ID")
		return
	}

	var req dto.AlertRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := system.UpdateAlertRule(uint(id), &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "更新成功", nil)
}

// ListAlertRecords 告警记录列表
func (ctrl *SystemController) ListAlertRecords(c *gin.Context) {
	var query dto.AlertRecordQuery
	c.ShouldBindQuery(&query)
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 || query.PageSize > 100 {
		query.PageSize = 10
	}

	records, total, err := system.ListAlertRecords(&query)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, records, total, query.Page, query.PageSize)
}
