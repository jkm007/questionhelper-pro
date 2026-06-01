package statistics

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/service/statistics"
	"questionhelper-server/pkg/response"
)

// GetQuestionDifficulty 题目难度分析
func (ctrl *StatisticsController) GetQuestionDifficulty(c *gin.Context) {
	var req dto.QuestionDifficultyRequest
	c.ShouldBindQuery(&req)

	items, err := statistics.GetQuestionDifficulty(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, items)
}

// GetQuestionDiscrimination 题目区分度分析
func (ctrl *StatisticsController) GetQuestionDiscrimination(c *gin.Context) {
	var req dto.QuestionDiscriminationRequest
	c.ShouldBindQuery(&req)

	items, err := statistics.GetQuestionDiscrimination(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, items)
}
