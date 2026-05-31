package model

import (
	"time"

	"gorm.io/gorm"
)

// Tag 用户标签表
type Tag struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Name      string         `gorm:"uniqueIndex;size:50;not null" json:"name"`
	Code      string         `gorm:"uniqueIndex;size:50;not null" json:"code"`
	Color     string         `gorm:"size:20" json:"color"`           // 标签颜色
	Icon      string         `gorm:"size:100" json:"icon"`           // 标签图标
	Type      int8           `gorm:"default:1" json:"type"`          // 类型:1=系统标签,2=自定义标签
	Sort      int            `gorm:"default:0" json:"sort"`          // 排序
	Status    int8           `gorm:"default:1" json:"status"`        // 状态:0=禁用,1=启用
	UserCount int            `gorm:"default:0" json:"user_count"`    // 关联用户数
	CreatedBy uint           `json:"created_by"`                     // 创建者ID
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	Users []User `gorm:"many2many:user_tags;" json:"users,omitempty"`
}

func (Tag) TableName() string {
	return "tags"
}

// UserTag 用户标签关联表
type UserTag struct {
	UserID    uint      `gorm:"primaryKey" json:"user_id"`
	TagID     uint      `gorm:"primaryKey" json:"tag_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (UserTag) TableName() string {
	return "user_tags"
}
