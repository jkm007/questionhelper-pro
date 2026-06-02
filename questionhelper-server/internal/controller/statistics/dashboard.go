package statistics

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/service/statistics"
	"questionhelper-server/pkg/response"
)

// GetDashboard 仪表盘（管理员）
// @Summary      获取仪表盘数据
// @Description  获取管理员仪表盘概览数据
// @Tags         数据统计
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /statistics/overview [get]
// @Security     BearerAuth
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
// @Summary      获取移动端个人概览
// @Description  获取移动端当前用户的个人学习概览数据
// @Tags         数据统计
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /statistics/mobile/overview [get]
// @Security     BearerAuth
func (ctrl *StatisticsController) GetMobileOverview(c *gin.Context) {
	userID := c.GetUint("user_id")

	overview, err := statistics.GetMobileOverview(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, overview)
}
