package dto

import "time"

// FolderInfo 收藏夹信息
type FolderInfo struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Icon        string    `json:"icon"`
	Sort        int       `json:"sort"`
	Count       int       `json:"count"`
	IsDefault   bool      `json:"is_default"`
	CreatedAt   time.Time `json:"created_at"`
}

// CreateFolderRequest 创建收藏夹请求
type CreateFolderRequest struct {
	Name        string `json:"name" binding:"required,max=100"`
	Description string `json:"description" binding:"omitempty,max=200"`
	Icon        string `json:"icon" binding:"omitempty,max=50"`
}

// UpdateFolderRequest 更新收藏夹请求
type UpdateFolderRequest struct {
	Name        string `json:"name" binding:"omitempty,max=100"`
	Description string `json:"description" binding:"omitempty,max=200"`
	Icon        string `json:"icon" binding:"omitempty,max=50"`
	Sort        int    `json:"sort"`
}

// FavoriteInfo 收藏信息
type FavoriteInfo struct {
	ID         uint         `json:"id"`
	QuestionID uint         `json:"question_id"`
	Question   *QuestionInfo `json:"question,omitempty"`
	FolderID   uint         `json:"folder_id"`
	FolderName string       `json:"folder_name"`
	Note       string       `json:"note"`
	CreatedAt  time.Time    `json:"created_at"`
}

// FavoriteRequest 收藏请求
type FavoriteRequest struct {
	FolderID uint   `json:"folder_id"`
	Note     string `json:"note" binding:"omitempty,max=500"`
}

// FavoriteListRequest 收藏列表请求
type FavoriteListRequest struct {
	PageRequest
	FolderID *uint `form:"folder_id"`
}
