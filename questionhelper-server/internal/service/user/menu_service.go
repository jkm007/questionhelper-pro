package user

import (
	"errors"
	"fmt"

	"questionhelper-server/internal/model"
	userRepo "questionhelper-server/internal/repository/user"
	"questionhelper-server/pkg/logger"
)

// RouteItem 前端路由项
type RouteItem struct {
	Name      string     `json:"name"`
	Path      string     `json:"path"`
	Component string     `json:"component,omitempty"`
	Redirect  string     `json:"redirect,omitempty"`
	Meta      *RouteMeta `json:"meta,omitempty"`
	Children  []RouteItem `json:"children,omitempty"`
}

// RouteMeta 路由元信息
type RouteMeta struct {
	Title     string `json:"title"`
	Icon      string `json:"icon,omitempty"`
	Hidden    bool   `json:"hidden,omitempty"`
	KeepAlive bool   `json:"keepAlive,omitempty"`
	Affix     bool   `json:"affix,omitempty"`
}

// GetUserMenus 获取用户的菜单树（根据角色）
func GetUserMenus(userID uint) ([]model.Menu, error) {
	u, err := userRepo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	// 收集所有角色的菜单ID
	menuMap := make(map[uint]*model.Menu)
	for _, role := range u.Roles {
		menus, err := userRepo.FindMenusByRoleID(role.ID)
		if err != nil {
			continue
		}
		for i := range menus {
			menu := menus[i]
			if menu.Type != 3 { // 排除按钮类型
				menuMap[menu.ID] = &menu
			}
		}
	}

	// 构建菜单树
	allMenus := make([]model.Menu, 0, len(menuMap))
	for _, m := range menuMap {
		allMenus = append(allMenus, *m)
	}

	return buildMenuTree(allMenus, nil), nil
}

// GetUserRoutes 获取用户的路由树（前端格式）
func GetUserRoutes(userID uint) ([]RouteItem, error) {
	menus, err := GetUserMenus(userID)
	if err != nil {
		return nil, err
	}

	return transformToRoutes(menus), nil
}

// transformToRoutes 将菜单树转换为路由树
func transformToRoutes(menus []model.Menu) []RouteItem {
	routes := make([]RouteItem, 0, len(menus))
	for _, menu := range menus {
		route := RouteItem{
			Name:      menu.Name,
			Path:      menu.Path,
			Component: menu.Component,
			Redirect:  menu.Redirect,
			Meta: &RouteMeta{
				Title:     menu.Title,
				Icon:      menu.Icon,
				Hidden:    menu.Hidden,
				KeepAlive: false,
			},
		}

		if len(menu.Children) > 0 {
			route.Children = transformToRoutes(menu.Children)
		}

		routes = append(routes, route)
	}
	return routes
}

// GetUserButtons 获取用户的按钮权限
func GetUserButtons(userID uint) ([]string, error) {
	u, err := userRepo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	permMap := make(map[string]bool)
	for _, role := range u.Roles {
		menus, err := userRepo.FindMenusByRoleID(role.ID)
		if err != nil {
			continue
		}
		for _, menu := range menus {
			if menu.Type == 3 && menu.Permission != "" {
				permMap[menu.Permission] = true
			}
		}
	}

	permissions := make([]string, 0, len(permMap))
	for p := range permMap {
		permissions = append(permissions, p)
	}
	return permissions, nil
}

// GetMenuTree 获取完整菜单树（管理员）
func GetMenuTree() ([]model.Menu, error) {
	return userRepo.FindMenuTree()
}

// ListMenus 菜单列表（平铺）
func ListMenus() ([]model.Menu, error) {
	return userRepo.FindAllMenus()
}

// GetMenu 获取菜单详情
func GetMenu(id uint) (*model.Menu, error) {
	menus, err := userRepo.FindAllMenus()
	if err != nil {
		return nil, fmt.Errorf("查询菜单失败: %w", err)
	}
	for _, m := range menus {
		if m.ID == id {
			return &m, nil
		}
	}
	return nil, errors.New("菜单不存在")
}

// CreateMenu 创建菜单
func CreateMenu(menu *model.Menu) error {
	if err := userRepo.CreateMenu(menu); err != nil {
		return fmt.Errorf("创建菜单失败: %w", err)
	}
	logger.Infof("创建菜单成功: %s", menu.Name)
	return nil
}

// UpdateMenu 更新菜单
func UpdateMenu(menu *model.Menu) error {
	// 检查菜单是否存在
	_, err := userRepo.FindAllMenus()
	if err != nil {
		return fmt.Errorf("查询菜单失败: %w", err)
	}

	if err := userRepo.UpdateMenu(menu); err != nil {
		return fmt.Errorf("更新菜单失败: %w", err)
	}
	logger.Infof("更新菜单 %d 成功", menu.ID)
	return nil
}

// DeleteMenu 删除菜单
func DeleteMenu(id uint) error {
	// 检查是否有子菜单
	menus, err := userRepo.FindAllMenus()
	if err != nil {
		return fmt.Errorf("查询菜单失败: %w", err)
	}
	for _, m := range menus {
		if m.ParentID != nil && *m.ParentID == id {
			return errors.New("存在子菜单，无法删除")
		}
	}

	if err := userRepo.DeleteMenuByID(id); err != nil {
		return fmt.Errorf("删除菜单失败: %w", err)
	}
	logger.Infof("删除菜单 %d 成功", id)
	return nil
}

// buildMenuTree 构建菜单树
func buildMenuTree(menus []model.Menu, parentID *uint) []model.Menu {
	tree := make([]model.Menu, 0)
	for _, menu := range menus {
		if (parentID == nil && menu.ParentID == nil) ||
			(parentID != nil && menu.ParentID != nil && *menu.ParentID == *parentID) {
			children := buildMenuTree(menus, &menu.ID)
			menu.Children = children
			tree = append(tree, menu)
		}
	}
	return tree
}
