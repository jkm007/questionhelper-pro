package model

import (
	"time"

	"gorm.io/gorm"
)

// NotificationTemplate 通知模板表
type NotificationTemplate struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Code      string    `gorm:"uniqueIndex;size:50;not null" json:"code"`  // 模板编码（唯一标识）
	Name      string    `gorm:"size:100;not null" json:"name"`            // 模板名称
	TitleTpl  string    `gorm:"size:200;not null" json:"title_tpl"`       // 标题模板，支持 {{.Var}} 变量
	ContentTpl string   `gorm:"type:text;not null" json:"content_tpl"`    // 内容模板，支持 {{.Var}} 变量
	Type      int8      `gorm:"not null;comment:类型:1=系统,2=考试,3=作业,4=班级,5=评论" json:"type"`
	Channel   string    `gorm:"size:20;default:app;comment:通知渠道:app/email/sms/wechat" json:"channel"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`            // 是否启用
	Remark    string    `gorm:"size:200" json:"remark"`                   // 备注
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (NotificationTemplate) TableName() string {
	return "notification_templates"
}

// NotificationSetting 通知设置表（用户通知偏好）
type NotificationSetting struct {
	ID              uint      `gorm:"primarykey" json:"id"`
	UserID          uint      `gorm:"uniqueIndex:idx_user_type_channel;not null" json:"user_id"`
	Type            int8      `gorm:"uniqueIndex:idx_user_type_channel;not null;comment:类型:1=系统,2=考试,3=作业,4=班级,5=评论" json:"type"`
	Channel         string    `gorm:"uniqueIndex:idx_user_type_channel;size:20;default:app;comment:通知渠道:app/email/sms/wechat" json:"channel"`
	Enabled         bool      `gorm:"default:true" json:"enabled"`         // 是否开启该类型通知
	DoNotDisturb    bool      `gorm:"default:false" json:"do_not_disturb"` // 免打扰
	DisturbStart    string    `gorm:"size:5" json:"disturb_start"`         // 免打扰开始时间 HH:mm
	DisturbEnd      string    `gorm:"size:5" json:"disturb_end"`           // 免打扰结束时间 HH:mm
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (NotificationSetting) TableName() string {
	return "notification_settings"
}

// ScheduledNotification 定时通知表
type ScheduledNotification struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	SenderID   uint           `gorm:"index;not null" json:"sender_id"`
	Title      string         `gorm:"size:200;not null" json:"title"`
	Content    string         `gorm:"type:text;not null" json:"content"`
	Type       int8           `gorm:"not null;comment:类型:1=系统,2=考试,3=作业,4=班级,5=评论" json:"type"`
	Channel    string         `gorm:"size:20;default:app;comment:通知渠道:app/email/sms/wechat" json:"channel"`
	TargetType string         `gorm:"size:50;comment:目标类型:all/role/class/group" json:"target_type"`
	TargetIDs  string         `gorm:"type:text;comment:目标ID列表(JSON数组)" json:"target_ids"`
	ScheduledAt time.Time     `gorm:"not null;index;comment:计划发送时间" json:"scheduled_at"`
	Status     int8           `gorm:"default:0;comment:状态:0=待发送,1=已发送,2=已取消,3=发送失败" json:"status"`
	BatchID    string         `gorm:"size:36;index" json:"batch_id"`
	ErrorMsg   string         `gorm:"type:text" json:"error_msg"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ScheduledNotification) TableName() string {
	return "scheduled_notifications"
}
