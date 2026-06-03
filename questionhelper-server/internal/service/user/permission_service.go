package user

import (
	"context"
	"fmt"

	"questionhelper-server/internal/dto"
	userRepo "questionhelper-server/internal/repository/user"
	"questionhelper-server/pkg/database"
	"questionhelper-server/pkg/logger"
)

// ListPermissions 权限列表
func ListPermissions() ([]dto.PermissionInfo, error) {
	permissions, err := userRepo.FindAllPermissions()
	if err != nil {
		return nil, fmt.Errorf("查询权限列表失败: %w", err)
	}

	list := make([]dto.PermissionInfo, 0, len(permissions))
	for _, p := range permissions {
		list = append(list, dto.PermissionInfo{
			ID:          p.ID,
			Name:        p.Name,
			Code:        p.Code,
			Description: p.Description,
			Type:        p.Type,
		})
	}

	return list, nil
}

// GetPermissionTree 按 Type 分组构建权限树
func GetPermissionTree() ([]dto.PermissionTreeNode, error) {
	permissions, err := userRepo.FindAllPermissions()
	if err != nil {
		return nil, fmt.Errorf("查询权限列表失败: %w", err)
	}

	// 按 Type 分组: 1=菜单权限, 2=按钮权限, 3=接口权限
	typeNames := map[int8]string{
		1: "菜单权限",
		2: "按钮权限",
		3: "接口权限",
	}

	groups := make(map[int8][]dto.PermissionTreeNode)
	for _, p := range permissions {
		groups[p.Type] = append(groups[p.Type], dto.PermissionTreeNode{
			PermissionInfo: dto.PermissionInfo{
				ID:          p.ID,
				Name:        p.Name,
				Code:        p.Code,
				Description: p.Description,
				Type:        p.Type,
			},
		})
	}

	tree := make([]dto.PermissionTreeNode, 0, len(typeNames))
	for _, t := range []int8{1, 2, 3} {
		children := groups[t]
		if children == nil {
			children = make([]dto.PermissionTreeNode, 0)
		}
		tree = append(tree, dto.PermissionTreeNode{
			PermissionInfo: dto.PermissionInfo{
				Name: typeNames[t],
				Type: t,
			},
			Children: children,
		})
	}

	return tree, nil
}

// GetRolePermissions 获取角色权限
func GetRolePermissions(roleID uint) ([]dto.PermissionInfo, error) {
	permissions, err := userRepo.FindPermissionsByRoleID(roleID)
	if err != nil {
		return nil, fmt.Errorf("查询角色权限失败: %w", err)
	}

	list := make([]dto.PermissionInfo, 0, len(permissions))
	for _, p := range permissions {
		list = append(list, dto.PermissionInfo{
			ID:          p.ID,
			Name:        p.Name,
			Code:        p.Code,
			Description: p.Description,
			Type:        p.Type,
		})
	}

	return list, nil
}

// AssignRoleMenus 分配角色菜单
func AssignRoleMenus(roleID uint, menuIDs []uint) error {
	if err := userRepo.AssignRoleMenus(roleID, menuIDs); err != nil {
		return fmt.Errorf("分配角色菜单失败: %w", err)
	}

	// 清除该角色下所有用户的权限缓存
	clearRoleUsersPermissionCache(roleID)

	logger.Infof("分配角色 %d 菜单成功, menuIDs: %v", roleID, menuIDs)
	return nil
}

// AssignRolePermissions 分配角色权限
func AssignRolePermissions(roleID uint, permissionIDs []uint) error {
	if err := userRepo.AssignRolePermissions(roleID, permissionIDs); err != nil {
		return fmt.Errorf("分配角色权限失败: %w", err)
	}

	// 清除该角色下所有用户的权限缓存
	clearRoleUsersPermissionCache(roleID)

	logger.Infof("分配角色 %d 权限成功, permissionIDs: %v", roleID, permissionIDs)
	return nil
}

// clearRoleUsersPermissionCache 清除角色关联的所有用户的权限缓存
func clearRoleUsersPermissionCache(roleID uint) {
	userIDs, err := userRepo.FindUserIDsByRoleID(roleID)
	if err != nil {
		logger.Errorf("查询角色 %d 关联用户失败: %v", roleID, err)
		return
	}

	ctx := context.Background()
	for _, uid := range userIDs {
		cacheKey := fmt.Sprintf("user:permissions:%d", uid)
		database.RDB.Del(ctx, cacheKey)
	}
}
