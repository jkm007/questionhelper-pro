package model

import (
	"time"

	"gorm.io/gorm"
)

// DataPermission 数据权限规则表
type DataPermission struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Name      string         `gorm:"size:50;not null" json:"name"`            // 规则名称
	Code      string         `gorm:"uniqueIndex;size:50;not null" json:"code"` // 规则编码
	Type      int8           `gorm:"not null" json:"type"`                    // 类型:1=全部,2=本部门及以下,3=本部门,4=仅自己
	DeptID    *uint          `gorm:"index" json:"dept_id"`                    // 部门ID(类型2/3时需要)
	UserID    *uint          `gorm:"index" json:"user_id"`                    // 用户ID(类型4时需要)
	Status    int8           `gorm:"default:1" json:"status"`                 // 状态:0=禁用,1=启用
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (DataPermission) TableName() string {
	return "data_permissions"
}

// Dept 部门表(用于数据权限)
type Dept struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	ParentID  *uint          `gorm:"index" json:"parent_id"`
	Name      string         `gorm:"size:50;not null" json:"name"`
	Code      string         `gorm:"uniqueIndex;size:50" json:"code"`
	Leader    string         `gorm:"size:50" json:"leader"`                   // 负责人
	Phone     string         `gorm:"size:20" json:"phone"`
	Email     string         `gorm:"size:100" json:"email"`
	Sort      int            `gorm:"default:0" json:"sort"`
	Status    int8           `gorm:"default:1" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	Children []Dept `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

func (Dept) TableName() string {
	return "depts"
}
