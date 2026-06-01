package exam

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/service/exam"
	"questionhelper-server/pkg/response"
)

// ListScores 成绩列表
func (ctrl *ExamController) ListScores(c *gin.Context) {
	var req dto.PageRequest
	c.ShouldBindQuery(&req)

	var examID *uint
	if eid := c.Query("exam_id"); eid != "" {
		id, err := strconv.ParseUint(eid, 10, 32)
		if err == nil {
			uid := uint(id)
			examID = &uid
		}
	}

	list, total, err := exam.ListScores(examID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// GetScore 获取成绩详情
func (ctrl *ExamController) GetScore(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的记录ID")
		return
	}

	result, err := exam.GetScore(uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, result)
}

// GetScoreAnalysis 成绩分析
func (ctrl *ExamController) GetScoreAnalysis(c *gin.Context) {
	examIDStr := c.Query("exam_id")
	if examIDStr == "" {
		response.Error(c, http.StatusBadRequest, "请输入考试ID")
		return
	}

	examID, err := strconv.ParseUint(examIDStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的考试ID")
		return
	}

	result, err := exam.GetScoreAnalysis(uint(examID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, result)
}
