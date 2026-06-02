package key

import "fmt"

// 菜单相关缓存 Key（统一使用 qh: 前缀）

// UserMenuKey 用户菜单缓存 Key
func UserMenuKey(userID uint) string {
	return fmt.Sprintf("qh:menu:user:%d", userID)
}

// UserPermissionKey 用户权限缓存 Key
func UserPermissionKey(userID uint) string {
	return fmt.Sprintf("qh:permission:user:%d", userID)
}

// RoleMenuKey 角色菜单缓存 Key
func RoleMenuKey(roleID uint) string {
	return fmt.Sprintf("qh:menu:role:%d", roleID)
}

// MenuTreeKey 菜单树缓存 Key
const MenuTreeKey = "qh:menu:tree"

// MenuExpire 菜单缓存过期时间（24小时）
const MenuExpire = 86400
