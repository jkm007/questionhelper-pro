package router

import (
	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/controller/file"
)

func SetupFileRoutes(r *gin.RouterGroup, ctrl *file.FileController) {
	// 兼容前端 /files 接口
	f := r.Group("/files")
	{
		f.POST("", ctrl.UploadFile)
		f.DELETE("", ctrl.DeleteFile)
	}
}
