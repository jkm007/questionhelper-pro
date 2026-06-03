package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/service/user"
	"questionhelper-server/pkg/response"
)

// ListRoles 角色列表
// @Summary      角色列表
// @Description  分页查询角色列表（管理员）
// @Tags         角色管理
// @Accept       json
// @Produce      json
// @Param        page       query     int     false  "页码"       default(1)
// @Param        page_size  query     int     false  "每页数量"   default(10)
// @Param        keyword    query     string  false  "搜索关键词"
// @Success      200  {object}  response.Response{data=dto.PageResponse{list=[]dto.RoleInfo}}  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/roles [get]
// @Security     BearerAuth
func (ctrl *UserController) ListRoles(c *gin.Context) {
	var req dto.RoleListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	list, total, err := user.ListRoles(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// GetRole 获取角色详情
// @Summary      获取角色详情
// @Description  根据ID获取角色详情（管理员）
// @Tags         角色管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "角色ID"
// @Success      200  {object}  response.Response{data=dto.RoleInfo}  "成功"
// @Failure      400  {object}  response.Response  "无效的角色ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/roles/{id} [get]
// @Security     BearerAuth
func (ctrl *UserController) GetRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的角色ID")
		return
	}

	role, err := user.GetRole(uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, role)
}

// CreateRole 创建角色
// @Summary      创建角色
// @Description  创建新角色（管理员）
// @Tags         角色管理
// @Accept       json
// @Produce      json
// @Param        req  body      dto.CreateRoleRequest  true  "角色信息"
// @Success      200  {object}  response.Response  "创建成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/roles [post]
// @Security     BearerAuth
func (ctrl *UserController) CreateRole(c *gin.Context) {
	var req dto.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := user.CreateRole(&req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "创建成功", nil)
}

// UpdateRole 更新角色
// @Summary      更新角色
// @Description  更新角色信息（管理员）
// @Tags         角色管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint                   true  "角色ID"
// @Param        req  body      dto.UpdateRoleRequest  true  "角色信息"
// @Success      200  {object}  response.Response  "更新成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/roles/{id} [put]
// @Security     BearerAuth
func (ctrl *UserController) UpdateRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的角色ID")
		return
	}

	var req dto.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := user.UpdateRole(uint(id), &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "更新成功", nil)
}

// DeleteRole 删除角色
// @Summary      删除角色
// @Description  删除指定角色（管理员）
// @Tags         角色管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "角色ID"
// @Success      200  {object}  response.Response  "删除成功"
// @Failure      400  {object}  response.Response  "无效的角色ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/roles/{id} [delete]
// @Security     BearerAuth
func (ctrl *UserController) DeleteRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的角色ID")
		return
	}

	if err := user.DeleteRole(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "删除成功", nil)
}

// AssignRoleMenus 分配角色菜单
// @Summary      分配角色菜单
// @Description  为指定角色分配菜单（管理员）
// @Tags         角色管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "角色ID"
// @Param        req  body      object{menu_ids=[]uint}  true  "菜单ID列表"
// @Success      200  {object}  response.Response  "分配成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/roles/{id}/menus [put]
// @Security     BearerAuth
func (ctrl *UserController) AssignRoleMenus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的角色ID")
		return
	}

	var req struct {
		MenuIDs []uint `json:"menu_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := user.AssignRoleMenus(uint(id), req.MenuIDs); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "分配成功", nil)
}

// AssignRolePermissions 分配角色权限
// @Summary      分配角色权限
// @Description  为指定角色分配权限（管理员）
// @Tags         角色管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "角色ID"
// @Param        req  body      object{permission_ids=[]uint}  true  "权限ID列表"
// @Success      200  {object}  response.Response  "分配成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/roles/{id}/permissions [put]
// @Security     BearerAuth
func (ctrl *UserController) AssignRolePermissions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的角色ID")
		return
	}

	var req struct {
		PermissionIDs []uint `json:"permission_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := user.AssignRolePermissions(uint(id), req.PermissionIDs); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "分配成功", nil)
}
