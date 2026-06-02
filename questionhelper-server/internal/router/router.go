package router

import (
	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/controller/auth"
	"questionhelper-server/internal/controller/class"
	"questionhelper-server/internal/controller/comment"
	"questionhelper-server/internal/controller/content"
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
	"questionhelper-server/internal/ws"
	"questionhelper-server/pkg/config"
)

func Setup(cfg *config.Config, hub *ws.Hub) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CorsMiddleware())
	r.Use(middleware.RateLimitMiddleware())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 静态文件服务 —— 上传文件可通过 /uploads/xxx 访问
	r.Static("/uploads", "./uploads")

	// WebSocket 端点（JWT 认证通过 query token 参数完成）
	r.GET("/ws", ws.HandleWebSocket(hub))

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
	examExtCtrl := exam.NewExamExtController()
	classCtrl := class.NewClassController()
	practiceCtrl := practice.NewPracticeController()
	practiceAdminCtrl := practice.NewPracticeAdminController()
	wrongCtrl := wrong.NewWrongController()
	commentCtrl := comment.NewCommentController()
	commentAdminCtrl := comment.NewCommentAdminController()
	notificationCtrl := notification.NewNotificationController()
	notificationAdminCtrl := notification.NewNotificationAdminController()
	statisticsCtrl := statistics.NewStatisticsController()
	statisticsAdminCtrl := statistics.NewStatisticsAdminController()
	fileCtrl := file.NewFileController()
	fileAdminCtrl := file.NewFileAdminController()
	systemCtrl := system.NewSystemController()
	contentCtrl := content.NewContentController()

	v1 := r.Group("/api/v1")
	{
		public := v1.Group("")
		SetupAuthRoutes(public, authCtrl)

		authorized := v1.Group("")
		authorized.Use(middleware.AuthMiddleware())
		authorized.Use(middleware.SensitiveFilterMiddleware()) // T05: 敏感词过滤中间件
		{
			SetupUserRoutes(authorized, userCtrl, profileCtrl, tagCtrl, applyCtrl)
			SetupUserAuthRoutes(authorized, authCtrl)
			SetupTagRoutes(authorized, tagCtrl)
			SetupMenuRoutes(authorized, menuCtrl)
			SetupQuestionRoutes(authorized, questionCtrl, versionCtrl, shareCtrl)
			SetupExamRoutes(authorized, examCtrl, answerCtrl, examExtCtrl)
			SetupClassRoutes(authorized, classCtrl)
			SetupPracticeRoutes(authorized, practiceCtrl, practiceAdminCtrl)
			SetupWrongRoutes(authorized, wrongCtrl)
			SetupCommentRoutes(authorized, commentCtrl)
			SetupNotificationRoutes(authorized, notificationCtrl)
			SetupStatisticsRoutes(authorized, statisticsCtrl)
			SetupFileRoutes(authorized, fileCtrl)
			SetupContentRoutes(authorized, contentCtrl)
		}

		admin := v1.Group("/admin")
		admin.Use(middleware.AuthMiddleware())
		admin.Use(middleware.AdminOnly())
		{
			SetupAdminUserRoutes(admin, userCtrl, profileCtrl)
			SetupAdminTagRoutes(admin, tagCtrl)
			SetupAdminApplicationRoutes(admin, applyCtrl)
			SetupAdminLogRoutes(admin, logController)
			SetupAdminMenuRoutes(admin, menuCtrl)
			SetupAdminQuestionRoutes(admin, questionCtrl, versionCtrl, batchCtrl, shareCtrl, questionStatsCtrl)
			SetupAdminExamRoutes(admin, examCtrl, paperCtrl, monitorCtrl, examExtCtrl)
			SetupAdminClassRoutes(admin, classCtrl)
			SetupAdminSystemRoutes(admin, systemCtrl)
			SetupAdminNotificationRoutes(admin, notificationAdminCtrl)
			SetupAdminCommentRoutes(admin, commentAdminCtrl)
			SetupAdminContentRoutes(admin, contentCtrl)
			SetupAdminStatisticsRoutes(admin, statisticsAdminCtrl, statisticsCtrl)
			SetupAdminFileRoutes(admin, fileAdminCtrl)
			SetupAdminPracticeRoutes(admin, practiceAdminCtrl)
		}
	}

	return r
}
