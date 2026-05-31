package dto

import "time"

// CommentInfo 评论信息
type CommentInfo struct {
	ID             uint         `json:"id"`
	TargetType     int8         `json:"target_type"`
	TargetID       uint         `json:"target_id"`
	UserID         uint         `json:"user_id"`
	UserName       string       `json:"user_name"`
	UserAvatar     string       `json:"user_avatar"`
	Content        string       `json:"content"`
	ParentID       *uint        `json:"parent_id"`
	Images         string       `json:"images"`
	Mentions       string       `json:"mentions"`
	ReplyToUserID  *uint        `json:"reply_to_user_id"`
	ReplyToUserName string      `json:"reply_to_user_name,omitempty"`
	LikeCount      int          `json:"like_count"`
	IsPinned       bool         `json:"is_pinned"`
	IsFeatured     bool         `json:"is_featured"`
	IsOfficial     bool         `json:"is_official"`
	Status         int8         `json:"status"`
	CreatedAt      time.Time    `json:"created_at"`
	Children       []CommentInfo `json:"children,omitempty"`
}

// CreateCommentRequest 创建评论请求
type CreateCommentRequest struct {
	TargetType    int8   `json:"target_type" binding:"required,oneof=1 2 3"`
	TargetID      uint   `json:"target_id" binding:"required"`
	Content       string `json:"content" binding:"required,max=500"`
	ParentID      *uint  `json:"parent_id"`
	Images        string `json:"images"`
	Mentions      string `json:"mentions"`
	ReplyToUserID *uint  `json:"reply_to_user_id"`
}

// UpdateCommentRequest 编辑评论请求
type UpdateCommentRequest struct {
	Content string `json:"content" binding:"required,max=500"`
	Images  string `json:"images"`
}

// CommentListRequest 评论列表请求
type CommentListRequest struct {
	PageRequest
	TargetType int8   `form:"target_type" binding:"required,oneof=1 2 3"`
	TargetID   uint   `form:"target_id" binding:"required"`
	Status     *int8  `form:"status"`
	Keyword    string `form:"keyword"`
}

// CommentAdminListRequest 管理员评论列表请求
type CommentAdminListRequest struct {
	PageRequest
	TargetType *int8  `form:"target_type"`
	Status     *int8  `form:"status"`
	Keyword    string `form:"keyword"`
	UserID     *uint  `form:"user_id"`
	StartDate  *time.Time `form:"start_date"`
	EndDate    *time.Time `form:"end_date"`
}

// ReportCommentRequest 举报评论请求
type ReportCommentRequest struct {
	Reason     string `json:"reason" binding:"required,max=200"`
	ReasonType string `json:"reason_type" binding:"omitempty,oneof=spam abuse infringement other"`
}

// UploadImageResponse 上传图片响应
type UploadImageResponse struct {
	URL  string `json:"url"`
	Name string `json:"name"`
	Size int64  `json:"size"`
}

// StickerInfo 表情包信息
type StickerInfo struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Category  string `json:"category"`
	ImageURL  string `json:"image_url"`
	SortOrder int    `json:"sort_order"`
}

