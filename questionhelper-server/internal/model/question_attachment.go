package model

import (
	"time"
)

// QuestionAttachment 题目附件表
type QuestionAttachment struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	QuestionID  uint      `gorm:"index;not null" json:"question_id"`
	FileID      uint      `gorm:"index" json:"file_id"`                  // 关联文件表
	FileName    string    `gorm:"size:200;not null" json:"file_name"`
	FileType    string    `gorm:"size:20;not null" json:"file_type"`     // image/audio/video/document
	FileURL     string    `gorm:"size:500;not null" json:"file_url"`
	FileSize    int64     `gorm:"default:0" json:"file_size"`
	Sort        int       `gorm:"default:0" json:"sort"`
	CreatedAt   time.Time `json:"created_at"`
}

func (QuestionAttachment) TableName() string {
	return "question_attachments"
}
