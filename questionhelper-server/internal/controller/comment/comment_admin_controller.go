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
// @Summary      管理员评论列表
// @Description  管理员分页获取评论列表（含审核状态筛选）
// @Tags         评论审核
// @Accept       json
// @Produce      json
// @Param        page       query     int     false  "页码"
// @Param        page_size  query     int     false  "每页数量"
// @Param        status     query     string  false  "审核状态"
// @Param        keyword    query     string  false  "搜索关键词"
// @Success      200  {object}  response.Response{data=object{list=[]dto.CommentInfo,total=int64,page=int,page_size=int}}  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/comments [get]
// @Security     BearerAuth
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
// @Summary      置顶评论
// @Description  管理员置顶指定评论
// @Tags         评论审核
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "评论ID"
// @Success      200  {object}  response.Response  "置顶成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/comments/{id}/pin [post]
// @Security     BearerAuth
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
// @Summary      取消置顶
// @Description  管理员取消评论置顶
// @Tags         评论审核
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "评论ID"
// @Success      200  {object}  response.Response  "取消置顶成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/comments/{id}/pin [delete]
// @Security     BearerAuth
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
// @Summary      精选评论
// @Description  管理员将评论标记为精选
// @Tags         评论审核
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "评论ID"
// @Success      200  {object}  response.Response  "精选成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/comments/{id}/featured [post]
// @Security     BearerAuth
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
// @Summary      取消精选
// @Description  管理员取消评论精选标记
// @Tags         评论审核
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "评论ID"
// @Success      200  {object}  response.Response  "取消精选成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/comments/{id}/featured [delete]
// @Security     BearerAuth
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
// @Summary      标记官方解答
// @Description  管理员将评论标记为官方解答
// @Tags         评论审核
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "评论ID"
// @Success      200  {object}  response.Response  "标记官方解答成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/comments/{id}/official [post]
// @Security     BearerAuth
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
// @Summary      取消官方解答
// @Description  管理员取消评论的官方解答标记
// @Tags         评论审核
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "评论ID"
// @Success      200  {object}  response.Response  "取消官方解答成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/comments/{id}/official [delete]
// @Security     BearerAuth
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
// @Summary      黑名单列表
// @Description  分页获取评论黑名单列表
// @Tags         评论审核
// @Accept       json
// @Produce      json
// @Param        page       query     int  false  "页码"
// @Param        page_size  query     int  false  "每页数量"
// @Success      200  {object}  response.Response{data=object{list=[]dto.BlacklistInfo,total=int64,page=int,page_size=int}}  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/comments/blacklist [get]
// @Security     BearerAuth
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
// @Summary      添加黑名单
// @Description  将用户添加到评论黑名单
// @Tags         评论审核
// @Accept       json
// @Produce      json
// @Param        req  body      dto.AddBlacklistRequest  true  "黑名单信息"
// @Success      200  {object}  response.Response  "添加黑名单成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/comments/blacklist [post]
// @Security     BearerAuth
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
// @Summary      移除黑名单
// @Description  将用户从评论黑名单中移除
// @Tags         评论审核
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "黑名单ID"
// @Success      200  {object}  response.Response  "移除黑名单成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/comments/blacklist/{id} [delete]
// @Security     BearerAuth
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
// @Summary      审核规则列表
// @Description  获取所有评论审核规则
// @Tags         评论审核
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response{data=[]dto.AuditRuleInfo}  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/comments/audit-rules [get]
// @Security     BearerAuth
func (ctrl *CommentAdminController) ListAuditRules(c *gin.Context) {
	list, err := comment.ListAuditRules()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, list)
}

// CreateAuditRule 创建审核规则
// @Summary      创建审核规则
// @Description  创建一条新的评论审核规则
// @Tags         评论审核
// @Accept       json
// @Produce      json
// @Param        req  body      dto.CreateAuditRuleRequest  true  "审核规则信息"
// @Success      200  {object}  response.Response  "创建成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/comments/audit-rules [post]
// @Security     BearerAuth
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
// @Summary      更新审核规则
// @Description  更新指定审核规则
// @Tags         评论审核
// @Accept       json
// @Produce      json
// @Param        id   path      uint                       true  "规则ID"
// @Param        req  body      dto.UpdateAuditRuleRequest  true  "审核规则信息"
// @Success      200  {object}  response.Response  "更新成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/comments/audit-rules/{id} [put]
// @Security     BearerAuth
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
// @Summary      删除审核规则
// @Description  删除指定审核规则
// @Tags         评论审核
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "规则ID"
// @Success      200  {object}  response.Response  "删除成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/comments/audit-rules/{id} [delete]
// @Security     BearerAuth
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
// @Summary      批量审核
// @Description  批量审核评论（通过/拒绝）
// @Tags         评论审核
// @Accept       json
// @Produce      json
// @Param        req  body      dto.BatchAuditRequest  true  "批量审核信息"
// @Success      200  {object}  response.Response  "批量审核成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/comments/batch-audit [put]
// @Security     BearerAuth
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
// @Summary      批量删除评论
// @Description  批量删除指定评论
// @Tags         评论审核
// @Accept       json
// @Produce      json
// @Param        req  body      dto.BatchCommentDeleteRequest  true  "评论ID列表"
// @Success      200  {object}  response.Response  "批量删除成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/comments/batch-delete [post]
// @Security     BearerAuth
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
// @Summary      评论统计
// @Description  获取评论统计数据
// @Tags         评论审核
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response{data=dto.CommentStats}  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/comments/stats [get]
// @Security     BearerAuth
func (ctrl *CommentAdminController) GetCommentStats(c *gin.Context) {
	stats, err := comment.GetCommentStats()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, stats)
}

// ExportComments 导出评论
// @Summary      导出评论
// @Description  导出评论数据为 Excel 文件
// @Tags         评论审核
// @Accept       json
// @Produce      application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Param        page       query     int     false  "页码"
// @Param        page_size  query     int     false  "每页数量"
// @Param        status     query     string  false  "审核状态"
// @Success      200  {file}    binary  "Excel 文件"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/comments/export [get]
// @Security     BearerAuth
func (ctrl *CommentAdminController) ExportComments(c *gin.Context) {
	var req dto.CommentExportRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	f, err := comment.ExportComments(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "导出失败: "+err.Error())
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=comments.xlsx")
	if err := f.Write(c.Writer); err != nil {
		response.Error(c, http.StatusInternalServerError, "写入文件失败")
		return
	}
}
