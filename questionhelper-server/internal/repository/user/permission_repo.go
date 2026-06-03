package user

import (
	"questionhelper-server/internal/model"
	"questionhelper-server/pkg/database"
)

// FindAllPermissions 查询所有权限
func FindAllPermissions() ([]model.Permission, error) {
	var permissions []model.Permission
	err := database.DB.Order("type ASC, id ASC").Find(&permissions).Error
	return permissions, err
}

// FindPermissionsByRoleID 查询角色的权限
func FindPermissionsByRoleID(roleID uint) ([]model.Permission, error) {
	var role model.Role
	if err := database.DB.Preload("Permissions").First(&role, roleID).Error; err != nil {
		return nil, err
	}
	return role.Permissions, nil
}

// AssignRoleMenus 替换角色的菜单关联
func AssignRoleMenus(roleID uint, menuIDs []uint) error {
	var role model.Role
	if err := database.DB.First(&role, roleID).Error; err != nil {
		return err
	}

	menus := make([]model.Menu, 0, len(menuIDs))
	for _, menuID := range menuIDs {
		menus = append(menus, model.Menu{ID: menuID})
	}

	return database.DB.Model(&role).Association("Menus").Replace(menus)
}

// AssignRolePermissions 替换角色的权限关联
func AssignRolePermissions(roleID uint, permissionIDs []uint) error {
	var role model.Role
	if err := database.DB.First(&role, roleID).Error; err != nil {
		return err
	}

	permissions := make([]model.Permission, 0, len(permissionIDs))
	for _, permissionID := range permissionIDs {
		permissions = append(permissions, model.Permission{ID: permissionID})
	}

	return database.DB.Model(&role).Association("Permissions").Replace(permissions)
}

// FindUserIDsByRoleID 查询角色关联的所有用户ID
func FindUserIDsByRoleID(roleID uint) ([]uint, error) {
	var userIDs []uint
	err := database.DB.Table("user_roles").Where("role_id = ?", roleID).Pluck("user_id", &userIDs).Error
	return userIDs, err
}
