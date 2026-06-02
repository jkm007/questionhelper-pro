package statistics

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/service/statistics"
	"questionhelper-server/pkg/response"
)

// GetExamStatistics 考试统计
// @Summary      获取考试统计
// @Description  获取考试统计数据概览
// @Tags         数据统计
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /statistics/exam [get]
// @Security     BearerAuth
func (ctrl *StatisticsController) GetExamStatistics(c *gin.Context) {
	stats, err := statistics.GetExamStatistics()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, stats)
}

// GetScorePrediction 成绩预测
// @Summary      获取成绩预测
// @Description  根据查询条件获取成绩预测数据
// @Tags         数据统计
// @Accept       json
// @Produce      json
// @Param        req  query     dto.ScorePredictionRequest  false  "查询参数"
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /statistics/score/prediction [get]
// @Security     BearerAuth
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
// @Summary      获取成绩预警
// @Description  根据查询条件获取成绩预警数据
// @Tags         数据统计
// @Accept       json
// @Produce      json
// @Param        req  query     dto.ScoreAlertRequest  false  "查询参数"
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /statistics/score/alert [get]
// @Security     BearerAuth
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
