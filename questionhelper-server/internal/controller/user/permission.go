package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/service/user"
	"questionhelper-server/pkg/response"
)

// ListPermissions 获取权限列表
// @Summary      获取权限列表
// @Description  获取系统权限列表（管理员）
// @Tags         权限管理
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response{data=[]dto.PermissionInfo}  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/permissions [get]
// @Security     BearerAuth
func (ctrl *UserController) ListPermissions(c *gin.Context) {
	list, err := user.ListPermissions()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, list)
}

// GetPermissionTree 获取权限树
// @Summary      获取权限树
// @Description  获取系统权限树结构，按类型分组（管理员）
// @Tags         权限管理
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response{data=[]dto.PermissionTreeNode}  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/permissions/tree [get]
// @Security     BearerAuth
func (ctrl *UserController) GetPermissionTree(c *gin.Context) {
	tree, err := user.GetPermissionTree()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, tree)
}
