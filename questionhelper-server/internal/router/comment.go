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
		c.DELETE("/:id", ctrl.DeleteComment)
		c.POST("/:id/like", ctrl.LikeComment)
		c.POST("/:id/report", ctrl.ReportComment)
	}
}
