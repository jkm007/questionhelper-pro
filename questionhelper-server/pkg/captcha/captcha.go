package captcha

import (
	"context"
	"strings"
	"time"

	"github.com/mojocn/base64Captcha"
	"questionhelper-server/pkg/cache/key"
	"questionhelper-server/pkg/database"
)

// CaptchaType 验证码类型
type CaptchaType string

const (
	CaptchaTypeDigit  CaptchaType = "digit"  // 纯数字验证码
	CaptchaTypeLetter CaptchaType = "letter" // 字母+数字验证码
	CaptchaTypeMath   CaptchaType = "math"   // 数学运算验证码
)

// RedisStore 使用 Redis 存储验证码（T03: 支持多实例部署）
type RedisStore struct{}

func (s *RedisStore) Set(id string, value string) error {
	ctx := context.Background()
	return database.RDB.Set(
		ctx,
		key.CaptchaKey(id),
		value,
		time.Duration(key.CaptchaExpire)*time.Second,
	).Err()
}

func (s *RedisStore) Get(id string, clear bool) string {
	ctx := context.Background()
	codeKey := key.CaptchaKey(id)
	val, err := database.RDB.Get(ctx, codeKey).Result()
	if err != nil {
		return ""
	}
	if clear {
		database.RDB.Del(ctx, codeKey)
	}
	return val
}

func (s *RedisStore) Verify(id, answer string, clear bool) bool {
	stored := s.Get(id, clear)
	if stored == "" {
		return false
	}
	// 忽略大小写
	return strings.EqualFold(stored, answer)
}

var store base64Captcha.Store = &RedisStore{}

// GenerateCaptcha 生成验证码（支持 digit/letter/math 三种类型）
func GenerateCaptcha(captchaType CaptchaType) (id string, b64s string, err error) {
	var driver base64Captcha.Driver
	switch captchaType {
	case CaptchaTypeLetter:
		// 字母+数字验证码：6位，中等干扰
		driver = base64Captcha.NewDriverString(80, 240, 5, base64Captcha.OptionShowHollowLine, 6, "abcdefghjkmnpqrstuvwxyz23456789", nil, nil, nil)
	case CaptchaTypeMath:
		// 数学运算验证码：如 "3+5=?"
		driver = base64Captcha.NewDriverMath(80, 240, 5, base64Captcha.OptionShowHollowLine, nil, nil, nil)
	default:
		// 默认纯数字验证码（向后兼容）
		driver = base64Captcha.NewDriverDigit(80, 240, 5, 0.7, 80)
	}
	c := base64Captcha.NewCaptcha(driver, store)
	id, b64s, _, err = c.Generate()
	return
}

// maxCaptchaAttempts 验证码最大尝试次数
const maxCaptchaAttempts = 5

// VerifyCaptcha 验证验证码（T03: 使用 Redis 存储，忽略大小写；T30: 支持尝试次数跟踪）
// 当 attemptsKey 非空时，启用尝试次数跟踪：失败超过 maxCaptchaAttempts 次后需重新获取验证码
func VerifyCaptcha(id string, answer string, attemptsKey string) (ok bool, tooManyAttempts bool) {
	ctx := context.Background()

	// T30: 检查尝试次数（如果启用了跟踪）
	if attemptsKey != "" {
		attemptsRedisKey := key.CaptchaAttemptsKey(attemptsKey)
		count, _ := database.RDB.Get(ctx, attemptsRedisKey).Int64()
		if count >= maxCaptchaAttempts {
			// 超过最大尝试次数，删除验证码并要求重新获取
			store.Get(id, true)
			return false, true
		}
	}

	// 验证答案（不清除验证码，成功或达到上限后再清除）
	stored := store.Get(id, false)
	if stored == "" {
		return false, false
	}
	ok = strings.EqualFold(stored, answer)

	if ok {
		// 验证成功：清除验证码和尝试次数
		store.Get(id, true)
		if attemptsKey != "" {
			database.RDB.Del(ctx, key.CaptchaAttemptsKey(attemptsKey))
		}
		return true, false
	}

	// 验证失败：记录尝试次数
	if attemptsKey != "" {
		attemptsRedisKey := key.CaptchaAttemptsKey(attemptsKey)
		database.RDB.Incr(ctx, attemptsRedisKey)
		database.RDB.Expire(ctx, attemptsRedisKey, time.Duration(key.CaptchaExpire)*time.Second)
	}

	return false, false
}
