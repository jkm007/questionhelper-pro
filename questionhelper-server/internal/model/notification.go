package model

import (
	"time"

	"gorm.io/gorm"
)

type Notification struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	UserID     uint           `gorm:"index;not null" json:"user_id"`
	Type       int8           `gorm:"not null;comment:类型:1=系统,2=考试,3=作业,4=班级,5=评论" json:"type"`
	Title      string         `gorm:"size:200;not null" json:"title"`
	Content    string         `gorm:"type:text" json:"content"`
	TargetType string         `gorm:"size:50" json:"target_type"`
	TargetID   uint           `gorm:"index" json:"target_id"`
	IsRead     bool           `gorm:"default:false" json:"is_read"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Notification) TableName() string {
	return "notifications"
}
