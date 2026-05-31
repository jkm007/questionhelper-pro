package exam

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/service/exam"
	"questionhelper-server/pkg/response"
)

type ExamController struct{}

func NewExamController() *ExamController {
	return &ExamController{}
}

// ListExams 可用考试列表（学生）
func (ctrl *ExamController) ListExams(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.PageRequest
	c.ShouldBindQuery(&req)

	list, total, err := exam.ListAvailableExams(userID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// GetExam 获取考试详情
func (ctrl *ExamController) GetExam(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的考试ID")
		return
	}

	info, err := exam.GetExam(uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, info)
}

// StartExam 开始考试
func (ctrl *ExamController) StartExam(c *gin.Context) {
	userID := c.GetUint("user_id")
	examID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的考试ID")
		return
	}

	ip := c.ClientIP()
	record, err := exam.StartExam(uint(examID), userID, ip)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, record)
}

// SubmitExam 提交考试
func (ctrl *ExamController) SubmitExam(c *gin.Context) {
	recordID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的记录ID")
		return
	}

	var req dto.SubmitExamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := exam.SubmitExam(uint(recordID), &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "提交成功", nil)
}

// GetExamResult 获取考试结果
func (ctrl *ExamController) GetExamResult(c *gin.Context) {
	userID := c.GetUint("user_id")
	examID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的考试ID")
		return
	}

	result, err := exam.GetExamResult(uint(examID), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, result)
}

// GetExamHistory 考试历史
func (ctrl *ExamController) GetExamHistory(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.PageRequest
	c.ShouldBindQuery(&req)

	list, total, err := exam.GetExamHistory(userID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// ListPapers 试卷列表（管理员）
func (ctrl *ExamController) ListPapers(c *gin.Context) {
	var req dto.PageRequest
	c.ShouldBindQuery(&req)

	list, total, err := exam.ListPapers(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// GetPaper 获取试卷详情
func (ctrl *ExamController) GetPaper(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的试卷ID")
		return
	}

	info, err := exam.GetPaper(uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, info)
}

// CreatePaper 创建试卷
func (ctrl *ExamController) CreatePaper(c *gin.Context) {
	creatorID := c.GetUint("user_id")

	var req dto.CreatePaperRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := exam.CreatePaper(creatorID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "创建成功", nil)
}

// UpdatePaper 更新试卷
func (ctrl *ExamController) UpdatePaper(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的试卷ID")
		return
	}

	var req dto.CreatePaperRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := exam.UpdatePaper(uint(id), &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "更新成功", nil)
}

// DeletePaper 删除试卷
func (ctrl *ExamController) DeletePaper(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的试卷ID")
		return
	}

	if err := exam.DeletePaper(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "删除成功", nil)
}

// AdminListExams 考试列表（管理员）
func (ctrl *ExamController) AdminListExams(c *gin.Context) {
	var req dto.ExamListRequest
	c.ShouldBindQuery(&req)

	list, total, err := exam.ListExams(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// AdminGetExam 获取考试详情（管理员）
func (ctrl *ExamController) AdminGetExam(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的考试ID")
		return
	}

	info, err := exam.GetExam(uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, info)
}

// CreateExam 创建考试
func (ctrl *ExamController) CreateExam(c *gin.Context) {
	creatorID := c.GetUint("user_id")

	var req dto.CreateExamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := exam.CreateExam(creatorID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "创建成功", nil)
}

// UpdateExam 更新考试
func (ctrl *ExamController) UpdateExam(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的考试ID")
		return
	}

	var req dto.UpdateExamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := exam.UpdateExam(uint(id), &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "更新成功", nil)
}

// DeleteExam 删除考试
func (ctrl *ExamController) DeleteExam(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的考试ID")
		return
	}

	if err := exam.DeleteExam(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "删除成功", nil)
}

// PublishExam 发布考试
func (ctrl *ExamController) PublishExam(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的考试ID")
		return
	}

	if err := exam.PublishExam(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "发布成功", nil)
}

// CloseExam 结束考试
func (ctrl *ExamController) CloseExam(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的考试ID")
		return
	}

	if err := exam.CloseExam(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "考试已结束", nil)
}

// ListScores 成绩列表
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
