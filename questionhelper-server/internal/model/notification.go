package model

import (
	"time"

	"gorm.io/gorm"
)

type Notification struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	UserID     uint           `gorm:"index;not null" json:"user_id"`
	SenderID   *uint          `gorm:"index" json:"sender_id"`                                                                         // 发送者ID（系统通知可为空）
	Type       int8           `gorm:"not null;comment:类型:1=系统,2=考试,3=作业,4=班级,5=评论" json:"type"`
	Title      string         `gorm:"size:200;not null" json:"title"`
	Content    string         `gorm:"type:text" json:"content"`
	TargetType string         `gorm:"size:50" json:"target_type"`
	TargetID   uint           `gorm:"index" json:"target_id"`
	IsRead     bool           `gorm:"default:false" json:"is_read"`
	Extra      string         `gorm:"type:text" json:"extra"`                                                                          // 扩展数据(JSON)
	IsRecalled bool           `gorm:"default:false" json:"is_recalled"`                                                                // 是否已撤回
	RecalledAt *time.Time     `json:"recalled_at"`                                                                                     // 撤回时间
	Channel    string         `gorm:"size:20;default:app;comment:通知渠道:app/email/sms/wechat" json:"channel"`                          // 通知渠道
	BatchID    string         `gorm:"size:36;index" json:"batch_id"`                                                                   // 批次ID（群发通知用）
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Notification) TableName() string {
	return "notifications"
}
