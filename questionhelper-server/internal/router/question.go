package router

import (
	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/controller/question"
)

func SetupQuestionRoutes(r *gin.RouterGroup, ctrl *question.QuestionController,
	versionCtrl *question.VersionController, shareCtrl *question.ShareController) {
	// 题目列表/详情
	q := r.Group("/questions")
	{
		q.GET("", ctrl.ListQuestions)
		q.GET("/:id", ctrl.GetQuestion)
		q.POST("/:id/favorite", ctrl.FavoriteQuestion)
		q.POST("/:id/like", ctrl.LikeQuestion)

		// 版本管理
		q.GET("/:id/versions", versionCtrl.ListVersions)
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
}
