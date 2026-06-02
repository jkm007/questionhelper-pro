package user

import (
	"github.com/gin-gonic/gin"
	"questionhelper-server/pkg/response"
)

// ListPermissions 获取权限列表
// @Summary      获取权限列表
// @Description  获取系统权限列表（管理员）
// @Tags         角色管理
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/permissions [get]
// @Security     BearerAuth
func (ctrl *UserController) ListPermissions(c *gin.Context) {
	// TODO: 实现权限列表查询
	response.Error(c, 501, "接口未实现")
}

// GetPermissionTree 获取权限树
// @Summary      获取权限树
// @Description  获取系统权限树结构（管理员）
// @Tags         角色管理
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/permissions/tree [get]
// @Security     BearerAuth
func (ctrl *UserController) GetPermissionTree(c *gin.Context) {
	// TODO: 实现权限树查询
	response.Error(c, 501, "接口未实现")
}
