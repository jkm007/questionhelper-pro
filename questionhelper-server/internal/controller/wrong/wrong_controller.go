package wrong

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/service/wrong"
	"questionhelper-server/pkg/response"
)

type WrongController struct{}

func NewWrongController() *WrongController {
	return &WrongController{}
}

// ListWrongQuestions 错题列表
func (ctrl *WrongController) ListWrongQuestions(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.WrongListRequest
	c.ShouldBindQuery(&req)

	list, total, err := wrong.ListWrongQuestions(userID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// GetWrongQuestion 获取错题详情
func (ctrl *WrongController) GetWrongQuestion(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的错题ID")
		return
	}

	info, err := wrong.GetWrongQuestion(uint(id), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, info)
}

// ReviewWrongQuestion 复习错题
func (ctrl *WrongController) ReviewWrongQuestion(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的错题ID")
		return
	}

	var req dto.ReviewWrongRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	isCorrect, err := wrong.ReviewWrongQuestion(uint(id), userID, req.Answer)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	if isCorrect {
		response.SuccessWithMessage(c, "回答正确！", gin.H{"is_correct": true})
	} else {
		response.SuccessWithMessage(c, "回答错误", gin.H{"is_correct": false})
	}
}

// RemoveWrongQuestion 移除错题
func (ctrl *WrongController) RemoveWrongQuestion(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的错题ID")
		return
	}

	if err := wrong.RemoveWrongQuestion(uint(id), userID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "移除成功", nil)
}

// GetWrongAnalysis 错题分析
func (ctrl *WrongController) GetWrongAnalysis(c *gin.Context) {
	userID := c.GetUint("user_id")

	analysis, err := wrong.GetWrongAnalysis(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, analysis)
}
