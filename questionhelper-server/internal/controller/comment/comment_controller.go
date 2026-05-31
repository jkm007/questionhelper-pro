package comment

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/dto"
	fileService "questionhelper-server/internal/service/file"
	"questionhelper-server/internal/service/comment"
	"questionhelper-server/pkg/response"
)

type CommentController struct{}

func NewCommentController() *CommentController {
	return &CommentController{}
}

// ListComments 评论列表
func (ctrl *CommentController) ListComments(c *gin.Context) {
	var req dto.CommentListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	list, total, err := comment.ListComments(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// CreateComment 创建评论
func (ctrl *CommentController) CreateComment(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := comment.CreateComment(userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "评论成功", nil)
}

// EditComment 编辑评论
func (ctrl *CommentController) EditComment(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的评论ID")
		return
	}

	var req dto.UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := comment.EditComment(uint(id), userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "编辑成功", nil)
}

// DeleteComment 删除评论
func (ctrl *CommentController) DeleteComment(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的评论ID")
		return
	}

	if err := comment.DeleteComment(uint(id), userID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "删除成功", nil)
}

// LikeComment 点赞评论
func (ctrl *CommentController) LikeComment(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的评论ID")
		return
	}

	if err := comment.LikeComment(uint(id), userID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "操作成功", nil)
}

// ReportComment 举报评论
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

// UploadImage 上传评论图片
func (ctrl *CommentController) UploadImage(c *gin.Context) {
	userID := c.GetUint("user_id")

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "请选择文件")
		return
	}
	defer file.Close()

	result, err := fileService.UploadFile(userID, header.Filename, header.Size, header.Header.Get("Content-Type"), file)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, dto.UploadImageResponse{
		URL:  result.URL,
		Name: result.Name,
		Size: result.Size,
	})
}

// ListStickers 获取表情包列表
func (ctrl *CommentController) ListStickers(c *gin.Context) {
	category := c.Query("category")

	list, err := comment.ListStickers(category)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, list)
}

// ListStickerCategories 获取表情分类
func (ctrl *CommentController) ListStickerCategories(c *gin.Context) {
	categories, err := comment.ListStickerCategories()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, categories)
}

// SearchUsers 搜索用户（用于@功能）
func (ctrl *CommentController) SearchUsers(c *gin.Context) {
	var req dto.UserSearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	limit := req.PageSize
	if limit <= 0 {
		limit = 10
	}

	list, err := comment.SearchUsers(req.Keyword, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, list)
}
