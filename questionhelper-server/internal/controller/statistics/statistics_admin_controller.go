package statistics

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/service/statistics"
	"questionhelper-server/pkg/response"
)

type StatisticsAdminController struct{}

func NewStatisticsAdminController() *StatisticsAdminController {
	return &StatisticsAdminController{}
}

// ==================== 用户留存分析 ====================

// GetRetention 用户留存统计
func (ctrl *StatisticsAdminController) GetRetention(c *gin.Context) {
	var req dto.RetentionRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	items, err := statistics.GetRetention(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, items)
}

// ==================== 用户流失分析 ====================

// GetChurn 用户流失统计
func (ctrl *StatisticsAdminController) GetChurn(c *gin.Context) {
	var req dto.ChurnRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	items, err := statistics.GetChurn(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, items)
}

// ==================== 用户行为事件 ====================

// CreateEvent 上报用户行为事件
func (ctrl *StatisticsAdminController) CreateEvent(c *gin.Context) {
	var req dto.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	userID := c.GetUint("user_id")
	ip := c.ClientIP()

	if err := statistics.CreateEvent(userID, ip, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "上报成功", nil)
}

// AnalyzeEvents 行为事件分析
func (ctrl *StatisticsAdminController) AnalyzeEvents(c *gin.Context) {
	var req dto.EventAnalysisRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	result, err := statistics.AnalyzeEvents(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, result)
}

// ==================== 用户分群 ====================

// ListSegments 分群列表
func (ctrl *StatisticsAdminController) ListSegments(c *gin.Context) {
	var req dto.SegmentListRequest
	c.ShouldBindQuery(&req)

	list, total, err := statistics.ListSegments(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// CreateSegment 创建分群
func (ctrl *StatisticsAdminController) CreateSegment(c *gin.Context) {
	var req dto.CreateSegmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	creatorID := c.GetUint("user_id")
	if err := statistics.CreateSegment(creatorID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "创建成功", nil)
}

// GetSegment 分群详情
func (ctrl *StatisticsAdminController) GetSegment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的分群ID")
		return
	}

	detail, err := statistics.GetSegment(uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, detail)
}

// UpdateSegment 更新分群
func (ctrl *StatisticsAdminController) UpdateSegment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的分群ID")
		return
	}

	var req dto.UpdateSegmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := statistics.UpdateSegment(uint(id), &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "更新成功", nil)
}

// DeleteSegment 删除分群
func (ctrl *StatisticsAdminController) DeleteSegment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的分群ID")
		return
	}

	if err := statistics.DeleteSegment(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "删除成功", nil)
}

// ==================== 用户路径分析 ====================

// GetPathAnalysis 访问路径分析
func (ctrl *StatisticsAdminController) GetPathAnalysis(c *gin.Context) {
	var req dto.PathAnalysisRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	result, err := statistics.GetPathAnalysis(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, result)
}

// ==================== 转化漏斗 ====================

// ListFunnels 漏斗列表
func (ctrl *StatisticsAdminController) ListFunnels(c *gin.Context) {
	var req dto.PageRequest
	c.ShouldBindQuery(&req)

	list, total, err := statistics.ListFunnels(req.Page, req.PageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// CreateFunnel 创建漏斗
func (ctrl *StatisticsAdminController) CreateFunnel(c *gin.Context) {
	var req dto.CreateFunnelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	creatorID := c.GetUint("user_id")
	if err := statistics.CreateFunnel(creatorID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "创建成功", nil)
}

// GetFunnelStats 漏斗统计
func (ctrl *StatisticsAdminController) GetFunnelStats(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的漏斗ID")
		return
	}

	var req dto.FunnelStatsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	result, err := statistics.GetFunnelStats(uint(id), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, result)
}

// ==================== 数据预警 ====================

// ListAlertRules 预警规则列表
func (ctrl *StatisticsAdminController) ListAlertRules(c *gin.Context) {
	var req dto.AlertRuleListRequest
	c.ShouldBindQuery(&req)

	list, total, err := statistics.ListAlertRules(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// CreateAlertRule 创建预警规则
func (ctrl *StatisticsAdminController) CreateAlertRule(c *gin.Context) {
	var req dto.CreateAlertRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	creatorID := c.GetUint("user_id")
	if err := statistics.CreateAlertRule(creatorID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "创建成功", nil)
}

// ListAlertRecords 预警记录列表
func (ctrl *StatisticsAdminController) ListAlertRecords(c *gin.Context) {
	var req dto.AlertRecordListRequest
	c.ShouldBindQuery(&req)

	list, total, err := statistics.ListAlertRecords(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// ==================== 数据导出 ====================

// ExportData 导出统计数据
func (ctrl *StatisticsAdminController) ExportData(c *gin.Context) {
	var req dto.ExportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	userID := c.GetUint("user_id")
	info, err := statistics.ExportData(userID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, info)
}

// ==================== 数据对比 ====================

// CompareData 数据对比
func (ctrl *StatisticsAdminController) CompareData(c *gin.Context) {
	var req dto.CompareRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	result, err := statistics.CompareData(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, result)
}
