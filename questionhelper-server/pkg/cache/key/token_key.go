package key

import "fmt"

// Token 相关缓存 Key（统一使用 qh: 前缀）

// TokenBlacklistKey Token 黑名单 Key（使用 JTI）
func TokenBlacklistKey(jti string) string {
	return fmt.Sprintf("qh:token:bl:%s", jti)
}

// CaptchaAttemptsKey 验证码错误次数 Key
func CaptchaAttemptsKey(key string) string {
	return fmt.Sprintf("qh:captcha:attempts:%s", key)
}

// SmsKey 短信验证码 Key
func SmsKey(phone string) string {
	return fmt.Sprintf("qh:sms:%s", phone)
}

// SmsLimitKey 短信发送频率限制 Key
func SmsLimitKey(phone string) string {
	return fmt.Sprintf("qh:sms:limit:%s", phone)
}

// SmsDailyKey 每日短信限制 Key
func SmsDailyKey(phone string) string {
	return fmt.Sprintf("qh:sms:daily:%s", phone)
}

// EmailKey 邮箱验证码 Key
func EmailKey(email string) string {
	return fmt.Sprintf("qh:email:%s", email)
}

// EmailLimitKey 邮箱发送频率限制 Key
func EmailLimitKey(email string) string {
	return fmt.Sprintf("qh:email:limit:%s", email)
}

// LoginLimitKey IP 登录限制 Key
func LoginLimitKey(ip string) string {
	return fmt.Sprintf("qh:login:limit:%s", ip)
}

// ResetTokenKey 密码重置令牌 Key
func ResetTokenKey(uid uint) string {
	return fmt.Sprintf("qh:reset:token:%d", uid)
}

// DeactivateCodeKey 注销确认验证码 Key
func DeactivateCodeKey(uid uint) string {
	return fmt.Sprintf("qh:deactivate:code:%d", uid)
}
