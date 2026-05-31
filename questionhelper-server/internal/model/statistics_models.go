package model

import (
	"time"

	"gorm.io/gorm"
)

// UserRetention 用户留存表
type UserRetention struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	Date          time.Time      `gorm:"type:date;not null;uniqueIndex:uk_date_period" json:"date"`
	NewUsers      int            `gorm:"not null;default:0" json:"new_users"`
	RetainedUsers int            `gorm:"not null;default:0" json:"retained_users"`
	RetentionRate float64        `gorm:"type:decimal(5,2);not null;default:0" json:"retention_rate"`
	Period        string         `gorm:"size:10;not null;uniqueIndex:uk_date_period" json:"period"` // 统计周期:day/week/month
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (UserRetention) TableName() string {
	return "user_retention"
}

// UserChurn 用户流失表
type UserChurn struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	Date         time.Time      `gorm:"type:date;not null;uniqueIndex:uk_date_period" json:"date"`
	ChurnedUsers int            `gorm:"not null;default:0" json:"churned_users"`
	ChurnRate    float64        `gorm:"type:decimal(5,2);not null;default:0" json:"churn_rate"`
	ChurnReasons string         `gorm:"type:text" json:"churn_reasons"` // 流失原因(JSON)
	Period       string         `gorm:"size:10;not null;uniqueIndex:uk_date_period" json:"period"` // 统计周期:day/week/month
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (UserChurn) TableName() string {
	return "user_churn"
}

// UserEvent 用户行为事件表
type UserEvent struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	UserID     uint      `gorm:"index;not null" json:"user_id"`
	EventType  string    `gorm:"size:50;not null;index" json:"event_type"` // 事件类型:click/view/submit/search/share
	EventName  string    `gorm:"size:100;not null" json:"event_name"`
	Page       string    `gorm:"size:200;index" json:"page"`
	Element    string    `gorm:"size:200" json:"element"`
	ExtraData  string    `gorm:"type:text" json:"extra_data"` // 扩展数据(JSON)
	SessionID  string    `gorm:"size:64;index" json:"session_id"`
	DeviceType string    `gorm:"size:20" json:"device_type"` // 设备类型:web/ios/android/miniapp
	IP         string    `gorm:"size:50" json:"ip"`
	CreatedAt  time.Time `json:"created_at"`
}

func (UserEvent) TableName() string {
	return "user_events"
}

// UserSegment 用户分群表
type UserSegment struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Description string         `gorm:"size:500" json:"description"`
	Rules       string         `gorm:"type:text;not null" json:"rules"` // 分群规则(JSON)
	UserCount   int            `gorm:"default:0" json:"user_count"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatorID   uint           `gorm:"index;not null" json:"creator_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (UserSegment) TableName() string {
	return "user_segments"
}

// UserSegmentMember 用户分群成员表
type UserSegmentMember struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	SegmentID uint      `gorm:"uniqueIndex:uk_segment_user;not null" json:"segment_id"`
	UserID    uint      `gorm:"uniqueIndex:uk_segment_user;not null" json:"user_id"`
	JoinedAt  time.Time `json:"joined_at"`
}

func (UserSegmentMember) TableName() string {
	return "user_segment_members"
}

// UserPageView 用户访问路径表
type UserPageView struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	UserID     uint      `gorm:"index;not null" json:"user_id"`
	Page       string    `gorm:"size:200;not null;index" json:"page"`
	Referrer   string    `gorm:"size:200" json:"referrer"`
	Duration   int       `gorm:"default:0;comment:停留时长(秒)" json:"duration"`
	SessionID  string    `gorm:"size:64;index" json:"session_id"`
	DeviceType string    `gorm:"size:20" json:"device_type"` // 设备类型:web/ios/android/miniapp
	IP         string    `gorm:"size:50" json:"ip"`
	CreatedAt  time.Time `json:"created_at"`
}

func (UserPageView) TableName() string {
	return "user_page_views"
}

// ConversionFunnel 转化漏斗表
type ConversionFunnel struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Description string         `gorm:"size:500" json:"description"`
	Steps       string         `gorm:"type:text;not null" json:"steps"` // 漏斗步骤(JSON)
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatorID   uint           `gorm:"index;not null" json:"creator_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ConversionFunnel) TableName() string {
	return "conversion_funnels"
}

// ConversionFunnelStat 转化漏斗统计表
type ConversionFunnelStat struct {
	ID             uint      `gorm:"primarykey" json:"id"`
	FunnelID       uint      `gorm:"uniqueIndex:uk_funnel_date_step;not null" json:"funnel_id"`
	Date           time.Time `gorm:"type:date;not null;uniqueIndex:uk_funnel_date_step" json:"date"`
	StepIndex      int       `gorm:"not null;uniqueIndex:uk_funnel_date_step" json:"step_index"`
	StepName       string    `gorm:"size:100;not null" json:"step_name"`
	UserCount      int       `gorm:"not null;default:0" json:"user_count"`
	ConversionRate float64   `gorm:"type:decimal(5,2);not null;default:0" json:"conversion_rate"`
	CreatedAt      time.Time `json:"created_at"`
}

func (ConversionFunnelStat) TableName() string {
	return "conversion_funnel_stats"
}

// AlertRule 数据预警规则表
type AlertRule struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	MetricName  string         `gorm:"size:100;not null;index" json:"metric_name"`
	Condition   string         `gorm:"size:10;not null" json:"condition"` // 条件:gt/lt/eq/gte/lte
	Threshold   float64        `gorm:"not null" json:"threshold"`
	Duration    int            `gorm:"default:0;comment:持续时间(分钟)" json:"duration"`
	NotifyType  string         `gorm:"size:20;not null;default:system" json:"notify_type"` // 通知方式:system/email/sms
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatorID   uint           `gorm:"index;not null" json:"creator_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (AlertRule) TableName() string {
	return "alert_rules"
}

