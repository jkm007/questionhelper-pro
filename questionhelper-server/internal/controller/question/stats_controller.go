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
func (ctrl *QuestionStatsController) GetStats(c *gin.Context) {
	stats, err := question.GetQuestionStats()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, stats)
}
