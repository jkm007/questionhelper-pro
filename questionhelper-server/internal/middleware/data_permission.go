package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DataPermissionType 数据权限类型（与数据库设计文档 02-权限管理 保持一致）
type DataPermissionType int

const (
	DataPermissionAll     DataPermissionType = 1 // 全部数据
	DataPermissionDeptSub DataPermissionType = 2 // 本部门及下级
	DataPermissionDept    DataPermissionType = 3 // 本部门
	DataPermissionSelf    DataPermissionType = 4 // 仅本人数据
)

// DataPermissionMiddleware 数据权限中间件
func DataPermissionMiddleware(tableName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.Next()
			return
		}

		// 管理员拥有全部权限
		if isAdmin(c) {
			c.Next()
			return
		}

		// 设置数据权限条件
		condition := fmt.Sprintf("%s.created_by = ?", tableName)
		c.Set("data_permission_condition", condition)
		c.Set("data_permission_args", []interface{}{userID})

		c.Next()
	}
}

// isAdmin 判断是否为管理员（super_admin 或 admin 角色）
// 优先从 JWT 上下文中的 role_codes 判断，避免每次查数据库
func isAdmin(c *gin.Context) bool {
	if codes, exists := c.Get("role_codes"); exists {
		if roleCodes, ok := codes.([]string); ok {
			for _, code := range roleCodes {
				if code == "super_admin" || code == "admin" {
					return true
				}
			}
		}
	}
	return false
}

// ApplyDataPermission 应用数据权限到查询
func ApplyDataPermission(db *gorm.DB, c *gin.Context, tableName string) *gorm.DB {
	condition, exists := c.Get("data_permission_condition")
	if !exists {
		return db
	}

	args, _ := c.Get("data_permission_args")
	return db.Where(condition.(string), args.([]interface{})...)
}
