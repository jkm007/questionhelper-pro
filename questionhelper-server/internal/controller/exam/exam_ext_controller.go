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
// @Summary      试卷共享
// @Description  将试卷共享给其他用户
// @Tags         考试扩展
// @Accept       json
// @Produce      json
// @Param        id   path      uint                      true  "试卷ID"
// @Param        req  body      dto.SharePaperRequest     true  "共享信息"
// @Success      200  {object}  response.Response  "共享成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/papers/{id}/share [post]
// @Security     BearerAuth
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
// @Summary      收藏试卷
// @Description  收藏或取消收藏试卷
// @Tags         考试扩展
// @Accept       json
// @Produce      json
// @Param        id   path      uint                        true  "试卷ID"
// @Param        req  body      dto.FavoritePaperRequest    true  "收藏信息"
// @Success      200  {object}  response.Response  "操作成功"
// @Failure      400  {object}  response.Response  "无效的试卷ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/papers/{id}/favorite [post]
// @Security     BearerAuth
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
// @Summary      导入试卷
// @Description  从JSON数据导入试卷
// @Tags         考试扩展
// @Accept       json
// @Produce      json
// @Param        req  body      dto.ExportPaperResponse  true  "试卷数据"
// @Success      200  {object}  response.Response{data=dto.PaperInfo}  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/papers/import [post]
// @Security     BearerAuth
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
// @Summary      延长考试时间
// @Description  管理员延长指定考试的时间
// @Tags         考试扩展
// @Accept       json
// @Produce      json
// @Param        id   path      uint                      true  "考试ID"
// @Param        req  body      dto.ExtendExamRequest     true  "延长时间信息"
// @Success      200  {object}  response.Response  "延长成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/exams/{id}/extend [post]
// @Security     BearerAuth
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
// @Summary      暂停考试
// @Description  管理员暂停指定考试
// @Tags         考试扩展
// @Accept       json
// @Produce      json
// @Param        id   path      uint                     true  "考试ID"
// @Param        req  body      dto.PauseExamRequest     true  "暂停原因"
// @Success      200  {object}  response.Response  "考试已暂停"
// @Failure      400  {object}  response.Response  "无效的考试ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/exams/{id}/pause [post]
// @Security     BearerAuth
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
// @Summary      恢复考试
// @Description  管理员恢复已暂停的考试
// @Tags         考试扩展
// @Accept       json
// @Produce      json
// @Param        id   path      uint                      true  "考试ID"
// @Param        req  body      dto.ResumeExamRequest     true  "恢复原因"
// @Success      200  {object}  response.Response  "考试已恢复"
// @Failure      400  {object}  response.Response  "无效的考试ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/exams/{id}/resume [post]
// @Security     BearerAuth
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
// @Summary      获取考试成绩列表
// @Description  获取指定考试的成绩列表，支持分页（管理员）
// @Tags         成绩管理
// @Accept       json
// @Produce      json
// @Param        id         path      uint  true   "考试ID"
// @Param        page       query     int   false  "页码"
// @Param        page_size  query     int   false  "每页数量"
// @Success      200  {object}  response.PageResponse{data=[]dto.ExamScoreInfo}  "成功"
// @Failure      400  {object}  response.Response  "无效的考试ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/exams/{id}/scores [get]
// @Security     BearerAuth
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
// @Summary      导出考试成绩
// @Description  导出指定考试的成绩数据（管理员）
// @Tags         成绩管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "考试ID"
// @Success      200  {object}  dto.ExportScoresResponse  "成功"
// @Failure      400  {object}  response.Response  "无效的考试ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/exams/{id}/scores/export [get]
// @Security     BearerAuth
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
// @Summary      获取考试统计
// @Description  获取指定考试的统计数据（管理员）
// @Tags         成绩管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "考试ID"
// @Success      200  {object}  response.Response{data=dto.ExamStatisticsInfo}  "成功"
// @Failure      400  {object}  response.Response  "无效的考试ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/exams/{id}/statistics [get]
// @Security     BearerAuth
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
// @Summary      申请成绩复核
// @Description  学生申请对指定考试成绩进行复核
// @Tags         成绩管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint                        true  "考试ID"
// @Param        req  body      dto.ScoreReviewRequest      true  "复核申请信息"
// @Success      200  {object}  response.Response  "复核申请已提交"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /exam/{id}/score-review [post]
// @Security     BearerAuth
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
// @Summary      获取复核申请列表
// @Description  获取成绩复核申请列表，支持分页（管理员）
// @Tags         成绩管理
// @Accept       json
// @Produce      json
// @Param        page       query     int  false  "页码"
// @Param        page_size  query     int  false  "每页数量"
// @Success      200  {object}  response.PageResponse{data=[]dto.ScoreReviewInfo}  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/score-reviews [get]
// @Security     BearerAuth
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
// @Summary      处理成绩复核
// @Description  管理员处理成绩复核申请
// @Tags         成绩管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint                            true  "复核ID"
// @Param        req  body      dto.HandleScoreReviewRequest    true  "处理结果"
// @Success      200  {object}  response.Response  "处理成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/score-reviews/{id} [put]
// @Security     BearerAuth
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
// @Summary      创建考试公告
// @Description  管理员创建考试公告通知
// @Tags         考试扩展
// @Accept       json
// @Produce      json
// @Param        id   path      uint                           true  "考试ID"
// @Param        req  body      dto.CreateExamNoticeRequest    true  "公告信息"
// @Success      200  {object}  response.Response  "公告发布成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/exams/{id}/notice [post]
// @Security     BearerAuth
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
// @Summary      获取考试公告列表
// @Description  获取指定考试的公告列表，支持分页
// @Tags         考试扩展
// @Accept       json
// @Produce      json
// @Param        id         path      uint  true   "考试ID"
// @Param        page       query     int   false  "页码"
// @Param        page_size  query     int   false  "每页数量"
// @Success      200  {object}  response.PageResponse{data=[]dto.ExamNoticeInfo}  "成功"
// @Failure      400  {object}  response.Response  "无效的考试ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/exams/{id}/notices [get]
// @Security     BearerAuth
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
// @Summary      上报切屏事件
// @Description  学生在考试过程中切换屏幕时上报
// @Tags         考试扩展
// @Accept       json
// @Produce      json
// @Param        id   path      uint                        true  "考试ID"
// @Param        req  body      dto.SwitchScreenRequest     true  "切屏详情"
// @Success      200  {object}  response.Response  "上报成功"
// @Failure      400  {object}  response.Response  "无效的考试ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /exam/{id}/switch-screen [post]
// @Security     BearerAuth
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
// @Summary      断线续考
// @Description  学生断线后恢复考试状态继续答题
// @Tags         考试扩展
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "考试ID"
// @Success      200  {object}  response.Response{data=dto.ExamRecordInfo}  "成功"
// @Failure      400  {object}  response.Response  "无效的考试ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /exam/{id}/resume [get]
// @Security     BearerAuth
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
// @Summary      获取即将开始的考试
// @Description  获取当前学生即将开始的考试列表，支持分页
// @Tags         考试扩展
// @Accept       json
// @Produce      json
// @Param        page       query     int  false  "页码"
// @Param        page_size  query     int  false  "每页数量"
// @Success      200  {object}  response.PageResponse{data=[]dto.ExamInfo}  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /exams/upcoming [get]
// @Security     BearerAuth
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
// @Summary      获取考试成绩排名
// @Description  获取指定考试的成绩排名，支持分页
// @Tags         考试扩展
// @Accept       json
// @Produce      json
// @Param        id         path      uint  true   "考试ID"
// @Param        page       query     int   false  "页码"
// @Param        page_size  query     int   false  "每页数量"
// @Success      200  {object}  response.PageResponse{data=[]dto.ExamRankingInfo}  "成功"
// @Failure      400  {object}  response.Response  "无效的考试ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /exams/{id}/rankings [get]
// @Security     BearerAuth
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
// @Summary      提交考试反馈
// @Description  学生提交考试反馈信息
// @Tags         考试扩展
// @Accept       json
// @Produce      json
// @Param        id   path      uint                  true  "考试ID"
// @Param        req  body      dto.FeedbackRequest   true  "反馈信息"
// @Success      200  {object}  response.Response  "反馈提交成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /exam/{id}/feedback [post]
// @Security     BearerAuth
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
