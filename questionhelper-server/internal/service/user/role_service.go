package user

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	userRepo "questionhelper-server/internal/repository/user"
	"questionhelper-server/pkg/logger"
)

// ListRoles 角色列表
func ListRoles(req *dto.RoleListRequest) ([]dto.RoleInfo, int64, error) {
	roles, total, err := userRepo.ListRoles(req)
	if err != nil {
		return nil, 0, fmt.Errorf("查询角色列表失败: %w", err)
	}

	list := make([]dto.RoleInfo, 0, len(roles))
	for _, r := range roles {
		list = append(list, dto.RoleInfo{
			ID:          r.ID,
			Name:        r.Name,
			Code:        r.Code,
			Description: r.Description,
		})
	}

	return list, total, nil
}

// GetRole 获取角色详情
func GetRole(id uint) (*dto.RoleInfo, error) {
	role, err := userRepo.FindRoleByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("角色不存在")
		}
		return nil, fmt.Errorf("查询角色失败: %w", err)
	}

	return &dto.RoleInfo{
		ID:          role.ID,
		Name:        role.Name,
		Code:        role.Code,
		Description: role.Description,
	}, nil
}

// CreateRole 创建角色
func CreateRole(req *dto.CreateRoleRequest) error {
	role := &model.Role{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Status:      1,
	}

	// 关联菜单
	if len(req.MenuIDs) > 0 {
		menus := make([]model.Menu, 0, len(req.MenuIDs))
		for _, menuID := range req.MenuIDs {
			menus = append(menus, model.Menu{ID: menuID})
		}
		role.Menus = menus
	}

	if err := userRepo.CreateRole(role); err != nil {
		return fmt.Errorf("创建角色失败: %w", err)
	}

	logger.Infof("创建角色成功: %s", req.Name)
	return nil
}

// UpdateRole 更新角色
func UpdateRole(id uint, req *dto.UpdateRoleRequest) error {
	role, err := userRepo.FindRoleByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("角色不存在")
		}
		return fmt.Errorf("查询角色失败: %w", err)
	}

	if req.Name != "" {
		role.Name = req.Name
	}
	if req.Description != "" {
		role.Description = req.Description
	}
	if req.Status != nil {
		role.Status = *req.Status
	}

	// 更新菜单关联
	if req.MenuIDs != nil {
		menus := make([]model.Menu, 0, len(req.MenuIDs))
		for _, menuID := range req.MenuIDs {
			menus = append(menus, model.Menu{ID: menuID})
		}
		role.Menus = menus
	}

	if err := userRepo.UpdateRole(role); err != nil {
		return fmt.Errorf("更新角色失败: %w", err)
	}

	logger.Infof("更新角色 %d 成功", id)
	return nil
}

// DeleteRole 删除角色
func DeleteRole(id uint) error {
	_, err := userRepo.FindRoleByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("角色不存在")
		}
		return fmt.Errorf("查询角色失败: %w", err)
	}

	if err := userRepo.DeleteRoleByID(id); err != nil {
		return fmt.Errorf("删除角色失败: %w", err)
	}

	logger.Infof("删除角色 %d 成功", id)
	return nil
}

// GetRoleMenus 获取角色的菜单权限
func GetRoleMenus(roleID uint) ([]model.Menu, error) {
	menus, err := userRepo.FindMenusByRoleID(roleID)
	if err != nil {
		return nil, fmt.Errorf("查询角色菜单失败: %w", err)
	}
	return menus, nil
}
