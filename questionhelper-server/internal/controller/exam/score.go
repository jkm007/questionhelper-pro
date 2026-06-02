package exam

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/service/exam"
	"questionhelper-server/pkg/response"
)

// ListScores 成绩列表
// @Summary      获取成绩列表
// @Description  获取成绩列表，支持按考试筛选和分页（管理员）
// @Tags         成绩管理
// @Accept       json
// @Produce      json
// @Param        page       query     int  false  "页码"
// @Param        page_size  query     int  false  "每页数量"
// @Param        exam_id    query     uint false  "考试ID"
// @Success      200  {object}  response.PageResponse{data=[]dto.ExamRecordInfo}  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/scores [get]
// @Security     BearerAuth
func (ctrl *ExamController) ListScores(c *gin.Context) {
	var req dto.PageRequest
	c.ShouldBindQuery(&req)

	var examID *uint
	if eid := c.Query("exam_id"); eid != "" {
		id, err := strconv.ParseUint(eid, 10, 32)
		if err == nil {
			uid := uint(id)
			examID = &uid
		}
	}

	list, total, err := exam.ListScores(examID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// GetScore 获取成绩详情
// @Summary      获取成绩详情
// @Description  根据记录ID获取成绩详细信息（管理员）
// @Tags         成绩管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "成绩记录ID"
// @Success      200  {object}  response.Response{data=dto.ExamRecordInfo}  "成功"
// @Failure      400  {object}  response.Response  "无效的记录ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/scores/{id} [get]
// @Security     BearerAuth
func (ctrl *ExamController) GetScore(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的记录ID")
		return
	}

	result, err := exam.GetScore(uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, result)
}

// GetScoreAnalysis 成绩分析
// @Summary      获取成绩分析
// @Description  获取指定考试的成绩分析数据（管理员）
// @Tags         成绩管理
// @Accept       json
// @Produce      json
// @Param        exam_id  query     uint  true  "考试ID"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "请输入考试ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/scores/analysis [get]
// @Security     BearerAuth
func (ctrl *ExamController) GetScoreAnalysis(c *gin.Context) {
	examIDStr := c.Query("exam_id")
	if examIDStr == "" {
		response.Error(c, http.StatusBadRequest, "请输入考试ID")
		return
	}

	examID, err := strconv.ParseUint(examIDStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的考试ID")
		return
	}

	result, err := exam.GetScoreAnalysis(uint(examID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, result)
}
