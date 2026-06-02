package practice

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	"questionhelper-server/internal/service/practice"
	"questionhelper-server/pkg/database"
	"questionhelper-server/pkg/response"
)

type PracticeController struct{}

func NewPracticeController() *PracticeController {
	return &PracticeController{}
}

// GetPracticeHistory 练习历史
// @Summary      获取练习历史
// @Description  分页获取当前用户的练习历史记录
// @Tags         练习管理
// @Accept       json
// @Produce      json
// @Param        page       query     int  false  "页码"
// @Param        page_size  query     int  false  "每页数量"
// @Success      200  {object}  response.Response{data=dto.PracticeRecordInfo}  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /practice [get]
// @Security     BearerAuth
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
// @Summary      获取练习统计
// @Description  获取当前用户的练习统计数据
// @Tags         练习管理
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /practice/stats [get]
// @Security     BearerAuth
func (ctrl *PracticeController) GetPracticeStats(c *gin.Context) {
	userID := c.GetUint("user_id")

	stats, err := practice.GetPracticeStats(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, stats)
}

// ==================== 模拟考试 ====================

// GetMockExamHistory 模拟考试历史
// @Summary      获取模拟考试历史
// @Description  分页获取当前用户的模拟考试历史记录
// @Tags         模拟考试
// @Accept       json
// @Produce      json
// @Param        page       query     int  false  "页码"
// @Param        page_size  query     int  false  "每页数量"
// @Success      200  {object}  response.Response{data=dto.MockExamHistoryItem}  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /practice/mock/history [get]
// @Security     BearerAuth
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
// @Summary      获取模拟考试详情
// @Description  根据ID获取模拟考试详细信息
// @Tags         模拟考试
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "模拟考试ID"
// @Success      200  {object}  response.Response{data=dto.MockExamResultInfo}  "成功"
// @Failure      400  {object}  response.Response  "无效的模拟考试ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /practice/mock/{id}/detail [get]
// @Security     BearerAuth
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
// @Summary      获取练习计划列表
// @Description  分页获取当前用户的练习计划列表
// @Tags         练习计划
// @Accept       json
// @Produce      json
// @Param        page       query     int  false  "页码"
// @Param        page_size  query     int  false  "每页数量"
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /practice/plans [get]
// @Security     BearerAuth
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
// @Summary      创建练习计划
// @Description  创建一个新的练习计划
// @Tags         练习计划
// @Accept       json
// @Produce      json
// @Param        req  body      dto.CreatePlanRequest  true  "计划信息"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /practice/plans [post]
// @Security     BearerAuth
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
// @Summary      获取计划详情
// @Description  根据ID获取练习计划详细信息
// @Tags         练习计划
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "计划ID"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "无效的计划ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /practice/plans/{id} [get]
// @Security     BearerAuth
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
// @Summary      更新练习计划
// @Description  根据ID更新练习计划信息
// @Tags         练习计划
// @Accept       json
// @Produce      json
// @Param        id   path      uint                   true  "计划ID"
// @Param        req  body      dto.UpdatePlanRequest   true  "更新信息"
// @Success      200  {object}  response.Response  "更新成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /practice/plans/{id} [put]
// @Security     BearerAuth
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
// @Summary      删除练习计划
// @Description  根据ID删除练习计划
// @Tags         练习计划
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "计划ID"
// @Success      200  {object}  response.Response  "删除成功"
// @Failure      400  {object}  response.Response  "无效的计划ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /practice/plans/{id} [delete]
// @Security     BearerAuth
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
// @Summary      执行练习计划
// @Description  根据ID执行练习计划
// @Tags         练习计划
// @Accept       json
// @Produce      json
// @Param        id   path      uint                    true  "计划ID"
// @Param        req  body      dto.ExecutePlanRequest   true  "执行参数"
// @Success      200  {object}  response.Response  "执行成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /practice/plans/{id}/execute [post]
// @Security     BearerAuth
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
// @Summary      获取今日练习
// @Description  获取当前用户的今日练习信息
// @Tags         练习管理
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /practice/daily/today [get]
// @Security     BearerAuth
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
// @Summary      完成今日练习
// @Description  标记今日练习为已完成
// @Tags         练习管理
// @Accept       json
// @Produce      json
// @Param        req  body      dto.CompleteDailyRequest  true  "完成信息"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /practice/daily/complete [post]
// @Security     BearerAuth
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
// @Summary      练习打卡
// @Description  用户进行练习打卡
// @Tags         练习管理
// @Accept       json
// @Produce      json
// @Param        req  body      dto.CheckinRequest  true  "打卡信息"
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /practice/checkin [post]
// @Security     BearerAuth
func (ctrl *PracticeController) Checkin(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.CheckinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	info, err := practice.Checkin(userID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, info)
}

