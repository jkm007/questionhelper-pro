package dto

import "time"

// LoginRequest 登录请求
type LoginRequest struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	CaptchaID  string `json:"captchaId"`
	Captcha    string `json:"captchaCode"`
	DeviceID   string `json:"deviceId"`   // 设备唯一标识
	DeviceInfo string `json:"deviceInfo"` // 设备信息（User-Agent）
	RememberMe bool   `json:"rememberMe"` // 记住我，Refresh Token延长至30天
}

// LoginResponse 登录响应
type LoginResponse struct {
	AccessToken  string   `json:"accessToken"`
	RefreshToken string   `json:"refreshToken"`
	TokenType    string   `json:"tokenType"`
	ExpiresIn    int      `json:"expiresIn"`
	User         UserInfo `json:"user,omitempty"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username        string `json:"username" binding:"required,min=2,max=50"`
	Password        string `json:"password" binding:"required,min=6,max=20"`
	ConfirmPassword string `json:"confirmPassword" binding:"required,eqfield=Password"`
	Nickname        string `json:"nickname" binding:"omitempty,max=50"`
	Phone           string `json:"phone" binding:"omitempty"`
	Email           string `json:"email" binding:"omitempty,email"`
	CaptchaID       string `json:"captchaId"`
	CaptchaCode     string `json:"captchaCode"`
	Agreement       bool   `json:"agreement" binding:"required"`
	Source          string `json:"source"` // 注册来源:h5/miniapp/app/web
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
	NewPassword     string `json:"newPassword" binding:"required,min=6,max=20"`
	ConfirmPassword string `json:"confirmPassword" binding:"required,eqfield=NewPassword"`
}

// ResetPasswordRequest 重置密码请求
type ResetPasswordRequest struct {
	Account     string `json:"account" binding:"required"` // 用户名/手机号/邮箱
	Code        string `json:"code" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=6,max=20"`
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
