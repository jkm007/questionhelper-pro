package model

import (
	"time"
)

// ExamWarning 考试异常行为记录表
type ExamWarning struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	RecordID   uint      `gorm:"index;not null" json:"record_id"`
	ExamID     uint      `gorm:"index;not null" json:"exam_id"`
	UserID     uint      `gorm:"index;not null" json:"user_id"`
	Type       string    `gorm:"size:30;not null" json:"type"`       // switch_screen/ip_change/copy_paste/tab_change/dev_tools
	Detail     string    `gorm:"size:500" json:"detail"`
	IP         string    `gorm:"size:50" json:"ip"`
	UserAgent  string    `gorm:"size:500" json:"user_agent"`
	CreatedAt  time.Time `json:"created_at"`
}

func (ExamWarning) TableName() string {
	return "exam_warnings"
}
