package comment

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/service/comment"
	"questionhelper-server/pkg/response"
)

// ReportComment 举报评论
// @Summary      举报评论
// @Description  举报一条评论
// @Tags         举报管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint                      true  "评论ID"
// @Param        req  body      dto.ReportCommentRequest   true  "举报信息"
// @Success      200  {object}  response.Response  "举报成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /comments/{id}/report [post]
// @Security     BearerAuth
func (ctrl *CommentController) ReportComment(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的评论ID")
		return
	}

	var req dto.ReportCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := comment.ReportComment(uint(id), userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "举报成功", nil)
}

// ListReports 举报列表
// @Summary      举报列表
// @Description  管理员分页获取举报列表
// @Tags         举报管理
// @Accept       json
// @Produce      json
// @Param        page       query     int     false  "页码"
// @Param        page_size  query     int     false  "每页数量"
// @Param        status     query     string  false  "处理状态"
// @Success      200  {object}  response.Response{data=object{list=[]dto.ReportInfo,total=int64,page=int,page_size=int}}  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/comments/reports [get]
// @Security     BearerAuth
func (ctrl *CommentAdminController) ListReports(c *gin.Context) {
	var req dto.ReportListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	list, total, err := comment.ListReports(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// HandleReport 处理举报
// @Summary      处理举报
// @Description  管理员处理指定举报
// @Tags         举报管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint                     true  "举报ID"
// @Param        req  body      dto.HandleReportRequest  true  "处理信息"
// @Success      200  {object}  response.Response  "处理成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/comments/reports/{id} [put]
// @Security     BearerAuth
func (ctrl *CommentAdminController) HandleReport(c *gin.Context) {
	handlerID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的举报ID")
		return
	}

	var req dto.HandleReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := comment.HandleReport(uint(id), handlerID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "处理成功", nil)
}
