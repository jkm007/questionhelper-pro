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
// @Summary      获取当前用户的菜单
// @Description  获取当前登录用户的菜单列表
// @Tags         菜单管理
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /menu/user [get]
// @Security     BearerAuth
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
// @Summary      获取当前用户的路由
// @Description  获取当前登录用户的路由列表（前端格式）
// @Tags         菜单管理
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /menus/routes [get]
// @Security     BearerAuth
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
// @Summary      获取当前用户的按钮权限
// @Description  获取当前登录用户的按钮权限列表
// @Tags         菜单管理
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /menu/buttons [get]
// @Security     BearerAuth
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
// @Summary      获取菜单列表
// @Description  获取所有菜单列表（管理员）
// @Tags         菜单管理
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/menus [get]
// @Security     BearerAuth
func (ctrl *MenuController) ListMenus(c *gin.Context) {
	menus, err := user.ListMenus()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, menus)
}

// GetMenuTree 菜单树
// @Summary      获取菜单树
// @Description  获取菜单树形结构（管理员）
// @Tags         菜单管理
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/menus/tree [get]
// @Security     BearerAuth
func (ctrl *MenuController) GetMenuTree(c *gin.Context) {
	menus, err := user.GetMenuTree()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, menus)
}

// GetMenu 获取菜单详情
// @Summary      获取菜单详情
// @Description  根据ID获取菜单详情（管理员）
// @Tags         菜单管理
// @Accept       json
// @Produce      json
// @Param        id  path      uint  true  "菜单ID"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "无效的菜单ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/menus/{id} [get]
// @Security     BearerAuth
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
// @Summary      创建菜单
// @Description  创建一个新的菜单项（管理员）
// @Tags         菜单管理
// @Accept       json
// @Produce      json
// @Param        menu  body      model.Menu  true  "菜单数据"
// @Success      200  {object}  response.Response  "创建成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/menus [post]
// @Security     BearerAuth
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
// @Summary      更新菜单
// @Description  根据ID更新菜单信息（管理员）
// @Tags         菜单管理
// @Accept       json
// @Produce      json
// @Param        id    path      uint        true  "菜单ID"
// @Param        menu  body      model.Menu  true  "菜单数据"
// @Success      200  {object}  response.Response  "更新成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/menus/{id} [put]
// @Security     BearerAuth
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
// @Summary      删除菜单
// @Description  根据ID删除菜单（管理员）
// @Tags         菜单管理
// @Accept       json
// @Produce      json
// @Param        id  path      uint  true  "菜单ID"
// @Success      200  {object}  response.Response  "删除成功"
// @Failure      400  {object}  response.Response  "无效的菜单ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/menus/{id} [delete]
// @Security     BearerAuth
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
