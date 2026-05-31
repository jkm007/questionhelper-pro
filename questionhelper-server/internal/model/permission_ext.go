package model

import "time"

// ButtonPermission 按钮权限表
type ButtonPermission struct {
	ID             uint      `gorm:"primarykey" json:"id"`
	MenuID         uint      `gorm:"index;not null;comment:所属菜单ID" json:"menu_id"`
	PermissionCode string    `gorm:"uniqueIndex;size:100;not null;comment:权限编码" json:"permission_code"`
	PermissionName string    `gorm:"size:100;not null;comment:权限名称" json:"permission_name"`
	Sort           int       `gorm:"default:0;comment:排序" json:"sort"`
	Status         bool      `gorm:"default:true;comment:状态" json:"status"`
	Remark         string    `gorm:"size:500;comment:备注" json:"remark"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (ButtonPermission) TableName() string {
	return "button_permissions"
}
