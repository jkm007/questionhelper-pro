package menu

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/model"
	"questionhelper-server/internal/service/user"
	"questionhelper-server/pkg/response"
)

type MenuController struct{}

func NewMenuController() *MenuController {
	return &MenuController{}
}

// GetUserMenus 获取当前用户的菜单
func (ctrl *MenuController) GetUserMenus(c *gin.Context) {
	userID := c.GetUint("user_id")

	menus, err := user.GetUserMenus(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, menus)
}

// GetUserRoutes 获取当前用户的路由（前端格式）
func (ctrl *MenuController) GetUserRoutes(c *gin.Context) {
	userID := c.GetUint("user_id")

	routes, err := user.GetUserRoutes(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, routes)
}

// GetUserButtons 获取当前用户的按钮权限
func (ctrl *MenuController) GetUserButtons(c *gin.Context) {
	userID := c.GetUint("user_id")

	buttons, err := user.GetUserButtons(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, buttons)
}

// ListMenus 菜单列表
func (ctrl *MenuController) ListMenus(c *gin.Context) {
	menus, err := user.ListMenus()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, menus)
}

// GetMenuTree 菜单树
func (ctrl *MenuController) GetMenuTree(c *gin.Context) {
	menus, err := user.GetMenuTree()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, menus)
}

// GetMenu 获取菜单详情
func (ctrl *MenuController) GetMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的菜单ID")
		return
	}

	menu, err := user.GetMenu(uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, menu)
}

// CreateMenu 创建菜单
func (ctrl *MenuController) CreateMenu(c *gin.Context) {
	var menu model.Menu
	if err := c.ShouldBindJSON(&menu); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := user.CreateMenu(&menu); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "创建成功", nil)
}

// UpdateMenu 更新菜单
func (ctrl *MenuController) UpdateMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的菜单ID")
		return
	}

	var menu model.Menu
	if err := c.ShouldBindJSON(&menu); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	menu.ID = uint(id)

	if err := user.UpdateMenu(&menu); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "更新成功", nil)
}

// DeleteMenu 删除菜单
func (ctrl *MenuController) DeleteMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的菜单ID")
		return
	}

	if err := user.DeleteMenu(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "删除成功", nil)
}
