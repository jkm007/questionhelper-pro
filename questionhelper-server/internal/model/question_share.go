package model

import (
	"time"
)

// QuestionShare 题目分享表
type QuestionShare struct {
	ID          uint       `gorm:"primarykey" json:"id"`
	ShareCode   string     `gorm:"uniqueIndex;size:32;not null" json:"share_code"` // 分享码
	QuestionID  uint       `gorm:"index;not null" json:"question_id"`
	UserID      uint       `gorm:"index;not null" json:"user_id"`                 // 分享者
	ShareType   int8       `gorm:"default:1" json:"share_type"`                   // 1=链接,2=二维码
	Password    string     `gorm:"size:20" json:"-"`                               // 访问密码
	ViewCount   int        `gorm:"default:0" json:"view_count"`
	Status      int8       `gorm:"default:1" json:"status"`                       // 0=已撤销,1=有效
	ExpiresAt   *time.Time `json:"expires_at"`                                    // 过期时间
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (QuestionShare) TableName() string {
	return "question_shares"
}
