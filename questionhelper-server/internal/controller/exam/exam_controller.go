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
// @Summary      获取可用考试列表
// @Description  获取当前学生可参加的考试列表，支持分页
// @Tags         考试管理
// @Accept       json
// @Produce      json
// @Param        page       query     int  false  "页码"
// @Param        page_size  query     int  false  "每页数量"
// @Success      200  {object}  response.PageResponse{data=[]dto.ExamInfo}  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /exam [get]
// @Security     BearerAuth
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
// @Summary      获取考试详情
// @Description  根据考试ID获取考试详细信息
// @Tags         考试管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "考试ID"
// @Success      200  {object}  response.Response{data=dto.ExamInfo}  "成功"
// @Failure      400  {object}  response.Response  "无效的考试ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /exam/{id} [get]
// @Security     BearerAuth
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
// @Summary      开始考试
// @Description  学生开始指定考试，创建考试记录
// @Tags         考试管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "考试ID"
// @Success      200  {object}  response.Response{data=dto.ExamRecordInfo}  "成功"
// @Failure      400  {object}  response.Response  "无效的考试ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /exam/{id}/start [post]
// @Security     BearerAuth
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
// @Summary      提交考试
// @Description  学生提交考试答卷
// @Tags         考试管理
// @Accept       json
// @Produce      json
// @Param        recordId  path      uint                    true  "考试记录ID"
// @Param        req       body      dto.SubmitExamRequest    true  "提交信息"
// @Success      200  {object}  response.Response  "提交成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /exam/{recordId}/submit [post]
// @Security     BearerAuth
func (ctrl *ExamController) SubmitExam(c *gin.Context) {
	recordID, err := strconv.ParseUint(c.Param("recordId"), 10, 32)
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
// @Summary      获取考试结果
// @Description  获取指定考试的成绩结果
// @Tags         考试管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "考试ID"
// @Success      200  {object}  response.Response{data=dto.ExamResultInfo}  "成功"
// @Failure      400  {object}  response.Response  "无效的考试ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /exam/{id}/result [get]
// @Security     BearerAuth
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
// @Summary      获取考试历史
// @Description  获取当前学生的考试历史记录，支持分页
// @Tags         考试管理
// @Accept       json
// @Produce      json
// @Param        page       query     int  false  "页码"
// @Param        page_size  query     int  false  "每页数量"
// @Success      200  {object}  response.PageResponse{data=[]dto.ExamRecordInfo}  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /exam/history [get]
// @Security     BearerAuth
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
// @Summary      获取试卷列表
// @Description  获取试卷列表，支持分页（管理员）
// @Tags         试卷管理
// @Accept       json
// @Produce      json
// @Param        page       query     int  false  "页码"
// @Param        page_size  query     int  false  "每页数量"
// @Success      200  {object}  response.PageResponse{data=[]dto.PaperInfo}  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/papers [get]
// @Security     BearerAuth
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
// @Summary      获取试卷详情
// @Description  根据试卷ID获取试卷详细信息（管理员）
// @Tags         试卷管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "试卷ID"
// @Success      200  {object}  response.Response{data=dto.PaperInfo}  "成功"
// @Failure      400  {object}  response.Response  "无效的试卷ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/papers/{id} [get]
// @Security     BearerAuth
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
// @Summary      创建试卷
// @Description  管理员创建新试卷
// @Tags         试卷管理
// @Accept       json
// @Produce      json
// @Param        req  body      dto.CreatePaperRequest  true  "试卷信息"
// @Success      200  {object}  response.Response  "创建成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/papers [post]
// @Security     BearerAuth
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
// @Summary      更新试卷
// @Description  管理员更新指定试卷信息
// @Tags         试卷管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint                    true  "试卷ID"
// @Param        req  body      dto.CreatePaperRequest  true  "试卷信息"
// @Success      200  {object}  response.Response  "更新成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/papers/{id} [put]
// @Security     BearerAuth
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
// @Summary      删除试卷
// @Description  管理员删除指定试卷
// @Tags         试卷管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "试卷ID"
// @Success      200  {object}  response.Response  "删除成功"
// @Failure      400  {object}  response.Response  "无效的试卷ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/papers/{id} [delete]
// @Security     BearerAuth
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
// @Summary      获取考试列表
// @Description  管理员获取考试列表，支持分页和筛选
// @Tags         考试管理
// @Accept       json
// @Produce      json
// @Param        page       query     int     false  "页码"
// @Param        page_size  query     int     false  "每页数量"
// @Param        status     query     string  false  "考试状态"
// @Param        keyword    query     string  false  "搜索关键词"
// @Success      200  {object}  response.PageResponse{data=[]dto.ExamInfo}  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/exams [get]
// @Security     BearerAuth
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
// @Summary      获取考试详情
// @Description  管理员根据考试ID获取考试详细信息
// @Tags         考试管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "考试ID"
// @Success      200  {object}  response.Response{data=dto.ExamInfo}  "成功"
// @Failure      400  {object}  response.Response  "无效的考试ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/exams/{id} [get]
// @Security     BearerAuth
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
// @Summary      创建考试
// @Description  管理员创建新考试
// @Tags         考试管理
// @Accept       json
// @Produce      json
// @Param        req  body      dto.CreateExamRequest  true  "考试信息"
// @Success      200  {object}  response.Response  "创建成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/exams [post]
// @Security     BearerAuth
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
// @Summary      更新考试
// @Description  管理员更新指定考试信息
// @Tags         考试管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint                    true  "考试ID"
// @Param        req  body      dto.UpdateExamRequest   true  "考试信息"
// @Success      200  {object}  response.Response  "更新成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/exams/{id} [put]
// @Security     BearerAuth
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
// @Summary      删除考试
// @Description  管理员删除指定考试
// @Tags         考试管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "考试ID"
// @Success      200  {object}  response.Response  "删除成功"
// @Failure      400  {object}  response.Response  "无效的考试ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/exams/{id} [delete]
// @Security     BearerAuth
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
// @Summary      发布考试
// @Description  管理员发布指定考试，使其对学生可见
// @Tags         考试管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "考试ID"
// @Success      200  {object}  response.Response  "发布成功"
// @Failure      400  {object}  response.Response  "无效的考试ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/exams/{id}/publish [put]
// @Security     BearerAuth
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
// @Summary      结束考试
// @Description  管理员结束指定考试
// @Tags         考试管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "考试ID"
// @Success      200  {object}  response.Response  "考试已结束"
// @Failure      400  {object}  response.Response  "无效的考试ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/exams/{id}/close [put]
// @Security     BearerAuth
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

