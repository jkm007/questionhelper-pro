package statistics

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/service/statistics"
	"questionhelper-server/pkg/response"
)

// GetClassStatistics 班级统计
// @Summary      获取班级统计
// @Description  获取班级统计数据概览
// @Tags         数据统计
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /statistics/classes [get]
// @Security     BearerAuth
func (ctrl *StatisticsController) GetClassStatistics(c *gin.Context) {
	stats, err := statistics.GetClassStatistics()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, stats)
}

// GetClassOverview 班级概览
// @Summary      获取班级概览
// @Description  根据班级ID获取班级概览数据
// @Tags         数据统计
// @Accept       json
// @Produce      json
// @Param        id  path      uint  true  "班级ID"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "无效的班级ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /statistics/class/{id}/overview [get]
// @Security     BearerAuth
func (ctrl *StatisticsController) GetClassOverview(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	overview, err := statistics.GetClassOverview(uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, overview)
}

// GetClassStudents 班级学生成绩列表
// @Summary      获取班级学生成绩列表
// @Description  根据班级ID获取学生成绩列表，支持分页
// @Tags         数据统计
// @Accept       json
// @Produce      json
// @Param        id         path      uint  true   "班级ID"
// @Param        page       query     int   false  "页码"
// @Param        page_size  query     int   false  "每页数量"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "无效的班级ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /statistics/class/{id}/students [get]
// @Security     BearerAuth
func (ctrl *StatisticsController) GetClassStudents(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	var req dto.ClassStudentListRequest
	c.ShouldBindQuery(&req)

	list, total, err := statistics.GetClassStudents(uint(id), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// GetClassPracticeStats 班级练习统计
// @Summary      获取班级练习统计
// @Description  根据班级ID获取班级练习统计数据
// @Tags         数据统计
// @Accept       json
// @Produce      json
// @Param        id   path      uint                       true  "班级ID"
// @Param        req  query     dto.ClassPracticeRequest  false  "查询参数"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "无效的班级ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /statistics/class/{id}/practice [get]
// @Security     BearerAuth
func (ctrl *StatisticsController) GetClassPracticeStats(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	var req dto.ClassPracticeRequest
	c.ShouldBindQuery(&req)

	result, err := statistics.GetClassPracticeStats(uint(id), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, result)
}

// GetClassExamStats 班级考试统计
// @Summary      获取班级考试统计
// @Description  根据班级ID获取班级考试统计数据
// @Tags         数据统计
// @Accept       json
// @Produce      json
// @Param        id   path      uint                   true  "班级ID"
// @Param        req  query     dto.ClassExamRequest  false  "查询参数"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "无效的班级ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /statistics/class/{id}/exam [get]
// @Security     BearerAuth
func (ctrl *StatisticsController) GetClassExamStats(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	var req dto.ClassExamRequest
	c.ShouldBindQuery(&req)

	result, err := statistics.GetClassExamStats(uint(id), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, result)
}

// GetClassQuestionStats 班级题目统计
// @Summary      获取班级题目统计
// @Description  根据班级ID获取班级题目统计数据
// @Tags         数据统计
// @Accept       json
// @Produce      json
// @Param        id  path      uint  true  "班级ID"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "无效的班级ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /statistics/class/{id}/questions [get]
// @Security     BearerAuth
func (ctrl *StatisticsController) GetClassQuestionStats(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	result, err := statistics.GetClassQuestionStats(uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, result)
}
