package router

import (
	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/controller/file"
)

// SetupFileRoutes 用户文件路由
func SetupFileRoutes(r *gin.RouterGroup, ctrl *file.FileController) {
	// 兼容前端 /files 接口
	f := r.Group("/files")
	{
		f.POST("", ctrl.UploadFile)
		f.DELETE("/:id", ctrl.DeleteFile)
	}

	// 文件上传
	upload := r.Group("/file/upload")
	{
		upload.POST("/image", ctrl.UploadImage)
		upload.POST("/batch", ctrl.BatchUpload)
	}

	// 文件下载
	r.GET("/file/:id/download", ctrl.DownloadFile)

	// 缩略图
	r.GET("/file/:id/thumbnail/:size", ctrl.GetThumbnail)

	// 文件引用管理
	ref := r.Group("/file/:id/reference")
	{
		ref.POST("", ctrl.AddReference)
		ref.DELETE("/:refId", ctrl.DeleteReference)
	}
	r.GET("/file/:id/references", ctrl.GetReferences)

	// 文件访问日志
	r.GET("/file/:id/access-logs", ctrl.GetAccessLogs)
}

// SetupAdminFileRoutes 管理员文件路由
func SetupAdminFileRoutes(r *gin.RouterGroup, ctrl *file.FileAdminController) {
	// 防盗链管理
	hotlink := r.Group("/file/hotlink-rules")
	{
		hotlink.GET("", ctrl.ListHotlinkRules)
		hotlink.POST("", ctrl.CreateHotlinkRule)
		hotlink.PUT("/:id", ctrl.UpdateHotlinkRule)
		hotlink.DELETE("/:id", ctrl.DeleteHotlinkRule)
	}

	// 水印管理
	watermark := r.Group("/file/watermark-configs")
	{
		watermark.GET("", ctrl.ListWatermarkConfigs)
		watermark.POST("", ctrl.CreateWatermarkConfig)
		watermark.PUT("/:id", ctrl.UpdateWatermarkConfig)
	}

	// 孤立文件清理
	cleanup := r.Group("/file/cleanup")
	{
		cleanup.POST("/orphan", ctrl.CleanupOrphanFiles)
		cleanup.GET("/logs", ctrl.GetCleanupLogs)
	}

	// 存储统计
	r.GET("/file/storage/statistics", ctrl.GetStorageStatistics)
}
