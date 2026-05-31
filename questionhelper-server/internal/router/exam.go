package router

import (
	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/controller/exam"
)

func SetupExamRoutes(r *gin.RouterGroup, ctrl *exam.ExamController,
	answerCtrl *exam.AnswerController) {
	// 考试列表/详情
	e := r.Group("/exam")
	{
		e.GET("", ctrl.ListExams)
		e.GET("/:id", ctrl.GetExam)
		e.POST("/:id/start", ctrl.StartExam)
		e.POST("/:id/submit", ctrl.SubmitExam)
		e.GET("/:id/result", ctrl.GetExamResult)
		e.GET("/history", ctrl.GetExamHistory)

		// 答案管理
		e.GET("/:id/standard-answers", answerCtrl.GetStandardAnswers)
		e.GET("/:id/guide", answerCtrl.GetExamGuide)
		e.POST("/:id/feedback", answerCtrl.SubmitFeedback)
	}

	// 考试记录相关
	record := r.Group("/exam-records")
	{
		record.POST("/:recordId/save-answer", answerCtrl.SaveAnswer)
		record.POST("/:recordId/save-answers", answerCtrl.SaveAnswers)
		record.POST("/:recordId/mark", answerCtrl.MarkQuestion)
		record.GET("/:recordId/marked", answerCtrl.GetMarkedQuestions)
		record.POST("/:recordId/warning", answerCtrl.ReportWarning)
	}
}

func SetupAdminExamRoutes(r *gin.RouterGroup, ctrl *exam.ExamController,
	paperCtrl *exam.PaperController, monitorCtrl *exam.MonitorController) {
	// 试卷管理
	paper := r.Group("/papers")
	{
		paper.GET("", ctrl.ListPapers)
		paper.GET("/:id", ctrl.GetPaper)
		paper.POST("", ctrl.CreatePaper)
		paper.PUT("/:id", ctrl.UpdatePaper)
		paper.DELETE("/:id", ctrl.DeletePaper)

		// 新增试卷功能
		paper.GET("/:id/preview", paperCtrl.PreviewPaper)
		paper.POST("/:id/copy", paperCtrl.CopyPaper)
		paper.PUT("/:id/publish", paperCtrl.PublishPaper)
		paper.POST("/:id/save-template", paperCtrl.SaveAsTemplate)
		paper.GET("/:id/export", paperCtrl.ExportPaper)
		paper.GET("/:id/stats", paperCtrl.GetPaperStats)
	}

	// 模板管理
	template := r.Group("/templates")
	{
		template.GET("", paperCtrl.ListTemplates)
		template.POST("/create", paperCtrl.CreateFromTemplate)
	}

	// 考试管理
	examGroup := r.Group("/exams")
	{
		examGroup.GET("", ctrl.AdminListExams)
		examGroup.GET("/:id", ctrl.AdminGetExam)
		examGroup.POST("", ctrl.CreateExam)
		examGroup.PUT("/:id", ctrl.UpdateExam)
		examGroup.DELETE("/:id", ctrl.DeleteExam)
		examGroup.PUT("/:id/publish", ctrl.PublishExam)
		examGroup.PUT("/:id/close", ctrl.CloseExam)

		// 监控/阅卷/分析
		examGroup.GET("/:id/monitor", monitorCtrl.GetExamMonitor)
		examGroup.POST("/:id/review", monitorCtrl.ReviewExam)
		examGroup.GET("/:id/analysis", monitorCtrl.GetExamAnalysis)
	}

	// 成绩管理
	score := r.Group("/scores")
	{
		score.GET("", ctrl.ListScores)
		score.GET("/:id", ctrl.GetScore)
		score.GET("/analysis", ctrl.GetScoreAnalysis)
		score.GET("/:id/export", monitorCtrl.ExportScores)
	}
}
