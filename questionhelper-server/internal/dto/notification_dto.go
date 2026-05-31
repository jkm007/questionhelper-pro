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
	IsRecalled bool      `json:"is_recalled"`
	Channel    string    `json:"channel"`
	BatchID    string    `json:"batch_id"`
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

// ==================== 通知模板 ====================

// TemplateInfo 通知模板信息
type TemplateInfo struct {
	ID         uint      `json:"id"`
	Code       string    `json:"code"`
	Name       string    `json:"name"`
	TitleTpl   string    `json:"title_tpl"`
	ContentTpl string    `json:"content_tpl"`
	Type       int8      `json:"type"`
	Channel    string    `json:"channel"`
	IsActive   bool      `json:"is_active"`
	Remark     string    `json:"remark"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// CreateTemplateRequest 创建模板请求
type CreateTemplateRequest struct {
	Code       string `json:"code" binding:"required,max=50"`
	Name       string `json:"name" binding:"required,max=100"`
	TitleTpl   string `json:"title_tpl" binding:"required,max=200"`
	ContentTpl string `json:"content_tpl" binding:"required"`
	Type       int8   `json:"type" binding:"required,oneof=1 2 3 4 5"`
	Channel    string `json:"channel" binding:"required,oneof=app email sms wechat"`
	IsActive   *bool  `json:"is_active"`
	Remark     string `json:"remark" binding:"max=200"`
}

// UpdateTemplateRequest 更新模板请求
type UpdateTemplateRequest struct {
	Name       string `json:"name" binding:"max=100"`
	TitleTpl   string `json:"title_tpl" binding:"max=200"`
	ContentTpl string `json:"content_tpl"`
	Type       int8   `json:"type" binding:"omitempty,oneof=1 2 3 4 5"`
	Channel    string `json:"channel" binding:"omitempty,oneof=app email sms wechat"`
	IsActive   *bool  `json:"is_active"`
	Remark     string `json:"remark" binding:"max=200"`
}

// NotificationTemplateListRequest 通知模板列表请求
type NotificationTemplateListRequest struct {
	PageRequest
	Keyword  string `form:"keyword"`
	Type     *int8  `form:"type"`
	Channel  string `form:"channel"`
	IsActive *bool  `form:"is_active"`
}

// ==================== 群发通知 ====================

// BatchSendRequest 群发通知请求
type BatchSendRequest struct {
	Title      string   `json:"title" binding:"required,max=200"`
	Content    string   `json:"content" binding:"required"`
	Type       int8     `json:"type" binding:"required,oneof=1 2 3 4 5"`
	Channel    string   `json:"channel" binding:"required,oneof=app email sms wechat"`
	TargetType string   `json:"target_type" binding:"required,oneof=all role class group"`
	TargetIDs  []uint   `json:"target_ids"` // 角色ID/班级ID/用户组ID列表，target_type=all时可为空
	UserIDs    []uint   `json:"user_ids"`   // 直接指定用户ID列表
}

// BatchSendResponse 群发通知响应
type BatchSendResponse struct {
	BatchID   string `json:"batch_id"`
	Total     int    `json:"total"`
	Success   int    `json:"success"`
	Failed    int    `json:"failed"`
}

// ==================== 撤回通知 ====================

// RecallNotificationRequest 撤回通知请求（无参数，仅路径ID）
type RecallNotificationRequest struct{}

// ==================== 定时通知 ====================

// ScheduledNotificationInfo 定时通知信息
type ScheduledNotificationInfo struct {
	ID          uint      `json:"id"`
	SenderID    uint      `json:"sender_id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Type        int8      `json:"type"`
	Channel     string    `json:"channel"`
	TargetType  string    `json:"target_type"`
	TargetIDs   string    `json:"target_ids"`
	ScheduledAt time.Time `json:"scheduled_at"`
	Status      int8      `json:"status"`
	BatchID     string    `json:"batch_id"`
	ErrorMsg    string    `json:"error_msg"`
	CreatedAt   time.Time `json:"created_at"`
}

