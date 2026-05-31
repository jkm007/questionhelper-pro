package model

import (
	"time"

	"gorm.io/gorm"
)

// RoleApplication 角色申请表
type RoleApplication struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	UserID      uint           `gorm:"index;not null" json:"user_id"`          // 申请人ID
	RoleID      uint           `gorm:"index;not null" json:"role_id"`          // 申请角色ID
	Reason      string         `gorm:"size:500" json:"reason"`                 // 申请理由
	Status      int8           `gorm:"default:0" json:"status"`                // 状态:0=待审核,1=通过,2=驳回
	ReviewNote  string         `gorm:"size:500" json:"review_note"`            // 审核备注
	ReviewedBy  *uint          `json:"reviewed_by"`                            // 审核人ID
	ReviewedAt  *time.Time     `json:"reviewed_at"`                            // 审核时间
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	User     User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Role     Role `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	Reviewer User `gorm:"foreignKey:ReviewedBy" json:"reviewer,omitempty"`
}

func (RoleApplication) TableName() string {
	return "role_applications"
}
