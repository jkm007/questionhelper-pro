package encrypt

import (
	"golang.org/x/crypto/bcrypt"
	"questionhelper-server/pkg/config"
)

// HashPassword 使用 bcrypt 加密密码
func HashPassword(password string) (string, error) {
	cost := bcrypt.DefaultCost // 默认 10
	if config.Cfg != nil && config.Cfg.Auth.BcryptCost > 0 {
		cost = config.Cfg.Auth.BcryptCost
		if cost > bcrypt.MaxCost {
			cost = bcrypt.MaxCost
		}
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(bytes), err
}

// CheckPassword 验证密码
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
