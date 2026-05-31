package practice

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/service/practice"
	"questionhelper-server/pkg/response"
)

type PracticeController struct{}

func NewPracticeController() *PracticeController {
	return &PracticeController{}
}

// StartPractice 开始练习
func (ctrl *PracticeController) StartPractice(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.StartPracticeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	session, err := practice.StartPractice(userID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, session)
}

// SubmitPractice 提交练习答案
func (ctrl *PracticeController) SubmitPractice(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.SubmitPracticeAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	// 从请求中获取session_id
	sessionIDStr := c.Query("session_id")
	if sessionIDStr == "" {
		response.Error(c, http.StatusBadRequest, "请提供session_id")
		return
	}

	sessionID, err := strconv.ParseUint(sessionIDStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的session_id")
		return
	}

	if err := practice.SubmitPractice(uint(sessionID), userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "提交成功", nil)
}

// GetPracticeResult 获取练习结果
func (ctrl *PracticeController) GetPracticeResult(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的练习ID")
		return
	}

	result, err := practice.GetPracticeResult(uint(id), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, result)
}

// GetPracticeHistory 练习历史
func (ctrl *PracticeController) GetPracticeHistory(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.PracticeListRequest
	c.ShouldBindQuery(&req)

	list, total, err := practice.GetPracticeHistory(userID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// GetPracticeStats 练习统计
func (ctrl *PracticeController) GetPracticeStats(c *gin.Context) {
	userID := c.GetUint("user_id")

	stats, err := practice.GetPracticeStats(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, stats)
}

// FinishPractice 完成练习
func (ctrl *PracticeController) FinishPractice(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的练习ID")
		return
	}

	if err := practice.FinishPractice(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "练习完成", nil)
}

// ==================== 模拟考试 ====================

// StartMockExam 开始模拟考试
func (ctrl *PracticeController) StartMockExam(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.StartMockExamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	session, err := practice.StartMockExam(userID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, session)
}

// SubmitMockExam 提交模拟考试
func (ctrl *PracticeController) SubmitMockExam(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的模拟考试ID")
		return
	}

	var req dto.SubmitMockExamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	result, err := practice.SubmitMockExam(uint(id), userID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, result)
}

// GetMockExamHistory 模拟考试历史
func (ctrl *PracticeController) GetMockExamHistory(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.MockExamHistoryRequest
	c.ShouldBindQuery(&req)

	list, total, err := practice.GetMockExamHistory(userID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// GetMockExamDetail 模拟考试详情
func (ctrl *PracticeController) GetMockExamDetail(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的模拟考试ID")
		return
	}

	result, err := practice.GetMockExamDetail(uint(id), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, result)
}

// ==================== 练习计划 ====================

// GetPlans 获取练习计划列表
func (ctrl *PracticeController) GetPlans(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.PlanListRequest
	c.ShouldBindQuery(&req)

	list, total, err := practice.ListPlans(userID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// CreatePlan 创建练习计划
func (ctrl *PracticeController) CreatePlan(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.CreatePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	plan, err := practice.CreatePlan(userID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, plan)
}

// GetPlan 获取计划详情
func (ctrl *PracticeController) GetPlan(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的计划ID")
		return
	}

	plan, err := practice.GetPlan(uint(id), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, plan)
}

// UpdatePlan 更新计划
func (ctrl *PracticeController) UpdatePlan(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的计划ID")
		return
	}

	var req dto.UpdatePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := practice.UpdatePlan(uint(id), userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "更新成功", nil)
}

// DeletePlan 删除计划
func (ctrl *PracticeController) DeletePlan(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的计划ID")
		return
	}

	if err := practice.DeletePlan(uint(id), userID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "删除成功", nil)
}

// ExecutePlan 执行计划
func (ctrl *PracticeController) ExecutePlan(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的计划ID")
		return
	}

	var req dto.ExecutePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := practice.ExecutePlan(uint(id), userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "执行成功", nil)
}

// ==================== 每日练习 ====================

// GetTodayPractice 获取今日练习
func (ctrl *PracticeController) GetTodayPractice(c *gin.Context) {
	userID := c.GetUint("user_id")

	info, err := practice.GetTodayPractice(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, info)
}

// CompleteDailyPractice 完成今日练习
func (ctrl *PracticeController) CompleteDailyPractice(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.CompleteDailyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	info, err := practice.CompleteDailyPractice(userID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, info)
}

// ==================== 练习打卡 ====================

// Checkin 打卡
func (ctrl *PracticeController) Checkin(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.CheckinRequest
	c.ShouldBindJSON(&req)

	info, err := practice.Checkin(userID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, info)
}

// GetCheckinCalendar 打卡日历
func (ctrl *PracticeController) GetCheckinCalendar(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.CheckinCalendarRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	items, err := practice.GetCheckinCalendar(userID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, items)
}

// ==================== 排行榜 ====================

// GetLeaderboard 排行榜
func (ctrl *PracticeController) GetLeaderboard(c *gin.Context) {
	var req dto.LeaderboardRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	items, err := practice.GetLeaderboard(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, items)
}

// ==================== 闯关模式 ====================

// GetChallengeLevels 关卡列表
func (ctrl *PracticeController) GetChallengeLevels(c *gin.Context) {
	userID := c.GetUint("user_id")

	levels, err := practice.GetChallengeLevels(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, levels)
}

// GetChallengeLevel 关卡详情
func (ctrl *PracticeController) GetChallengeLevel(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的关卡ID")
		return
	}

	level, err := practice.GetChallengeLevel(uint(id), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, level)
}

// StartChallenge 开始闯关
func (ctrl *PracticeController) StartChallenge(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的关卡ID")
		return
	}

	session, err := practice.StartChallenge(userID, uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, session)
}

// SubmitChallenge 提交闯关
func (ctrl *PracticeController) SubmitChallenge(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的关卡ID")
		return
	}

	var req dto.SubmitChallengeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	result, err := practice.SubmitChallenge(uint(id), userID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, result)
}

// GetChallengeProgress 闯关进度
func (ctrl *PracticeController) GetChallengeProgress(c *gin.Context) {
	userID := c.GetUint("user_id")

	progress, err := practice.GetChallengeProgress(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, progress)
}
