package router

import (
	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/controller/practice"
)

func SetupPracticeRoutes(r *gin.RouterGroup, ctrl *practice.PracticeController) {
	// 与设计文档保持一致使用 /practice（单数）
	p := r.Group("/practice")
	{
		// 基础练习
		p.GET("", ctrl.GetPracticeHistory)
		p.GET("/:id", ctrl.GetPracticeResult)
		p.GET("/stats", ctrl.GetPracticeStats)
		p.POST("/start", ctrl.StartPractice)
		p.POST("/submit", ctrl.SubmitPractice)
		p.POST("/:id/finish", ctrl.FinishPractice)

		// 模拟考试
		mock := p.Group("/mock")
		{
			mock.POST("/start", ctrl.StartMockExam)
			mock.GET("/history", ctrl.GetMockExamHistory)
			mock.GET("/:id/detail", ctrl.GetMockExamDetail)
			mock.POST("/:id/submit", ctrl.SubmitMockExam)
		}

		// 练习计划
		plans := p.Group("/plans")
		{
			plans.GET("", ctrl.GetPlans)
			plans.POST("", ctrl.CreatePlan)
			plans.GET("/:id", ctrl.GetPlan)
			plans.PUT("/:id", ctrl.UpdatePlan)
			plans.DELETE("/:id", ctrl.DeletePlan)
			plans.POST("/:id/execute", ctrl.ExecutePlan)
		}

		// 每日练习
		daily := p.Group("/daily")
		{
			daily.GET("/today", ctrl.GetTodayPractice)
			daily.POST("/complete", ctrl.CompleteDailyPractice)
		}

		// 练习打卡
		p.POST("/checkin", ctrl.Checkin)
		p.GET("/checkin/calendar", ctrl.GetCheckinCalendar)

		// 排行榜
		p.GET("/leaderboard", ctrl.GetLeaderboard)

		// 闯关模式
		challenge := p.Group("/challenge")
		{
			challenge.GET("/levels", ctrl.GetChallengeLevels)
			challenge.GET("/levels/:id", ctrl.GetChallengeLevel)
			challenge.POST("/levels/:id/start", ctrl.StartChallenge)
			challenge.POST("/levels/:id/submit", ctrl.SubmitChallenge)
			challenge.GET("/progress", ctrl.GetChallengeProgress)
		}
	}
}
