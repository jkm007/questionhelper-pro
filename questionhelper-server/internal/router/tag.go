package router

import (
	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/controller/user"
)

func SetupTagRoutes(r *gin.RouterGroup, ctrl *user.TagController) {
	tag := r.Group("/tags")
	{
		tag.GET("", ctrl.GetAllTags)
	}
}

func SetupAdminTagRoutes(r *gin.RouterGroup, ctrl *user.TagController) {
	tag := r.Group("/tags")
	{
		tag.GET("", ctrl.ListTags)
		tag.GET("/:id", ctrl.GetTag)
		tag.POST("", ctrl.CreateTag)
		tag.PUT("/:id", ctrl.UpdateTag)
		tag.DELETE("/:id", ctrl.DeleteTag)
	}
}