// AlertRecord 预警记录表
type AlertRecord struct {
	ID          uint       `gorm:"primarykey" json:"id"`
	RuleID      uint       `gorm:"index;not null" json:"rule_id"`
	MetricValue float64    `gorm:"not null" json:"metric_value"`
	Threshold   float64    `gorm:"not null" json:"threshold"`
	Level       string     `gorm:"size:20;not null;index" json:"level"` // 预警级别:info/warning/critical
	Message     string     `gorm:"size:500" json:"message"`
	HandledAt   *time.Time `json:"handled_at"`
	HandlerID   *uint      `json:"handler_id"`
	CreatedAt   time.Time  `json:"created_at"`
}

func (AlertRecord) TableName() string {
	return "alert_records"
}

// DataSubscription 数据订阅表
type DataSubscription struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	UserID     uint           `gorm:"index;not null" json:"user_id"`
	ReportType string         `gorm:"size:50;not null" json:"report_type"` // 报表类型
	Frequency  string         `gorm:"size:10;not null" json:"frequency"`   // 订阅频率:daily/weekly/monthly
	Channels   string         `gorm:"type:text;not null" json:"channels"`  // 通知渠道(JSON)
	IsActive   bool           `gorm:"default:true" json:"is_active"`
	LastSentAt *time.Time     `json:"last_sent_at"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (DataSubscription) TableName() string {
	return "data_subscriptions"
}

// ExportRecord 数据导出记录表
type ExportRecord struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	UserID     uint      `gorm:"index;not null" json:"user_id"`
	ExportType string    `gorm:"size:50;not null" json:"export_type"` // 导出类型
	FileFormat string    `gorm:"size:10;not null" json:"file_format"` // 文件格式:xlsx/csv/pdf
	FilePath   string    `gorm:"size:500" json:"file_path"`
	FileSize   int64     `gorm:"default:0" json:"file_size"`           // 文件大小(字节)
	Status     string    `gorm:"size:20;not null;default:pending;index" json:"status"` // 状态:pending/processing/completed/failed
	Filters    string    `gorm:"type:text" json:"filters"`             // 筛选条件(JSON)
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (ExportRecord) TableName() string {
	return "export_records"
}

// ComparisonSnapshot 数据对比快照表
type ComparisonSnapshot struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	SnapshotType string        `gorm:"size:50;not null;index" json:"snapshot_type"` // 快照类型
	PeriodStart time.Time      `gorm:"type:date;not null" json:"period_start"`
	PeriodEnd   time.Time      `gorm:"type:date;not null" json:"period_end"`
	Data        string         `gorm:"type:text;not null" json:"data"` // 快照数据(JSON)
	CreatorID   uint           `gorm:"index;not null" json:"creator_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ComparisonSnapshot) TableName() string {
	return "comparison_snapshots"
}

// ScoreAlertConfig 成绩预警配置表
type ScoreAlertConfig struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	UserID     uint           `gorm:"uniqueIndex:uk_user_exam;not null" json:"user_id"`
	ExamID     uint           `gorm:"uniqueIndex:uk_user_exam;index;not null" json:"exam_id"`
	MinScore   float64        `gorm:"type:decimal(5,2);not null" json:"min_score"`
	NotifyType string         `gorm:"size:20;not null;default:system" json:"notify_type"` // 通知方式:system/email/sms
	IsActive   bool           `gorm:"default:true" json:"is_active"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ScoreAlertConfig) TableName() string {
	return "score_alert_configs"
}

// DataShare 数据分享表
type DataShare struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	UserID      uint           `gorm:"index;not null" json:"user_id"`
	ShareType   string         `gorm:"size:50;not null" json:"share_type"` // 分享类型
	TargetID    uint           `gorm:"index;not null" json:"target_id"`    // 目标对象ID
	Permissions string         `gorm:"type:text" json:"permissions"`       // 权限配置(JSON)
	ExpireAt    *time.Time     `gorm:"index" json:"expire_at"`
	ShareCode   string         `gorm:"size:64;uniqueIndex" json:"share_code"` // 分享码
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (DataShare) TableName() string {
	return "data_shares"
}

// StatisticsSnapshot 统计快照表
type StatisticsSnapshot struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	SnapshotType string    `gorm:"size:50;not null;uniqueIndex:uk_type_date_period" json:"snapshot_type"` // 快照类型
	Date         time.Time `gorm:"type:date;not null;uniqueIndex:uk_type_date_period" json:"date"`
	Data         string    `gorm:"type:text;not null" json:"data"` // 快照数据(JSON)
	Period       string    `gorm:"size:10;not null;uniqueIndex:uk_type_date_period" json:"period"` // 统计周期:daily/weekly/monthly
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (StatisticsSnapshot) TableName() string {
	return "statistics_snapshots"
}
