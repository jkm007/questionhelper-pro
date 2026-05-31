package dto

import "time"

// OperationLogListRequest 操作日志列表请求
type OperationLogListRequest struct {
	PageRequest
	UserID    *uint      `form:"user_id"`
	Module    string     `form:"module"`
	Action    string     `form:"action"`
	Status    *int8      `form:"status"`
	IP        string     `form:"ip"`
	StartTime *time.Time `form:"start_time"`
	EndTime   *time.Time `form:"end_time"`
}

// OperationLogInfo 操作日志信息
type OperationLogInfo struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"user_id"`
	Username    string    `json:"username"`
	Module      string    `json:"module"`
	Action      string    `json:"action"`
	Resource    string    `json:"resource"`
	ResourceID  string    `json:"resource_id"`
	Description string    `json:"description"`
	IP          string    `json:"ip"`
	UserAgent   string    `json:"user_agent"`
	Status      int8      `json:"status"`
	ErrorMsg    string    `json:"error_msg"`
	CreatedAt   time.Time `json:"created_at"`
}

// LoginLogListRequest 登录日志列表请求
type LoginLogListRequest struct {
	PageRequest
	UserID    *uint      `form:"user_id"`
	Status    *int8      `form:"status"`
	IP        string     `form:"ip"`
	StartTime *time.Time `form:"start_time"`
	EndTime   *time.Time `form:"end_time"`
}

// LoginLogInfo 登录日志信息
type LoginLogInfo struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Username  string    `json:"username"`
	IP        string    `json:"ip"`
	Location  string    `json:"location"`
	Browser   string    `json:"browser"`
	OS        string    `json:"os"`
	Status    int8      `json:"status"`
	Msg       string    `json:"msg"`
	CreatedAt time.Time `json:"created_at"`
}
