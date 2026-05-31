package router

import (
	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/controller/auth"
	"questionhelper-server/internal/middleware"
)

// SetupAuthRoutes 设置认证相关路由（公开接口）
func SetupAuthRoutes(r *gin.RouterGroup, ctrl *auth.AuthController) {
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", ctrl.Login)
		authGroup.POST("/register", ctrl.Register)
		authGroup.POST("/refresh", ctrl.RefreshToken)
		authGroup.GET("/captcha", ctrl.GetCaptcha)
		authGroup.POST("/captcha/verify", ctrl.VerifyCaptcha)
		authGroup.POST("/password/reset-request", ctrl.RequestPasswordReset)
		authGroup.POST("/password/reset", ctrl.ResetPassword)
		authGroup.POST("/sms/send", ctrl.SendSmsCode)
		authGroup.POST("/email/send", ctrl.SendEmailCode)
		authGroup.GET("/oauth/:provider/url", ctrl.GetOAuthURL)
		authGroup.POST("/oauth/:provider", ctrl.OAuthLogin)
	}
}

// SetupUserAuthRoutes 设置用户认证相关路由（需要认证）
func SetupUserAuthRoutes(r *gin.RouterGroup, ctrl *auth.AuthController) {
	userGroup := r.Group("/user")
	userGroup.Use(middleware.AuthMiddleware())
	{
		// 退出登录
		userGroup.POST("/logout", ctrl.Logout)
		userGroup.POST("/logout-all", ctrl.LogoutAll)

		// 修改密码
		userGroup.PUT("/password", ctrl.ChangePassword)

		// 设备管理
		userGroup.GET("/devices", ctrl.GetUserDevices)
		userGroup.DELETE("/devices/:id", ctrl.KickDevice)
		userGroup.DELETE("/devices", ctrl.KickAllDevices)

		// 安全日志
		userGroup.GET("/security/logs", ctrl.GetSecurityLogs)

		// 安全设置
		userGroup.PUT("/security/settings", ctrl.UpdateSecuritySettings)

		// 第三方账号
		userGroup.GET("/oauth/status", ctrl.GetOAuthStatus)
		userGroup.POST("/oauth/bind/:provider", ctrl.OAuthBind)
		userGroup.DELETE("/oauth/unbind/:provider", ctrl.OAuthUnbind)

		// 账号注销
		userGroup.POST("/account/deactivate", ctrl.DeactivateAccount)
		userGroup.POST("/account/deactivate/confirm", ctrl.ConfirmDeactivate)
		userGroup.POST("/account/cancel-deactivate", ctrl.CancelDeactivate)

		// 数据导出
		userGroup.GET("/account/export", ctrl.ExportUserData)
	}
}
