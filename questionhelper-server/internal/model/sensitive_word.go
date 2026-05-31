package model

import (
	"time"

	"gorm.io/gorm"
)

// SensitiveWord 敏感词表
type SensitiveWord struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Word      string         `gorm:"uniqueIndex;size:100;not null" json:"word"`
	Category  string         `gorm:"size:50;index" json:"category"` // 分类:政治/色情/暴力/广告等
	Level     int8           `gorm:"default:1" json:"level"`        // 级别:1=替换,2=禁止
	ReplaceTo string         `gorm:"size:100" json:"replace_to"`    // 替换文本
	Status    int8           `gorm:"default:1" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (SensitiveWord) TableName() string {
	return "sensitive_words"
}
