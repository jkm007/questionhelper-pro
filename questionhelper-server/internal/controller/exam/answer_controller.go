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
// @Summary      保存单个答案
// @Description  保存学生对单个题目的作答
// @Tags         答案管理
// @Accept       json
// @Produce      json
// @Param        recordId  path      uint                   true  "考试记录ID"
// @Param        req       body      dto.SaveAnswerRequest  true  "答案信息"
// @Success      200  {object}  response.Response  "保存成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /exam-records/{recordId}/save-answer [post]
// @Security     BearerAuth
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
// @Summary      批量保存答案
// @Description  批量保存学生对多个题目的作答
// @Tags         答案管理
// @Accept       json
// @Produce      json
// @Param        recordId  path      uint                        true  "考试记录ID"
// @Param        req       body      dto.SaveAnswerBatchRequest  true  "批量答案信息"
// @Success      200  {object}  response.Response  "保存成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /exam-records/{recordId}/save-answers [post]
// @Security     BearerAuth
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
// @Summary      获取标准答案
// @Description  获取指定考试的标准答案
// @Tags         答案管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "考试ID"
// @Success      200  {object}  response.Response{data=[]dto.StandardAnswer}  "成功"
// @Failure      400  {object}  response.Response  "无效的考试ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /exam/{id}/standard-answers [get]
// @Security     BearerAuth
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
// @Summary      标记题目
// @Description  标记或取消标记考试中的题目
// @Tags         答案管理
// @Accept       json
// @Produce      json
// @Param        recordId  path      uint                      true  "考试记录ID"
// @Param        req       body      dto.MarkQuestionRequest   true  "标记信息"
// @Success      200  {object}  response.Response  "操作成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /exam-records/{recordId}/mark [post]
// @Security     BearerAuth
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
// @Summary      获取标记题目列表
// @Description  获取学生在考试中标记的题目列表
// @Tags         答案管理
// @Accept       json
// @Produce      json
// @Param        recordId  path      uint  true  "考试记录ID"
// @Success      200  {object}  response.Response{data=[]dto.SaveAnswerRequest}  "成功"
// @Failure      400  {object}  response.Response  "无效的考试记录ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /exam-records/{recordId}/marked [get]
// @Security     BearerAuth
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
// @Summary      获取答题指引
// @Description  获取指定考试的答题指引信息
// @Tags         答案管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "考试ID"
// @Success      200  {object}  response.Response{data=dto.ExamGuideResponse}  "成功"
// @Failure      400  {object}  response.Response  "无效的考试ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /exam/{id}/guide [get]
// @Security     BearerAuth
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
// @Summary      提交考后反馈
// @Description  学生提交考试后的反馈信息
// @Tags         答案管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint                  true  "考试ID"
// @Param        req  body      dto.FeedbackRequest   true  "反馈信息"
// @Success      200  {object}  response.Response  "反馈提交成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /exam/{id}/feedback [post]
// @Security     BearerAuth
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
// @Summary      上报异常行为
// @Description  上报学生在考试中的异常行为（如切屏、多设备登录等）
// @Tags         答案管理
// @Accept       json
// @Produce      json
// @Param        recordId  path      uint                          true  "考试记录ID"
// @Param        req       body      object{type=string,detail=string}  true  "异常信息"
// @Success      200  {object}  response.Response  "上报成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /exam-records/{recordId}/warning [post]
// @Security     BearerAuth
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
