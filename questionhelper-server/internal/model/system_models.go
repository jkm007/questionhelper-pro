package model

import (
	"time"

	"gorm.io/gorm"
)

// SystemLog 系统日志表
type SystemLog struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Level     string         `gorm:"size:20;not null;index;comment:日志级别:debug/info/warn/error" json:"level"`
	Module    string         `gorm:"size:50;index;comment:模块" json:"module"`
	Action    string         `gorm:"size:50;comment:操作" json:"action"`
	Message   string         `gorm:"type:text;comment:日志消息" json:"message"`
	Details   string         `gorm:"type:text;comment:详细信息" json:"details"`
	UserID    *uint          `gorm:"index;comment:操作用户ID" json:"user_id"`
	IP        string         `gorm:"size:50;comment:IP地址" json:"ip"`
	UserAgent string         `gorm:"size:500;comment:用户代理" json:"user_agent"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (SystemLog) TableName() string {
	return "system_logs"
}

// ErrorLog 错误日志表
type ErrorLog struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Level       string         `gorm:"size:20;not null;comment:错误级别" json:"level"`
	ErrorMessage string        `gorm:"type:text;not null;comment:错误信息" json:"error_message"`
	StackTrace  string         `gorm:"type:text;comment:堆栈跟踪" json:"stack_trace"`
	URL         string         `gorm:"size:500;comment:请求URL" json:"url"`
	Method      string         `gorm:"size:10;comment:请求方法" json:"method"`
	UserID      *uint          `gorm:"index;comment:用户ID" json:"user_id"`
	IP          string         `gorm:"size:50;comment:IP地址" json:"ip"`
	UserAgent   string         `gorm:"size:500;comment:用户代理" json:"user_agent"`
	OccurredAt  time.Time      `gorm:"not null;index;comment:发生时间" json:"occurred_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ErrorLog) TableName() string {
	return "error_logs"
}

// FeatureFlag 功能开关表
type FeatureFlag struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Key         string         `gorm:"uniqueIndex;size:100;not null;comment:功能标识" json:"key"`
	Name        string         `gorm:"size:100;not null;comment:功能名称" json:"name"`
	Description string         `gorm:"size:500;comment:功能描述" json:"description"`
	IsEnabled   bool           `gorm:"default:false;comment:是否启用" json:"is_enabled"`
	Config      string         `gorm:"type:json;comment:配置数据" json:"config"`
	Environment string         `gorm:"size:50;default:production;comment:环境:development/staging/production" json:"environment"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (FeatureFlag) TableName() string {
	return "feature_flags"
}

// SecurityConfig 安全配置表
type SecurityConfig struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	ConfigKey   string         `gorm:"uniqueIndex;size:100;not null;comment:配置键" json:"config_key"`
	ConfigValue string         `gorm:"type:text;comment:配置值" json:"config_value"`
	Description string         `gorm:"size:500;comment:配置说明" json:"description"`
	IsActive    bool           `gorm:"default:true;comment:是否启用" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (SecurityConfig) TableName() string {
	return "security_configs"
}

// StorageConfig 存储配置表
type StorageConfig struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Name      string         `gorm:"size:100;not null;comment:配置名称" json:"name"`
	Type      string         `gorm:"size:20;not null;comment:存储类型:local/minio/oss/s3" json:"type"`
	Endpoint  string         `gorm:"size:500;comment:服务端点" json:"endpoint"`
	Bucket    string         `gorm:"size:100;comment:存储桶名称" json:"bucket"`
	AccessKey string         `gorm:"size:200;comment:访问密钥ID" json:"access_key"`
	SecretKey string         `gorm:"size:200;comment:访问密钥" json:"secret_key"`
	Region    string         `gorm:"size:50;comment:区域" json:"region"`
	BaseURL   string         `gorm:"size:500;comment:基础访问URL" json:"base_url"`
	IsDefault bool           `gorm:"default:false;comment:是否默认" json:"is_default"`
	IsActive  bool           `gorm:"default:true;comment:是否启用" json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (StorageConfig) TableName() string {
	return "storage_configs"
}

// EmailConfig 邮件配置表
type EmailConfig struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `gorm:"size:100;not null;comment:配置名称" json:"name"`
	SMTPHost    string         `gorm:"size:200;not null;comment:SMTP服务器地址" json:"smtp_host"`
	SMTPPort    int            `gorm:"not null;comment:SMTP端口" json:"smtp_port"`
	Username    string         `gorm:"size:200;comment:用户名" json:"username"`
	Password    string         `gorm:"size:200;comment:密码" json:"password"`
	FromAddress string         `gorm:"size:200;comment:发件人地址" json:"from_address"`
	FromName    string         `gorm:"size:100;comment:发件人名称" json:"from_name"`
	IsDefault   bool           `gorm:"default:false;comment:是否默认" json:"is_default"`
	IsActive    bool           `gorm:"default:true;comment:是否启用" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (EmailConfig) TableName() string {
	return "email_configs"
}

