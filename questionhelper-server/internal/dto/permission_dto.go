package dto

// ButtonPermissionInfo 按钮权限信息
type ButtonPermissionInfo struct {
	ID             uint   `json:"id"`
	MenuID         uint   `json:"menu_id"`
	PermissionCode string `json:"permission_code"`
	PermissionName string `json:"permission_name"`
	Sort           int    `json:"sort"`
	Status         bool   `json:"status"`
	Remark         string `json:"remark"`
}

// ButtonPermissionListRequest 按钮权限列表请求
type ButtonPermissionListRequest struct {
	PageRequest
	MenuID         *uint  `form:"menu_id" json:"menu_id"`
	PermissionCode string `form:"permission_code" json:"permission_code"`
	PermissionName string `form:"permission_name" json:"permission_name"`
	Status         *bool  `form:"status" json:"status"`
}

// CreateButtonPermissionRequest 创建按钮权限请求
type CreateButtonPermissionRequest struct {
	MenuID         uint   `json:"menu_id" binding:"required"`
	PermissionCode string `json:"permission_code" binding:"required,max=100"`
	PermissionName string `json:"permission_name" binding:"required,max=100"`
	Sort           int    `json:"sort"`
	Remark         string `json:"remark" binding:"omitempty,max=500"`
}

// UpdateButtonPermissionRequest 更新按钮权限请求
type UpdateButtonPermissionRequest struct {
	MenuID         *uint  `json:"menu_id"`
	PermissionCode string `json:"permission_code" binding:"omitempty,max=100"`
	PermissionName string `json:"permission_name" binding:"omitempty,max=100"`
	Sort           *int   `json:"sort"`
	Status         *bool  `json:"status"`
	Remark         string `json:"remark" binding:"omitempty,max=500"`
}
