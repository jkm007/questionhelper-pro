package jwt

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte
var jwtIssuer string = "questionhelper"

// Claims JWT Claims
type Claims struct {
	UserID    uint     `json:"user_id"`
	Username  string   `json:"username"`
	RoleIDs   []uint   `json:"role_ids"`
	RoleCodes []string `json:"role_codes"` // 角色编码，用于中间件权限判断
	JTI       string   `json:"jti"`        // JWT Unique Identifier
	Type      string   `json:"type"`       // Token 类型: "access" 或 "refresh"（T28）
	jwt.RegisteredClaims
}

// Init 初始化 JWT 密钥和 Issuer
func Init(secret string, issuer ...string) {
	jwtSecret = []byte(secret)
	if len(issuer) > 0 && issuer[0] != "" {
		jwtIssuer = issuer[0]
	}
}

// GenerateJTI 生成唯一的 JTI
func GenerateJTI() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// GenerateToken 生成 Token（向后兼容，默认生成 access token）
func GenerateToken(userID uint, username string, roleIDs []uint, roleCodes []string, expire int) (string, error) {
	jti := GenerateJTI()
	return GenerateTokenWithJTI(userID, username, roleIDs, roleCodes, jti, expire, "access")
}

// GenerateTokenWithJTI 生成带 JTI 的 Token（T27: 添加 Issuer 字段，T28: 添加 Type 字段）
func GenerateTokenWithJTI(userID uint, username string, roleIDs []uint, roleCodes []string, jti string, expire int, tokenType string) (string, error) {
	claims := Claims{
		UserID:    userID,
		Username:  username,
		RoleIDs:   roleIDs,
		RoleCodes: roleCodes,
		JTI:       jti,
		Type:      tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expire) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        jti,
			Issuer:    jwtIssuer,
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
