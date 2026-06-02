package statistics

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/service/statistics"
	"questionhelper-server/pkg/response"
)

// GetUserStatistics 用户统计
// @Summary      获取用户统计
// @Description  获取用户统计数据概览
// @Tags         数据统计
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /statistics/user [get]
// @Security     BearerAuth
func (ctrl *StatisticsController) GetUserStatistics(c *gin.Context) {
	stats, err := statistics.GetUserStatistics()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, stats)
}

// GetPracticeStatistics 练习统计
// @Summary      获取练习统计
// @Description  获取练习统计数据概览
// @Tags         数据统计
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /statistics/practice [get]
// @Security     BearerAuth
func (ctrl *StatisticsController) GetPracticeStatistics(c *gin.Context) {
	stats, err := statistics.GetPracticeStatistics()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, stats)
}

// GetRanking 排行榜
// @Summary      获取排行榜
// @Description  获取用户排行榜数据，支持分页
// @Tags         数据统计
// @Accept       json
// @Produce      json
// @Param        page       query     int  false  "页码"
// @Param        page_size  query     int  false  "每页数量"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /statistics/ranking [get]
// @Security     BearerAuth
func (ctrl *StatisticsController) GetRanking(c *gin.Context) {
	var req dto.RankingRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	ranks, total, err := statistics.GetRanking(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, ranks, total, req.Page, req.PageSize)
}

// ==================== 用户行为事件(普通用户上报) ====================

// CreateEvent 上报用户行为事件
// @Summary      上报用户行为事件
// @Description  普通用户上报行为事件数据
// @Tags         数据统计
// @Accept       json
// @Produce      json
// @Param        req  body      dto.CreateEventRequest  true  "事件数据"
// @Success      200  {object}  response.Response  "上报成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /statistics/events [post]
// @Security     BearerAuth
func (ctrl *StatisticsController) CreateEvent(c *gin.Context) {
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

// ==================== 数据订阅 ====================

// ListSubscriptions 订阅列表
// @Summary      获取订阅列表
// @Description  获取当前用户的数据订阅列表，支持分页
// @Tags         数据统计
// @Accept       json
// @Produce      json
// @Param        page       query     int  false  "页码"
// @Param        page_size  query     int  false  "每页数量"
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /statistics/subscriptions [get]
// @Security     BearerAuth
func (ctrl *StatisticsController) ListSubscriptions(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.PageRequest
	c.ShouldBindQuery(&req)

	list, total, err := statistics.ListSubscriptions(userID, req.Page, req.PageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// CreateSubscription 创建订阅
// @Summary      创建数据订阅
// @Description  创建一条数据订阅
// @Tags         数据统计
// @Accept       json
// @Produce      json
// @Param        req  body      dto.CreateSubscriptionRequest  true  "订阅数据"
// @Success      200  {object}  response.Response  "创建成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /statistics/subscriptions [post]
// @Security     BearerAuth
func (ctrl *StatisticsController) CreateSubscription(c *gin.Context) {
	var req dto.CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	userID := c.GetUint("user_id")
	if err := statistics.CreateSubscription(userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "创建成功", nil)
}

// DeleteSubscription 取消订阅
// @Summary      取消数据订阅
// @Description  根据ID取消一条数据订阅
// @Tags         数据统计
// @Accept       json
// @Produce      json
// @Param        id  path      uint  true  "订阅ID"
// @Success      200  {object}  response.Response  "取消成功"
// @Failure      400  {object}  response.Response  "无效的订阅ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /statistics/subscriptions/{id} [delete]
// @Security     BearerAuth
func (ctrl *StatisticsController) DeleteSubscription(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的订阅ID")
		return
	}

	userID := c.GetUint("user_id")
	if err := statistics.DeleteSubscription(userID, uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "取消成功", nil)
}

// ==================== 移动端练习统计 ====================

// GetMobilePracticeStats 移动端练习统计
// @Summary      获取移动端练习统计
// @Description  获取移动端当前用户的练习统计数据
// @Tags         数据统计
// @Accept       json
// @Produce      json
// @Param        req  query     dto.MobilePracticeRequest  false  "查询参数"
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /statistics/mobile/practice [get]
// @Security     BearerAuth
func (ctrl *StatisticsController) GetMobilePracticeStats(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.MobilePracticeRequest
	c.ShouldBindQuery(&req)

	result, err := statistics.GetMobilePracticeStats(userID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, result)
}

// GetMobileWrongStats 移动端错题统计
// @Summary      获取移动端错题统计
// @Description  获取移动端当前用户的错题统计数据
// @Tags         数据统计
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /statistics/mobile/wrong [get]
// @Security     BearerAuth
func (ctrl *StatisticsController) GetMobileWrongStats(c *gin.Context) {
	userID := c.GetUint("user_id")

	result, err := statistics.GetMobileWrongStats(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, result)
}

// GetMobileTrend 移动端学习趋势
// @Summary      获取移动端学习趋势
// @Description  获取移动端当前用户的学习趋势数据
// @Tags         数据统计
// @Accept       json
// @Produce      json
// @Param        req  query     dto.MobileTrendRequest  false  "查询参数"
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /statistics/mobile/trend [get]
// @Security     BearerAuth
func (ctrl *StatisticsController) GetMobileTrend(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.MobileTrendRequest
	c.ShouldBindQuery(&req)

	items, err := statistics.GetMobileTrend(userID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, items)
}
