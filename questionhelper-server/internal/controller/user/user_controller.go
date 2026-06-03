package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/service/user"
	"questionhelper-server/pkg/response"
)

type UserController struct{}

func NewUserController() *UserController {
	return &UserController{}
}

// GetProfile 获取个人信息
// @Summary      获取个人信息
// @Description  获取当前登录用户的个人信息
// @Tags         个人中心
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response{data=dto.UserInfo}  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /user/profile [get]
// @Security     BearerAuth
func (ctrl *UserController) GetProfile(c *gin.Context) {
	ctrl.getProfile(c)
}

// GetMe 获取个人信息（兼容前端 /users/me 接口）
// @Summary      获取个人信息
// @Description  获取当前登录用户的个人信息（兼容前端 /users/me 接口）
// @Tags         个人中心
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response{data=dto.UserInfo}  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /users/me [get]
// @Security     BearerAuth
func (ctrl *UserController) GetMe(c *gin.Context) {
	ctrl.getProfile(c)
}

func (ctrl *UserController) getProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	info, err := user.GetProfile(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, info)
}

// UpdateProfile 更新个人信息
// @Summary      更新个人信息
// @Description  更新当前登录用户的个人信息
// @Tags         个人中心
// @Accept       json
// @Produce      json
// @Param        req  body      dto.UpdateProfileRequest  true  "更新信息"
// @Success      200  {object}  response.Response  "更新成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /user/profile [put]
// @Security     BearerAuth
func (ctrl *UserController) UpdateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := user.UpdateProfile(userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "更新成功", nil)
}

// ChangePassword 修改密码
// @Summary      修改密码
// @Description  修改当前登录用户的密码
// @Tags         个人中心
// @Accept       json
// @Produce      json
// @Param        req  body      dto.ChangePasswordRequest  true  "密码信息"
// @Success      200  {object}  response.Response  "密码修改成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /user/password [put]
// @Security     BearerAuth
func (ctrl *UserController) ChangePassword(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := user.ChangePassword(userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "密码修改成功", nil)
}

// UploadAvatar 上传头像
// @Summary      上传头像
// @Description  上传并更新用户头像
// @Tags         个人中心
// @Accept       multipart/form-data
// @Produce      json
// @Param        file  formData  file  true  "头像文件"
// @Success      200   {object}  response.Response  "上传成功"
// @Failure      500   {object}  response.Response  "服务器内部错误"
// @Router       /user/avatar [post]
// @Security     BearerAuth
func (ctrl *UserController) UploadAvatar(c *gin.Context) {
	// TODO: 实现文件上传
	response.Error(c, 501, "接口未实现")
}

// RealNameAuth 实名认证
// @Summary      实名认证
// @Description  提交实名认证信息
// @Tags         个人中心
// @Accept       json
// @Produce      json
// @Param        req  body      dto.RealNameAuthRequest  true  "实名认证信息"
// @Success      200  {object}  response.Response  "实名认证成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /user/realname-auth [post]
// @Security     BearerAuth
func (ctrl *UserController) RealNameAuth(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.RealNameAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := user.RealNameAuth(userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "实名认证成功", nil)
}

// GetFavorites 获取收藏列表
// @Summary      获取收藏列表
// @Description  获取当前用户的收藏列表
// @Tags         个人中心
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /user/favorites [get]
// @Security     BearerAuth
func (ctrl *UserController) GetFavorites(c *gin.Context) {
	// TODO: 实现收藏功能
	response.Error(c, 501, "接口未实现")
}

// AddFavorite 添加收藏
// @Summary      添加收藏
// @Description  添加题目到收藏列表
// @Tags         个人中心
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response  "收藏成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /user/favorites [post]
// @Security     BearerAuth
func (ctrl *UserController) AddFavorite(c *gin.Context) {
	// TODO: 实现收藏功能
	response.Error(c, 501, "接口未实现")
}

// RemoveFavorite 取消收藏
// @Summary      取消收藏
// @Description  取消收藏指定题目
// @Tags         个人中心
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "收藏ID"
// @Success      200  {object}  response.Response  "取消收藏成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /user/favorites/{id} [delete]
// @Security     BearerAuth
func (ctrl *UserController) RemoveFavorite(c *gin.Context) {
	// TODO: 实现收藏功能
	response.Error(c, 501, "接口未实现")
}

// ListUsers 用户列表（管理员）
// @Summary      用户列表
// @Description  分页查询用户列表（管理员）
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        page       query     int     false  "页码"       default(1)
// @Param        page_size  query     int     false  "每页数量"   default(10)
// @Param        keyword    query     string  false  "搜索关键词"
// @Param        status     query     int     false  "用户状态"
// @Success      200  {object}  response.Response{data=dto.PageResponse{list=[]dto.UserInfo}}  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/users [get]
// @Security     BearerAuth
func (ctrl *UserController) ListUsers(c *gin.Context) {
	var req dto.UserListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	list, total, err := user.ListUsers(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// GetUser 获取用户详情（管理员）
// @Summary      获取用户详情
// @Description  根据ID获取用户详情（管理员）
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "用户ID"
// @Success      200  {object}  response.Response{data=dto.UserInfo}  "成功"
// @Failure      400  {object}  response.Response  "无效的用户ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/users/{id} [get]
// @Security     BearerAuth
func (ctrl *UserController) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	info, err := user.GetUser(uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, info)
}

// CreateUser 创建用户（管理员）
// @Summary      创建用户
// @Description  创建新用户（管理员）
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        req  body      dto.CreateUserRequest  true  "用户信息"
// @Success      200  {object}  response.Response  "创建成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/users [post]
// @Security     BearerAuth
func (ctrl *UserController) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := user.CreateUser(&req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "创建成功", nil)
}

