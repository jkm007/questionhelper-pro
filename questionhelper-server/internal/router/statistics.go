package router

import (
	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/controller/statistics"
)

func SetupStatisticsRoutes(r *gin.RouterGroup, ctrl *statistics.StatisticsController) {
	// 统计数据 (设计文档路径: /statistics)
	s := r.Group("/statistics")
	{
		s.GET("/overview", ctrl.GetDashboard)
		s.GET("/user", ctrl.GetUserStatistics)      // 设计文档: /statistics/user
		s.GET("/practice", ctrl.GetPracticeStatistics) // 设计文档: /statistics/practice
		s.GET("/exam", ctrl.GetExamStatistics)       // 设计文档: /statistics/exam
		s.GET("/classes", ctrl.GetClassStatistics)
		s.GET("/ranking", ctrl.GetRanking)
	}
}
