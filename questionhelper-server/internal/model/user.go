package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID                uint           `gorm:"primarykey" json:"id"`
	Username          string         `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Password          string         `gorm:"size:100;not null" json:"-"`
	Nickname          string         `gorm:"size:50" json:"nickname"`
	Email             string         `gorm:"size:100" json:"email"`
	Phone             string         `gorm:"size:20" json:"phone"`
	Avatar            string         `gorm:"size:255" json:"avatar"`
	Gender            int8           `gorm:"default:0" json:"gender"` // 0:未知 1:男 2:女
	Birthday          *time.Time     `json:"birthday"`
	Bio               string         `gorm:"size:500" json:"bio"`
	Status            int8           `gorm:"default:1" json:"status"` // 0:禁用 1:正常 2:注销中
	RealName          string         `gorm:"size:50" json:"real_name"`
	IDCard            string         `gorm:"size:100" json:"-"`
	IsReal            bool           `gorm:"default:false" json:"is_real"`
	RegisterSource    string         `gorm:"size:20;default:web" json:"register_source"` // 注册来源:h5/miniapp/app/web
	LastLoginAt       *time.Time     `json:"last_login_at"`
	LastLoginIP       string         `gorm:"size:50" json:"last_login_ip"`
	LastLoginDevice   string         `gorm:"size:200" json:"last_login_device"`
	LoginFailCount    int            `gorm:"default:0" json:"-"` // 连续登录失败次数
	LockUntil         *time.Time     `json:"-"`                  // 锁定截止时间
	PasswordChangedAt *time.Time     `json:"password_changed_at"`
	LogoutAt          *time.Time     `json:"logout_at"` // 注销申请时间
	Roles             []Role         `gorm:"many2many:user_roles;" json:"roles"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	Privacy    *UserPrivacy `gorm:"foreignKey:UserID" json:"privacy,omitempty"`
	OAuthUsers []OAuthUser  `gorm:"foreignKey:UserID" json:"oauth_users,omitempty"`
	Tags       []Tag        `gorm:"many2many:user_tags;" json:"tags,omitempty"`
}

func (User) TableName() string {
	return "users"
}

