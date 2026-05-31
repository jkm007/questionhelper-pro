package model

import (
	"time"

	"gorm.io/gorm"
)

type PracticeSession struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	UserID     uint           `gorm:"index;not null" json:"user_id"`
	CategoryID *uint          `gorm:"index" json:"category_id"`
	TotalCount int            `gorm:"default:0" json:"total_count"`
	CorrectCount int          `gorm:"default:0" json:"correct_count"`
	Accuracy   float64        `gorm:"default:0" json:"accuracy"`
	Duration   int            `gorm:"default:0;comment:练习时长(秒)" json:"duration"`
	Status     int8           `gorm:"default:0;comment:状态:0=进行中,1=已完成" json:"status"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (PracticeSession) TableName() string {
	return "practice_sessions"
}

type PracticeRecord struct {
	ID         uint   `gorm:"primarykey" json:"id"`
	SessionID  uint   `gorm:"index;not null" json:"session_id"`
	QuestionID uint   `gorm:"index;not null" json:"question_id"`
	Answer     string `gorm:"type:text" json:"answer"`
	IsCorrect  bool   `json:"is_correct"`
	Duration   int    `gorm:"default:0;comment:用时(秒)" json:"duration"`
}

func (PracticeRecord) TableName() string {
	return "practice_records"
}
