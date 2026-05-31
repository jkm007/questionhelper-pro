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

	if err := practice.SubmitPractice(uint(sessionID), &req); err != nil {
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
