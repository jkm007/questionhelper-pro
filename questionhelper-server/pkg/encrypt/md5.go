package encrypt

import (
	"crypto/md5"
	"encoding/hex"
)

// MD5 计算 MD5 哈希
func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// MD5WithSalt 带盐的 MD5
func MD5WithSalt(str, salt string) string {
	return MD5(str + salt)
}
