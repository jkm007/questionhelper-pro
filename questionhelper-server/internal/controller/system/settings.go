package system

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/service/system"
	"questionhelper-server/pkg/response"
)

// GetSettings 获取系统设置
// @Summary      获取系统设置
// @Description  获取系统全局设置
// @Tags         系统设置
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/settings [get]
// @Security     BearerAuth
func (ctrl *SystemController) GetSettings(c *gin.Context) {
	settings, err := system.GetSettings()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, settings)
}

// UpdateSettings 更新系统设置
// @Summary      更新系统设置
// @Description  更新系统全局设置
// @Tags         系统设置
// @Accept       json
// @Produce      json
// @Param        settings  body      map[string]string  true  "设置键值对"
// @Success      200  {object}  response.Response  "更新成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/settings [put]
// @Security     BearerAuth
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
// @Summary      获取班级设置
// @Description  获取系统班级相关设置
// @Tags         系统设置
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/settings/class [get]
// @Security     BearerAuth
func (ctrl *SystemController) GetClassSettings(c *gin.Context) {
	settings, err := system.GetClassSettings()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, settings)
}

// UpdateClassSettings 更新班级设置
// @Summary      更新班级设置
// @Description  更新系统班级相关设置
// @Tags         系统设置
// @Accept       json
// @Produce      json
// @Param        req  body      dto.ClassSettings  true  "班级设置数据"
// @Success      200  {object}  response.Response  "更新成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/settings/class [put]
// @Security     BearerAuth
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
// @Summary      获取资源设置
// @Description  获取系统资源相关设置
// @Tags         系统设置
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/settings/resource [get]
// @Security     BearerAuth
func (ctrl *SystemController) GetResourceSettings(c *gin.Context) {
	settings, err := system.GetResourceSettings()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, settings)
}

// UpdateResourceSettings 更新资源设置
// @Summary      更新资源设置
// @Description  更新系统资源相关设置
// @Tags         系统设置
// @Accept       json
// @Produce      json
// @Param        req  body      dto.ResourceSettings  true  "资源设置数据"
// @Success      200  {object}  response.Response  "更新成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/settings/resource [put]
// @Security     BearerAuth
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
