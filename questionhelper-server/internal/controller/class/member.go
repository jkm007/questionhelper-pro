package class

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/service/class"
	"questionhelper-server/pkg/response"
)

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

// AdminListMembers 成员列表（管理员）
// @Summary      获取班级成员列表
// @Description  管理员获取班级成员列表
// @Tags         班级管理
// @Accept       json
// @Produce      json
// @Param        id path int true "班级ID"
// @Param        page query int false "页码"
// @Param        page_size query int false "每页数量"
// @Security     BearerAuth
// @Success      200 {object} response.Response
// @Router       /admin/class/{id}/members [get]
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
// @Summary      移除班级成员
// @Description  管理员移除班级成员
// @Tags         班级管理
// @Accept       json
// @Produce      json
// @Param        id path int true "班级ID"
// @Param        uid path int true "用户ID"
// @Security     BearerAuth
// @Success      200 {object} response.Response
// @Router       /admin/class/{id}/members/{uid} [delete]
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
