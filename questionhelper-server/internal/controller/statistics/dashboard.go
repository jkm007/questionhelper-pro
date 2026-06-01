package statistics

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/service/statistics"
	"questionhelper-server/pkg/response"
)

// GetDashboard 仪表盘（管理员）
func (ctrl *StatisticsController) GetDashboard(c *gin.Context) {
	dashboard, err := statistics.GetDashboard()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, dashboard)
}

// ==================== 移动端统计 ====================

// GetMobileOverview 移动端个人概览
func (ctrl *StatisticsController) GetMobileOverview(c *gin.Context) {
	userID := c.GetUint("user_id")

	overview, err := statistics.GetMobileOverview(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, overview)
}
