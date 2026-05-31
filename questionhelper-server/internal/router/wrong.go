package router

import (
	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/controller/wrong"
)

func SetupWrongRoutes(r *gin.RouterGroup, ctrl *wrong.WrongController) {
	// 错题本 (设计文档路径: /wrong)
	w := r.Group("/wrong")
	{
		// 基础CRUD
		w.GET("", ctrl.ListWrongQuestions)
		w.GET("/:id", ctrl.GetWrongQuestion)
		w.POST("/:id/review", ctrl.ReviewWrongQuestion)
		w.DELETE("/:id", ctrl.RemoveWrongQuestion)

		// 错题搜索
		w.GET("/search", ctrl.SearchWrongQuestions)

		// 批量操作
		w.POST("/batch/delete", ctrl.BatchDeleteWrongQuestions)
		w.POST("/batch/export", ctrl.BatchExportWrongQuestions)
		w.POST("/batch/tag", ctrl.BatchTagWrongQuestions)

		// 错题标签管理
		w.GET("/tags", ctrl.ListWrongTags)
		w.POST("/tags", ctrl.CreateWrongTag)
		w.PUT("/tags/:id", ctrl.UpdateWrongTag)
		w.DELETE("/tags/:id", ctrl.DeleteWrongTag)
		w.POST("/:id/tags", ctrl.AddTagToWrongQuestion)
		w.DELETE("/:id/tags/:tagId", ctrl.RemoveTagFromWrongQuestion)

		// 错题备注
		w.GET("/:id/notes", ctrl.ListWrongNotes)
		w.POST("/:id/notes", ctrl.CreateWrongNote)
		w.PUT("/:id/notes/:noteId", ctrl.UpdateWrongNote)
		w.DELETE("/:id/notes/:noteId", ctrl.DeleteWrongNote)

		// 错题附件
		w.POST("/:id/attachments", ctrl.UploadWrongAttachment)
		w.GET("/:id/attachments", ctrl.ListWrongAttachments)
		w.DELETE("/:id/attachments/:attachmentId", ctrl.DeleteWrongAttachment)

		// 错题收藏
		w.POST("/:id/favorite", ctrl.FavoriteWrongQuestion)
		w.DELETE("/:id/favorite", ctrl.UnfavoriteWrongQuestion)
		w.GET("/favorites", ctrl.ListWrongFavorites)

		// 错题复习
		w.GET("/review/today", ctrl.GetTodayReviewQuestions)
		w.POST("/:id/review/record", ctrl.RecordReviewResult)
		w.GET("/review/history", ctrl.GetReviewHistory)

		// 错题导出
		w.POST("/export", ctrl.ExportWrongQuestions)

		// 错题分析扩展
		w.GET("/analysis", ctrl.GetWrongAnalysis)
		w.GET("/analysis/trend", ctrl.GetWrongTrendAnalysis)
		w.GET("/analysis/category", ctrl.GetWrongCategoryAnalysis)
		w.GET("/analysis/accuracy", ctrl.GetWrongAccuracyAnalysis)
	}
}
