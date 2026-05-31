package router

import (
	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/controller/class"
)

func SetupClassRoutes(r *gin.RouterGroup, ctrl *class.ClassController) {
	// 与设计文档保持一致使用 /class（单数）
	c := r.Group("/class")
	{
		c.GET("", ctrl.ListClasses)
		c.GET("/:id", ctrl.GetClass)
		c.POST("", ctrl.CreateClass)
		c.PUT("/:id", ctrl.UpdateClass)
		c.DELETE("/:id", ctrl.DeleteClass)
		c.POST("/:id/join", ctrl.JoinClass)
		c.POST("/:id/leave", ctrl.LeaveClass)
		c.GET("/:id/members", ctrl.ListMembers)
		c.GET("/:id/notices", ctrl.ListNotices)
		c.GET("/:id/homework", ctrl.ListHomework)
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
