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
