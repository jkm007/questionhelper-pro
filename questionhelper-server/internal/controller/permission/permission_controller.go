package permission

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/service/permission"
	"questionhelper-server/pkg/response"
)

type PermissionController struct{}

func NewPermissionController() *PermissionController {
	return &PermissionController{}
}

// ListButtonPermissions 按钮权限列表
// @Summary      获取按钮权限列表
// @Description  获取按钮权限列表，支持分页
// @Tags         按钮权限
// @Accept       json
// @Produce      json
// @Param        page       query     int  false  "页码"
// @Param        page_size  query     int  false  "每页数量"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/button-permissions [get]
// @Security     BearerAuth
func (ctrl *PermissionController) ListButtonPermissions(c *gin.Context) {
	var req dto.ButtonPermissionListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	list, total, err := permission.ListButtonPermissions(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// GetButtonPermission 获取按钮权限详情
// @Summary      获取按钮权限详情
// @Description  根据ID获取按钮权限详情
// @Tags         按钮权限
// @Accept       json
// @Produce      json
// @Param        id  path      uint  true  "权限ID"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "无效的ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/button-permissions/{id} [get]
// @Security     BearerAuth
func (ctrl *PermissionController) GetButtonPermission(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的ID")
		return
	}

	info, err := permission.GetButtonPermission(uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, info)
}

// CreateButtonPermission 创建按钮权限
// @Summary      创建按钮权限
// @Description  创建一条新的按钮权限
// @Tags         按钮权限
// @Accept       json
// @Produce      json
// @Param        req  body      dto.CreateButtonPermissionRequest  true  "权限数据"
// @Success      200  {object}  response.Response  "创建成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/button-permissions [post]
// @Security     BearerAuth
func (ctrl *PermissionController) CreateButtonPermission(c *gin.Context) {
	var req dto.CreateButtonPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := permission.CreateButtonPermission(&req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "创建成功", nil)
}

// UpdateButtonPermission 更新按钮权限
// @Summary      更新按钮权限
// @Description  根据ID更新按钮权限
// @Tags         按钮权限
// @Accept       json
// @Produce      json
// @Param        id   path      uint                                true  "权限ID"
// @Param        req  body      dto.UpdateButtonPermissionRequest   true  "权限数据"
// @Success      200  {object}  response.Response  "更新成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/button-permissions/{id} [put]
// @Security     BearerAuth
func (ctrl *PermissionController) UpdateButtonPermission(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的ID")
		return
	}

	var req dto.UpdateButtonPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := permission.UpdateButtonPermission(uint(id), &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "更新成功", nil)
}

// DeleteButtonPermission 删除按钮权限
// @Summary      删除按钮权限
// @Description  根据ID删除按钮权限
// @Tags         按钮权限
// @Accept       json
// @Produce      json
// @Param        id  path      uint  true  "权限ID"
// @Success      200  {object}  response.Response  "删除成功"
// @Failure      400  {object}  response.Response  "无效的ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/button-permissions/{id} [delete]
// @Security     BearerAuth
func (ctrl *PermissionController) DeleteButtonPermission(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的ID")
		return
	}

	if err := permission.DeleteButtonPermission(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "删除成功", nil)
}

// ListMenuButtonPermissions 获取菜单下的按钮权限
// @Summary      获取菜单下的按钮权限
// @Description  根据菜单ID获取其下的所有按钮权限
// @Tags         按钮权限
// @Accept       json
// @Produce      json
// @Param        menu_id  path      uint  true  "菜单ID"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "无效的菜单ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/button-permissions/menu/{menu_id} [get]
// @Security     BearerAuth
func (ctrl *PermissionController) ListMenuButtonPermissions(c *gin.Context) {
	menuID, err := strconv.ParseUint(c.Param("menu_id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的菜单ID")
		return
	}

	list, err := permission.ListButtonPermissionsByMenuID(uint(menuID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, list)
}