// GetCheckinCalendar 打卡日历
// @Summary      获取打卡日历
// @Description  获取当前用户的打卡日历记录
// @Tags         练习管理
// @Accept       json
// @Produce      json
// @Param        req  query     dto.CheckinCalendarRequest  true  "查询参数"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /practice/checkin/calendar [get]
// @Security     BearerAuth
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
// @Summary      获取排行榜
// @Description  获取练习排行榜数据
// @Tags         练习管理
// @Accept       json
// @Produce      json
// @Param        req  query     dto.LeaderboardRequest  true  "查询参数"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /practice/leaderboard [get]
// @Security     BearerAuth
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
// @Summary      获取关卡列表
// @Description  获取当前用户的闯关关卡列表
// @Tags         闯关模式
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /practice/challenge/levels [get]
// @Security     BearerAuth
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
// @Summary      获取关卡详情
// @Description  根据ID获取闯关关卡详细信息
// @Tags         闯关模式
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "关卡ID"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "无效的关卡ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /practice/challenge/levels/{id} [get]
// @Security     BearerAuth
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

// GetChallengeProgress 闯关进度
// @Summary      获取闯关进度
// @Description  获取当前用户的闯关进度信息
// @Tags         闯关模式
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /practice/challenge/progress [get]
// @Security     BearerAuth
func (ctrl *PracticeController) GetChallengeProgress(c *gin.Context) {
	userID := c.GetUint("user_id")

	progress, err := practice.GetChallengeProgress(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, progress)
}

// ==================== 练习记录导出/搜索 ====================

// ExportPracticeRecords 导出练习记录
func (ctrl *PracticeController) ExportPracticeRecords(c *gin.Context) {
	userID := c.GetUint("user_id")
	var sessions []model.PracticeSession
	database.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&sessions)
	response.Success(c, sessions)
}

// SearchPracticeRecords 搜索练习记录
func (ctrl *PracticeController) SearchPracticeRecords(c *gin.Context) {
	userID := c.GetUint("user_id")
	keyword := c.Query("keyword")
	query := database.DB.Where("user_id = ?", userID)
	if keyword != "" {
		query = query.Joins("LEFT JOIN questions q ON q.id = practice_sessions.category_id").
			Where("q.title LIKE ?", "%"+keyword+"%")
	}
	var sessions []model.PracticeSession
	query.Order("created_at DESC").Find(&sessions)
	response.Success(c, sessions)
}

// ==================== 打卡状态/连续打卡 ====================

// GetCheckinStatus 获取打卡状态
func (ctrl *PracticeController) GetCheckinStatus(c *gin.Context) {
	userID := c.GetUint("user_id")
	var checkin model.PracticeCheckin
	today := time.Now().Format("2006-01-02")
	database.DB.Where("user_id = ? AND date = ?", userID, today).First(&checkin)
	response.Success(c, gin.H{
		"is_checkin": checkin.ID > 0,
		"streak":     checkin.Streak,
		"count":      checkin.QuestionCount,
	})
}

// GetCheckinStreak 获取连续打卡天数
func (ctrl *PracticeController) GetCheckinStreak(c *gin.Context) {
	userID := c.GetUint("user_id")
	var checkin model.PracticeCheckin
	database.DB.Where("user_id = ?", userID).Order("date DESC").First(&checkin)
	response.Success(c, gin.H{"streak": checkin.Streak})
}

// ==================== 排行榜"我的排名" ====================

// GetMyRank 获取我的排名
func (ctrl *PracticeController) GetMyRank(c *gin.Context) {
	userID := c.GetUint("user_id")
	rankType := c.DefaultQuery("rank_type", "1")

	var entry model.PracticeLeaderboard
	today := time.Now().Format("2006-01-02")
	database.DB.Where("user_id = ? AND rank_type = ? AND rank_date = ?", userID, rankType, today).First(&entry)

	var totalRank int64
	database.DB.Model(&model.PracticeLeaderboard{}).Where("rank_type = ? AND rank_date = ?", rankType, today).Count(&totalRank)

	response.Success(c, gin.H{
		"rank_pos":   entry.RankPos,
		"score":      entry.Score,
		"total_rank": totalRank,
	})
}
