package router

import (
	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/controller/comment"
)

func SetupCommentRoutes(r *gin.RouterGroup, ctrl *comment.CommentController) {
	// 兼容前端 /comments 接口
	c := r.Group("/comments")
	{
		c.GET("", ctrl.ListComments)
		c.POST("", ctrl.CreateComment)
		c.PUT("/:id", ctrl.EditComment)
		c.DELETE("/:id", ctrl.DeleteComment)
		c.POST("/:id/like", ctrl.LikeComment)
		c.POST("/:id/report", ctrl.ReportComment)
		c.POST("/upload-image", ctrl.UploadImage)
	}

	// 表情包
	sticker := r.Group("/sticker")
	{
		sticker.GET("/list", ctrl.ListStickers)
		sticker.GET("/categories", ctrl.ListStickerCategories)
	}

	// 用户搜索(@功能)
	r.GET("/user/search", ctrl.SearchUsers)
}

// SetupAdminCommentRoutes 管理员评论路由
func SetupAdminCommentRoutes(r *gin.RouterGroup, ctrl *comment.CommentAdminController) {
	c := r.Group("/comments")
	{
		// 评论列表
		c.GET("", ctrl.ListComments)

		// 评论管理
		c.POST("/:id/pin", ctrl.PinComment)
		c.DELETE("/:id/pin", ctrl.UnpinComment)
		c.POST("/:id/featured", ctrl.FeatureComment)
		c.DELETE("/:id/featured", ctrl.UnfeatureComment)
		c.POST("/:id/official", ctrl.MarkOfficial)
		c.DELETE("/:id/official", ctrl.UnmarkOfficial)

		// 批量操作
		c.PUT("/batch-audit", ctrl.BatchAudit)
		c.POST("/batch-delete", ctrl.BatchDelete)

		// 评论统计与导出
		c.GET("/stats", ctrl.GetCommentStats)
		c.GET("/export", ctrl.ExportComments)
	}

	// 举报管理
	report := r.Group("/comments/reports")
	{
		report.GET("", ctrl.ListReports)
		report.PUT("/:id", ctrl.HandleReport)
	}

	// 黑名单管理
	blacklist := r.Group("/comments/blacklist")
	{
		blacklist.GET("", ctrl.ListBlacklists)
		blacklist.POST("", ctrl.AddBlacklist)
		blacklist.DELETE("/:id", ctrl.RemoveBlacklist)
	}

	// 审核规则管理
	auditRules := r.Group("/comments/audit-rules")
	{
		auditRules.GET("", ctrl.ListAuditRules)
		auditRules.POST("", ctrl.CreateAuditRule)
		auditRules.PUT("/:id", ctrl.UpdateAuditRule)
		auditRules.DELETE("/:id", ctrl.DeleteAuditRule)
	}
}
