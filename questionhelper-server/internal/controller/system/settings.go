package system

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/service/system"
	"questionhelper-server/pkg/response"
)

// GetSettings 获取系统设置
func (ctrl *SystemController) GetSettings(c *gin.Context) {
	settings, err := system.GetSettings()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, settings)
}

// UpdateSettings 更新系统设置
func (ctrl *SystemController) UpdateSettings(c *gin.Context) {
	var settings map[string]string
	if err := c.ShouldBindJSON(&settings); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := system.UpdateSettings(settings); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "更新成功", nil)
}

// GetClassSettings 获取班级设置
func (ctrl *SystemController) GetClassSettings(c *gin.Context) {
	settings, err := system.GetClassSettings()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, settings)
}

// UpdateClassSettings 更新班级设置
func (ctrl *SystemController) UpdateClassSettings(c *gin.Context) {
	var req dto.ClassSettings
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	if err := system.UpdateClassSettings(&req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "更新成功", nil)
}

// GetResourceSettings 获取资源设置
func (ctrl *SystemController) GetResourceSettings(c *gin.Context) {
	settings, err := system.GetResourceSettings()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, settings)
}

// UpdateResourceSettings 更新资源设置
func (ctrl *SystemController) UpdateResourceSettings(c *gin.Context) {
	var req dto.ResourceSettings
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	if err := system.UpdateResourceSettings(&req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "更新成功", nil)
}
