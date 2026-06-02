package practice

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/service/practice"
	"questionhelper-server/pkg/response"
)

// StartPractice 开始练习
// @Summary      开始练习
// @Description  创建一个新的练习会话
// @Tags         练习管理
// @Accept       json
// @Produce      json
// @Param        req  body      dto.StartPracticeRequest  true  "练习配置"
// @Success      200  {object}  response.Response{data=dto.PracticeSessionInfo}  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /practice/start [post]
// @Security     BearerAuth
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
// @Summary      提交练习答案
// @Description  提交练习会话的答案
// @Tags         练习管理
// @Accept       json
// @Produce      json
// @Param        session_id  query     uint                            true  "会话ID"
// @Param        req         body      dto.SubmitPracticeAnswerRequest  true  "答案信息"
// @Success      200  {object}  response.Response  "提交成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /practice/submit [post]
// @Security     BearerAuth
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
// @Summary      获取练习结果
// @Description  根据ID获取练习结果详情
// @Tags         练习管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "练习ID"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "无效的练习ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /practice/{id} [get]
// @Security     BearerAuth
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

// FinishPractice 完成练习
// @Summary      完成练习
// @Description  根据ID完成练习会话
// @Tags         练习管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "练习ID"
// @Success      200  {object}  response.Response  "练习完成"
// @Failure      400  {object}  response.Response  "无效的练习ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /practice/{id}/finish [post]
// @Security     BearerAuth
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

// StartMockExam 开始模拟考试
// @Summary      开始模拟考试
// @Description  创建一个新的模拟考试会话
// @Tags         模拟考试
// @Accept       json
// @Produce      json
// @Param        req  body      dto.StartMockExamRequest  true  "考试配置"
// @Success      200  {object}  response.Response{data=dto.MockExamSessionInfo}  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /practice/mock/start [post]
// @Security     BearerAuth
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
// @Summary      提交模拟考试
// @Description  根据ID提交模拟考试答案
// @Tags         模拟考试
// @Accept       json
// @Produce      json
// @Param        id   path      uint                        true  "模拟考试ID"
// @Param        req  body      dto.SubmitMockExamRequest    true  "考试答案"
// @Success      200  {object}  response.Response{data=dto.MockExamResultInfo}  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /practice/mock/{id}/submit [post]
// @Security     BearerAuth
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

// StartChallenge 开始闯关
// @Summary      开始闯关
// @Description  根据关卡ID开始闯关挑战
// @Tags         闯关模式
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "关卡ID"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "无效的关卡ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /practice/challenge/levels/{id}/start [post]
// @Security     BearerAuth
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
// @Summary      提交闯关答案
// @Description  根据关卡ID提交闯关挑战答案
// @Tags         闯关模式
// @Accept       json
// @Produce      json
// @Param        id   path      uint                          true  "关卡ID"
// @Param        req  body      dto.SubmitChallengeRequest    true  "闯关答案"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /practice/challenge/levels/{id}/submit [post]
// @Security     BearerAuth
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
