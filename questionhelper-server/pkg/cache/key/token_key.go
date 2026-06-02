package key

import (
	"fmt"
	"time"
)

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

// DeactivateCooldownKey 注销取消冷却期 Key（30天内不可再次申请注销）
func DeactivateCooldownKey(uid uint) string {
	return fmt.Sprintf("qh:deactivate:cooldown:%d", uid)
}

// DeactivateCooldownDuration 注销取消冷却期时长（30天）
const DeactivateCooldownDuration = 30 * 24 * time.Hour

// ---- 通用验证码过期时间常量 ----

// CodeExpire 验证码过期时间（5分钟）
const CodeExpire = 300

// CodeLimitExpire 发送限制时间（1分钟）
const CodeLimitExpire = 60

// ---- 每日发送限制 ----

// EmailDailyKey 邮箱每日发送次数 Key
func EmailDailyKey(email string) string {
	return fmt.Sprintf("qh:email:daily:%s", email)
}

// EmailIPKey 邮箱验证码 IP 限制 Key
func EmailIPKey(ip string) string {
	return fmt.Sprintf("qh:email:ip:%s", ip)
}

// EmailAttemptsKey 邮箱验证码验证尝试次数 Key
func EmailAttemptsKey(email, code string) string {
	return fmt.Sprintf("qh:email:attempts:%s:%s", email, code)
}

// DailyLimitExpire 每日限制过期时间（计算到当天结束的秒数）
func DailyLimitExpire() time.Duration {
	now := time.Now()
	endOfDay := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
	return time.Until(endOfDay)
}

// EmailIPLimitExpire 邮箱 IP 限制过期时间（1小时）
const EmailIPLimitExpire = 1 * time.Hour

// MaxEmailDaily 每日邮箱发送上限
const MaxEmailDaily = 10

// MaxSmsDaily 每日短信发送上限
const MaxSmsDaily = 10

// MaxEmailIPHourly 每小时邮箱 IP 发送上限
const MaxEmailIPHourly = 50

// MaxEmailCodeAttempts 邮箱验证码验证尝试上限
const MaxEmailCodeAttempts = 5
