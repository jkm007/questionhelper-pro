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
