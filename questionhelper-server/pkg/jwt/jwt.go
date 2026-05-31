package jwt

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

// Claims JWT Claims
type Claims struct {
	UserID   uint     `json:"user_id"`
	Username string   `json:"username"`
	RoleIDs  []uint   `json:"role_ids"`
	JTI      string   `json:"jti"` // JWT Unique Identifier
	jwt.RegisteredClaims
}

// Init 初始化 JWT 密钥
func Init(secret string) {
	jwtSecret = []byte(secret)
}

// GenerateJTI 生成唯一的 JTI
func GenerateJTI() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// GenerateToken 生成 Token（向后兼容）
func GenerateToken(userID uint, username string, roleIDs []uint, expire int) (string, error) {
	jti := GenerateJTI()
	return GenerateTokenWithJTI(userID, username, roleIDs, jti, expire)
}

// GenerateTokenWithJTI 生成带 JTI 的 Token
func GenerateTokenWithJTI(userID uint, username string, roleIDs []uint, jti string, expire int) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		RoleIDs:  roleIDs,
		JTI:      jti,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expire) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        jti,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseToken 解析 Token
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
