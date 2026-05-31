package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// 默认密钥(生产环境应该从配置或环境变量读取)
var defaultKey = []byte("questionhelper2024aes256key!!") // 32字节 for AES-256

// SetAESKey 设置AES密钥
func SetAESKey(key []byte) error {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return errors.New("AES密钥长度必须为16/24/32字节")
	}
	defaultKey = key
	return nil
}

// AESEncrypt AES加密(返回base64编码)
func AESEncrypt(plaintext string) (string, error) {
	if len(defaultKey) == 0 {
		return "", errors.New("AES密钥未设置")
	}

	block, err := aes.NewCipher(defaultKey)
	if err != nil {
		return "", err
	}

	// 使用GCM模式
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 生成随机nonce
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// 加密
	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)

	// 返回base64编码
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// AESDecrypt AES解密(输入base64编码)
func AESDecrypt(ciphertext string) (string, error) {
	if len(defaultKey) == 0 {
		return "", errors.New("AES密钥未设置")
	}

	// base64解码
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(defaultKey)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("密文长度无效")
	}

	nonce, encryptedData := data[:nonceSize], data[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// MaskIDCard 身份证号脱敏(保留前3后4)
func MaskIDCard(idCard string) string {
	if len(idCard) < 8 {
		return idCard
	}
	return idCard[:3] + "***********" + idCard[14:]
}

// MaskPhone 手机号脱敏(保留前3后4)
func MaskPhone(phone string) string {
	if len(phone) < 7 {
		return phone
	}
	return phone[:3] + "****" + phone[7:]
}

// MaskEmail 邮箱脱敏
func MaskEmail(email string) string {
	for i, c := range email {
		if c == '@' {
			if i <= 2 {
				return email
			}
			return email[:2] + "***" + email[i:]
		}
	}
	return email
}
