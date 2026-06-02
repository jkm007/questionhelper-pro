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
// @Summary      获取用户留存统计
// @Description  根据查询条件获取用户留存统计数据
// @Tags         统计管理
// @Accept       json
// @Produce      json
// @Param        req  query     dto.RetentionRequest  true  "查询参数"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/statistics/retention [get]
// @Security     BearerAuth
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
// @Summary      获取用户流失统计
// @Description  根据查询条件获取用户流失统计数据
// @Tags         统计管理
// @Accept       json
// @Produce      json
// @Param        req  query     dto.ChurnRequest  true  "查询参数"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/statistics/churn [get]
// @Security     BearerAuth
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
// @Summary      管理员上报用户行为事件
// @Description  管理员端上报用户行为事件数据
// @Tags         统计管理
// @Accept       json
// @Produce      json
// @Param        req  body      dto.CreateEventRequest  true  "事件数据"
// @Success      200  {object}  response.Response  "上报成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/statistics/events [post]
// @Security     BearerAuth
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
// @Summary      行为事件分析
// @Description  根据查询条件分析用户行为事件数据
// @Tags         统计管理
// @Accept       json
// @Produce      json
// @Param        req  query     dto.EventAnalysisRequest  true  "查询参数"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/statistics/events/analysis [get]
// @Security     BearerAuth
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
// @Summary      获取分群列表
// @Description  获取用户分群列表，支持分页
// @Tags         统计管理
// @Accept       json
// @Produce      json
// @Param        page       query     int  false  "页码"
// @Param        page_size  query     int  false  "每页数量"
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/statistics/segments [get]
// @Security     BearerAuth
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
// @Summary      创建用户分群
// @Description  创建一个新的用户分群
// @Tags         统计管理
// @Accept       json
// @Produce      json
// @Param        req  body      dto.CreateSegmentRequest  true  "分群数据"
// @Success      200  {object}  response.Response  "创建成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/statistics/segments [post]
// @Security     BearerAuth
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
// @Summary      获取分群详情
// @Description  根据ID获取用户分群详情
// @Tags         统计管理
// @Accept       json
// @Produce      json
// @Param        id  path      uint  true  "分群ID"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "无效的分群ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/statistics/segments/{id} [get]
// @Security     BearerAuth
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
// @Summary      更新用户分群
// @Description  根据ID更新用户分群信息
// @Tags         统计管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint                      true  "分群ID"
// @Param        req  body      dto.UpdateSegmentRequest  true  "分群数据"
// @Success      200  {object}  response.Response  "更新成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/statistics/segments/{id} [put]
// @Security     BearerAuth
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
// @Summary      删除用户分群
// @Description  根据ID删除用户分群
// @Tags         统计管理
// @Accept       json
// @Produce      json
// @Param        id  path      uint  true  "分群ID"
// @Success      200  {object}  response.Response  "删除成功"
// @Failure      400  {object}  response.Response  "无效的分群ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/statistics/segments/{id} [delete]
// @Security     BearerAuth
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
// @Summary      获取访问路径分析
// @Description  根据查询条件获取用户访问路径分析数据
// @Tags         统计管理
// @Accept       json
// @Produce      json
// @Param        req  query     dto.PathAnalysisRequest  true  "查询参数"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/statistics/paths [get]
// @Security     BearerAuth
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
// @Summary      获取转化漏斗列表
// @Description  获取转化漏斗列表，支持分页
// @Tags         统计管理
// @Accept       json
// @Produce      json
// @Param        page       query     int  false  "页码"
// @Param        page_size  query     int  false  "每页数量"
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/statistics/funnels [get]
// @Security     BearerAuth
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
// @Summary      创建转化漏斗
// @Description  创建一个新的转化漏斗
// @Tags         统计管理
// @Accept       json
// @Produce      json
// @Param        req  body      dto.CreateFunnelRequest  true  "漏斗数据"
// @Success      200  {object}  response.Response  "创建成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/statistics/funnels [post]
// @Security     BearerAuth
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
// @Summary      获取漏斗统计数据
// @Description  根据ID获取转化漏斗统计数据
// @Tags         统计管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint                    true  "漏斗ID"
// @Param        req  query     dto.FunnelStatsRequest  false  "查询参数"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "无效的漏斗ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/statistics/funnels/{id}/stats [get]
// @Security     BearerAuth
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
// @Summary      获取预警规则列表
// @Description  获取数据预警规则列表，支持分页
// @Tags         统计管理
// @Accept       json
// @Produce      json
// @Param        page       query     int  false  "页码"
// @Param        page_size  query     int  false  "每页数量"
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/statistics/alerts/rules [get]
// @Security     BearerAuth
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
// @Summary      创建预警规则
// @Description  创建一条新的数据预警规则
// @Tags         统计管理
// @Accept       json
// @Produce      json
// @Param        req  body      dto.CreateAlertRuleRequest  true  "规则数据"
// @Success      200  {object}  response.Response  "创建成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/statistics/alerts/rules [post]
// @Security     BearerAuth
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
// @Summary      获取预警记录列表
// @Description  获取数据预警记录列表，支持分页
// @Tags         统计管理
// @Accept       json
// @Produce      json
// @Param        page       query     int  false  "页码"
// @Param        page_size  query     int  false  "每页数量"
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/statistics/alerts/records [get]
// @Security     BearerAuth
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
// @Summary      导出统计数据
// @Description  根据请求条件导出统计数据
// @Tags         统计管理
// @Accept       json
// @Produce      json
// @Param        req  body      dto.ExportRequest  true  "导出参数"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/statistics/export [post]
// @Security     BearerAuth
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
// @Summary      数据对比
// @Description  根据查询条件进行数据对比分析
// @Tags         统计管理
// @Accept       json
// @Produce      json
// @Param        req  query     dto.CompareRequest  true  "查询参数"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/statistics/compare [get]
// @Security     BearerAuth
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
