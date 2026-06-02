package question

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/service/question"
	"questionhelper-server/pkg/response"
)

type QuestionStatsController struct{}

func NewQuestionStatsController() *QuestionStatsController {
	return &QuestionStatsController{}
}

// GetStats 获取题目统计
// @Summary      获取题目统计
// @Description  管理员获取题目统计数据，包括总量、题型分布、分类分布等
// @Tags         题目管理
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response{data=dto.QuestionStatsResponse}  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/questions/stats [get]
// @Security     BearerAuth
func (ctrl *QuestionStatsController) GetStats(c *gin.Context) {
	stats, err := question.GetQuestionStats()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, stats)
}
