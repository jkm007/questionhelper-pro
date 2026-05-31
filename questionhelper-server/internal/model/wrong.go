package model

import (
	"time"

	"gorm.io/gorm"
)

type WrongQuestion struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	UserID        uint           `gorm:"index;not null" json:"user_id"`
	QuestionID    uint           `gorm:"index;not null" json:"question_id"`
	Question      Question       `json:"question,omitempty"`
	Source        int8           `gorm:"default:1;comment:来源:1=练习,2=考试" json:"source"`
	SourceID      uint           `gorm:"index" json:"source_id"`
	WrongCount    int            `gorm:"default:1" json:"wrong_count"`
	LastAnswer    string         `gorm:"type:text" json:"last_answer"`
	Mastered      bool           `gorm:"default:false" json:"mastered"`
	CorrectStreak int            `gorm:"default:0" json:"correct_streak"`              // 连续正确次数
	LastReviewAt  *time.Time     `json:"last_review_at"`                               // 上次复习时间
	NextReviewAt  *time.Time     `gorm:"index" json:"next_review_at"`                  // 下次复习时间
	ReviewCount   int            `gorm:"default:0" json:"review_count"`                // 复习次数
	IsFavorite    bool           `gorm:"default:false" json:"is_favorite"`             // 是否收藏
	IsExported    bool           `gorm:"default:false" json:"is_exported"`             // 是否已导出
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (WrongQuestion) TableName() string {
	return "wrong_questions"
}
