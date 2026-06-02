package exam

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/service/exam"
	"questionhelper-server/pkg/response"
)

type MonitorController struct{}

func NewMonitorController() *MonitorController {
	return &MonitorController{}
}

// GetExamMonitor 获取考试监控
// @Summary      获取考试监控数据
// @Description  获取指定考试的实时监控数据（管理员）
// @Tags         考试监控
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "考试ID"
// @Success      200  {object}  response.Response{data=dto.ExamMonitorInfo}  "成功"
// @Failure      400  {object}  response.Response  "无效的考试ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/exams/{id}/monitor [get]
// @Security     BearerAuth
func (ctrl *MonitorController) GetExamMonitor(c *gin.Context) {
	examID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的考试ID")
		return
	}

	result, err := exam.GetExamMonitor(uint(examID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, result)
}

// ExportScores 导出成绩
// @Summary      导出成绩
// @Description  导出指定考试的成绩数据
// @Tags         考试监控
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "考试ID"
// @Success      200  {object}  dto.ExportScoresResponse  "成功"
// @Failure      400  {object}  response.Response  "无效的考试ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/exams/{id}/scores/export [get]
// @Security     BearerAuth
func (ctrl *MonitorController) ExportScores(c *gin.Context) {
	examID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的考试ID")
		return
	}

	result, err := exam.ExportScores(uint(examID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 设置下载头
	c.Header("Content-Type", "application/json")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=scores_%d.json", examID))
	c.JSON(http.StatusOK, result)
}

// ReviewExam 阅卷
// @Summary      阅卷
// @Description  管理员对考试进行阅卷操作
// @Tags         考试监控
// @Accept       json
// @Produce      json
// @Param        id   path      uint               true  "考试ID"
// @Param        req  body      dto.ReviewRequest  true  "阅卷信息"
// @Success      200  {object}  response.Response  "阅卷成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/exams/{id}/review [post]
// @Security     BearerAuth
func (ctrl *MonitorController) ReviewExam(c *gin.Context) {
	examID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的考试ID")
		return
	}

	var req dto.ReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := exam.ReviewExam(uint(examID), &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "阅卷成功", nil)
}

// GetExamAnalysis 获取考试分析
// @Summary      获取考试分析
// @Description  获取指定考试的详细分析数据（管理员）
// @Tags         考试监控
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "考试ID"
// @Success      200  {object}  response.Response{data=dto.ExamAnalysisInfo}  "成功"
// @Failure      400  {object}  response.Response  "无效的考试ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/exams/{id}/analysis [get]
// @Security     BearerAuth
func (ctrl *MonitorController) GetExamAnalysis(c *gin.Context) {
	examID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的考试ID")
		return
	}

	result, err := exam.GetExamAnalysis(uint(examID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, result)
}
