package key

import "fmt"

// 短信验证码相关缓存 Key

// SMSCodeKey 短信验证码 Key
func SMSCodeKey(phone string) string {
	return fmt.Sprintf("sms:code:%s", phone)
}

// SMSLimitKey 短信发送限制 Key
func SMSLimitKey(phone string) string {
	return fmt.Sprintf("sms:limit:%s", phone)
}

// SMSExpire 验证码过期时间（5分钟）
const SMSExpire = 300

// SMSLimitExpire 发送限制时间（1分钟）
const SMSLimitExpire = 60
