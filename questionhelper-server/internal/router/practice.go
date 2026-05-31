package router

import (
	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/controller/practice"
)

func SetupPracticeRoutes(r *gin.RouterGroup, ctrl *practice.PracticeController) {
	// 兼容前端 /practices 接口
	p := r.Group("/practices")
	{
		p.GET("", ctrl.GetPracticeHistory)
		p.GET("/:id", ctrl.GetPracticeResult)
		p.GET("/stats", ctrl.GetPracticeStats)
		p.POST("/start", ctrl.StartPractice)
		p.POST("/submit", ctrl.SubmitPractice)
	}
}
