package exam

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/service/exam"
	"questionhelper-server/pkg/response"
)

type AnswerController struct{}

func NewAnswerController() *AnswerController {
	return &AnswerController{}
}

// SaveAnswer 保存答案
func (ctrl *AnswerController) SaveAnswer(c *gin.Context) {
	recordID, err := strconv.ParseUint(c.Param("recordId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的考试记录ID")
		return
	}

	var req dto.SaveAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := exam.SaveAnswer(uint(recordID), &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "保存成功", nil)
}

// SaveAnswers 批量保存答案
func (ctrl *AnswerController) SaveAnswers(c *gin.Context) {
	recordID, err := strconv.ParseUint(c.Param("recordId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的考试记录ID")
		return
	}

	var req dto.SaveAnswerBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := exam.SaveAnswers(uint(recordID), &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "保存成功", nil)
}

// GetStandardAnswers 获取标准答案
func (ctrl *AnswerController) GetStandardAnswers(c *gin.Context) {
	userID := c.GetUint("user_id")

	examID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的考试ID")
		return
	}

	result, err := exam.GetStandardAnswers(uint(examID), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, result)
}

// MarkQuestion 标记题目
func (ctrl *AnswerController) MarkQuestion(c *gin.Context) {
	recordID, err := strconv.ParseUint(c.Param("recordId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的考试记录ID")
		return
	}

	var req dto.MarkQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := exam.MarkQuestion(uint(recordID), req.QuestionID, req.IsMarked); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "操作成功", nil)
}

// GetMarkedQuestions 获取标记题目
func (ctrl *AnswerController) GetMarkedQuestions(c *gin.Context) {
	recordID, err := strconv.ParseUint(c.Param("recordId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的考试记录ID")
		return
	}

	result, err := exam.GetMarkedQuestions(uint(recordID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, result)
}

// GetExamGuide 获取答题指引
func (ctrl *AnswerController) GetExamGuide(c *gin.Context) {
	examID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的考试ID")
		return
	}

	result, err := exam.GetExamGuide(uint(examID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, result)
}

// SubmitFeedback 提交考后反馈
func (ctrl *AnswerController) SubmitFeedback(c *gin.Context) {
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

	if err := exam.SubmitFeedback(uint(examID), userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "反馈提交成功", nil)
}

// ReportWarning 上报异常行为
func (ctrl *AnswerController) ReportWarning(c *gin.Context) {
	userID := c.GetUint("user_id")

	recordID, err := strconv.ParseUint(c.Param("recordId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的考试记录ID")
		return
	}

	var req struct {
		Type   string `json:"type" binding:"required"`
		Detail string `json:"detail"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	// 获取考试ID
	record, err := exam.FindRecordByID(uint(recordID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "考试记录不存在")
		return
	}

	if err := exam.RecordWarning(uint(recordID), record.ExamID, userID, req.Type, req.Detail, c.ClientIP(), c.GetHeader("User-Agent")); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "上报成功", nil)
}
