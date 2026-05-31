package dto

// PageRequest 分页请求
type PageRequest struct {
	Page     int `form:"page" json:"page"`
	PageSize int `form:"page_size" json:"page_size"`
}

// GetOffset 获取偏移量
func (p *PageRequest) GetOffset() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
	return (p.Page - 1) * p.PageSize
}

// GetLimit 获取每页数量
func (p *PageRequest) GetLimit() int {
	if p.PageSize <= 0 {
		return 10
	}
	if p.PageSize > 100 {
		return 100
	}
	return p.PageSize
}

// PageResponse 分页响应
type PageResponse struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

// IDRequest ID请求
type IDRequest struct {
	ID uint `uri:"id" json:"id" binding:"required"`
}

// IDsRequest 批量ID请求
type IDsRequest struct {
	IDs []uint `json:"ids" binding:"required"`
}

// StatusRequest 状态更新请求
type StatusRequest struct {
	Status int8 `json:"status" binding:"required"`
}
