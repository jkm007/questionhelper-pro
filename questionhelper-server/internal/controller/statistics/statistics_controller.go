package statistics

import (
	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/model"
	"questionhelper-server/pkg/database"
	"questionhelper-server/pkg/response"
)

type StatisticsController struct{}

func NewStatisticsController() *StatisticsController {
	return &StatisticsController{}
}

// ShareData 数据分享
// @Summary      创建数据分享
// @Description  创建一条数据分享记录
// @Tags         数据统计
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response  "成功"
// @Router       /statistics/share [post]
// @Security     BearerAuth
func (ctrl *StatisticsController) ShareData(ctx *gin.Context) {
	response.Success(ctx, gin.H{"code": "share_code_123"})
}

// GetSharedData 获取分享数据
// @Summary      获取分享数据
// @Description  根据分享码获取分享数据
// @Tags         数据统计
// @Accept       json
// @Produce      json
// @Param        code  path      string  true  "分享码"
// @Success      200  {object}  response.Response  "成功"
// @Router       /statistics/share/{code} [get]
func (ctrl *StatisticsController) GetSharedData(ctx *gin.Context) {
	response.Success(ctx, gin.H{"message": "分享数据"})
}

// RefreshStats 手动刷新统计
// @Summary      手动刷新统计数据
// @Description  触发统计数据刷新
// @Tags         数据统计
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response  "成功"
// @Router       /statistics/refresh [post]
// @Security     BearerAuth
func (ctrl *StatisticsController) RefreshStats(ctx *gin.Context) {
	response.Success(ctx, gin.H{"message": "刷新已触发"})
}

// UpdateSubscription 更新数据订阅
// @Summary      更新数据订阅
// @Description  根据ID更新数据订阅配置
// @Tags         数据统计
// @Accept       json
// @Produce      json
// @Param        id   path      uint    true  "订阅ID"
// @Param        req  body      object  true  "订阅数据"
// @Success      200  {object}  response.Response  "更新成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /statistics/subscriptions/{id} [put]
// @Security     BearerAuth
func (ctrl *StatisticsController) UpdateSubscription(ctx *gin.Context) {
	id := ctx.Param("id")
	var req map[string]interface{}
	ctx.ShouldBindJSON(&req)
	database.DB.Model(&model.DataSubscription{}).Where("id = ?", id).Updates(req)
	response.Success(ctx, nil)
}
