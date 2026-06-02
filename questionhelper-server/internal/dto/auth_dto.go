package dto

import "time"

// LoginRequest 登录请求
type LoginRequest struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required_if=LoginType password"`
	CaptchaID  string `json:"captchaId"`
	Captcha    string `json:"captchaCode"`
	DeviceID   string `json:"deviceId"`   // 设备唯一标识
	DeviceInfo string `json:"deviceInfo"` // 设备信息（User-Agent）
	UserAgent  string `json:"-"`          // User-Agent 请求头（由控制器自动填充）
	RememberMe bool   `json:"rememberMe"` // 记住我，Refresh Token延长至30天
	LoginType  string `json:"loginType" binding:"omitempty,oneof=password email_code"` // 登录类型，默认 password
	EmailCode  string `json:"emailCode"`  // 邮箱验证码（邮箱登录时必填）
}

// LoginResponse 登录响应
type LoginResponse struct {
	AccessToken  string   `json:"accessToken"`
	RefreshToken string   `json:"refreshToken"`
	TokenType    string   `json:"tokenType"`
	ExpiresIn    int      `json:"expiresIn"`
	User         UserInfo `json:"user,omitempty"`
}

// RegisterRequest 注册请求（T08: 邮箱必填，T17: 用户名4-20，T18: 昵称2-20，T19: 验证码必填）
type RegisterRequest struct {
	Username        string `json:"username" binding:"required,min=4,max=20"`          // T17: 4-20字符
	Password        string `json:"password" binding:"required,min=6,max=20,password"` // T09: 需包含字母+数字（自定义验证器）
	ConfirmPassword string `json:"confirmPassword" binding:"required,eqfield=Password"`
	Nickname        string `json:"nickname" binding:"omitempty,min=2,max=20"` // T18: 2-20字符
	Phone           string `json:"phone" binding:"omitempty"`
	Email           string `json:"email" binding:"required,email"` // T08: 邮箱必填
	EmailCode       string `json:"emailCode"`                      // T08: 邮箱验证码
	CaptchaID       string `json:"captchaId" binding:"required"`   // T19: 必填
	CaptchaCode     string `json:"captchaCode" binding:"required"` // T19: 必填
	Agreement       bool   `json:"agreement" binding:"required"`
	Source          string `json:"source"` // 注册来源:h5/miniapp/app/web
	DeviceID        string `json:"-"`      // 设备标识（由控制器自动填充）
	DeviceInfo      string `json:"-"`      // 设备IP（由控制器自动填充）
	UserAgent       string `json:"-"`      // User-Agent（由控制器自动填充）
}

// RegisterResponse 注册响应
type RegisterResponse struct {
	AccessToken  string   `json:"accessToken"`
	RefreshToken string   `json:"refreshToken"`
	TokenType    string   `json:"tokenType"`
	ExpiresIn    int      `json:"expiresIn"`
	User         UserInfo `json:"user,omitempty"`
}

// RefreshTokenRequest 刷新Token请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword     string `json:"oldPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required,min=6,max=20,password"`
	ConfirmPassword string `json:"confirmPassword" binding:"required,eqfield=NewPassword"`
}

// RequestPasswordResetRequest 请求重置密码请求（T23: 验证码必填）
type RequestPasswordResetRequest struct {
	Account     string `json:"account" binding:"required"`     // 用户名/手机号/邮箱
	CaptchaID   string `json:"captchaId" binding:"required"`   // T23: 图形验证码ID
	CaptchaCode string `json:"captchaCode" binding:"required"` // T23: 图形验证码
}

// ResetPasswordRequest 重置密码请求
type ResetPasswordRequest struct {
	Account     string `json:"account" binding:"required"` // 用户名/手机号/邮箱
	Code        string `json:"code" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=6,max=20,password"`
}

// CaptchaResponse 验证码响应
type CaptchaResponse struct {
	CaptchaID     string `json:"captchaId"`
	CaptchaBase64 string `json:"captchaBase64"`
}

// DeviceResponse 设备信息响应
type DeviceResponse struct {
	ID           uint   `json:"id"`
	DeviceType   string `json:"deviceType"`
	DeviceName   string `json:"deviceName"`
	Browser      string `json:"browser"`
	OS           string `json:"os"`
	IP           string `json:"ip"`
	Location     string `json:"location"`
	LastActiveAt string `json:"lastActiveAt"`
	IsCurrent    bool   `json:"isCurrent"`
	CreatedAt    string `json:"createdAt"`
}

// SecurityLogResponse 安全日志响应
type SecurityLogResponse struct {
	ID          uint   `json:"id"`
	EventType   string `json:"eventType"`
	EventDetail string `json:"eventDetail"`
	IP          string `json:"ip"`
	Status      int8   `json:"status"`
	CreatedAt   string `json:"createdAt"`
}

// DeactivateRequest 注销账号请求
type DeactivateRequest struct {
	Password string `json:"password" binding:"required"`
	Reason   string `json:"reason"`
}

// DeactivateConfirmRequest 确认注销请求
type DeactivateConfirmRequest struct {
	Code string `json:"code" binding:"required"`
}

// UserDataExport 用户数据导出
type UserDataExport struct {
	ExportAt       time.Time   `json:"exportAt"`
	User           interface{} `json:"user"`
	Exams          interface{} `json:"exams"`
	Practices      interface{} `json:"practices"`
	WrongQuestions interface{} `json:"wrongQuestions"`
	Favorites      interface{} `json:"favorites"`
	Comments       interface{} `json:"comments"`
}

// SwitchRoleRequest 角色切换请求
type SwitchRoleRequest struct {
	RoleID uint `json:"role_id" binding:"required"`
}

// SecuritySettingsRequest 安全设置更新请求
type SecuritySettingsRequest struct {
	LoginNotification    *bool `json:"loginNotification"`    // 新登录通知
	PasswordChangeNotify *bool `json:"passwordChangeNotify"` // 密码变更通知
	DeviceManageNotify   *bool `json:"deviceManageNotify"`   // 设备管理通知
}

// SecuritySettingsResponse 安全设置响应
type SecuritySettingsResponse struct {
	LoginNotification    bool `json:"loginNotification"`
	PasswordChangeNotify bool `json:"passwordChangeNotify"`
	DeviceManageNotify   bool `json:"deviceManageNotify"`
}

// OAuthBindRequest 绑定第三方账号请求
type OAuthBindRequest struct {
	Code  string `json:"code" binding:"required"`
	State string `json:"state"`
}

// OAuthStatusResponse 第三方绑定状态响应
type OAuthStatusResponse struct {
	Wechat *OAuthProviderStatus `json:"wechat,omitempty"`
	Github *OAuthProviderStatus `json:"github,omitempty"`
	Google *OAuthProviderStatus `json:"google,omitempty"`
}

// OAuthProviderStatus 第三方账号绑定状态
type OAuthProviderStatus struct {
	Bound    bool   `json:"bound"`
	Nickname string `json:"nickname,omitempty"`
	BoundAt  string `json:"boundAt,omitempty"`
}
