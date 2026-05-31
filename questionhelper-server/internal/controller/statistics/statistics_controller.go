package statistics

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/service/statistics"
	"questionhelper-server/pkg/response"
)

type StatisticsController struct{}

func NewStatisticsController() *StatisticsController {
	return &StatisticsController{}
}

// GetUserStatistics 用户统计
func (ctrl *StatisticsController) GetUserStatistics(c *gin.Context) {
	stats, err := statistics.GetUserStatistics()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, stats)
}

// GetPracticeStatistics 练习统计
func (ctrl *StatisticsController) GetPracticeStatistics(c *gin.Context) {
	stats, err := statistics.GetPracticeStatistics()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, stats)
}

// GetExamStatistics 考试统计
func (ctrl *StatisticsController) GetExamStatistics(c *gin.Context) {
	stats, err := statistics.GetExamStatistics()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, stats)
}

// GetRanking 排行榜
func (ctrl *StatisticsController) GetRanking(c *gin.Context) {
	var req dto.RankingRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	ranks, total, err := statistics.GetRanking(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, ranks, total, req.Page, req.PageSize)
}

// GetClassStatistics 班级统计
func (ctrl *StatisticsController) GetClassStatistics(c *gin.Context) {
	stats, err := statistics.GetClassStatistics()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, stats)
}

// GetDashboard 仪表盘（管理员）
func (ctrl *StatisticsController) GetDashboard(c *gin.Context) {
	dashboard, err := statistics.GetDashboard()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, dashboard)
}
