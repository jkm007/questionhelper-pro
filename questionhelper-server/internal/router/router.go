package router

import (
	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/controller/auth"
	"questionhelper-server/internal/controller/class"
	"questionhelper-server/internal/controller/comment"
	"questionhelper-server/internal/controller/exam"
	"questionhelper-server/internal/controller/file"
	logCtrl "questionhelper-server/internal/controller/log"
	"questionhelper-server/internal/controller/menu"
	"questionhelper-server/internal/controller/notification"
	"questionhelper-server/internal/controller/practice"
	"questionhelper-server/internal/controller/question"
	"questionhelper-server/internal/controller/statistics"
	"questionhelper-server/internal/controller/system"
	"questionhelper-server/internal/controller/user"
	"questionhelper-server/internal/controller/wrong"
	"questionhelper-server/internal/middleware"
	"questionhelper-server/pkg/config"
)

func Setup(cfg *config.Config) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CorsMiddleware())
	r.Use(middleware.RateLimitMiddleware())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 创建控制器
	authCtrl := auth.NewAuthController(&cfg.JWT)
	userCtrl := user.NewUserController()
	profileCtrl := user.NewProfileController()
	tagCtrl := user.NewTagController()
	applyCtrl := user.NewApplyController()
	logController := logCtrl.NewLogController()
	menuCtrl := menu.NewMenuController()
	questionCtrl := question.NewQuestionController()
	versionCtrl := question.NewVersionController()
	batchCtrl := question.NewBatchController()
	shareCtrl := question.NewShareController()
	questionStatsCtrl := question.NewQuestionStatsController()
	examCtrl := exam.NewExamController()
	paperCtrl := exam.NewPaperController()
	answerCtrl := exam.NewAnswerController()
	monitorCtrl := exam.NewMonitorController()
	classCtrl := class.NewClassController()
	practiceCtrl := practice.NewPracticeController()
	wrongCtrl := wrong.NewWrongController()
	commentCtrl := comment.NewCommentController()
	notificationCtrl := notification.NewNotificationController()
	statisticsCtrl := statistics.NewStatisticsController()
	fileCtrl := file.NewFileController()
	systemCtrl := system.NewSystemController()

	v1 := r.Group("/api/v1")
	{
		public := v1.Group("")
		SetupAuthRoutes(public, authCtrl)

		authorized := v1.Group("")
		authorized.Use(middleware.AuthMiddleware())
		{
			SetupUserRoutes(authorized, userCtrl, profileCtrl, tagCtrl, applyCtrl)
			SetupUserAuthRoutes(authorized, authCtrl)
			SetupTagRoutes(authorized, tagCtrl)
			SetupMenuRoutes(authorized, menuCtrl)
			SetupQuestionRoutes(authorized, questionCtrl, versionCtrl, shareCtrl)
			SetupExamRoutes(authorized, examCtrl, answerCtrl)
			SetupClassRoutes(authorized, classCtrl)
			SetupPracticeRoutes(authorized, practiceCtrl)
			SetupWrongRoutes(authorized, wrongCtrl)
			SetupCommentRoutes(authorized, commentCtrl)
			SetupNotificationRoutes(authorized, notificationCtrl)
			SetupStatisticsRoutes(authorized, statisticsCtrl)
			SetupFileRoutes(authorized, fileCtrl)
		}

		admin := v1.Group("/admin")
		admin.Use(middleware.AuthMiddleware())
		{
			SetupAdminUserRoutes(admin, userCtrl, profileCtrl)
			SetupAdminTagRoutes(admin, tagCtrl)
			SetupAdminApplicationRoutes(admin, applyCtrl)
			SetupAdminLogRoutes(admin, logController)
			SetupAdminMenuRoutes(admin, menuCtrl)
			SetupAdminQuestionRoutes(admin, questionCtrl, versionCtrl, batchCtrl, shareCtrl, questionStatsCtrl)
			SetupAdminExamRoutes(admin, examCtrl, paperCtrl, monitorCtrl)
			SetupAdminClassRoutes(admin, classCtrl)
			SetupAdminSystemRoutes(admin, systemCtrl)
		}
	}

	return r
}
