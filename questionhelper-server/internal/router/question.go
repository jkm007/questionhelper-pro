package router

import (
	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/controller/question"
)

func SetupQuestionRoutes(r *gin.RouterGroup, ctrl *question.QuestionController,
	versionCtrl *question.VersionController, shareCtrl *question.ShareController) {
	// 题目列表/详情/搜索
	q := r.Group("/questions")
	{
		q.GET("", ctrl.ListQuestions)
		q.GET("/search", ctrl.SearchQuestions)
		q.GET("/:id", ctrl.GetQuestion)
		q.POST("/:id/favorite", ctrl.FavoriteQuestion)
		q.POST("/:id/like", ctrl.LikeQuestion)

		// 版本管理
		q.GET("/:id/versions", versionCtrl.ListVersions)

		// 题目笔记
		q.GET("/:id/notes", ctrl.GetQuestionNotes)
		q.POST("/:id/notes", ctrl.CreateNote)
		q.PUT("/:id/notes/:noteId", ctrl.UpdateNote)
		q.DELETE("/:id/notes/:noteId", ctrl.DeleteNote)

		// 题目评价
		q.POST("/:id/difficulty-rating", ctrl.RateDifficulty)
		q.POST("/:id/quality-rating", ctrl.RateQuality)

		// 题目纠错
		q.POST("/:id/corrections", ctrl.CreateCorrection)
		q.GET("/:id/corrections", ctrl.GetCorrections)
	}

	// 分享
	share := r.Group("/shares")
	{
		share.GET("/:code", shareCtrl.GetShare)
	}

	cat := r.Group("/categories")
	{
		cat.GET("", ctrl.ListCategories)
		cat.GET("/tree", ctrl.GetCategoryTree)
	}

	kp := r.Group("/knowledge-points")
	{
		kp.GET("", ctrl.ListKnowledgePoints)
	}

	// 收藏夹管理
	favoriteGroup := r.Group("/favorites")
	{
		favoriteGroup.POST("/folders", ctrl.CreateFavoriteFolder)
		favoriteGroup.PUT("/folders/:id", ctrl.UpdateFavoriteFolder)
		favoriteGroup.DELETE("/folders/:id", ctrl.DeleteFavoriteFolder)
		favoriteGroup.GET("/folders", ctrl.ListFavoriteFolders)
	}
}

func SetupAdminQuestionRoutes(r *gin.RouterGroup, ctrl *question.QuestionController,
	versionCtrl *question.VersionController, batchCtrl *question.BatchController,
	shareCtrl *question.ShareController, statsCtrl *question.QuestionStatsController) {
	// 题目CRUD
	q := r.Group("/questions")
	{
		q.GET("", ctrl.AdminListQuestions)
		q.GET("/:id", ctrl.AdminGetQuestion)
		q.POST("", ctrl.CreateQuestion)
		q.PUT("/:id", ctrl.UpdateQuestion)
		q.DELETE("/:id", ctrl.DeleteQuestion)
		q.PUT("/:id/status", ctrl.UpdateQuestionStatus)

		// 版本管理
		q.GET("/:id/versions", versionCtrl.ListVersions)
		q.GET("/:id/versions/:versionId", versionCtrl.GetVersionDetail)
		q.POST("/:id/versions/:version/rollback", versionCtrl.RollbackVersion)

		// 分享
		q.POST("/:id/share", shareCtrl.CreateShare)
		q.GET("/:id/shares", shareCtrl.ListMyShares)

		// 导入导出
		q.POST("/import", ctrl.ImportQuestions)
		q.GET("/export", ctrl.ExportQuestions)
	}

	// 批量操作
	batch := r.Group("/questions/batch")
	{
		batch.POST("/publish", batchCtrl.BatchPublish)
		batch.POST("/archive", batchCtrl.BatchArchive)
		batch.POST("/delete", batchCtrl.BatchDelete)
		batch.POST("/move", batchCtrl.BatchMoveCategory)
	}

	// 统计
	stats := r.Group("/questions/stats")
	{
		stats.GET("", statsCtrl.GetStats)
	}

	// 分享管理
	shares := r.Group("/shares")
	{
		shares.GET("", shareCtrl.ListMyShares)
		shares.DELETE("/:id", shareCtrl.RevokeShare)
	}

	// 分类管理（管理员）
	categories := r.Group("/categories")
	{
		categories.POST("", ctrl.CreateCategory)
		categories.PUT("/:id", ctrl.UpdateCategory)
		categories.DELETE("/:id", ctrl.DeleteCategory)
	}

	// 知识点管理（管理员）
	knowledge := r.Group("/knowledge-points")
	{
		knowledge.POST("", ctrl.CreateKnowledgePoint)
		knowledge.PUT("/:id", ctrl.UpdateKnowledgePoint)
		knowledge.DELETE("/:id", ctrl.DeleteKnowledgePoint)
	}

	// 内容审核（管理员）
	review := r.Group("/reviews")
	{
		review.GET("", ctrl.ListPendingReviews)
		review.GET("/:id", ctrl.GetReviewDetail)
		review.POST("/:id/approve", ctrl.ApproveReview)
		review.POST("/:id/reject", ctrl.RejectReview)
	}

	// 敏感词管理（管理员）
	sensitive := r.Group("/sensitive-words")
	{
		sensitive.GET("", ctrl.ListSensitiveWords)
		sensitive.POST("", ctrl.CreateSensitiveWord)
		sensitive.DELETE("/:id", ctrl.DeleteSensitiveWord)
		sensitive.POST("/import", ctrl.ImportSensitiveWords)
		sensitive.POST("/test", ctrl.TestSensitiveWord)
	}
}
