package dto

import "time"

// ==================== 分类设置 ====================

// ClassSettings 班级设置
type ClassSettings struct {
	MaxStudentsPerClass int    `json:"max_students_per_class"` // 每班最大学生数
	AllowSelfJoin       bool   `json:"allow_self_join"`        // 允许自主加入
	RequireApproval     bool   `json:"require_approval"`       // 需要审批
	DefaultClassType    string `json:"default_class_type"`     // 默认班级类型
	EnableClassCode     bool   `json:"enable_class_code"`      // 启用班级码
	ClassCodeLength     int    `json:"class_code_length"`      // 班级码长度
}

// ResourceSettings 资源设置
type ResourceSettings struct {
	MaxUploadSize    int64    `json:"max_upload_size"`    // 最大上传大小(MB)
	AllowedFileTypes []string `json:"allowed_file_types"` // 允许的文件类型
	EnablePreview    bool     `json:"enable_preview"`     // 启用预览
	StorageQuota     int64    `json:"storage_quota"`      // 存储配额(GB)
	EnableCompress   bool     `json:"enable_compress"`    // 启用压缩
	WatermarkEnabled bool     `json:"watermark_enabled"`  // 启用水印
}

// ==================== 系统日志 ====================

// SystemLogQuery 系统日志查询
type SystemLogQuery struct {
	PageRequest
	Level    string `form:"level"`     // 日志级别
	Module   string `form:"module"`    // 模块
	Keyword  string `form:"keyword"`   // 关键词
	StartAt  string `form:"start_at"`  // 开始时间
	EndAt    string `form:"end_at"`    // 结束时间
}

// ErrorLogQuery 错误日志查询
type ErrorLogQuery struct {
	PageRequest
	Level    string `form:"level"`
	Keyword  string `form:"keyword"`
	StartAt  string `form:"start_at"`
	EndAt    string `form:"end_at"`
}

// LogSearchRequest 日志搜索请求
type LogSearchRequest struct {
	PageRequest
	Keyword   string `form:"keyword" json:"keyword"`
	LogType   string `form:"log_type" json:"log_type"` // system/error/operation
	Level     string `form:"level" json:"level"`
	Module    string `form:"module" json:"module"`
	StartAt   string `form:"start_at" json:"start_at"`
	EndAt     string `form:"end_at" json:"end_at"`
}

// LogArchiveRequest 日志归档请求
type LogArchiveRequest struct {
	LogType  string `json:"log_type" binding:"required"` // system/error
	BeforeAt string `json:"before_at" binding:"required"` // 归档此时间之前的日志
}

// LogStatsResponse 日志统计响应
type LogStatsResponse struct {
	TotalCount    int64            `json:"total_count"`
	TodayCount    int64            `json:"today_count"`
	ErrorCount    int64            `json:"error_count"`
	WarnCount     int64            `json:"warn_count"`
	LevelStats    []LevelStatItem  `json:"level_stats"`
	ModuleStats   []ModuleStatItem `json:"module_stats"`
	TrendData     []TrendItem      `json:"trend_data"`
}

// LevelStatItem 级别统计项
type LevelStatItem struct {
	Level string `json:"level"`
	Count int64  `json:"count"`
}

// ModuleStatItem 模块统计项
type ModuleStatItem struct {
	Module string `json:"module"`
	Count  int64  `json:"count"`
}

// TrendItem 趋势数据项
type TrendItem struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

// ==================== 通知渠道 ====================

// UpdateChannelRequest 更新渠道请求
type UpdateChannelRequest struct {
	Name      string `json:"name"`
	Config    string `json:"config"`
	IsEnabled *bool  `json:"is_enabled"`
	Priority  *int   `json:"priority"`
}

// ==================== 数据备份 ====================

// CreateBackupRequest 创建备份请求
type CreateBackupRequest struct {
	Type       string `json:"type" binding:"required"` // full/incremental
	ConfigID   uint   `json:"config_id"`
}

// BackupListQuery 备份列表查询
type BackupListQuery struct {
	PageRequest
	Status string `form:"status"`
	Type   string `form:"type"`
}

// BackupConfigRequest 备份配置请求
type BackupConfigRequest struct {
	Name        string `json:"name" binding:"required"`
	Type        string `json:"type" binding:"required"`
	Schedule    string `json:"schedule" binding:"required"`
	StoragePath string `json:"storage_path" binding:"required"`
	RetainDays  int    `json:"retain_days"`
	IsActive    *bool  `json:"is_active"`
}

// ==================== 功能开关 ====================

// UpdateFeatureRequest 更新功能开关请求
type UpdateFeatureRequest struct {
	IsEnabled *bool  `json:"is_enabled"`
	Config    string `json:"config"`
}

// ==================== 安全配置 ====================

// SecurityConfigRequest 安全配置请求
type SecurityConfigRequest struct {
	Configs []SecurityConfigItem `json:"configs" binding:"required"`
}

