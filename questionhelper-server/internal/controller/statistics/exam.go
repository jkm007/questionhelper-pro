package statistics

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/service/statistics"
	"questionhelper-server/pkg/response"
)

// GetExamStatistics 考试统计
func (ctrl *StatisticsController) GetExamStatistics(c *gin.Context) {
	stats, err := statistics.GetExamStatistics()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, stats)
}

// GetScorePrediction 成绩预测
func (ctrl *StatisticsController) GetScorePrediction(c *gin.Context) {
	var req dto.ScorePredictionRequest
	c.ShouldBindQuery(&req)

	userID := c.GetUint("user_id")
	result, err := statistics.GetScorePrediction(userID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, result)
}

// GetScoreAlert 成绩预警
func (ctrl *StatisticsController) GetScoreAlert(c *gin.Context) {
	var req dto.ScoreAlertRequest
	c.ShouldBindQuery(&req)

	userID := c.GetUint("user_id")
	items, err := statistics.GetScoreAlert(userID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, items)
}
