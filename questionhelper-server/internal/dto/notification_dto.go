package dto

import "time"

// NotificationInfo 通知信息
type NotificationInfo struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"user_id"`
	Type       int8      `json:"type"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	TargetType string    `json:"target_type"`
	TargetID   uint      `json:"target_id"`
	IsRead     bool      `json:"is_read"`
	CreatedAt  time.Time `json:"created_at"`
}

// NotificationListRequest 通知列表请求
type NotificationListRequest struct {
	PageRequest
	Type   *int8  `form:"type"`
	IsRead *bool  `form:"is_read"`
}

// UnreadCountResponse 未读数量响应
type UnreadCountResponse struct {
	Count int `json:"count"`
}