// UpdateUser 更新用户（管理员）
// @Summary      更新用户
// @Description  更新用户信息（管理员）
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint                   true  "用户ID"
// @Param        req  body      dto.UpdateUserRequest   true  "用户信息"
// @Success      200  {object}  response.Response  "更新成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/users/{id} [put]
// @Security     BearerAuth
func (ctrl *UserController) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := user.UpdateUser(uint(id), &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "更新成功", nil)
}

// DeleteUser 删除用户（管理员）
// @Summary      删除用户
// @Description  删除指定用户（管理员）
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "用户ID"
// @Success      200  {object}  response.Response  "删除成功"
// @Failure      400  {object}  response.Response  "无效的用户ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/users/{id} [delete]
// @Security     BearerAuth
func (ctrl *UserController) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	if err := user.DeleteUser(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "删除成功", nil)
}

// UpdateUserStatus 更新用户状态（管理员）
// @Summary      更新用户状态
// @Description  启用或禁用用户（管理员）
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint                true  "用户ID"
// @Param        req  body      dto.StatusRequest   true  "状态信息"
// @Success      200  {object}  response.Response  "状态更新成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/users/{id}/status [put]
// @Security     BearerAuth
func (ctrl *UserController) UpdateUserStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	var req dto.StatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := user.UpdateUserStatus(uint(id), req.Status); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "状态更新成功", nil)
}

// BatchUpdateStatus 批量更新用户状态
// @Summary      批量更新用户状态
// @Description  批量启用或禁用用户（管理员）
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        req  body      dto.BatchStatusRequest  true  "批量状态请求"
// @Success      200  {object}  response.Response  "批量更新成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/users/batch-status [post]
// @Security     BearerAuth
func (ctrl *UserController) BatchUpdateStatus(c *gin.Context) {
	var req dto.BatchStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := user.BatchUpdateStatus(&req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "批量更新成功", nil)
}

// BatchDeleteUsers 批量删除用户
// @Summary      批量删除用户
// @Description  批量删除多个用户（管理员）
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        req  body      dto.BatchDeleteRequest  true  "批量删除请求"
// @Success      200  {object}  response.Response  "批量删除成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/users/batch-delete [post]
// @Security     BearerAuth
func (ctrl *UserController) BatchDeleteUsers(c *gin.Context) {
	var req dto.BatchDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := user.BatchDeleteUsers(&req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "批量删除成功", nil)
}

// BatchAssignRoles 批量分配角色
// @Summary      批量分配角色
// @Description  批量为用户分配角色（管理员）
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        req  body      dto.BatchRoleRequest  true  "批量角色分配请求"
// @Success      200  {object}  response.Response  "批量分配角色成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/users/batch-roles [post]
// @Security     BearerAuth
func (ctrl *UserController) BatchAssignRoles(c *gin.Context) {
	var req dto.BatchRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := user.BatchAssignRoles(&req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "批量分配角色成功", nil)
}

// ResetPassword 重置用户密码
// @Summary      重置用户密码
// @Description  管理员重置指定用户的密码
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint                          true  "用户ID"
// @Param        req  body      dto.AdminResetPasswordRequest  true  "新密码信息"
// @Success      200  {object}  response.Response  "密码重置成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/users/{id}/reset-password [post]
// @Security     BearerAuth
func (ctrl *UserController) ResetPassword(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	var req dto.AdminResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := user.ResetPassword(uint(id), &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "密码重置成功", nil)
}

// ExportUsers 导出用户
// @Summary      导出用户
// @Description  导出用户列表为CSV文件（管理员）
// @Tags         用户管理
// @Accept       json
// @Produce      text/csv
// @Param        page       query     int     false  "页码"
// @Param        page_size  query     int     false  "每页数量"
// @Param        keyword    query     string  false  "搜索关键词"
// @Param        status     query     int     false  "用户状态"
// @Success      200  {string}  string  "CSV文件"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/users/export [get]
// @Security     BearerAuth
func (ctrl *UserController) ExportUsers(c *gin.Context) {
	var req dto.UserListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	buf, err := user.ExportUsers(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=users.csv")
	c.Data(http.StatusOK, "text/csv; charset=utf-8", buf.Bytes())
}

// AssignRoles 分配用户角色
// @Summary      分配用户角色
// @Description  为指定用户分配角色（管理员）
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "用户ID"
// @Param        req  body      object{role_ids=[]uint}  true  "角色ID列表"
// @Success      200  {object}  response.Response  "角色分配成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/users/{id}/roles [put]
// @Security     BearerAuth
func (ctrl *UserController) AssignRoles(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	var req struct {
		RoleIDs []uint `json:"role_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := user.AssignRoles(uint(id), req.RoleIDs); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "角色分配成功", nil)
}

// ListRoleMenus 获取角色菜单
// @Summary      获取角色菜单
// @Description  获取指定角色的菜单列表（管理员）
// @Tags         角色管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "角色ID"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "无效的角色ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/roles/{id}/menus [get]
// @Security     BearerAuth
func (ctrl *UserController) ListRoleMenus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的角色ID")
		return
	}

	menus, err := user.GetRoleMenus(uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, menus)
}

// ListRolePermissions 获取角色权限
// @Summary      获取角色权限
// @Description  获取指定角色的权限列表（管理员）
// @Tags         角色管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "角色ID"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "无效的角色ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/roles/{id}/permissions [get]
// @Security     BearerAuth
func (ctrl *UserController) ListRolePermissions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的角色ID")
		return
	}

	permissions, err := user.GetRolePermissions(uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, permissions)
}

