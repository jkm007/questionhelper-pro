package router

import (
	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/controller/content"
)

func SetupContentRoutes(r *gin.RouterGroup, ctrl *content.ContentController) {
	// 创作者申请
	userApply := r.Group("/user/apply")
	{
		userApply.POST("/creator", ctrl.ApplyCreator)
	}

	// 创作者信息
	creator := r.Group("/creator")
	{
		creator.GET("/profile", ctrl.GetCreatorProfile)
		creator.GET("/level", ctrl.GetCreatorLevel)
		creator.GET("/points", ctrl.GetCreatorPoints)
		creator.GET("/points/logs", ctrl.GetCreatorPointLogs)
	}

	// 创作者协议
	agreement := r.Group("/creator/agreement")
	{
		agreement.GET("", ctrl.GetCreatorAgreement)
		agreement.POST("/sign", ctrl.SignAgreement)
	}

	// 创作者作品集
	portfolios := r.Group("/creator/portfolios")
	{
		portfolios.GET("", ctrl.ListPortfolios)
		portfolios.POST("", ctrl.CreatePortfolio)
		portfolios.GET("/:id", ctrl.GetPortfolio)
		portfolios.PUT("/:id", ctrl.UpdatePortfolio)
		portfolios.DELETE("/:id", ctrl.DeletePortfolio)
	}

	// 内容版本管理
	contentGroup := r.Group("/content")
	{
		contentGroup.GET("/:type/:id/versions", ctrl.ListContentVersions)
		contentGroup.GET("/:type/:id/versions/:versionId", ctrl.GetContentVersion)
		contentGroup.POST("/:type/:id/versions/rollback", ctrl.RollbackContentVersion)
	}

	// 内容标签
	contentTags := r.Group("/content")
	{
		contentTags.GET("/:type/:id/tags", ctrl.GetContentTags)
		contentTags.POST("/:type/:id/tags", ctrl.AddContentTag)
		contentTags.DELETE("/:type/:id/tags/:tagId", ctrl.RemoveContentTag)
		contentTags.GET("/tags/hot", ctrl.GetHotTags)
	}

	// 内容收藏
	contentFavorites := r.Group("/content")
	{
		contentFavorites.GET("/:type/:id/favorites", ctrl.CheckContentFavorite)
		contentFavorites.POST("/:type/:id/favorites", ctrl.AddContentFavorite)
		contentFavorites.DELETE("/:type/:id/favorites", ctrl.RemoveContentFavorite)
	}

	// 我的收藏列表
	favorites := r.Group("/favorites")
	{
		favorites.GET("", ctrl.ListContentFavorites)
	}

	// 内容预览
	contentPreview := r.Group("/content")
	{
		contentPreview.GET("/:type/:id/preview", ctrl.GetContentPreview)
	}

	// 综合搜索
	search := r.Group("/search")
	{
		search.GET("", ctrl.Search)
		search.GET("/suggestions", ctrl.GetSearchSuggestions)
		search.GET("/hot", ctrl.GetHotSearches)
		search.GET("/history", ctrl.GetSearchHistory)
	}

	// 提交内容审核（用户提交审核，保留在 authorized 组）
	contentReviews := r.Group("/content")
	{
		contentReviews.POST("/reviews", ctrl.SubmitReview)
	}
}

// SetupAdminContentRoutes 管理员内容路由
func SetupAdminContentRoutes(r *gin.RouterGroup, ctrl *content.ContentController) {
	// 审核流程（仅管理员可操作）
	reviews := r.Group("/reviews")
	{
		reviews.GET("/pending", ctrl.ListPendingReviews)
		reviews.GET("/:id", ctrl.GetReviewDetail)
		reviews.POST("/:id/approve", ctrl.ApproveReview)
		reviews.POST("/:id/reject", ctrl.RejectReview)
		reviews.POST("/:id/comment", ctrl.AddReviewComment)
	}
}
