package class

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/service/class"
	"questionhelper-server/pkg/response"
)

// ==================== Peer Review ====================

// AssignPeerReview 分配互评
func (ctrl *ClassController) AssignPeerReview(c *gin.Context) {
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

	var req dto.AssignPeerReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := class.AssignPeerReview(uint(classID), uint(homeworkID), userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "分配成功", nil)
}

// ListPeerReviews 互评列表
func (ctrl *ClassController) ListPeerReviews(c *gin.Context) {
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

	list, total, err := class.ListPeerReviews(uint(classID), uint(homeworkID), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// GetMyPeerReviews 我的互评任务
func (ctrl *ClassController) GetMyPeerReviews(c *gin.Context) {
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

	list, err := class.GetMyPeerReviews(uint(classID), uint(homeworkID), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, list)
}

// SubmitPeerReview 提交互评
func (ctrl *ClassController) SubmitPeerReview(c *gin.Context) {
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

	reviewID, err := strconv.ParseUint(c.Param("reviewId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的互评ID")
		return
	}

	var req dto.SubmitPeerReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := class.SubmitPeerReview(uint(classID), uint(homeworkID), uint(reviewID), userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "提交成功", nil)
}

// GetPeerReviewResult 互评结果
func (ctrl *ClassController) GetPeerReviewResult(c *gin.Context) {
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

	result, err := class.GetPeerReviewResult(uint(classID), uint(homeworkID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, result)
}

// ==================== Group Management ====================

// ListGroups 分组列表
func (ctrl *ClassController) ListGroups(c *gin.Context) {
	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	list, err := class.ListGroups(uint(classID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, list)
}

// CreateGroup 创建分组
func (ctrl *ClassController) CreateGroup(c *gin.Context) {
	userID := c.GetUint("user_id")

	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	var req dto.CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := class.CreateGroup(uint(classID), userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "创建成功", nil)
}

// UpdateGroup 更新分组
func (ctrl *ClassController) UpdateGroup(c *gin.Context) {
	userID := c.GetUint("user_id")

	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	groupID, err := strconv.ParseUint(c.Param("groupId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的分组ID")
		return
	}

	var req dto.UpdateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := class.UpdateGroup(uint(classID), uint(groupID), userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "更新成功", nil)
}

// DeleteGroup 删除分组
func (ctrl *ClassController) DeleteGroup(c *gin.Context) {
	userID := c.GetUint("user_id")

	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	groupID, err := strconv.ParseUint(c.Param("groupId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的分组ID")
		return
	}

	if err := class.DeleteGroup(uint(classID), uint(groupID), userID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "删除成功", nil)
}

// AddGroupMember 添加分组成员
func (ctrl *ClassController) AddGroupMember(c *gin.Context) {
	userID := c.GetUint("user_id")

	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	groupID, err := strconv.ParseUint(c.Param("groupId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的分组ID")
		return
	}

	var req dto.AddGroupMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := class.AddGroupMember(uint(classID), uint(groupID), userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "添加成功", nil)
}

// RemoveGroupMember 移除分组成员
func (ctrl *ClassController) RemoveGroupMember(c *gin.Context) {
	userID := c.GetUint("user_id")

	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	groupID, err := strconv.ParseUint(c.Param("groupId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的分组ID")
		return
	}

	memberUserID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	if err := class.RemoveGroupMember(uint(classID), uint(groupID), uint(memberUserID), userID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "移除成功", nil)
}

// ==================== Application ====================

// ApplyClass 提交申请
func (ctrl *ClassController) ApplyClass(c *gin.Context) {
	userID := c.GetUint("user_id")

	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	var req dto.ClassApplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := class.ApplyClass(uint(classID), userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "申请已提交", nil)
}

// ListApplications 申请列表
func (ctrl *ClassController) ListApplications(c *gin.Context) {
	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	var req dto.PageRequest
	c.ShouldBindQuery(&req)

	list, total, err := class.ListApplications(uint(classID), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// ApproveApplication 审批通过
func (ctrl *ClassController) ApproveApplication(c *gin.Context) {
	userID := c.GetUint("user_id")

	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	appID, err := strconv.ParseUint(c.Param("appId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的申请ID")
		return
	}

	var req dto.ClassReviewApplicationRequest
	c.ShouldBindJSON(&req)

	if err := class.ApproveApplication(uint(classID), uint(appID), userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "审批通过", nil)
}

// RejectApplication 审批驳回
func (ctrl *ClassController) RejectApplication(c *gin.Context) {
	userID := c.GetUint("user_id")

	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	appID, err := strconv.ParseUint(c.Param("appId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的申请ID")
		return
	}

	var req dto.ClassReviewApplicationRequest
	c.ShouldBindJSON(&req)

	if err := class.RejectApplication(uint(classID), uint(appID), userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "已驳回", nil)
}

// ==================== Attendance ====================

// ListAttendances 考勤列表
func (ctrl *ClassController) ListAttendances(c *gin.Context) {
	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	var req dto.PageRequest
	c.ShouldBindQuery(&req)

	list, total, err := class.ListAttendances(uint(classID), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// CreateAttendance 创建考勤
func (ctrl *ClassController) CreateAttendance(c *gin.Context) {
	userID := c.GetUint("user_id")

	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	var req dto.CreateAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := class.CreateAttendance(uint(classID), userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "创建成功", nil)
}

// UpdateAttendance 编辑考勤
func (ctrl *ClassController) UpdateAttendance(c *gin.Context) {
	userID := c.GetUint("user_id")

	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	attID, err := strconv.ParseUint(c.Param("attId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的考勤ID")
		return
	}

	var req dto.UpdateAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := class.UpdateAttendance(uint(classID), uint(attID), userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "更新成功", nil)
}

// DeleteAttendance 删除考勤
func (ctrl *ClassController) DeleteAttendance(c *gin.Context) {
	userID := c.GetUint("user_id")

	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	attID, err := strconv.ParseUint(c.Param("attId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的考勤ID")
		return
	}

	if err := class.DeleteAttendance(uint(classID), uint(attID), userID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "删除成功", nil)
}

// Checkin 学生签到
func (ctrl *ClassController) Checkin(c *gin.Context) {
	userID := c.GetUint("user_id")

	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	attID, err := strconv.ParseUint(c.Param("attId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的考勤ID")
		return
	}

	var req dto.AttendanceCheckinRequest
	c.ShouldBindJSON(&req)

	// 获取客户端IP
	if req.IP == "" {
		req.IP = c.ClientIP()
	}

	if err := class.Checkin(uint(classID), uint(attID), userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "签到成功", nil)
}

// Checkout 学生签退
func (ctrl *ClassController) Checkout(c *gin.Context) {
	userID := c.GetUint("user_id")

	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	attID, err := strconv.ParseUint(c.Param("attId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的考勤ID")
		return
	}

	var req dto.AttendanceCheckinRequest
	c.ShouldBindJSON(&req)

	if req.IP == "" {
		req.IP = c.ClientIP()
	}

	if err := class.Checkout(uint(classID), uint(attID), userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "签退成功", nil)
}

// ListAttendanceRecords 考勤记录
func (ctrl *ClassController) ListAttendanceRecords(c *gin.Context) {
	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	attID, err := strconv.ParseUint(c.Param("attId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的考勤ID")
		return
	}

	var req dto.PageRequest
	c.ShouldBindQuery(&req)

	list, total, err := class.ListAttendanceRecords(uint(classID), uint(attID), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// ExportAttendance 导出考勤
func (ctrl *ClassController) ExportAttendance(c *gin.Context) {
	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	attID, err := strconv.ParseUint(c.Param("attId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的考勤ID")
		return
	}

	list, err := class.ExportAttendance(uint(classID), uint(attID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, list)
}

// ==================== Study Plan ====================

// ListStudyPlans 计划列表
func (ctrl *ClassController) ListStudyPlans(c *gin.Context) {
	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	var req dto.PageRequest
	c.ShouldBindQuery(&req)

	list, total, err := class.ListStudyPlans(uint(classID), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// CreateStudyPlan 创建计划
func (ctrl *ClassController) CreateStudyPlan(c *gin.Context) {
	userID := c.GetUint("user_id")

	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	var req dto.CreateStudyPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := class.CreateStudyPlan(uint(classID), userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "创建成功", nil)
}

// GetStudyPlan 计划详情
func (ctrl *ClassController) GetStudyPlan(c *gin.Context) {
	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	planID, err := strconv.ParseUint(c.Param("planId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的计划ID")
		return
	}

	info, err := class.GetStudyPlan(uint(classID), uint(planID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, info)
}

// UpdateStudyPlan 更新计划
func (ctrl *ClassController) UpdateStudyPlan(c *gin.Context) {
	userID := c.GetUint("user_id")

	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	planID, err := strconv.ParseUint(c.Param("planId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的计划ID")
		return
	}

	var req dto.UpdateStudyPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := class.UpdateStudyPlan(uint(classID), uint(planID), userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "更新成功", nil)
}

// DeleteStudyPlan 删除计划
func (ctrl *ClassController) DeleteStudyPlan(c *gin.Context) {
	userID := c.GetUint("user_id")

	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	planID, err := strconv.ParseUint(c.Param("planId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的计划ID")
		return
	}

	if err := class.DeleteStudyPlan(uint(classID), uint(planID), userID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "删除成功", nil)
}

// AddStudyPlanItem 添加任务
func (ctrl *ClassController) AddStudyPlanItem(c *gin.Context) {
	userID := c.GetUint("user_id")

	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	planID, err := strconv.ParseUint(c.Param("planId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的计划ID")
		return
	}

	var req dto.CreateStudyPlanItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := class.AddStudyPlanItem(uint(classID), uint(planID), userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "添加成功", nil)
}

// UpdateStudyPlanItem 更新任务
func (ctrl *ClassController) UpdateStudyPlanItem(c *gin.Context) {
	userID := c.GetUint("user_id")

	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	planID, err := strconv.ParseUint(c.Param("planId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的计划ID")
		return
	}

	itemID, err := strconv.ParseUint(c.Param("itemId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的任务ID")
		return
	}

	var req dto.UpdateStudyPlanItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := class.UpdateStudyPlanItem(uint(classID), uint(planID), uint(itemID), userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "更新成功", nil)
}

// DeleteStudyPlanItem 删除任务
func (ctrl *ClassController) DeleteStudyPlanItem(c *gin.Context) {
	userID := c.GetUint("user_id")

	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	planID, err := strconv.ParseUint(c.Param("planId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的计划ID")
		return
	}

	itemID, err := strconv.ParseUint(c.Param("itemId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的任务ID")
		return
	}

	if err := class.DeleteStudyPlanItem(uint(classID), uint(planID), uint(itemID), userID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "删除成功", nil)
}

// CompleteStudyPlanItem 完成任务
func (ctrl *ClassController) CompleteStudyPlanItem(c *gin.Context) {
	userID := c.GetUint("user_id")

	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	planID, err := strconv.ParseUint(c.Param("planId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的计划ID")
		return
	}

	itemID, err := strconv.ParseUint(c.Param("itemId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的任务ID")
		return
	}

	if err := class.CompleteStudyPlanItem(uint(classID), uint(planID), uint(itemID), userID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "已完成", nil)
}

// GetStudyPlanProgress 进度查看
func (ctrl *ClassController) GetStudyPlanProgress(c *gin.Context) {
	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	planID, err := strconv.ParseUint(c.Param("planId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的计划ID")
		return
	}

	info, err := class.GetStudyPlanProgress(uint(classID), uint(planID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, info)
}

// ==================== Class File ====================

// ListClassFiles 文件列表
func (ctrl *ClassController) ListClassFiles(c *gin.Context) {
	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	var req dto.PageRequest
	c.ShouldBindQuery(&req)

	list, total, err := class.ListClassFiles(uint(classID), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// UploadClassFile 上传文件
func (ctrl *ClassController) UploadClassFile(c *gin.Context) {
	userID := c.GetUint("user_id")

	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "请选择文件")
		return
	}

	// 保存文件到本地（简化实现，实际应使用 OSS）
	savePath := "uploads/class/" + strconv.FormatUint(classID, 10) + "/" + file.Filename
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		response.Error(c, http.StatusInternalServerError, "文件保存失败")
		return
	}

	if err := class.UploadClassFile(uint(classID), userID, file.Filename, savePath, file.Header.Get("Content-Type"), file.Size); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "上传成功", nil)
}

// DeleteClassFile 删除文件
func (ctrl *ClassController) DeleteClassFile(c *gin.Context) {
	userID := c.GetUint("user_id")

	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	fileID, err := strconv.ParseUint(c.Param("fileId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的文件ID")
		return
	}

	if err := class.DeleteClassFile(uint(classID), uint(fileID), userID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "删除成功", nil)
}

// DownloadClassFile 下载文件
func (ctrl *ClassController) DownloadClassFile(c *gin.Context) {
	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	fileID, err := strconv.ParseUint(c.Param("fileId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的文件ID")
		return
	}

	info, err := class.GetClassFile(uint(classID), uint(fileID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.File(info.Path)
}

// ==================== Ranking ====================

// ListRanking 排名列表
func (ctrl *ClassController) ListRanking(c *gin.Context) {
	userID := c.GetUint("user_id")

	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	rankingType := c.DefaultQuery("type", "overall")

	list, err := class.ListRanking(uint(classID), userID, rankingType)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, list)
}

// CalculateRanking 触发排名计算
func (ctrl *ClassController) CalculateRanking(c *gin.Context) {
	userID := c.GetUint("user_id")

	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	var req dto.CalculateRankingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := class.CalculateRanking(uint(classID), userID, req.RankingType); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "排名计算完成", nil)
}

// ==================== Tag Management ====================

// ListTags 标签列表
func (ctrl *ClassController) ListTags(c *gin.Context) {
	list, err := class.ListTags()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, list)
}

// CreateTag 创建标签
func (ctrl *ClassController) CreateTag(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.ClassTagCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := class.CreateTag(userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "创建成功", nil)
}

// UpdateTag 更新标签
func (ctrl *ClassController) UpdateTag(c *gin.Context) {
	userID := c.GetUint("user_id")

	tagID, err := strconv.ParseUint(c.Param("tagId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的标签ID")
		return
	}

	var req dto.ClassTagUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := class.UpdateTag(uint(tagID), userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "更新成功", nil)
}

// DeleteTag 删除标签
func (ctrl *ClassController) DeleteTag(c *gin.Context) {
	tagID, err := strconv.ParseUint(c.Param("tagId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的标签ID")
		return
	}

	if err := class.DeleteTag(uint(tagID)); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "删除成功", nil)
}

// AddClassTag 为班级添加标签
func (ctrl *ClassController) AddClassTag(c *gin.Context) {
	userID := c.GetUint("user_id")

	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	var req dto.AddClassTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := class.AddClassTag(uint(classID), userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "添加成功", nil)
}

// RemoveClassTag 移除班级标签
func (ctrl *ClassController) RemoveClassTag(c *gin.Context) {
	userID := c.GetUint("user_id")

	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	tagID, err := strconv.ParseUint(c.Param("tagId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的标签ID")
		return
	}

	if err := class.RemoveClassTag(uint(classID), uint(tagID), userID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "移除成功", nil)
}

// ==================== Template Management ====================

// ListTemplates 模板列表
func (ctrl *ClassController) ListTemplates(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.PageRequest
	c.ShouldBindQuery(&req)

	list, total, err := class.ListTemplates(userID, &req, false)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// CreateTemplate 创建模板
func (ctrl *ClassController) CreateTemplate(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.ClassTemplateCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := class.CreateTemplate(userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "创建成功", nil)
}

// GetTemplate 模板详情
func (ctrl *ClassController) GetTemplate(c *gin.Context) {
	templateID, err := strconv.ParseUint(c.Param("templateId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的模板ID")
		return
	}

	info, err := class.GetTemplate(uint(templateID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, info)
}

// UpdateTemplate 更新模板
func (ctrl *ClassController) UpdateTemplate(c *gin.Context) {
	userID := c.GetUint("user_id")

	templateID, err := strconv.ParseUint(c.Param("templateId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的模板ID")
		return
	}

	var req dto.ClassTemplateUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := class.UpdateTemplate(uint(templateID), userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "更新成功", nil)
}

// DeleteTemplate 删除模板
func (ctrl *ClassController) DeleteTemplate(c *gin.Context) {
	userID := c.GetUint("user_id")

	templateID, err := strconv.ParseUint(c.Param("templateId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的模板ID")
		return
	}

	if err := class.DeleteTemplate(uint(templateID), userID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "删除成功", nil)
}

// CreateClassFromTemplate 从模板创建班级
func (ctrl *ClassController) CreateClassFromTemplate(c *gin.Context) {
	userID := c.GetUint("user_id")

	templateID, err := strconv.ParseUint(c.Param("templateId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的模板ID")
		return
	}

	var req dto.CreateClassFromTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := class.CreateClassFromTemplate(uint(templateID), userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "创建成功", nil)
}

// ==================== Discussion Management ====================

// ListDiscussions 讨论列表
func (ctrl *ClassController) ListDiscussions(c *gin.Context) {
	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	var req dto.PageRequest
	c.ShouldBindQuery(&req)

	list, total, err := class.ListDiscussions(uint(classID), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// CreateDiscussion 发布讨论
func (ctrl *ClassController) CreateDiscussion(c *gin.Context) {
	userID := c.GetUint("user_id")

	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	var req dto.CreateDiscussionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := class.CreateDiscussion(uint(classID), userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "发布成功", nil)
}

// GetDiscussion 讨论详情
func (ctrl *ClassController) GetDiscussion(c *gin.Context) {
	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	discussionID, err := strconv.ParseUint(c.Param("discussionId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的讨论ID")
		return
	}

	info, err := class.GetDiscussion(uint(classID), uint(discussionID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, info)
}

// UpdateDiscussion 编辑讨论
func (ctrl *ClassController) UpdateDiscussion(c *gin.Context) {
	userID := c.GetUint("user_id")

	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	discussionID, err := strconv.ParseUint(c.Param("discussionId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的讨论ID")
		return
	}

	var req dto.UpdateDiscussionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := class.UpdateDiscussion(uint(classID), uint(discussionID), userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "编辑成功", nil)
}

// DeleteDiscussion 删除讨论
func (ctrl *ClassController) DeleteDiscussion(c *gin.Context) {
	userID := c.GetUint("user_id")

	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	discussionID, err := strconv.ParseUint(c.Param("discussionId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的讨论ID")
		return
	}

	if err := class.DeleteDiscussion(uint(classID), uint(discussionID), userID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "删除成功", nil)
}

// ToggleDiscussionPin 置顶/取消置顶
func (ctrl *ClassController) ToggleDiscussionPin(c *gin.Context) {
	userID := c.GetUint("user_id")

	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	discussionID, err := strconv.ParseUint(c.Param("discussionId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的讨论ID")
		return
	}

	if err := class.ToggleDiscussionPin(uint(classID), uint(discussionID), userID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "操作成功", nil)
}

// ==================== Management Enhancement ====================

// ArchiveClass 归档班级
func (ctrl *ClassController) ArchiveClass(c *gin.Context) {
	userID := c.GetUint("user_id")

	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	if err := class.ArchiveClass(uint(classID), userID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "归档成功", nil)
}

// UnarchiveClass 取消归档
func (ctrl *ClassController) UnarchiveClass(c *gin.Context) {
	userID := c.GetUint("user_id")

	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	if err := class.UnarchiveClass(uint(classID), userID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "取消归档成功", nil)
}

// PinClass 置顶班级
func (ctrl *ClassController) PinClass(c *gin.Context) {
	userID := c.GetUint("user_id")

	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	if err := class.PinClass(uint(classID), userID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "置顶成功", nil)
}

// UnpinClass 取消置顶
func (ctrl *ClassController) UnpinClass(c *gin.Context) {
	userID := c.GetUint("user_id")

	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	if err := class.UnpinClass(uint(classID), userID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "取消置顶成功", nil)
}

// SearchClasses 搜索班级
func (ctrl *ClassController) SearchClasses(c *gin.Context) {
	var req dto.SearchClassRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	list, total, err := class.SearchClasses(req.Keyword, &req.PageRequest)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// GenerateQRCode 生成二维码
func (ctrl *ClassController) GenerateQRCode(c *gin.Context) {
	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	info, err := class.GenerateQRCode(uint(classID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, info)
}

// SetClassExpire 设置有效期
func (ctrl *ClassController) SetClassExpire(c *gin.Context) {
	userID := c.GetUint("user_id")

	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	var req dto.SetExpireRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := class.SetClassExpire(uint(classID), userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "设置成功", nil)
}

// ListClassExams 班级考试列表
func (ctrl *ClassController) ListClassExams(c *gin.Context) {
	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	var req dto.PageRequest
	c.ShouldBindQuery(&req)

	list, total, err := class.ListClassExams(uint(classID), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// GetClassNotice 班级公告
func (ctrl *ClassController) GetClassNotice(c *gin.Context) {
	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	info, err := class.GetClassNotice(uint(classID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, info)
}
