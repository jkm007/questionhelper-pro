package model

import (
	"time"

	"gorm.io/gorm"
)

// UserRealName 实名认证表
type UserRealName struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	UserID     uint           `gorm:"uniqueIndex;not null" json:"user_id"`    // 用户ID
	RealName   string         `gorm:"size:50;not null" json:"real_name"`      // 真实姓名
	IDCard     string         `gorm:"size:200;not null" json:"-"`             // 身份证号(AES加密)
	IDCardHash string         `gorm:"size:64;index" json:"-"`                 // 身份证号哈希(用于查询)
	Status     int8           `gorm:"default:0" json:"status"`                // 状态:0=待审核,1=已认证,2=认证失败
	RejectReason string       `gorm:"size:200" json:"reject_reason"`          // 驳回原因
	ReviewedBy *uint          `json:"reviewed_by"`                            // 审核人ID
	ReviewedAt *time.Time     `json:"reviewed_at"`                            // 审核时间
	SubmittedAt time.Time     `json:"submitted_at"`                           // 提交时间
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	User     User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Reviewer User `gorm:"foreignKey:ReviewedBy" json:"reviewer,omitempty"`
}

func (UserRealName) TableName() string {
	return "user_real_names"
}
