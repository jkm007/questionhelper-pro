package key

import "fmt"

// 图形验证码相关缓存 Key

// CaptchaKey 图形验证码 Key（统一使用 qh: 前缀）
func CaptchaKey(captchaID string) string {
	return fmt.Sprintf("qh:captcha:%s", captchaID)
}

// CaptchaExpire 验证码过期时间（5分钟）
const CaptchaExpire = 300