// CreateScheduledRequest 创建定时通知请求
type CreateScheduledRequest struct {
	Title       string    `json:"title" binding:"required,max=200"`
	Content     string    `json:"content" binding:"required"`
	Type        int8      `json:"type" binding:"required,oneof=1 2 3 4 5"`
	Channel     string    `json:"channel" binding:"required,oneof=app email sms wechat"`
	TargetType  string    `json:"target_type" binding:"required,oneof=all role class group"`
	TargetIDs   []uint    `json:"target_ids"`
	ScheduledAt time.Time `json:"scheduled_at" binding:"required"`
}

// ScheduledListRequest 定时通知列表请求
type ScheduledListRequest struct {
	PageRequest
	Status *int8 `form:"status"`
}

// ==================== 通知设置 ====================

// NotificationSettingInfo 通知设置信息
type NotificationSettingInfo struct {
	ID           uint   `json:"id"`
	Type         int8   `json:"type"`
	Channel      string `json:"channel"`
	Enabled      bool   `json:"enabled"`
	DoNotDisturb bool   `json:"do_not_disturb"`
	DisturbStart string `json:"disturb_start"`
	DisturbEnd   string `json:"disturb_end"`
}

// UpdateNotificationSettingRequest 更新通知设置请求
type UpdateNotificationSettingRequest struct {
	Settings []NotificationSettingItem `json:"settings" binding:"required"`
}

// NotificationSettingItem 单条通知设置
type NotificationSettingItem struct {
	Type         int8   `json:"type" binding:"required,oneof=1 2 3 4 5"`
	Channel      string `json:"channel" binding:"required,oneof=app email sms wechat"`
	Enabled      *bool  `json:"enabled"`
	DoNotDisturb *bool  `json:"do_not_disturb"`
	DisturbStart string `json:"disturb_start"`
	DisturbEnd   string `json:"disturb_end"`
}

// ==================== 通知统计 ====================

// NotificationStats 通知统计信息
type NotificationStats struct {
	TotalCount   int64            `json:"total_count"`
	UnreadCount  int64            `json:"unread_count"`
	ReadCount    int64            `json:"read_count"`
	TypeStats    []TypeStatItem   `json:"type_stats"`
	ChannelStats []ChannelStatItem `json:"channel_stats"`
	DailyStats   []DailyStatItem  `json:"daily_stats"`
}

// TypeStatItem 按类型统计
type TypeStatItem struct {
	Type  int8  `json:"type"`
	Count int64 `json:"count"`
}

// ChannelStatItem 按渠道统计
type ChannelStatItem struct {
	Channel string `json:"channel"`
	Count   int64  `json:"count"`
}

// DailyStatItem 每日统计
type DailyStatItem struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

// ==================== 批量操作 ====================

// NotificationBatchReadRequest 批量标记已读请求
type NotificationBatchReadRequest struct {
	IDs []uint `json:"ids" binding:"required"`
}

// NotificationBatchDeleteRequest 批量删除请求
type NotificationBatchDeleteRequest struct {
	IDs []uint `json:"ids" binding:"required"`
}

// ==================== 通知渠道管理 ====================

// ChannelInfo 通知渠道信息
type ChannelInfo struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Config    string    `json:"config"`
	IsEnabled bool      `json:"is_enabled"`
	Priority  int       `json:"priority"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ChannelListRequest 渠道列表请求
type ChannelListRequest struct {
	PageRequest
	Type      string `form:"type"`
	IsEnabled *bool  `form:"is_enabled"`
}

// UpdateChannelRequest 更新渠道配置请求
type UpdateChannelRequest struct {
	Name      string `json:"name" binding:"max=100"`
	Config    string `json:"config"`
	IsEnabled *bool  `json:"is_enabled"`
	Priority  *int   `json:"priority"`
}
