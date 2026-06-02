package router

import (
	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/controller/system"
)

func SetupAdminSystemRoutes(r *gin.RouterGroup, ctrl *system.SystemController) {
	// 系统设置 (原有)
	settings := r.Group("/settings")
	{
		settings.GET("", ctrl.GetSettings)
		settings.PUT("", ctrl.UpdateSettings)
	}

	// 操作日志 (原有)
	logs := r.Group("/logs")
	{
		logs.GET("/operation", ctrl.ListOperationLogs)
		logs.GET("/login", ctrl.ListLoginLogs)
	}

	// 分类设置
	settings.GET("/class", ctrl.GetClassSettings)
	settings.PUT("/class", ctrl.UpdateClassSettings)
	settings.GET("/resource", ctrl.GetResourceSettings)
	settings.PUT("/resource", ctrl.UpdateResourceSettings)

	// 系统日志
	adminLogs := r.Group("/logs")
	{
		adminLogs.GET("/system", ctrl.ListSystemLogs)
		adminLogs.GET("/error", ctrl.ListErrorLogs)
		adminLogs.GET("/search", ctrl.SearchLogs)
		adminLogs.POST("/archive", ctrl.ArchiveLogs)
		adminLogs.GET("/stats", ctrl.GetLogStats)
	}

	// 数据备份
	backup := r.Group("/backup")
	{
		backup.POST("/create", ctrl.CreateBackup)
		backup.GET("/list", ctrl.ListBackupRecords)
		backup.POST("/:id/restore", ctrl.RestoreBackup)
		backup.DELETE("/:id", ctrl.DeleteBackup)
		backup.GET("/config", ctrl.GetBackupConfigs)
		backup.PUT("/config", ctrl.UpdateBackupConfig)
	}

	// 功能开关
	features := r.Group("/features")
	{
		features.GET("", ctrl.ListFeatureFlags)
		features.PUT("/:key", ctrl.UpdateFeatureFlag)
	}

	// 安全配置
	security := r.Group("/security")
	{
		security.GET("", ctrl.GetSecurityConfigs)
		security.PUT("", ctrl.UpdateSecurityConfigs)
	}

	// 存储配置
	storage := r.Group("/storage")
	{
		storage.GET("", ctrl.ListStorageConfigs)
		storage.POST("", ctrl.CreateStorageConfig)
		storage.PUT("/:id", ctrl.UpdateStorageConfig)
	}

	// 邮件配置
	email := r.Group("/email")
	{
		email.GET("/config", ctrl.GetEmailConfig)
		email.PUT("/config", ctrl.UpdateEmailConfig)
		email.GET("/templates", ctrl.ListEmailTemplates)
		email.POST("/templates", ctrl.CreateEmailTemplate)
	}

	// 短信配置
	sms := r.Group("/sms")
	{
		sms.GET("/config", ctrl.GetSMSConfig)
		sms.PUT("/config", ctrl.UpdateSMSConfig)
		sms.GET("/templates", ctrl.ListSMSTemplates)
		sms.POST("/templates", ctrl.CreateSMSTemplate)
	}

	// 缓存管理
	cache := r.Group("/cache")
	{
		cache.GET("/stats", ctrl.GetCacheStats)
		cache.POST("/clear", ctrl.ClearCache)
	}

	// 主题配置
	theme := r.Group("/theme")
	{
		theme.GET("", ctrl.GetThemeConfig)
		theme.PUT("", ctrl.UpdateThemeConfig)
	}

	// 告警管理
	alerts := r.Group("/alerts")
	{
		alerts.GET("/rules", ctrl.ListAlertRules)
		alerts.POST("/rules", ctrl.CreateAlertRule)
		alerts.PUT("/rules/:id", ctrl.UpdateAlertRule)
		alerts.GET("/records", ctrl.ListAlertRecords)
	}
}
