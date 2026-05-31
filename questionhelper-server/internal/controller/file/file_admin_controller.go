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
func (ctrl *FileAdminController) ListHotlinkRules(c *gin.Context) {
	rules, err := fileService.ListHotlinkRules()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, rules)
}

// CreateHotlinkRule 创建防盗链规则
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
func (ctrl *FileAdminController) ListWatermarkConfigs(c *gin.Context) {
	configs, err := fileService.ListWatermarkConfigs()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, configs)
}

// CreateWatermarkConfig 创建水印配置
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
func (ctrl *FileAdminController) CleanupOrphanFiles(c *gin.Context) {
	result, err := fileService.CleanupOrphanFiles()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, result)
}

// GetCleanupLogs 清理日志列表
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
func (ctrl *FileAdminController) GetStorageStatistics(c *gin.Context) {
	stats, err := fileService.GetStorageStatistics()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, stats)
}