// UserPrivacy 用户隐私设置表
type UserPrivacy struct {
	ID              uint      `gorm:"primarykey" json:"id"`
	UserID          uint      `gorm:"uniqueIndex" json:"user_id"`
	ProfileVisible  int8      `gorm:"default:1" json:"profile_visible"`  // 个人主页可见性:1=所有人,2=仅班级成员,3=仅自己
	RealnameVisible int8      `gorm:"default:1" json:"realname_visible"` // 真实姓名可见性
	EmailVisible    int8      `gorm:"default:1" json:"email_visible"`    // 邮箱可见性
	StatsVisible    int8      `gorm:"default:1" json:"stats_visible"`    // 学习统计可见性
	ClassVisible    int8      `gorm:"default:1" json:"class_visible"`    // 班级信息可见性
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (UserPrivacy) TableName() string {
	return "user_privacy"
}

// OAuthUser 第三方登录绑定表
type OAuthUser struct {
	ID               uint       `gorm:"primarykey" json:"id"`
	UserID           uint       `gorm:"index;not null" json:"user_id"`
	Provider         string     `gorm:"size:20;not null" json:"provider"`          // 提供商:wechat/github/google
	ProviderType     string     `gorm:"size:20" json:"provider_type"`              // 提供商类型:miniapp/official/scan
	ProviderUserID   string     `gorm:"size:100;not null" json:"provider_user_id"` // 第三方用户ID
	ProviderUsername string     `gorm:"size:100" json:"provider_username"`
	ProviderAvatar   string     `gorm:"size:255" json:"provider_avatar"`
	AccessToken      string     `gorm:"size:500" json:"-"`
	RefreshToken     string     `gorm:"size:500" json:"-"`
	ExpiresAt        *time.Time `json:"expires_at"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

func (OAuthUser) TableName() string {
	return "oauth_users"
}

// LoginDevice 登录设备表
type LoginDevice struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	UserID       uint      `gorm:"index;not null" json:"user_id"`
	DeviceID     string    `gorm:"size:100;not null" json:"device_id"`
	DeviceType   string    `gorm:"size:20;not null" json:"device_type"` // 设备类型:web/ios/android/miniapp
	DeviceName   string    `gorm:"size:100" json:"device_name"`
	Browser      string    `gorm:"size:50" json:"browser"`
	OS           string    `gorm:"size:50" json:"os"`
	IP           string    `gorm:"size:50" json:"ip"`
	Location     string    `gorm:"size:100" json:"location"`
	TokenJTI     string    `gorm:"size:100;not null" json:"token_jti"`
	LastActiveAt time.Time `json:"last_active_at"`
	IsCurrent    bool      `gorm:"default:false" json:"is_current"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (LoginDevice) TableName() string {
	return "login_devices"
}

// SecurityLog 安全日志表
type SecurityLog struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	UserID      uint      `gorm:"index;not null" json:"user_id"`
	EventType   string    `gorm:"size:30;not null" json:"event_type"` // 事件类型:login/logout/password_change/bind/unbind/account_deactivate
	EventDetail string    `gorm:"size:500" json:"event_detail"`
	IP          string    `gorm:"size:50" json:"ip"`
	UserAgent   string    `gorm:"size:500" json:"user_agent"`
	Status      int8      `gorm:"default:1" json:"status"` // 0=失败,1=成功
	CreatedAt   time.Time `json:"created_at"`
}

func (SecurityLog) TableName() string {
	return "security_logs"
}

// PasswordHistory 密码历史表
type PasswordHistory struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	Password  string    `gorm:"size:100;not null" json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

func (PasswordHistory) TableName() string {
	return "password_history"
}

type Role struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `gorm:"size:50;not null" json:"name"`
	Code        string         `gorm:"uniqueIndex;size:50;not null" json:"code"`
	Description string         `gorm:"size:200" json:"description"`
	IsDefault   bool           `gorm:"default:false" json:"is_default"`
	IsSystem    bool           `gorm:"default:false" json:"is_system"` // 系统角色不可删除
	Sort        int            `gorm:"default:0" json:"sort"`
	Status      int8           `gorm:"default:1" json:"status"`
	Menus       []Menu         `gorm:"many2many:role_menus;" json:"menus"`
	Permissions []Permission   `gorm:"many2many:role_permissions;" json:"permissions"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Role) TableName() string {
	return "roles"
}

type Menu struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	ParentID   *uint          `gorm:"index" json:"parent_id"`
	Name       string         `gorm:"size:50;not null" json:"name"`
	Path       string         `gorm:"size:200" json:"path"`
	Component  string         `gorm:"size:200" json:"component"`
	Redirect   string         `gorm:"size:200" json:"redirect"`
	Title      string         `gorm:"size:50;not null" json:"title"`
	Icon       string         `gorm:"size:50" json:"icon"`
	Hidden     bool           `gorm:"default:false" json:"hidden"`
	Type       int8           `gorm:"not null;comment:类型:1=目录,2=菜单,3=按钮" json:"type"`
	Permission string         `gorm:"size:100" json:"permission"`
	Sort       int            `gorm:"default:0" json:"sort"`
	Status     int8           `gorm:"default:1" json:"status"`
	Children   []Menu         `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Menu) TableName() string {
	return "menus"
}

type Permission struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `gorm:"size:50;not null" json:"name"`
	Code        string         `gorm:"uniqueIndex;size:100;not null" json:"code"`
	Description string         `gorm:"size:200" json:"description"`
	Type        int8           `gorm:"default:1;comment:类型:1=菜单,2=按钮,3=接口" json:"type"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Permission) TableName() string {
	return "permissions"
}
