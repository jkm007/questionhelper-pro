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