// SecurityConfigItem 安全配置项
type SecurityConfigItem struct {
	ConfigKey   string `json:"config_key" binding:"required"`
	ConfigValue string `json:"config_value"`
	Description string `json:"description"`
	IsActive    *bool  `json:"is_active"`
}

// ==================== 存储配置 ====================

// StorageConfigRequest 存储配置请求
type StorageConfigRequest struct {
	Name      string `json:"name" binding:"required"`
	Type      string `json:"type" binding:"required"`
	Endpoint  string `json:"endpoint"`
	Bucket    string `json:"bucket"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Region    string `json:"region"`
	BaseURL   string `json:"base_url"`
	IsDefault *bool  `json:"is_default"`
	IsActive  *bool  `json:"is_active"`
}

// ==================== 邮件配置 ====================

// EmailConfigRequest 邮件配置请求
type EmailConfigRequest struct {
	Name        string `json:"name" binding:"required"`
	SMTPHost    string `json:"smtp_host" binding:"required"`
	SMTPPort    int    `json:"smtp_port" binding:"required"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	FromAddress string `json:"from_address"`
	FromName    string `json:"from_name"`
	IsDefault   *bool  `json:"is_default"`
	IsActive    *bool  `json:"is_active"`
}

// EmailTemplateRequest 邮件模板请求
type EmailTemplateRequest struct {
	Code      string `json:"code" binding:"required"`
	Name      string `json:"name" binding:"required"`
	Subject   string `json:"subject" binding:"required"`
	Body      string `json:"body" binding:"required"`
	Variables string `json:"variables"`
	IsActive  *bool  `json:"is_active"`
}

// ==================== 短信配置 ====================

// SMSConfigRequest 短信配置请求
type SMSConfigRequest struct {
	Provider        string `json:"provider" binding:"required"`
	AccessKeyID     string `json:"access_key_id" binding:"required"`
	AccessKeySecret string `json:"access_key_secret" binding:"required"`
	SignName        string `json:"sign_name" binding:"required"`
	IsDefault       *bool  `json:"is_default"`
	IsActive        *bool  `json:"is_active"`
}

// SMSTemplateRequest 短信模板请求
type SMSTemplateRequest struct {
	Code               string `json:"code" binding:"required"`
	Name               string `json:"name" binding:"required"`
	Content            string `json:"content" binding:"required"`
	Variables          string `json:"variables"`
	ProviderTemplateID string `json:"provider_template_id"`
	IsActive           *bool  `json:"is_active"`
}

// ==================== 缓存管理 ====================

// CacheStatsResponse 缓存统计响应
type CacheStatsResponse struct {
	TotalKeys      int64              `json:"total_keys"`
	UsedMemory     string             `json:"used_memory"`
	UsedMemoryPeak string             `json:"used_memory_peak"`
	HitRate        float64            `json:"hit_rate"`
	CacheDetails   []CacheDetailItem  `json:"cache_details"`
}

// CacheDetailItem 缓存详情项
type CacheDetailItem struct {
	CacheKey    string     `json:"cache_key"`
	TTL         int        `json:"ttl"`
	Description string     `json:"description"`
	IsEnabled   bool       `json:"is_enabled"`
	LastCleared *time.Time `json:"last_cleared"`
}

// ClearCacheRequest 清除缓存请求
type ClearCacheRequest struct {
	Pattern string `json:"pattern"` // 为空则清除全部，支持通配符
}

// ==================== 主题配置 ====================

// ThemeConfigRequest 主题配置请求
type ThemeConfigRequest struct {
	Name           string `json:"name" binding:"required"`
	PrimaryColor   string `json:"primary_color"`
	SecondaryColor string `json:"secondary_color"`
	LogoPath       string `json:"logo_path"`
	FaviconPath    string `json:"favicon_path"`
	IsDefault      *bool  `json:"is_default"`
	IsActive       *bool  `json:"is_active"`
	Config         string `json:"config"`
}

// ==================== 告警管理 ====================

// AlertRuleRequest 告警规则请求
type AlertRuleRequest struct {
	Name       string `json:"name" binding:"required"`
	Level      string `json:"level" binding:"required"`
	Module     string `json:"module"`
	Pattern    string `json:"pattern"`
	Threshold  int    `json:"threshold" binding:"required"`
	Duration   int    `json:"duration" binding:"required"`
	NotifyType string `json:"notify_type" binding:"required"`
	IsActive   *bool  `json:"is_active"`
}

// AlertRuleQuery 告警规则查询
type AlertRuleQuery struct {
	PageRequest
	Level    string `form:"level"`
	IsActive *bool  `form:"is_active"`
}

// AlertRecordQuery 告警记录查询
type AlertRecordQuery struct {
	PageRequest
	RuleID  uint   `form:"rule_id"`
	Level   string `form:"level"`
	StartAt string `form:"start_at"`
	EndAt   string `form:"end_at"`
}