// StickerCategoryInfo 表情分类信息
type StickerCategoryInfo struct {
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

// UserSearchRequest 用户搜索请求(@功能)
type UserSearchRequest struct {
	Keyword  string `form:"keyword" binding:"required"`
	PageSize int    `form:"page_size"`
}

// UserSearchInfo 用户搜索结果
type UserSearchInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

// ==================== Blacklist DTOs ====================

// BlacklistInfo 黑名单信息
type BlacklistInfo struct {
	ID         uint       `json:"id"`
	UserID     uint       `json:"user_id"`
	UserName   string     `json:"user_name"`
	UserAvatar string     `json:"user_avatar"`
	TargetType int8       `json:"target_type"`
	TargetID   uint       `json:"target_id"`
	Reason     string     `json:"reason"`
	OperatorID uint       `json:"operator_id"`
	ExpiresAt  *time.Time `json:"expires_at"`
	CreatedAt  time.Time  `json:"created_at"`
}

// BlacklistListRequest 黑名单列表请求
type BlacklistListRequest struct {
	PageRequest
	UserID     *uint `form:"user_id"`
	TargetType *int8 `form:"target_type"`
}

// AddBlacklistRequest 添加黑名单请求
type AddBlacklistRequest struct {
	UserID     uint       `json:"user_id" binding:"required"`
	TargetType int8       `json:"target_type" binding:"required,oneof=0 1 2 3"`
	TargetID   uint       `json:"target_id"`
	Reason     string     `json:"reason" binding:"omitempty,max=200"`
	Duration   int        `json:"duration"` // 拉黑时长（小时），0表示永久
}

// ==================== Audit Rule DTOs ====================

// AuditRuleInfo 审核规则信息
type AuditRuleInfo struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	RuleType    string `json:"rule_type"`
	Pattern     string `json:"pattern"`
	Action      int8   `json:"action"`
	Priority    int    `json:"priority"`
	Description string `json:"description"`
	Status      int8   `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// CreateAuditRuleRequest 创建审核规则请求
type CreateAuditRuleRequest struct {
	Name        string `json:"name" binding:"required,max=100"`
	RuleType    string `json:"rule_type" binding:"required,oneof=keyword regex length repeat"`
	Pattern     string `json:"pattern" binding:"required"`
	Action      int8   `json:"action" binding:"required,oneof=1 2 3 4"`
	Priority    int    `json:"priority"`
	Description string `json:"description" binding:"omitempty,max=200"`
}

// UpdateAuditRuleRequest 更新审核规则请求
type UpdateAuditRuleRequest struct {
	Name        string `json:"name" binding:"omitempty,max=100"`
	RuleType    string `json:"rule_type" binding:"omitempty,oneof=keyword regex length repeat"`
	Pattern     string `json:"pattern"`
	Action      int8   `json:"action" binding:"omitempty,oneof=1 2 3 4"`
	Priority    int    `json:"priority"`
	Description string `json:"description" binding:"omitempty,max=200"`
	Status      *int8  `json:"status" binding:"omitempty,oneof=0 1"`
}

// ==================== Report DTOs ====================

// ReportInfo 举报信息
type ReportInfo struct {
	ID           uint       `json:"id"`
	CommentID    uint       `json:"comment_id"`
	CommentContent string   `json:"comment_content,omitempty"`
	UserID       uint       `json:"user_id"`
	UserName     string     `json:"user_name,omitempty"`
	Reason       string     `json:"reason"`
	ReasonType   string     `json:"reason_type"`
	Status       int8       `json:"status"`
	HandlerID    *uint      `json:"handler_id"`
	HandleRemark string     `json:"handle_remark"`
	HandledAt    *time.Time `json:"handled_at"`
	CreatedAt    time.Time  `json:"created_at"`
}

// ReportListRequest 举报列表请求
type ReportListRequest struct {
	PageRequest
	Status     *int8  `form:"status"`
	ReasonType string `form:"reason_type"`
}

// HandleReportRequest 处理举报请求
type HandleReportRequest struct {
	Status       int8   `json:"status" binding:"required,oneof=1 2"`
	HandleRemark string `json:"handle_remark" binding:"omitempty,max=500"`
}

// ==================== Batch Operation DTOs ====================

// BatchAuditRequest 批量审核请求
type BatchAuditRequest struct {
	IDs    []uint `json:"ids" binding:"required,min=1"`
	Status int8   `json:"status" binding:"required,oneof=0 1"`
}

// BatchDeleteRequest 批量删除评论请求
type BatchCommentDeleteRequest struct {
	IDs []uint `json:"ids" binding:"required,min=1"`
}

// ==================== Stats & Export DTOs ====================

// CommentStats 评论统计
type CommentStats struct {
	TotalCount     int64            `json:"total_count"`
	TodayCount     int64            `json:"today_count"`
	PendingCount   int64            `json:"pending_count"`
	HiddenCount    int64            `json:"hidden_count"`
	ReportCount    int64            `json:"report_count"`
	DailyStats     []DailyCommentStat `json:"daily_stats,omitempty"`
	TopTargetStats []TargetCommentStat `json:"top_target_stats,omitempty"`
}

// DailyCommentStat 每日评论统计
type DailyCommentStat struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

// TargetCommentStat 目标评论统计
type TargetCommentStat struct {
	TargetType int8   `json:"target_type"`
	TargetID   uint   `json:"target_id"`
	Count      int64  `json:"count"`
}

// CommentExportRequest 导出评论请求
type CommentExportRequest struct {
	TargetType *int8      `form:"target_type"`
	Status     *int8      `form:"status"`
	StartDate  *time.Time `form:"start_date"`
	EndDate    *time.Time `form:"end_date"`
	Format     string     `form:"format"` // csv, xlsx
}
