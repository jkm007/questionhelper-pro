package router

import (
	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/controller/user"
)

func SetupUserRoutes(r *gin.RouterGroup, userCtrl *user.UserController, profileCtrl *user.ProfileController, tagCtrl *user.TagController, applyCtrl *user.ApplyController) {
	// 个人中心
	u := r.Group("/user")
	{
		u.GET("/profile", userCtrl.GetProfile)
		u.PUT("/profile", userCtrl.UpdateProfile)
		u.POST("/avatar", userCtrl.UploadAvatar)

		// 隐私设置
		u.GET("/privacy", profileCtrl.GetPrivacy)
		u.PUT("/privacy", profileCtrl.UpdatePrivacy)

		// 绑定手机/邮箱
		u.POST("/bind-phone", profileCtrl.BindPhone)
		u.POST("/bind-email", profileCtrl.BindEmail)

		// 实名认证
		u.GET("/realname", profileCtrl.GetRealName)
		u.POST("/realname", profileCtrl.SubmitRealName)

		// 第三方账号
		u.GET("/oauth", profileCtrl.GetOAuthBindings)
		u.DELETE("/oauth/:provider", profileCtrl.UnbindOAuth)

		// 收藏
		u.GET("/favorites", userCtrl.GetFavorites)
		u.POST("/favorites", userCtrl.AddFavorite)
		u.DELETE("/favorites/:id", userCtrl.RemoveFavorite)
	}

	// 兼容前端 /users/me 接口
	users := r.Group("/users")
	{
		users.GET("/me", userCtrl.GetMe)
	}

	// 角色申请
	SetupApplicationRoutes(r, applyCtrl)
}

func SetupAdminUserRoutes(r *gin.RouterGroup, userCtrl *user.UserController, profileCtrl *user.ProfileController) {
	// 用户管理
	u := r.Group("/users")
	{
		u.GET("", userCtrl.ListUsers)
		u.GET("/:id", userCtrl.GetUser)
		u.POST("", userCtrl.CreateUser)
		u.PUT("/:id", userCtrl.UpdateUser)
		u.DELETE("/:id", userCtrl.DeleteUser)
		u.PUT("/:id/status", userCtrl.UpdateUserStatus)
		u.POST("/:id/reset-password", userCtrl.ResetPassword)
		u.PUT("/:id/roles", userCtrl.AssignRoles)

		// 批量操作
		u.POST("/batch-status", userCtrl.BatchUpdateStatus)
		u.POST("/batch-delete", userCtrl.BatchDeleteUsers)
		u.POST("/batch-roles", userCtrl.BatchAssignRoles)

		// 导出
		u.GET("/export", userCtrl.ExportUsers)
	}

	// 角色管理
	role := r.Group("/roles")
	{
		role.GET("", userCtrl.ListRoles)
		role.GET("/:id", userCtrl.GetRole)
		role.POST("", userCtrl.CreateRole)
		role.PUT("/:id", userCtrl.UpdateRole)
		role.DELETE("/:id", userCtrl.DeleteRole)

		// 角色菜单管理
		role.GET("/:id/menus", userCtrl.ListRoleMenus)
		role.PUT("/:id/menus", userCtrl.AssignRoleMenus)

		// 角色权限管理
		role.GET("/:id/permissions", userCtrl.ListRolePermissions)
		role.PUT("/:id/permissions", userCtrl.AssignRolePermissions)
	}

	// 权限管理
	perm := r.Group("/permissions")
	{
		perm.GET("", userCtrl.ListPermissions)
		perm.GET("/tree", userCtrl.GetPermissionTree)
	}

	// 实名认证审核
	realname := r.Group("/realnames")
	{
		realname.GET("", profileCtrl.ListRealNames)
		realname.PUT("/:id/review", profileCtrl.ReviewRealName)
	}
}