// EmailTemplate 邮件模板表
type EmailTemplate struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Code      string         `gorm:"uniqueIndex;size:100;not null;comment:模板编码" json:"code"`
	Name      string         `gorm:"size:100;not null;comment:模板名称" json:"name"`
	Subject   string         `gorm:"size:200;not null;comment:邮件主题" json:"subject"`
	Body      string         `gorm:"type:text;not null;comment:邮件内容" json:"body"`
	Variables string         `gorm:"type:json;comment:变量列表" json:"variables"`
	IsActive  bool           `gorm:"default:true;comment:是否启用" json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (EmailTemplate) TableName() string {
	return "email_templates"
}

// SMSConfig 短信配置表
type SMSConfig struct {
	ID             uint           `gorm:"primarykey" json:"id"`
	Provider       string         `gorm:"size:20;not null;comment:短信服务商:aliyun/tencent" json:"provider"`
	AccessKeyID    string         `gorm:"size:200;not null;comment:AccessKey ID" json:"access_key_id"`
	AccessKeySecret string        `gorm:"size:200;not null;comment:AccessKey Secret" json:"access_key_secret"`
	SignName       string         `gorm:"size:100;not null;comment:短信签名" json:"sign_name"`
	IsDefault      bool           `gorm:"default:false;comment:是否默认" json:"is_default"`
	IsActive       bool           `gorm:"default:true;comment:是否启用" json:"is_active"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (SMSConfig) TableName() string {
	return "sms_configs"
}

// SMSTemplate 短信模板表
type SMSTemplate struct {
	ID               uint           `gorm:"primarykey" json:"id"`
	Code             string         `gorm:"uniqueIndex;size:100;not null;comment:模板编码" json:"code"`
	Name             string         `gorm:"size:100;not null;comment:模板名称" json:"name"`
	Content          string         `gorm:"type:text;not null;comment:模板内容" json:"content"`
	Variables        string         `gorm:"type:json;comment:变量列表" json:"variables"`
	ProviderTemplateID string       `gorm:"size:100;comment:服务商模板ID" json:"provider_template_id"`
	IsActive         bool           `gorm:"default:true;comment:是否启用" json:"is_active"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

func (SMSTemplate) TableName() string {
	return "sms_templates"
}

// NotificationChannel 通知渠道配置表
type NotificationChannel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Name      string         `gorm:"size:100;not null;comment:渠道名称" json:"name"`
	Type      string         `gorm:"size:20;not null;comment:渠道类型:app/email/sms/wechat" json:"type"`
	Config    string         `gorm:"type:json;comment:渠道配置" json:"config"`
	IsEnabled bool           `gorm:"default:true;comment:是否启用" json:"is_enabled"`
	Priority  int            `gorm:"default:0;comment:优先级(数值越大优先级越高)" json:"priority"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (NotificationChannel) TableName() string {
	return "notification_channels"
}

// NotificationRateLimit 通知频率限制配置表
type NotificationRateLimit struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	UserID       uint           `gorm:"not null;index;comment:用户ID" json:"user_id"`
	Channel      string         `gorm:"size:20;not null;comment:通知渠道" json:"channel"`
	MaxCount     int            `gorm:"not null;comment:最大发送次数" json:"max_count"`
	Period       string         `gorm:"size:20;not null;comment:限制周期:hourly/daily" json:"period"`
	CurrentCount int            `gorm:"default:0;comment:当前发送次数" json:"current_count"`
	ResetAt      time.Time      `gorm:"comment:重置时间" json:"reset_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (NotificationRateLimit) TableName() string {
	return "notification_rate_limits"
}

// LogAlertRule 日志告警规则表
type LogAlertRule struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	Name       string         `gorm:"size:100;not null;comment:规则名称" json:"name"`
	Level      string         `gorm:"size:20;not null;comment:日志级别:debug/info/warn/error" json:"level"`
	Module     string         `gorm:"size:50;comment:模块" json:"module"`
	Pattern    string         `gorm:"type:text;comment:匹配模式(正则)" json:"pattern"`
	Threshold  int            `gorm:"not null;comment:触发阈值" json:"threshold"`
	Duration   int            `gorm:"not null;comment:统计时长(分钟)" json:"duration"`
	NotifyType string         `gorm:"size:50;not null;comment:通知方式:email/sms/webhook" json:"notify_type"`
	IsActive   bool           `gorm:"default:true;comment:是否启用" json:"is_active"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (LogAlertRule) TableName() string {
	return "log_alert_rules"
}

// LogAlertRecord 日志告警记录表
type LogAlertRecord struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	RuleID    uint           `gorm:"not null;index;comment:规则ID" json:"rule_id"`
	LogIDs    string         `gorm:"type:json;comment:触发日志ID列表" json:"log_ids"`
	Count     int            `gorm:"not null;comment:触发次数" json:"count"`
	Level     string         `gorm:"size:20;not null;comment:日志级别" json:"level"`
	Message   string         `gorm:"type:text;comment:告警消息" json:"message"`
	HandledAt *time.Time     `gorm:"comment:处理时间" json:"handled_at"`
	HandlerID *uint          `gorm:"index;comment:处理人ID" json:"handler_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (LogAlertRecord) TableName() string {
	return "log_alert_records"
}

// BackupConfig 备份配置表
type BackupConfig struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	Name         string         `gorm:"size:100;not null;comment:配置名称" json:"name"`
	Type         string         `gorm:"size:20;not null;comment:备份类型:full/incremental" json:"type"`
	Schedule     string         `gorm:"size:50;not null;comment:备份计划(cron表达式)" json:"schedule"`
	StoragePath  string         `gorm:"size:500;not null;comment:存储路径" json:"storage_path"`
	RetainDays   int            `gorm:"not null;default:30;comment:保留天数" json:"retain_days"`
	IsActive     bool           `gorm:"default:true;comment:是否启用" json:"is_active"`
	LastBackupAt *time.Time     `gorm:"comment:最近备份时间" json:"last_backup_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (BackupConfig) TableName() string {
	return "backup_configs"
}

