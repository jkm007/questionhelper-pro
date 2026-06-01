package comment

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/service/comment"
	"questionhelper-server/pkg/response"
)

type CommentAdminController struct{}

func NewCommentAdminController() *CommentAdminController {
	return &CommentAdminController{}
}

// ==================== Comment Management ====================

// ListComments 管理员评论列表
func (ctrl *CommentAdminController) ListComments(c *gin.Context) {
	var req dto.CommentAdminListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	list, total, err := comment.ListCommentsAdmin(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// PinComment 置顶评论
func (ctrl *CommentAdminController) PinComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的评论ID")
		return
	}

	if err := comment.PinComment(uint(id), true); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "置顶成功", nil)
}

// UnpinComment 取消置顶
func (ctrl *CommentAdminController) UnpinComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的评论ID")
		return
	}

	if err := comment.PinComment(uint(id), false); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "取消置顶成功", nil)
}

// FeatureComment 精选评论
func (ctrl *CommentAdminController) FeatureComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的评论ID")
		return
	}

	if err := comment.FeatureComment(uint(id), true); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "精选成功", nil)
}

// UnfeatureComment 取消精选
func (ctrl *CommentAdminController) UnfeatureComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的评论ID")
		return
	}

	if err := comment.FeatureComment(uint(id), false); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "取消精选成功", nil)
}

// MarkOfficial 标记官方解答
func (ctrl *CommentAdminController) MarkOfficial(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的评论ID")
		return
	}

	if err := comment.OfficialComment(uint(id), true); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "标记官方解答成功", nil)
}

// UnmarkOfficial 取消官方解答
func (ctrl *CommentAdminController) UnmarkOfficial(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的评论ID")
		return
	}

	if err := comment.OfficialComment(uint(id), false); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "取消官方解答成功", nil)
}

// ==================== Blacklist ====================

// ListBlacklists 黑名单列表
func (ctrl *CommentAdminController) ListBlacklists(c *gin.Context) {
	var req dto.BlacklistListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	list, total, err := comment.ListBlacklists(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// AddBlacklist 添加黑名单
func (ctrl *CommentAdminController) AddBlacklist(c *gin.Context) {
	operatorID := c.GetUint("user_id")

	var req dto.AddBlacklistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := comment.AddBlacklist(operatorID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "添加黑名单成功", nil)
}

// RemoveBlacklist 移除黑名单
func (ctrl *CommentAdminController) RemoveBlacklist(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的黑名单ID")
		return
	}

	if err := comment.RemoveBlacklist(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "移除黑名单成功", nil)
}

// ==================== Audit Rules ====================

// ListAuditRules 审核规则列表
func (ctrl *CommentAdminController) ListAuditRules(c *gin.Context) {
	list, err := comment.ListAuditRules()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, list)
}

// CreateAuditRule 创建审核规则
func (ctrl *CommentAdminController) CreateAuditRule(c *gin.Context) {
	var req dto.CreateAuditRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := comment.CreateAuditRule(&req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "创建成功", nil)
}

// UpdateAuditRule 更新审核规则
func (ctrl *CommentAdminController) UpdateAuditRule(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的规则ID")
		return
	}

	var req dto.UpdateAuditRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := comment.UpdateAuditRule(uint(id), &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "更新成功", nil)
}

// DeleteAuditRule 删除审核规则
func (ctrl *CommentAdminController) DeleteAuditRule(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的规则ID")
		return
	}

	if err := comment.DeleteAuditRule(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "删除成功", nil)
}

// ==================== Batch Operations ====================

// BatchAudit 批量审核
func (ctrl *CommentAdminController) BatchAudit(c *gin.Context) {
	var req dto.BatchAuditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := comment.BatchAudit(&req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "批量审核成功", nil)
}

// BatchDelete 批量删除
func (ctrl *CommentAdminController) BatchDelete(c *gin.Context) {
	var req dto.BatchCommentDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := comment.BatchDelete(req.IDs); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "批量删除成功", nil)
}

// ==================== Stats & Export ====================

// GetCommentStats 评论统计
func (ctrl *CommentAdminController) GetCommentStats(c *gin.Context) {
	stats, err := comment.GetCommentStats()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, stats)
}

// ExportComments 导出评论
func (ctrl *CommentAdminController) ExportComments(c *gin.Context) {
	var req dto.CommentExportRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	list, err := comment.ExportComments(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, list)
}
