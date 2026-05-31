package dto

import "time"

// ShareInfo 分享信息
type ShareInfo struct {
	ID         uint       `json:"id"`
	ShareCode  string     `json:"share_code"`
	ShareURL   string     `json:"share_url"`
	QuestionID uint       `json:"question_id"`
	ShareType  int8       `json:"share_type"`
	HasPassword bool      `json:"has_password"`
	ViewCount  int        `json:"view_count"`
	Status     int8       `json:"status"`
	ExpiresAt  *time.Time `json:"expires_at"`
	CreatedAt  time.Time  `json:"created_at"`
}

// CreateShareRequest 创建分享请求
type CreateShareRequest struct {
	QuestionID uint   `json:"question_id" binding:"required"`
	ShareType  int8   `json:"share_type" binding:"omitempty,oneof=1 2"`
	Password   string `json:"password" binding:"omitempty,max=20"`
	ExpiresIn  int    `json:"expires_in"` // 过期时间(小时),0=永不过期
}
