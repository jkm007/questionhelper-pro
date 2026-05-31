package model

import (
	"time"
)

// QuestionReview 题目审核表
type QuestionReview struct {
	ID          uint       `gorm:"primarykey" json:"id"`
	QuestionID  uint       `gorm:"index;not null" json:"question_id"`
	ReviewerID  *uint      `gorm:"index" json:"reviewer_id"`                    // 审核人
	Status      int8       `gorm:"default:0" json:"status"`                     // 0=待审核,1=通过,2=驳回,3=需修改
	Opinion     string     `gorm:"size:500" json:"opinion"`                     // 审核意见
	Reply       string     `gorm:"size:500" json:"reply"`                       // 审核意见回复
	BatchID     string     `gorm:"size:50;index" json:"batch_id"`               // 批量审核批次号
	SubmittedBy uint       `gorm:"index" json:"submitted_by"`                   // 提交人
	SubmittedAt time.Time  `json:"submitted_at"`
	ReviewedAt  *time.Time `json:"reviewed_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (QuestionReview) TableName() string {
	return "question_reviews"
}
