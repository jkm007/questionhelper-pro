package model

import (
	"time"
)

// QuestionVersion 题目版本历史表
type QuestionVersion struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	QuestionID  uint      `gorm:"index;not null" json:"question_id"`
	Version     int       `gorm:"not null" json:"version"`
	Title       string    `gorm:"type:text;not null" json:"title"`
	Content     string    `gorm:"type:text" json:"content"`
	Type        int8      `gorm:"not null" json:"type"`
	Difficulty  int8      `gorm:"default:2" json:"difficulty"`
	Answer      string    `gorm:"type:text" json:"answer"`
	Analysis    string    `gorm:"type:text" json:"analysis"`
	CategoryID  uint      `gorm:"index" json:"category_id"`
	Options     string    `gorm:"type:text" json:"options"` // JSON格式的选项
	ChangeLog   string    `gorm:"size:500" json:"change_log"` // 修改说明
	CreatorID   uint      `gorm:"index" json:"creator_id"`
	CreatedAt   time.Time `json:"created_at"`
}

func (QuestionVersion) TableName() string {
	return "question_versions"
}
