package model

import (
	"time"

	"gorm.io/gorm"
)

type WrongQuestion struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	UserID     uint           `gorm:"index;not null" json:"user_id"`
	QuestionID uint           `gorm:"index;not null" json:"question_id"`
	Question   Question       `json:"question,omitempty"`
	Source     int8           `gorm:"default:1;comment:来源:1=练习,2=考试" json:"source"`
	SourceID   uint           `gorm:"index" json:"source_id"`
	WrongCount int            `gorm:"default:1" json:"wrong_count"`
	LastAnswer string         `gorm:"type:text" json:"last_answer"`
	Mastered   bool           `gorm:"default:false" json:"mastered"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (WrongQuestion) TableName() string {
	return "wrong_questions"
}
