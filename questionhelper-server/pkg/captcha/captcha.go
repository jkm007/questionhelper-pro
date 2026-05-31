package captcha

import (
	"context"
	"time"

	"github.com/mojocn/base64Captcha"
	"questionhelper-server/pkg/cache/key"
	"questionhelper-server/pkg/database"
)

var store = base64Captcha.DefaultMemStore

// GenerateCaptcha 生成验证码
func GenerateCaptcha() (id string, b64s string, err error) {
	driver := base64Captcha.NewDriverDigit(80, 240, 5, 0.7, 80)
	c := base64Captcha.NewCaptcha(driver, store)
	id, b64s, _, err = c.Generate()
	return
}

// VerifyCaptcha 验证验证码
func VerifyCaptcha(id string, answer string) bool {
	return store.Verify(id, answer, true)
}

// SaveCaptchaToRedis 保存验证码到 Redis
func SaveCaptchaToRedis(id string, answer string) error {
	ctx := context.Background()
	return database.RDB.Set(
		ctx,
		key.CaptchaKey(id),
		answer,
		time.Duration(key.CaptchaExpire)*time.Second,
	).Err()
}

// VerifyCaptchaFromRedis 从 Redis 验证验证码
func VerifyCaptchaFromRedis(id string, answer string) bool {
	ctx := context.Background()
	stored, err := database.RDB.Get(ctx, key.CaptchaKey(id)).Result()
	if err != nil {
		return false
	}
	database.RDB.Del(ctx, key.CaptchaKey(id))
	return stored == answer
}
