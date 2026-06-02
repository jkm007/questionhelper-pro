package file

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/dto"
	fileService "questionhelper-server/internal/service/file"
	"questionhelper-server/pkg/response"
)

type FileAdminController struct{}

func NewFileAdminController() *FileAdminController {
	return &FileAdminController{}
}

// ==================== Hotlink Protection ====================

// ListHotlinkRules 防盗链规则列表
// @Summary      获取防盗链规则列表
// @Description  获取所有防盗链规则
// @Tags         文件管理(管理员)
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/file/hotlink-rules [get]
// @Security     BearerAuth
func (ctrl *FileAdminController) ListHotlinkRules(c *gin.Context) {
	rules, err := fileService.ListHotlinkRules()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, rules)
}

// CreateHotlinkRule 创建防盗链规则
// @Summary      创建防盗链规则
// @Description  创建一条新的防盗链规则
// @Tags         文件管理(管理员)
// @Accept       json
// @Produce      json
// @Param        req  body      dto.HotlinkRuleRequest  true  "规则数据"
// @Success      200  {object}  response.Response  "创建成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/file/hotlink-rules [post]
// @Security     BearerAuth
func (ctrl *FileAdminController) CreateHotlinkRule(c *gin.Context) {
	var req dto.HotlinkRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := fileService.CreateHotlinkRule(&req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "创建成功", nil)
}

// UpdateHotlinkRule 更新防盗链规则
// @Summary      更新防盗链规则
// @Description  根据ID更新防盗链规则
// @Tags         文件管理(管理员)
// @Accept       json
// @Produce      json
// @Param        id   path      uint                    true  "规则ID"
// @Param        req  body      dto.HotlinkRuleRequest  true  "规则数据"
// @Success      200  {object}  response.Response  "更新成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/file/hotlink-rules/{id} [put]
// @Security     BearerAuth
func (ctrl *FileAdminController) UpdateHotlinkRule(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的规则ID")
		return
	}

	var req dto.HotlinkRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := fileService.UpdateHotlinkRule(uint(id), &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "更新成功", nil)
}

// DeleteHotlinkRule 删除防盗链规则
// @Summary      删除防盗链规则
// @Description  根据ID删除防盗链规则
// @Tags         文件管理(管理员)
// @Accept       json
// @Produce      json
// @Param        id  path      uint  true  "规则ID"
// @Success      200  {object}  response.Response  "删除成功"
// @Failure      400  {object}  response.Response  "无效的规则ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/file/hotlink-rules/{id} [delete]
// @Security     BearerAuth
func (ctrl *FileAdminController) DeleteHotlinkRule(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的规则ID")
		return
	}

	if err := fileService.DeleteHotlinkRule(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "删除成功", nil)
}

// ==================== Watermark Config ====================

// ListWatermarkConfigs 水印配置列表
// @Summary      获取水印配置列表
// @Description  获取所有水印配置
// @Tags         文件管理(管理员)
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/file/watermark-configs [get]
// @Security     BearerAuth
func (ctrl *FileAdminController) ListWatermarkConfigs(c *gin.Context) {
	configs, err := fileService.ListWatermarkConfigs()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, configs)
}

// CreateWatermarkConfig 创建水印配置
// @Summary      创建水印配置
// @Description  创建一条新的水印配置
// @Tags         文件管理(管理员)
// @Accept       json
// @Produce      json
// @Param        req  body      dto.WatermarkConfigRequest  true  "配置数据"
// @Success      200  {object}  response.Response  "创建成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/file/watermark-configs [post]
// @Security     BearerAuth
func (ctrl *FileAdminController) CreateWatermarkConfig(c *gin.Context) {
	var req dto.WatermarkConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := fileService.CreateWatermarkConfig(&req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "创建成功", nil)
}

// UpdateWatermarkConfig 更新水印配置
// @Summary      更新水印配置
// @Description  根据ID更新水印配置
// @Tags         文件管理(管理员)
// @Accept       json
// @Produce      json
// @Param        id   path      uint                          true  "配置ID"
// @Param        req  body      dto.WatermarkConfigRequest    true  "配置数据"
// @Success      200  {object}  response.Response  "更新成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/file/watermark-configs/{id} [put]
// @Security     BearerAuth
func (ctrl *FileAdminController) UpdateWatermarkConfig(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的配置ID")
		return
	}

	var req dto.WatermarkConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := fileService.UpdateWatermarkConfig(uint(id), &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "更新成功", nil)
}

// ==================== Orphan Cleanup ====================

// CleanupOrphanFiles 执行孤立文件清理
// @Summary      执行孤立文件清理
// @Description  清理未被引用的孤立文件
// @Tags         文件管理(管理员)
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/file/cleanup/orphan [post]
// @Security     BearerAuth
func (ctrl *FileAdminController) CleanupOrphanFiles(c *gin.Context) {
	result, err := fileService.CleanupOrphanFiles()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, result)
}

// GetCleanupLogs 清理日志列表
// @Summary      获取清理日志列表
// @Description  获取文件清理日志列表，支持分页
// @Tags         文件管理(管理员)
// @Accept       json
// @Produce      json
// @Param        page       query     int  false  "页码"
// @Param        page_size  query     int  false  "每页数量"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/file/cleanup/logs [get]
// @Security     BearerAuth
func (ctrl *FileAdminController) GetCleanupLogs(c *gin.Context) {
	var req dto.PageRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	logs, total, err := fileService.GetCleanupLogs(req.Page, req.PageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, logs, total, req.Page, req.PageSize)
}

// ==================== Storage Statistics ====================

// GetStorageStatistics 存储统计
// @Summary      获取存储统计
// @Description  获取文件存储统计数据
// @Tags         文件管理(管理员)
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/file/storage/statistics [get]
// @Security     BearerAuth
func (ctrl *FileAdminController) GetStorageStatistics(c *gin.Context) {
	stats, err := fileService.GetStorageStatistics()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, stats)
}
