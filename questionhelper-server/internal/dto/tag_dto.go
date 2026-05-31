package dto

// CreateTagRequest 创建标签请求
type CreateTagRequest struct {
	Name  string `json:"name" binding:"required,max=50"`
	Code  string `json:"code" binding:"required,max=50"`
	Color string `json:"color" binding:"omitempty,max=20"`
	Icon  string `json:"icon" binding:"omitempty,max=100"`
	Type  int8   `json:"type" binding:"omitempty,oneof=1 2"`
	Sort  int    `json:"sort"`
}

// UpdateTagRequest 更新标签请求
type UpdateTagRequest struct {
	Name   string `json:"name" binding:"omitempty,max=50"`
	Color  string `json:"color" binding:"omitempty,max=20"`
	Icon   string `json:"icon" binding:"omitempty,max=100"`
	Sort   int    `json:"sort"`
	Status *int8  `json:"status"`
}

// TagListRequest 标签列表请求
type TagListRequest struct {
	PageRequest
	Keyword string `form:"keyword"`
	Type    *int8  `form:"type"`
	Status  *int8  `form:"status"`
}

// TagInfo 标签信息
type TagInfo struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Code      string `json:"code"`
	Color     string `json:"color"`
	Icon      string `json:"icon"`
	Type      int8   `json:"type"`
	Sort      int    `json:"sort"`
	Status    int8   `json:"status"`
	UserCount int    `json:"user_count"`
}

// UserTagRequest 用户标签操作请求
type UserTagRequest struct {
	TagIDs []uint `json:"tag_ids" binding:"required,min=1"`
}
