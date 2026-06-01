package user

import (
	"github.com/gin-gonic/gin"
	"questionhelper-server/pkg/response"
)

// ListPermissions 获取权限列表
func (ctrl *UserController) ListPermissions(c *gin.Context) {
	// TODO: 实现权限列表查询
	response.Error(c, 501, "接口未实现")
}

// GetPermissionTree 获取权限树
func (ctrl *UserController) GetPermissionTree(c *gin.Context) {
	// TODO: 实现权限树查询
	response.Error(c, 501, "接口未实现")
}
