package class

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/service/class"
	"questionhelper-server/pkg/response"
)

// GetHomework 获取作业详情
func (ctrl *ClassController) GetHomework(c *gin.Context) {
	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	homeworkID, err := strconv.ParseUint(c.Param("homeworkId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的作业ID")
		return
	}

	info, err := class.GetHomework(uint(classID), uint(homeworkID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, info)
}

// CreateHomework 创建作业
func (ctrl *ClassController) CreateHomework(c *gin.Context) {
	userID := c.GetUint("user_id")

	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	var req dto.CreateHomeworkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := class.CreateHomework(uint(classID), userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "创建成功", nil)
}

// UpdateHomework 更新作业
func (ctrl *ClassController) UpdateHomework(c *gin.Context) {
	userID := c.GetUint("user_id")

	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	homeworkID, err := strconv.ParseUint(c.Param("homeworkId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的作业ID")
		return
	}

	var req dto.UpdateHomeworkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := class.UpdateHomework(uint(classID), uint(homeworkID), userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "更新成功", nil)
}

// DeleteHomework 删除作业
func (ctrl *ClassController) DeleteHomework(c *gin.Context) {
	userID := c.GetUint("user_id")

	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	homeworkID, err := strconv.ParseUint(c.Param("homeworkId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的作业ID")
		return
	}

	if err := class.DeleteHomework(uint(classID), uint(homeworkID), userID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "删除成功", nil)
}

// SubmitHomework 提交作业
func (ctrl *ClassController) SubmitHomework(c *gin.Context) {
	userID := c.GetUint("user_id")

	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	homeworkID, err := strconv.ParseUint(c.Param("homeworkId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的作业ID")
		return
	}

	var req dto.SubmitHomeworkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := class.SubmitHomework(uint(classID), uint(homeworkID), userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "提交成功", nil)
}

// GradeHomework 批改作业
func (ctrl *ClassController) GradeHomework(c *gin.Context) {
	userID := c.GetUint("user_id")

	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	homeworkID, err := strconv.ParseUint(c.Param("homeworkId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的作业ID")
		return
	}

	submissionID, err := strconv.ParseUint(c.Param("submissionId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的提交ID")
		return
	}

	var req dto.GradeHomeworkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := class.GradeHomework(uint(classID), uint(homeworkID), uint(submissionID), userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "批改成功", nil)
}

// ListHomeworkSubmissions 作业提交列表
func (ctrl *ClassController) ListHomeworkSubmissions(c *gin.Context) {
	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	homeworkID, err := strconv.ParseUint(c.Param("homeworkId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的作业ID")
		return
	}

	var req dto.PageRequest
	c.ShouldBindQuery(&req)

	list, total, err := class.ListHomeworkSubmissions(uint(classID), uint(homeworkID), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}
