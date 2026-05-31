package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"questionhelper-server/pkg/cache/key"
	"questionhelper-server/pkg/database"
	"questionhelper-server/pkg/jwt"
	"questionhelper-server/pkg/response"
)

// AuthMiddleware JWT 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "未提供认证令牌")
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			response.Error(c, http.StatusUnauthorized, "认证格式错误")
			c.Abort()
			return
		}

		claims, err := jwt.ParseToken(tokenString)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "认证令牌无效")
			c.Abort()
			return
		}

		// 检查 Token 是否在黑名单中
		ctx := context.Background()
		blacklistKey := key.TokenBlacklistKey(claims.JTI)
		exists, _ := database.RDB.Exists(ctx, blacklistKey).Result()
		if exists > 0 {
			response.Error(c, http.StatusUnauthorized, "令牌已失效")
			c.Abort()
			return
		}

		// 设置上下文信息
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role_ids", claims.RoleIDs)
		c.Set("jti", claims.JTI)

		c.Next()
	}
}
