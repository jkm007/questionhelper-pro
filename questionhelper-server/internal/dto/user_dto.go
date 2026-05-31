package dto

import "time"

// UserInfo 用户信息
type UserInfo struct {
	ID        uint       `json:"id"`
	Username  string     `json:"username"`
	Nickname  string     `json:"nickname"`
	Email     string     `json:"email"`
	Phone     string     `json:"phone"`
	Avatar    string     `json:"avatar"`
	Gender    int8       `json:"gender"`
	Birthday  *time.Time `json:"birthday"`
	Bio       string     `json:"bio"`
	Status    int8       `json:"status"`
	IsReal    bool       `json:"is_real"`
	Roles     []RoleInfo `json:"roles"`
	CreatedAt time.Time  `json:"created_at"`
}

// RoleInfo 角色信息
type RoleInfo struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

// UpdateProfileRequest 更新个人信息请求
type UpdateProfileRequest struct {
	Nickname string     `json:"nickname" binding:"omitempty,max=50"`
	Email    string     `json:"email" binding:"omitempty,email"`
	Gender   int8       `json:"gender" binding:"omitempty,oneof=0 1 2"`
	Birthday *time.Time `json:"birthday"`
	Bio      string     `json:"bio" binding:"omitempty,max=500"`
}

// RealNameAuthRequest 实名认证请求
type RealNameAuthRequest struct {
	RealName string `json:"real_name" binding:"required,max=50"`
	IDCard   string `json:"id_card" binding:"required"`
}

// UserListRequest 用户列表请求
type UserListRequest struct {
	PageRequest
	Keyword   string     `form:"keyword"`
	Status    *int8      `form:"status"`
	RoleID    *uint      `form:"role_id"`
	RoleCode  string     `form:"role_code"`  // 按角色编码筛选
	TagID     *uint      `form:"tag_id"`     // 按标签筛选
	StartDate *time.Time `form:"start_date"` // 注册开始时间
	EndDate   *time.Time `form:"end_date"`   // 注册结束时间
}

// BatchStatusRequest 批量更新状态请求
type BatchStatusRequest struct {
	IDs    []uint `json:"ids" binding:"required,min=1"`
	Status int8   `json:"status" binding:"required,oneof=0 1"`
}

// BatchDeleteRequest 批量删除请求
type BatchDeleteRequest struct {
	IDs []uint `json:"ids" binding:"required,min=1"`
}

// BatchRoleRequest 批量分配角色请求
type BatchRoleRequest struct {
	IDs     []uint `json:"ids" binding:"required,min=1"`
	RoleIDs []uint `json:"role_ids" binding:"required,min=1"`
}

// AdminResetPasswordRequest 管理员重置密码请求
type AdminResetPasswordRequest struct {
	Password string `json:"password" binding:"required,min=6,max=20"`
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=2,max=20"`
	Password string `json:"password" binding:"required,min=6,max=20"`
	Nickname string `json:"nickname" binding:"omitempty,max=50"`
	Phone    string `json:"phone" binding:"omitempty"`
	Email    string `json:"email" binding:"omitempty,email"`
	RoleIDs  []uint `json:"role_ids"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Nickname string `json:"nickname" binding:"omitempty,max=50"`
	Phone    string `json:"phone" binding:"omitempty"`
	Email    string `json:"email" binding:"omitempty,email"`
	Status   *int8  `json:"status"`
	RoleIDs  []uint `json:"role_ids"`
}

// RoleListRequest 角色列表请求
type RoleListRequest struct {
	PageRequest
	Keyword string `form:"keyword"`
}

// CreateRoleRequest 创建角色请求
type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required,max=50"`
	Code        string `json:"code" binding:"required,max=50"`
	Description string `json:"description" binding:"omitempty,max=200"`
	MenuIDs     []uint `json:"menu_ids"`
}

// UpdateRoleRequest 更新角色请求
type UpdateRoleRequest struct {
	Name        string `json:"name" binding:"omitempty,max=50"`
	Description string `json:"description" binding:"omitempty,max=200"`
	Status      *int8  `json:"status"`
	MenuIDs     []uint `json:"menu_ids"`
}
