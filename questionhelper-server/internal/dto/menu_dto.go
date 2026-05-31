package dto

// MenuInfo 菜单信息
type MenuInfo struct {
	ID        uint       `json:"id"`
	ParentID  *uint      `json:"parent_id"`
	Name      string     `json:"name"`
	Path      string     `json:"path"`
	Component string     `json:"component"`
	Title     string     `json:"title"`
	Icon      string     `json:"icon"`
	Hidden    bool       `json:"hidden"`
	Type      int8       `json:"type"`
	Permission string    `json:"permission"`
	Sort      int        `json:"sort"`
	Status    int8       `json:"status"`
	Children  []MenuInfo `json:"children,omitempty"`
}

// CreateMenuRequest 创建菜单请求
type CreateMenuRequest struct {
	ParentID   *uint  `json:"parent_id"`
	Name       string `json:"name" binding:"required,max=50"`
	Path       string `json:"path" binding:"omitempty,max=200"`
	Component  string `json:"component" binding:"omitempty,max=200"`
	Redirect   string `json:"redirect" binding:"omitempty,max=200"`
	Title      string `json:"title" binding:"required,max=50"`
	Icon       string `json:"icon" binding:"omitempty,max=50"`
	Hidden     bool   `json:"hidden"`
	Type       int8   `json:"type" binding:"required,oneof=1 2 3"`
	Permission string `json:"permission" binding:"omitempty,max=100"`
	Sort       int    `json:"sort"`
}

// UpdateMenuRequest 更新菜单请求
type UpdateMenuRequest struct {
	ParentID   *uint  `json:"parent_id"`
	Name       string `json:"name" binding:"omitempty,max=50"`
	Path       string `json:"path" binding:"omitempty,max=200"`
	Component  string `json:"component" binding:"omitempty,max=200"`
	Redirect   string `json:"redirect" binding:"omitempty,max=200"`
	Title      string `json:"title" binding:"omitempty,max=50"`
	Icon       string `json:"icon" binding:"omitempty,max=50"`
	Hidden     *bool  `json:"hidden"`
	Type       int8   `json:"type" binding:"omitempty,oneof=1 2 3"`
	Permission string `json:"permission" binding:"omitempty,max=100"`
	Sort       *int   `json:"sort"`
	Status     *int8  `json:"status"`
}
