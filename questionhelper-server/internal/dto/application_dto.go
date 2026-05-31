package dto

import "time"

// CreateApplicationRequest 创建角色申请请求
type CreateApplicationRequest struct {
	RoleID uint   `json:"role_id" binding:"required"`
	Reason string `json:"reason" binding:"required,max=500"`
}

// ApplicationListRequest 角色申请列表请求
type ApplicationListRequest struct {
	PageRequest
	UserID *uint  `form:"user_id"`
	RoleID *uint  `form:"role_id"`
	Status *int8  `form:"status"`
}

// ApplicationInfo 角色申请信息
type ApplicationInfo struct {
	ID         uint       `json:"id"`
	UserID     uint       `json:"user_id"`
	Username   string     `json:"username"`
	Nickname   string     `json:"nickname"`
	RoleID     uint       `json:"role_id"`
	RoleName   string     `json:"role_name"`
	Reason     string     `json:"reason"`
	Status     int8       `json:"status"`
	ReviewNote string     `json:"review_note"`
	ReviewedBy *uint      `json:"reviewed_by"`
	Reviewer   string     `json:"reviewer"`
	ReviewedAt *time.Time `json:"reviewed_at"`
	CreatedAt  time.Time  `json:"created_at"`
}

// ReviewApplicationRequest 审核角色申请请求
type ReviewApplicationRequest struct {
	Status int8   `json:"status" binding:"required,oneof=1 2"` // 1=通过,2=驳回
	Note   string `json:"note" binding:"omitempty,max=500"`
}
