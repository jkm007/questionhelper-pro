package exam

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/service/exam"
	"questionhelper-server/pkg/response"
)

type ExamExtController struct{}

func NewExamExtController() *ExamExtController {
	return &ExamExtController{}
}

// ==================== 试卷共享/收藏 ====================

// SharePaper 试卷共享
func (ctrl *ExamExtController) SharePaper(c *gin.Context) {
	userID := c.GetUint("user_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的试卷ID")
		return
	}

	var req dto.SharePaperRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := exam.SharePaper(uint(id), userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "共享成功", nil)
}

// FavoritePaper 收藏试卷
func (ctrl *ExamExtController) FavoritePaper(c *gin.Context) {
	userID := c.GetUint("user_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的试卷ID")
		return
	}

	var req dto.FavoritePaperRequest
	c.ShouldBindJSON(&req)

	msg, err := exam.FavoritePaper(userID, uint(id), req.Note)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, msg, nil)
}

// ImportPaper 导入试卷
func (ctrl *ExamExtController) ImportPaper(c *gin.Context) {
	creatorID := c.GetUint("user_id")

	var req dto.ExportPaperResponse
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	result, err := exam.ImportPaper(creatorID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, result)
}

// ==================== 考试操作 ====================

// ExtendExam 延长考试
func (ctrl *ExamExtController) ExtendExam(c *gin.Context) {
	operatorID := c.GetUint("user_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的考试ID")
		return
	}

	var req dto.ExtendExamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := exam.ExtendExam(uint(id), operatorID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "延长成功", nil)
}

// PauseExam 暂停考试
func (ctrl *ExamExtController) PauseExam(c *gin.Context) {
	operatorID := c.GetUint("user_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的考试ID")
		return
	}

	var req dto.PauseExamRequest
	c.ShouldBindJSON(&req)

	if err := exam.PauseExam(uint(id), operatorID, req.Reason); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "考试已暂停", nil)
}

// ResumeExam 恢复考试
func (ctrl *ExamExtController) ResumeExam(c *gin.Context) {
	operatorID := c.GetUint("user_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的考试ID")
		return
	}

	var req dto.ResumeExamRequest
	c.ShouldBindJSON(&req)

	if err := exam.ResumeExam(uint(id), operatorID, req.Reason); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "考试已恢复", nil)
}

// ==================== 成绩管理 ====================

// GetExamScores 考试成绩列表
func (ctrl *ExamExtController) GetExamScores(c *gin.Context) {
	examID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的考试ID")
		return
	}

	var req dto.PageRequest
	c.ShouldBindQuery(&req)

	list, total, err := exam.GetExamScores(uint(examID), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// ExportExamScores 导出考试成绩
func (ctrl *ExamExtController) ExportExamScores(c *gin.Context) {
	examID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的考试ID")
		return
	}

	result, err := exam.ExportExamScores(uint(examID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Header("Content-Type", "application/json")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=scores_%d.json", examID))
	c.JSON(http.StatusOK, result)
}

// GetExamStatistics 考试统计
func (ctrl *ExamExtController) GetExamStatistics(c *gin.Context) {
	examID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的考试ID")
		return
	}

	result, err := exam.GetExamStatistics(uint(examID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, result)
}

// ==================== 成绩复核 ====================

// SubmitScoreReview 申请成绩复核
func (ctrl *ExamExtController) SubmitScoreReview(c *gin.Context) {
	userID := c.GetUint("user_id")

	examID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的考试ID")
		return
	}

	var req dto.ScoreReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := exam.SubmitScoreReview(userID, uint(examID), &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "复核申请已提交", nil)
}

// ListScoreReviews 复核申请列表
func (ctrl *ExamExtController) ListScoreReviews(c *gin.Context) {
	var req dto.ScoreReviewListRequest
	c.ShouldBindQuery(&req)

	list, total, err := exam.ListScoreReviews(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// HandleScoreReview 处理复核
func (ctrl *ExamExtController) HandleScoreReview(c *gin.Context) {
	reviewerID := c.GetUint("user_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的复核ID")
		return
	}

	var req dto.HandleScoreReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := exam.HandleScoreReview(uint(id), reviewerID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "处理成功", nil)
}

// ==================== 考试公告 ====================

// CreateExamNotice 创建考试公告
func (ctrl *ExamExtController) CreateExamNotice(c *gin.Context) {
	creatorID := c.GetUint("user_id")

	examID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的考试ID")
		return
	}

	var req dto.CreateExamNoticeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := exam.CreateExamNotice(uint(examID), creatorID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "公告发布成功", nil)
}

// ListExamNotices 考试公告列表
func (ctrl *ExamExtController) ListExamNotices(c *gin.Context) {
	examID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的考试ID")
		return
	}

	var req dto.PageRequest
	c.ShouldBindQuery(&req)

	list, total, err := exam.ListExamNotices(uint(examID), req.Page, req.PageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// ==================== 考试防作弊 ====================

// ReportSwitchScreen 上报切屏
func (ctrl *ExamExtController) ReportSwitchScreen(c *gin.Context) {
	userID := c.GetUint("user_id")

	examID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的考试ID")
		return
	}

	var req dto.SwitchScreenRequest
	c.ShouldBindJSON(&req)

	// 获取用户的考试记录
	record, findErr := exam.FindRecordByUser(uint(examID), userID)
	if findErr != nil {
		response.Error(c, http.StatusInternalServerError, "未找到考试记录")
		return
	}

	if err := exam.ReportSwitchScreen(record.ID, uint(examID), userID, req.Detail, c.ClientIP(), c.GetHeader("User-Agent")); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "上报成功", nil)
}

// ResumeExamStudent 断线续考
func (ctrl *ExamExtController) ResumeExamStudent(c *gin.Context) {
	userID := c.GetUint("user_id")

	examID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的考试ID")
		return
	}

	result, err := exam.ResumeExamForStudent(uint(examID), userID, c.ClientIP())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, result)
}

// ==================== 考试查询增强 ====================

// ListUpcomingExams 即将开始的考试
func (ctrl *ExamExtController) ListUpcomingExams(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.PageRequest
	c.ShouldBindQuery(&req)

	list, total, err := exam.ListUpcomingExams(userID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// GetExamRankings 成绩排名
func (ctrl *ExamExtController) GetExamRankings(c *gin.Context) {
	examID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的考试ID")
		return
	}

	var req dto.PageRequest
	c.ShouldBindQuery(&req)

	list, total, err := exam.GetExamRankings(uint(examID), req.Page, req.PageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// SubmitExamFeedback 提交考试反馈
func (ctrl *ExamExtController) SubmitExamFeedback(c *gin.Context) {
	userID := c.GetUint("user_id")

	examID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的考试ID")
		return
	}

	var req dto.FeedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := exam.SubmitExamFeedback(uint(examID), userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "反馈提交成功", nil)
}
