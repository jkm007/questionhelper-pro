package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"questionhelper-server/pkg/database"
)

// DataPermissionType 数据权限类型
type DataPermissionType int

const (
	DataPermissionAll     DataPermissionType = 1 // 全部数据
	DataPermissionCustom  DataPermissionType = 2 // 自定义数据
	DataPermissionDept    DataPermissionType = 3 // 部门数据
	DataPermissionDeptSub DataPermissionType = 4 // 部门及以下数据
	DataPermissionSelf    DataPermissionType = 5 // 仅本人数据
)

// DataPermissionMiddleware 数据权限中间件
func DataPermissionMiddleware(tableName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.Next()
			return
		}

		// 获取用户角色
		roleIDs, _ := c.Get("role_ids")

		// 管理员拥有全部权限
		if isAdmin(roleIDs.([]uint)) {
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
func isAdmin(roleIDs []uint) bool {
	var count int64
	database.DB.Table("roles").
		Where("id IN ? AND code IN ?", roleIDs, []string{"super_admin", "admin"}).
		Count(&count)
	return count > 0
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
