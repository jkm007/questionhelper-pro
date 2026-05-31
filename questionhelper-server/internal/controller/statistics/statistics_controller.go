package statistics

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/service/statistics"
	"questionhelper-server/pkg/response"
)

type StatisticsController struct{}

func NewStatisticsController() *StatisticsController {
	return &StatisticsController{}
}

// ==================== 基础统计 ====================

// GetUserStatistics 用户统计
func (ctrl *StatisticsController) GetUserStatistics(c *gin.Context) {
	stats, err := statistics.GetUserStatistics()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, stats)
}

// GetPracticeStatistics 练习统计
func (ctrl *StatisticsController) GetPracticeStatistics(c *gin.Context) {
	stats, err := statistics.GetPracticeStatistics()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, stats)
}

// GetExamStatistics 考试统计
func (ctrl *StatisticsController) GetExamStatistics(c *gin.Context) {
	stats, err := statistics.GetExamStatistics()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, stats)
}

// GetRanking 排行榜
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

// GetClassStatistics 班级统计
func (ctrl *StatisticsController) GetClassStatistics(c *gin.Context) {
	stats, err := statistics.GetClassStatistics()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, stats)
}

// GetDashboard 仪表盘（管理员）
func (ctrl *StatisticsController) GetDashboard(c *gin.Context) {
	dashboard, err := statistics.GetDashboard()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, dashboard)
}

// ==================== 用户行为事件(普通用户上报) ====================

// CreateEvent 上报用户行为事件
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

// ==================== 题目分析 ====================

// GetQuestionDifficulty 题目难度分析
func (ctrl *StatisticsController) GetQuestionDifficulty(c *gin.Context) {
	var req dto.QuestionDifficultyRequest
	c.ShouldBindQuery(&req)

	items, err := statistics.GetQuestionDifficulty(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, items)
}

// GetQuestionDiscrimination 题目区分度分析
func (ctrl *StatisticsController) GetQuestionDiscrimination(c *gin.Context) {
	var req dto.QuestionDiscriminationRequest
	c.ShouldBindQuery(&req)

	items, err := statistics.GetQuestionDiscrimination(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, items)
}

// ==================== 成绩预测与预警 ====================

// GetScorePrediction 成绩预测
func (ctrl *StatisticsController) GetScorePrediction(c *gin.Context) {
	var req dto.ScorePredictionRequest
	c.ShouldBindQuery(&req)

	userID := c.GetUint("user_id")
	result, err := statistics.GetScorePrediction(userID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, result)
}

// GetScoreAlert 成绩预警
func (ctrl *StatisticsController) GetScoreAlert(c *gin.Context) {
	var req dto.ScoreAlertRequest
	c.ShouldBindQuery(&req)

	userID := c.GetUint("user_id")
	items, err := statistics.GetScoreAlert(userID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, items)
}

// ==================== 班级统计(教师视角) ====================

// GetClassOverview 班级概览
func (ctrl *StatisticsController) GetClassOverview(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	overview, err := statistics.GetClassOverview(uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, overview)
}

// GetClassStudents 班级学生成绩列表
func (ctrl *StatisticsController) GetClassStudents(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	var req dto.ClassStudentListRequest
	c.ShouldBindQuery(&req)

	list, total, err := statistics.GetClassStudents(uint(id), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// GetClassPracticeStats 班级练习统计
func (ctrl *StatisticsController) GetClassPracticeStats(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	var req dto.ClassPracticeRequest
	c.ShouldBindQuery(&req)

	result, err := statistics.GetClassPracticeStats(uint(id), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, result)
}

// GetClassExamStats 班级考试统计
func (ctrl *StatisticsController) GetClassExamStats(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	var req dto.ClassExamRequest
	c.ShouldBindQuery(&req)

	result, err := statistics.GetClassExamStats(uint(id), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, result)
}

// GetClassQuestionStats 班级题目统计
func (ctrl *StatisticsController) GetClassQuestionStats(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	result, err := statistics.GetClassQuestionStats(uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, result)
}

// ==================== 移动端统计 ====================

// GetMobileOverview 移动端个人概览
func (ctrl *StatisticsController) GetMobileOverview(c *gin.Context) {
	userID := c.GetUint("user_id")

	overview, err := statistics.GetMobileOverview(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, overview)
}

// GetMobilePracticeStats 移动端练习统计
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
