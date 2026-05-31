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
