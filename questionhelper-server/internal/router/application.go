package router

import (
	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/controller/user"
)

func SetupApplicationRoutes(r *gin.RouterGroup, ctrl *user.ApplyController) {
	app := r.Group("/applications")
	{
		app.POST("", ctrl.CreateApplication)
		app.GET("", ctrl.GetMyApplications)
		app.GET("/:id", ctrl.GetApplication)
	}
}

func SetupAdminApplicationRoutes(r *gin.RouterGroup, ctrl *user.ApplyController) {
	app := r.Group("/applications")
	{
		app.GET("", ctrl.ListApplications)
		app.GET("/:id", ctrl.GetApplication)
		app.PUT("/:id/review", ctrl.ReviewApplication)
	}
}
