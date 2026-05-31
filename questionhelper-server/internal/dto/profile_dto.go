package dto

// UpdatePrivacyRequest 更新隐私设置请求
type UpdatePrivacyRequest struct {
	ProfileVisible  *int8 `json:"profile_visible" binding:"omitempty,oneof=1 2 3"`
	RealnameVisible *int8 `json:"realname_visible" binding:"omitempty,oneof=1 2 3"`
	EmailVisible    *int8 `json:"email_visible" binding:"omitempty,oneof=1 2 3"`
	StatsVisible    *int8 `json:"stats_visible" binding:"omitempty,oneof=1 2 3"`
	ClassVisible    *int8 `json:"class_visible" binding:"omitempty,oneof=1 2 3"`
}

// PrivacyInfo 隐私设置信息
type PrivacyInfo struct {
	ProfileVisible  int8 `json:"profile_visible"`
	RealnameVisible int8 `json:"realname_visible"`
	EmailVisible    int8 `json:"email_visible"`
	StatsVisible    int8 `json:"stats_visible"`
	ClassVisible    int8 `json:"class_visible"`
}

// BindPhoneRequest 绑定手机号请求
type BindPhoneRequest struct {
	Phone string `json:"phone" binding:"required"`
	Code  string `json:"code" binding:"required,len=6"`
}

// BindEmailRequest 绑定邮箱请求
type BindEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required,len=6"`
}

// RealNameSubmitRequest 提交实名认证请求
type RealNameSubmitRequest struct {
	RealName string `json:"real_name" binding:"required,max=50"`
	IDCard   string `json:"id_card" binding:"required,len=18"`
}

// RealNameInfo 实名认证信息
type RealNameInfo struct {
	ID           uint   `json:"id"`
	RealName     string `json:"real_name"`
	IDCardMasked string `json:"id_card_masked"` // 脱敏后的身份证号
	Status       int8   `json:"status"`
	RejectReason string `json:"reject_reason"`
	SubmittedAt  string `json:"submitted_at"`
}

// ReviewRealNameRequest 审核实名认证请求
type ReviewRealNameRequest struct {
	Status int8   `json:"status" binding:"required,oneof=1 2"` // 1=通过,2=驳回`
	Reason string `json:"reason" binding:"omitempty,max=200"`
}

// OAuthInfo 第三方账号信息
type OAuthInfo struct {
	Provider       string `json:"provider"`
	ProviderType   string `json:"provider_type"`
	ProviderUserID string `json:"provider_user_id"`
	BindTime       string `json:"bind_time"`
}

// BindOAuthRequest 绑定第三方账号请求
type BindOAuthRequest struct {
	Provider       string `json:"provider" binding:"required"`
	ProviderUserID string `json:"provider_user_id" binding:"required"`
	AccessToken    string `json:"access_token"`
}
