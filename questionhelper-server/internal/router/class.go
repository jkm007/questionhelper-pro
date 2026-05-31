package router

import (
	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/controller/class"
)

func SetupClassRoutes(r *gin.RouterGroup, ctrl *class.ClassController) {
	// 与设计文档保持一致使用 /class（单数）
	c := r.Group("/class")
	{
		// ==================== 基础班级 ====================
		c.GET("", ctrl.ListClasses)
		c.GET("/:id", ctrl.GetClass)
		c.POST("", ctrl.CreateClass)
		c.PUT("/:id", ctrl.UpdateClass)
		c.DELETE("/:id", ctrl.DeleteClass)
		c.POST("/:id/join", ctrl.JoinClass)
		c.POST("/:id/leave", ctrl.LeaveClass)
		c.GET("/:id/members", ctrl.ListMembers)
		c.GET("/:id/notices", ctrl.ListNotices)

		// ==================== 班级管理增强 ====================
		c.POST("/:id/archive", ctrl.ArchiveClass)
		c.POST("/:id/unarchive", ctrl.UnarchiveClass)
		c.POST("/:id/pin", ctrl.PinClass)
		c.POST("/:id/unpin", ctrl.UnpinClass)
		c.GET("/search", ctrl.SearchClasses)
		c.GET("/:id/qrcode", ctrl.GenerateQRCode)
		c.PUT("/:id/expire", ctrl.SetClassExpire)
		c.GET("/:id/exams", ctrl.ListClassExams)
		c.GET("/:id/notice", ctrl.GetClassNotice)

		// ==================== 作业管理 ====================
		c.GET("/:id/homework", ctrl.ListHomework)
		c.POST("/:id/homework", ctrl.CreateHomework)
		c.GET("/:id/homework/:homeworkId", ctrl.GetHomework)
		c.PUT("/:id/homework/:homeworkId", ctrl.UpdateHomework)
		c.DELETE("/:id/homework/:homeworkId", ctrl.DeleteHomework)
		c.POST("/:id/homework/:homeworkId/submit", ctrl.SubmitHomework)
		c.PUT("/:id/homework/:homeworkId/submissions/:submissionId", ctrl.GradeHomework)
		c.GET("/:id/homework/:homeworkId/submissions", ctrl.ListHomeworkSubmissions)

		// ==================== 互评管理 ====================
		c.POST("/:id/homework/:homeworkId/peer-reviews/assign", ctrl.AssignPeerReview)
		c.GET("/:id/homework/:homeworkId/peer-reviews", ctrl.ListPeerReviews)
		c.GET("/:id/homework/:homeworkId/peer-reviews/mine", ctrl.GetMyPeerReviews)
		c.POST("/:id/homework/:homeworkId/peer-reviews/:reviewId", ctrl.SubmitPeerReview)
		c.GET("/:id/homework/:homeworkId/peer-reviews/result", ctrl.GetPeerReviewResult)

		// ==================== 分组管理 ====================
		c.GET("/:id/groups", ctrl.ListGroups)
		c.POST("/:id/groups", ctrl.CreateGroup)
		c.PUT("/:id/groups/:groupId", ctrl.UpdateGroup)
		c.DELETE("/:id/groups/:groupId", ctrl.DeleteGroup)
		c.POST("/:id/groups/:groupId/members", ctrl.AddGroupMember)
		c.DELETE("/:id/groups/:groupId/members/:userId", ctrl.RemoveGroupMember)

		// ==================== 加入申请 ====================
		c.POST("/:id/apply", ctrl.ApplyClass)
		c.GET("/:id/applications", ctrl.ListApplications)
		c.POST("/:id/applications/:appId/approve", ctrl.ApproveApplication)
		c.POST("/:id/applications/:appId/reject", ctrl.RejectApplication)

		// ==================== 考勤管理 ====================
		c.GET("/:id/attendances", ctrl.ListAttendances)
		c.POST("/:id/attendances", ctrl.CreateAttendance)
		c.PUT("/:id/attendances/:attId", ctrl.UpdateAttendance)
		c.DELETE("/:id/attendances/:attId", ctrl.DeleteAttendance)
		c.POST("/:id/attendances/:attId/checkin", ctrl.Checkin)
		c.POST("/:id/attendances/:attId/checkout", ctrl.Checkout)
		c.GET("/:id/attendances/:attId/records", ctrl.ListAttendanceRecords)
		c.GET("/:id/attendances/:attId/export", ctrl.ExportAttendance)

		// ==================== 学习计划 ====================
		c.GET("/:id/study-plans", ctrl.ListStudyPlans)
		c.POST("/:id/study-plans", ctrl.CreateStudyPlan)
		c.GET("/:id/study-plans/:planId", ctrl.GetStudyPlan)
		c.PUT("/:id/study-plans/:planId", ctrl.UpdateStudyPlan)
		c.DELETE("/:id/study-plans/:planId", ctrl.DeleteStudyPlan)
		c.POST("/:id/study-plans/:planId/items", ctrl.AddStudyPlanItem)
		c.PUT("/:id/study-plans/:planId/items/:itemId", ctrl.UpdateStudyPlanItem)
		c.DELETE("/:id/study-plans/:planId/items/:itemId", ctrl.DeleteStudyPlanItem)
		c.POST("/:id/study-plans/:planId/items/:itemId/complete", ctrl.CompleteStudyPlanItem)
		c.GET("/:id/study-plans/:planId/progress", ctrl.GetStudyPlanProgress)

		// ==================== 文件管理 ====================
		c.GET("/:id/files", ctrl.ListClassFiles)
		c.POST("/:id/files", ctrl.UploadClassFile)
		c.DELETE("/:id/files/:fileId", ctrl.DeleteClassFile)
		c.GET("/:id/files/:fileId/download", ctrl.DownloadClassFile)

		// ==================== 排名管理 ====================
		c.GET("/:id/ranking", ctrl.ListRanking)
		c.POST("/:id/ranking/calculate", ctrl.CalculateRanking)

		// ==================== 标签管理 ====================
		c.GET("/tags", ctrl.ListTags)
		c.POST("/tags", ctrl.CreateTag)
		c.PUT("/tags/:tagId", ctrl.UpdateTag)
		c.DELETE("/tags/:tagId", ctrl.DeleteTag)
		c.POST("/:id/tags", ctrl.AddClassTag)
		c.DELETE("/:id/tags/:tagId", ctrl.RemoveClassTag)

		// ==================== 模板管理 ====================
		c.GET("/templates", ctrl.ListTemplates)
		c.POST("/templates", ctrl.CreateTemplate)
		c.GET("/templates/:templateId", ctrl.GetTemplate)
		c.PUT("/templates/:templateId", ctrl.UpdateTemplate)
		c.DELETE("/templates/:templateId", ctrl.DeleteTemplate)
		c.POST("/templates/:templateId/create", ctrl.CreateClassFromTemplate)

		// ==================== 讨论管理 ====================
		c.GET("/:id/discussions", ctrl.ListDiscussions)
		c.POST("/:id/discussions", ctrl.CreateDiscussion)
		c.GET("/:id/discussions/:discussionId", ctrl.GetDiscussion)
		c.PUT("/:id/discussions/:discussionId", ctrl.UpdateDiscussion)
		c.DELETE("/:id/discussions/:discussionId", ctrl.DeleteDiscussion)
		c.POST("/:id/discussions/:discussionId/pin", ctrl.ToggleDiscussionPin)

		// ==================== 资源管理 ====================
		c.GET("/:id/resources", ctrl.ListResources)
		c.GET("/:id/resources/statistics", ctrl.GetResourceStatistics)
		c.POST("/:id/resources/import", ctrl.ImportResource)
		c.GET("/:id/resources/export", ctrl.ExportResource)

		// ==================== 创作者管理 ====================
		c.GET("/:id/creator-applications", ctrl.ListCreatorApplications)
		c.POST("/:id/creator-apply", ctrl.CreatorApply)
		c.POST("/:id/creator-applications/:appId/approve", ctrl.ApproveCreatorApplication)
		c.POST("/:id/creator-applications/:appId/reject", ctrl.RejectCreatorApplication)
		c.DELETE("/:id/creators/:userId", ctrl.RemoveCreator)
		c.GET("/:id/creators", ctrl.ListCreators)
	}
}

func SetupAdminClassRoutes(r *gin.RouterGroup, ctrl *class.ClassController) {
	c := r.Group("/class")
	{
		c.GET("", ctrl.AdminListClasses)
		c.GET("/:id", ctrl.AdminGetClass)
		c.POST("", ctrl.AdminCreateClass)
		c.PUT("/:id", ctrl.AdminUpdateClass)
		c.DELETE("/:id", ctrl.AdminDeleteClass)
		c.GET("/:id/members", ctrl.AdminListMembers)
		c.DELETE("/:id/members/:uid", ctrl.RemoveMember)
	}
}
