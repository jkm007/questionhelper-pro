package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/repository/user"
	"questionhelper-server/pkg/database"
	"questionhelper-server/pkg/logger"
	"questionhelper-server/pkg/response"
)

// PermissionMiddleware 权限验证中间件
func PermissionMiddleware(requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			response.Error(c, http.StatusForbidden, "无权限访问")
			c.Abort()
			return
		}

		uid, ok := userID.(uint)
		if !ok {
			response.Error(c, http.StatusForbidden, "无权限访问")
			c.Abort()
			return
		}

		// 从缓存获取用户权限
		permissions, err := getUserPermissions(uid)
		if err != nil {
			logger.Errorf("获取用户权限失败: %v", err)
			response.Error(c, http.StatusInternalServerError, "获取权限失败")
			c.Abort()
			return
		}

		// 检查是否有该权限
		for _, p := range permissions {
			if p == requiredPermission {
				c.Next()
				return
			}
		}

		logger.Warnf("用户 %d 无权限访问: %s", uid, requiredPermission)
		response.Error(c, http.StatusForbidden, "无权限访问")
		c.Abort()
	}
}

// AdminOnly 仅管理员可访问
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleIDs, exists := c.Get("role_ids")
		if !exists {
			response.Error(c, http.StatusForbidden, "无权限访问")
			c.Abort()
			return
		}

		ids, ok := roleIDs.([]uint)
		if !ok || len(ids) == 0 {
			response.Error(c, http.StatusForbidden, "无权限访问")
			c.Abort()
			return
		}

		// 通过数据库查询角色编码，而非硬编码 ID（与 data_permission.go 保持一致）
		if isAdmin(ids) {
			c.Next()
			return
		}

		response.Error(c, http.StatusForbidden, "需要管理员权限")
		c.Abort()
	}
}

// getUserPermissions 从缓存或数据库获取用户权限
func getUserPermissions(userID uint) ([]string, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("user:permissions:%d", userID)

	// 尝试从 Redis 缓存获取
	cached, err := database.RDB.SMembers(ctx, cacheKey).Result()
	if err == nil && len(cached) > 0 {
		return cached, nil
	}

	// 从数据库查询
	u, err := user.FindByID(userID)
	if err != nil {
		return nil, err
	}

	permSet := make(map[string]bool)
	for _, role := range u.Roles {
		menus, err := user.FindMenusByRoleID(role.ID)
		if err != nil {
			continue
		}
		for _, menu := range menus {
			if menu.Permission != "" {
				permSet[menu.Permission] = true
			}
		}
	}

	permissions := make([]string, 0, len(permSet))
	for p := range permSet {
		permissions = append(permissions, p)
	}

	// 缓存到 Redis（5分钟过期）
	if len(permissions) > 0 {
		members := make([]interface{}, len(permissions))
		for i, p := range permissions {
			members[i] = p
		}
		database.RDB.SAdd(ctx, cacheKey, members...)
		database.RDB.Expire(ctx, cacheKey, 300e9) // 5分钟
	}

	return permissions, nil
}

// ClearUserPermissionCache 清除用户权限缓存
func ClearUserPermissionCache(userID uint) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("user:permissions:%d", userID)
	database.RDB.Del(ctx, cacheKey)
}
