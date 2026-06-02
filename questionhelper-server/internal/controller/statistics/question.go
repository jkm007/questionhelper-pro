package statistics

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/service/statistics"
	"questionhelper-server/pkg/response"
)

// GetQuestionDifficulty 题目难度分析
// @Summary      获取题目难度分析
// @Description  根据查询条件获取题目难度分析数据
// @Tags         数据统计
// @Accept       json
// @Produce      json
// @Param        req  query     dto.QuestionDifficultyRequest  false  "查询参数"
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /statistics/question/difficulty [get]
// @Security     BearerAuth
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
// @Summary      获取题目区分度分析
// @Description  根据查询条件获取题目区分度分析数据
// @Tags         数据统计
// @Accept       json
// @Produce      json
// @Param        req  query     dto.QuestionDiscriminationRequest  false  "查询参数"
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /statistics/question/discrimination [get]
// @Security     BearerAuth
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
