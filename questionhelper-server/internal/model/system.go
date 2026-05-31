package model

import (
	"time"
)

type SystemSetting struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Key       string    `gorm:"uniqueIndex;size:100;not null" json:"key"`
	Value     string    `gorm:"type:text" json:"value"`
	Remark    string    `gorm:"size:200" json:"remark"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (SystemSetting) TableName() string {
	return "system_settings"
}

type OperationLog struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	UserID      uint      `gorm:"index" json:"user_id"`
	Username    string    `gorm:"size:50" json:"username"`
	Module      string    `gorm:"size:50" json:"module"`                       // 模块
	Action      string    `gorm:"size:50;not null" json:"action"`
	Resource    string    `gorm:"size:50" json:"resource"`
	ResourceID  string    `gorm:"size:50;index" json:"resource_id"`
	Description string    `gorm:"size:200" json:"description"`                 // 操作描述
	Method      string    `gorm:"size:10" json:"method"`
	Path        string    `gorm:"size:500" json:"path"`
	IP          string    `gorm:"size:50" json:"ip"`
	UserAgent   string    `gorm:"size:500" json:"user_agent"`
	Status      int8      `gorm:"default:1" json:"status"`                     // 状态:0=失败,1=成功
	ErrorMsg    string    `gorm:"size:500" json:"error_msg"`                   // 错误信息
	Latency     int64     `gorm:"comment:耗时(ms)" json:"latency"`
	CreatedAt   time.Time `json:"created_at"`
}

func (OperationLog) TableName() string {
	return "operation_logs"
}

type LoginLog struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	UserID      *uint     `gorm:"index" json:"user_id"` // 可空，登录失败时可能没有用户ID
	Username    string    `gorm:"size:50" json:"username"`
	IP          string    `gorm:"size:50" json:"ip"`
	Location    string    `gorm:"size:100" json:"location"`     // 登录地点
	UserAgent   string    `gorm:"size:500" json:"user_agent"`
	DeviceType  string    `gorm:"size:20" json:"device_type"`   // 设备类型:web/ios/android/miniapp
	DeviceName  string    `gorm:"size:100" json:"device_name"`  // 设备名称
	Browser     string    `gorm:"size:50" json:"browser"`       // 浏览器
	OS          string    `gorm:"size:50" json:"os"`            // 操作系统
	Status      int8      `gorm:"default:1;comment:状态:0=失败,1=成功" json:"status"`
	Msg         string    `gorm:"size:200" json:"msg"`          // 消息
	LoginType   string    `gorm:"size:20;default:password" json:"login_type"` // 登录类型:password/sms/oauth
	IsNewDevice bool      `gorm:"default:false" json:"is_new_device"`         // 是否新设备
	CreatedAt   time.Time `json:"created_at"`
}

func (LoginLog) TableName() string {
	return "login_logs"
}
