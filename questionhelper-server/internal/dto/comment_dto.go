package dto

import "time"

// CommentInfo 评论信息
type CommentInfo struct {
	ID         uint      `json:"id"`
	TargetType int8      `json:"target_type"`
	TargetID   uint      `json:"target_id"`
	UserID     uint      `json:"user_id"`
	UserName   string    `json:"user_name"`
	UserAvatar string    `json:"user_avatar"`
	Content    string    `json:"content"`
	ParentID   *uint     `json:"parent_id"`
	LikeCount  int       `json:"like_count"`
	Status     int8      `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	Children   []CommentInfo `json:"children,omitempty"`
}

// CreateCommentRequest 创建评论请求
type CreateCommentRequest struct {
	TargetType int8   `json:"target_type" binding:"required,oneof=1 2 3"`
	TargetID   uint   `json:"target_id" binding:"required"`
	Content    string `json:"content" binding:"required,max=500"`
	ParentID   *uint  `json:"parent_id"`
}

// CommentListRequest 评论列表请求
type CommentListRequest struct {
	PageRequest
	TargetType int8 `form:"target_type" binding:"required,oneof=1 2 3"`
	TargetID   uint `form:"target_id" binding:"required"`
}

// ReportCommentRequest 举报评论请求
type ReportCommentRequest struct {
	Reason string `json:"reason" binding:"required,max=200"`
}
