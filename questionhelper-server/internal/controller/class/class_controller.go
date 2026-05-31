package class

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/service/class"
	"questionhelper-server/pkg/response"
)

type ClassController struct{}

func NewClassController() *ClassController {
	return &ClassController{}
}

// ListClasses 我的班级列表
func (ctrl *ClassController) ListClasses(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.PageRequest
	c.ShouldBindQuery(&req)

	list, total, err := class.ListMyClasses(userID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// GetClass 获取班级详情
func (ctrl *ClassController) GetClass(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	info, err := class.GetClass(uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, info)
}

// CreateClass 创建班级
func (ctrl *ClassController) CreateClass(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.CreateClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := class.CreateClass(userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "创建成功", nil)
}

// UpdateClass 更新班级
func (ctrl *ClassController) UpdateClass(c *gin.Context) {
	userID := c.GetUint("user_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	var req dto.UpdateClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := class.UpdateClass(uint(id), userID, false, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "更新成功", nil)
}

// DeleteClass 删除班级
func (ctrl *ClassController) DeleteClass(c *gin.Context) {
	userID := c.GetUint("user_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	if err := class.DeleteClass(uint(id), userID, false); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "删除成功", nil)
}

// JoinClass 加入班级
func (ctrl *ClassController) JoinClass(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.JoinClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := class.JoinClass(userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "加入成功", nil)
}

// LeaveClass 离开班级
func (ctrl *ClassController) LeaveClass(c *gin.Context) {
	userID := c.GetUint("user_id")
	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	if err := class.LeaveClass(uint(classID), userID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "已离开班级", nil)
}

// ListMembers 成员列表
func (ctrl *ClassController) ListMembers(c *gin.Context) {
	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	var req dto.ClassMemberListRequest
	c.ShouldBindQuery(&req)

	list, total, err := class.ListMembers(uint(classID), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// ListNotices 公告列表
func (ctrl *ClassController) ListNotices(c *gin.Context) {
	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	var req dto.PageRequest
	c.ShouldBindQuery(&req)

	list, total, err := class.ListNotices(uint(classID), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// ListHomework 作业列表
func (ctrl *ClassController) ListHomework(c *gin.Context) {
	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	var req dto.PageRequest
	c.ShouldBindQuery(&req)

	list, total, err := class.ListHomework(uint(classID), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// AdminListClasses 班级列表（管理员）
func (ctrl *ClassController) AdminListClasses(c *gin.Context) {
	var req dto.ClassListRequest
	c.ShouldBindQuery(&req)

	list, total, err := class.ListClasses(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// AdminGetClass 获取班级详情（管理员）
func (ctrl *ClassController) AdminGetClass(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	info, err := class.GetClass(uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, info)
}

// AdminCreateClass 创建班级（管理员）
func (ctrl *ClassController) AdminCreateClass(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.CreateClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := class.CreateClass(userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "创建成功", nil)
}

// AdminUpdateClass 更新班级（管理员）
func (ctrl *ClassController) AdminUpdateClass(c *gin.Context) {
	userID := c.GetUint("user_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	var req dto.UpdateClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := class.UpdateClass(uint(id), userID, true, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "更新成功", nil)
}

// AdminDeleteClass 删除班级（管理员）
func (ctrl *ClassController) AdminDeleteClass(c *gin.Context) {
	userID := c.GetUint("user_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	if err := class.DeleteClass(uint(id), userID, true); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "删除成功", nil)
}

// AdminListMembers 成员列表（管理员）
func (ctrl *ClassController) AdminListMembers(c *gin.Context) {
	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	var req dto.ClassMemberListRequest
	c.ShouldBindQuery(&req)

	list, total, err := class.ListMembers(uint(classID), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// RemoveMember 移除成员（管理员）
func (ctrl *ClassController) RemoveMember(c *gin.Context) {
	operatorID := c.GetUint("user_id")

	classID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的班级ID")
		return
	}

	uid, err := strconv.ParseUint(c.Param("uid"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	if err := class.RemoveMember(uint(classID), uint(uid), operatorID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "移除成功", nil)
}
