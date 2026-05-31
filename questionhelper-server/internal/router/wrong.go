package router

import (
	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/controller/wrong"
)

func SetupWrongRoutes(r *gin.RouterGroup, ctrl *wrong.WrongController) {
	// 兼容前端 /wrong-questions 接口
	w := r.Group("/wrong-questions")
	{
		w.GET("", ctrl.ListWrongQuestions)
		w.GET("/:id", ctrl.GetWrongQuestion)
		w.POST("/:id/review", ctrl.ReviewWrongQuestion)
		w.DELETE("/:id", ctrl.RemoveWrongQuestion)
		w.GET("/analysis", ctrl.GetWrongAnalysis)
	}
}