// BackupRecord 备份记录表
type BackupRecord struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	ConfigID    uint           `gorm:"not null;index;comment:备份配置ID" json:"config_id"`
	Type        string         `gorm:"size:20;not null;comment:备份类型:full/incremental" json:"type"`
	FilePath    string         `gorm:"size:500;not null;comment:备份文件路径" json:"file_path"`
	FileSize    int64           `gorm:"comment:文件大小(字节)" json:"file_size"`
	Status      string         `gorm:"size:20;not null;default:pending;comment:状态:pending/running/completed/failed" json:"status"`
	StartedAt   *time.Time     `gorm:"comment:开始时间" json:"started_at"`
	CompletedAt *time.Time     `gorm:"comment:完成时间" json:"completed_at"`
	ErrorMessage string        `gorm:"type:text;comment:错误信息" json:"error_message"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (BackupRecord) TableName() string {
	return "backup_records"
}

// Language 语言表
type Language struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	Code       string         `gorm:"uniqueIndex;size:20;not null;comment:语言编码(如zh-CN)" json:"code"`
	Name       string         `gorm:"size:50;not null;comment:语言名称" json:"name"`
	NativeName string         `gorm:"size:50;not null;comment:本地名称" json:"native_name"`
	IsDefault  bool           `gorm:"default:false;comment:是否默认语言" json:"is_default"`
	IsActive   bool           `gorm:"default:true;comment:是否启用" json:"is_active"`
	Progress   float64        `gorm:"default:0;comment:翻译进度(百分比)" json:"progress"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Language) TableName() string {
	return "languages"
}

// ThemeConfig 主题配置表
type ThemeConfig struct {
	ID             uint           `gorm:"primarykey" json:"id"`
	Name           string         `gorm:"size:100;not null;comment:主题名称" json:"name"`
	PrimaryColor   string         `gorm:"size:20;comment:主色调" json:"primary_color"`
	SecondaryColor string         `gorm:"size:20;comment:副色调" json:"secondary_color"`
	LogoPath       string         `gorm:"size:500;comment:Logo路径" json:"logo_path"`
	FaviconPath    string         `gorm:"size:500;comment:Favicon路径" json:"favicon_path"`
	IsDefault      bool           `gorm:"default:false;comment:是否默认主题" json:"is_default"`
	IsActive       bool           `gorm:"default:true;comment:是否启用" json:"is_active"`
	Config         string         `gorm:"type:json;comment:主题配置数据" json:"config"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ThemeConfig) TableName() string {
	return "theme_configs"
}

// CacheConfig 缓存配置表
type CacheConfig struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	CacheKey      string         `gorm:"uniqueIndex;size:200;not null;comment:缓存键" json:"cache_key"`
	TTL           int            `gorm:"not null;default:3600;comment:过期时间(秒)" json:"ttl"`
	Description   string         `gorm:"size:500;comment:描述" json:"description"`
	IsEnabled     bool           `gorm:"default:true;comment:是否启用" json:"is_enabled"`
	LastClearedAt *time.Time     `gorm:"comment:最近清除时间" json:"last_cleared_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (CacheConfig) TableName() string {
	return "cache_configs"
}
