package router

import (
	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/controller/statistics"
)

func SetupStatisticsRoutes(r *gin.RouterGroup, ctrl *statistics.StatisticsController) {
	s := r.Group("/statistics")
	{
		s.GET("/overview", ctrl.GetDashboard)
		s.GET("/users", ctrl.GetUserStatistics)
		s.GET("/questions", ctrl.GetPracticeStatistics)
		s.GET("/exams", ctrl.GetExamStatistics)
		s.GET("/classes", ctrl.GetClassStatistics)
		s.GET("/ranking", ctrl.GetRanking)
	}
}
