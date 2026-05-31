package dto

import "time"

// ClassInfo 班级信息
type ClassInfo struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Cover       string    `json:"cover"`
	Code        string    `json:"code"`
	CreatorID   uint      `json:"creator_id"`
	CreatorName string    `json:"creator_name"`
	MemberCount int       `json:"member_count"`
	Status      int8      `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

// ClassMemberInfo 班级成员信息
type ClassMemberInfo struct {
	ID       uint      `json:"id"`
	ClassID  uint      `json:"class_id"`
	UserID   uint      `json:"user_id"`
	Role     int8      `json:"role"`
	JoinedAt time.Time `json:"joined_at"`
}

// CreateClassRequest 创建班级请求
type CreateClassRequest struct {
	Name        string `json:"name" binding:"required,max=100"`
	Description string `json:"description" binding:"omitempty,max=500"`
	Cover       string `json:"cover" binding:"omitempty,max=255"`
}

// UpdateClassRequest 更新班级请求
type UpdateClassRequest struct {
	Name        string `json:"name" binding:"omitempty,max=100"`
	Description string `json:"description" binding:"omitempty,max=500"`
	Cover       string `json:"cover" binding:"omitempty,max=255"`
}

// JoinClassRequest 加入班级请求
type JoinClassRequest struct {
	Code string `json:"code" binding:"required"`
}

// ClassListRequest 班级列表请求
type ClassListRequest struct {
	PageRequest
	Keyword string `form:"keyword"`
}

// ClassMemberListRequest 班级成员列表请求
type ClassMemberListRequest struct {
	PageRequest
	Role *int8 `form:"role"`
}

// HomeworkInfo 作业信息
type HomeworkInfo struct {
	ID          uint      `json:"id"`
	ClassID     uint      `json:"class_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Deadline    time.Time `json:"deadline"`
	CreatorID   uint      `json:"creator_id"`
}

// CreateHomeworkRequest 创建作业请求
type CreateHomeworkRequest struct {
	Title       string    `json:"title" binding:"required,max=200"`
	Description string    `json:"description"`
	Deadline    time.Time `json:"deadline" binding:"required"`
}

// ClassNoticeInfo 班级公告信息
type ClassNoticeInfo struct {
	ID        uint      `json:"id"`
	ClassID   uint      `json:"class_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatorID uint      `json:"creator_id"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateNoticeRequest 创建公告请求
type CreateNoticeRequest struct {
	Title   string `json:"title" binding:"required,max=200"`
	Content string `json:"content" binding:"required"`
}
