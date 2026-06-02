package router

import (
	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/controller/statistics"
)

// SetupStatisticsRoutes 注册用户端统计路由
func SetupStatisticsRoutes(r *gin.RouterGroup, ctrl *statistics.StatisticsController) {
	s := r.Group("/statistics")
	{
		// 基础统计
		s.GET("/overview", ctrl.GetDashboard)
		s.GET("/user", ctrl.GetUserStatistics)
		s.GET("/practice", ctrl.GetPracticeStatistics)
		s.GET("/exam", ctrl.GetExamStatistics)
		s.GET("/classes", ctrl.GetClassStatistics)
		s.GET("/ranking", ctrl.GetRanking)

		// 用户行为事件(普通用户上报)
		s.POST("/events", ctrl.CreateEvent)

		// 数据订阅
		s.GET("/subscriptions", ctrl.ListSubscriptions)
		s.POST("/subscriptions", ctrl.CreateSubscription)
		s.PUT("/subscriptions/:id", ctrl.UpdateSubscription)
		s.DELETE("/subscriptions/:id", ctrl.DeleteSubscription)

		// 题目分析
		s.GET("/question/difficulty", ctrl.GetQuestionDifficulty)
		s.GET("/question/discrimination", ctrl.GetQuestionDiscrimination)

		// 成绩预测与预警
		s.GET("/score/prediction", ctrl.GetScorePrediction)
		s.GET("/score/alert", ctrl.GetScoreAlert)

		// 班级统计(教师视角)
		s.GET("/class/:id/overview", ctrl.GetClassOverview)
		s.GET("/class/:id/students", ctrl.GetClassStudents)
		s.GET("/class/:id/practice", ctrl.GetClassPracticeStats)
		s.GET("/class/:id/exam", ctrl.GetClassExamStats)
		s.GET("/class/:id/questions", ctrl.GetClassQuestionStats)

		// 移动端统计
		s.GET("/mobile/overview", ctrl.GetMobileOverview)
		s.GET("/mobile/practice", ctrl.GetMobilePracticeStats)
		s.GET("/mobile/wrong", ctrl.GetMobileWrongStats)
		s.GET("/mobile/trend", ctrl.GetMobileTrend)

		// 数据分享
		s.POST("/share", ctrl.ShareData)
		s.GET("/share/:code", ctrl.GetSharedData)

		// 手动刷新
		s.POST("/refresh", ctrl.RefreshStats)
	}
}

// SetupAdminStatisticsRoutes 注册管理员统计路由
func SetupAdminStatisticsRoutes(r *gin.RouterGroup, adminCtrl *statistics.StatisticsAdminController, userCtrl *statistics.StatisticsController) {
	s := r.Group("/statistics")
	{
		// 管理端仪表盘
		s.GET("/dashboard", adminCtrl.GetAdminDashboard)

		// 用户留存分析
		s.GET("/retention", adminCtrl.GetRetention)

		// 留存/流失趋势
		s.GET("/retention/trend", adminCtrl.GetRetentionTrend)
		s.GET("/churn/trend", adminCtrl.GetChurnTrend)

		// 用户流失分析
		s.GET("/churn", adminCtrl.GetChurn)

		// 用户行为事件
		s.POST("/events", userCtrl.CreateEvent) // 管理员也可上报事件
		s.GET("/events/analysis", adminCtrl.AnalyzeEvents)

		// 用户分群 CRUD
		s.GET("/segments", adminCtrl.ListSegments)
		s.POST("/segments", adminCtrl.CreateSegment)
		s.GET("/segments/:id", adminCtrl.GetSegment)
		s.PUT("/segments/:id", adminCtrl.UpdateSegment)
		s.DELETE("/segments/:id", adminCtrl.DeleteSegment)

		// 分群成员
		s.GET("/segments/:id/members", adminCtrl.GetSegmentMembers)
		s.POST("/segments/:id/sync", adminCtrl.SyncSegmentMembers)

		// 用户路径分析
		s.GET("/paths", adminCtrl.GetPathAnalysis)

		// 转化漏斗
		s.GET("/funnels", adminCtrl.ListFunnels)
		s.POST("/funnels", adminCtrl.CreateFunnel)
		s.GET("/funnels/:id/stats", adminCtrl.GetFunnelStats)

		// 数据预警
		s.GET("/alerts/rules", adminCtrl.ListAlertRules)
		s.POST("/alerts/rules", adminCtrl.CreateAlertRule)
		s.PUT("/alerts/rules/:id", adminCtrl.UpdateAlertRule)
		s.DELETE("/alerts/rules/:id", adminCtrl.DeleteAlertRule)
		s.GET("/alerts/records", adminCtrl.ListAlertRecords)
		s.PUT("/alerts/records/:id/handle", adminCtrl.HandleAlertRecord)

		// 数据导出
		s.POST("/export", adminCtrl.ExportData)

		// 数据对比
		s.GET("/compare", adminCtrl.CompareData)

		// 对比指标列表
		s.GET("/compare/metrics", adminCtrl.GetCompareMetrics)
	}
}
