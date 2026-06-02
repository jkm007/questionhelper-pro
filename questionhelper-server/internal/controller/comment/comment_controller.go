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
// @Summary      获取评论列表
// @Description  分页获取评论列表
// @Tags         评论管理
// @Accept       json
// @Produce      json
// @Param        page       query     int  false  "页码"
// @Param        page_size  query     int  false  "每页数量"
// @Param        target_id  query     uint false  "目标ID"
// @Param        sort       query     string false "排序方式"
// @Success      200  {object}  response.Response{data=object{list=[]dto.CommentInfo,total=int64,page=int,page_size=int}}  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /comments [get]
// @Security     BearerAuth
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
// @Summary      创建评论
// @Description  创建一条新评论
// @Tags         评论管理
// @Accept       json
// @Produce      json
// @Param        req  body      dto.CreateCommentRequest  true  "评论内容"
// @Success      200  {object}  response.Response  "评论成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /comments [post]
// @Security     BearerAuth
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
// @Summary      编辑评论
// @Description  编辑指定评论内容
// @Tags         评论管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint                    true  "评论ID"
// @Param        req  body      dto.UpdateCommentRequest true  "更新内容"
// @Success      200  {object}  response.Response  "编辑成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /comments/{id} [put]
// @Security     BearerAuth
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
// @Summary      删除评论
// @Description  删除指定评论
// @Tags         评论管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "评论ID"
// @Success      200  {object}  response.Response  "删除成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /comments/{id} [delete]
// @Security     BearerAuth
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
// @Summary      点赞评论
// @Description  对指定评论进行点赞或取消点赞
// @Tags         评论管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "评论ID"
// @Success      200  {object}  response.Response  "操作成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /comments/{id}/like [post]
// @Security     BearerAuth
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

// UploadImage 上传评论图片
// @Summary      上传评论图片
// @Description  上传评论中使用的图片
// @Tags         评论管理
// @Accept       multipart/form-data
// @Produce      json
// @Param        file  formData  file  true  "图片文件"
// @Success      200   {object}  response.Response{data=dto.UploadImageResponse}  "上传成功"
// @Failure      400   {object}  response.Response  "参数错误"
// @Failure      500   {object}  response.Response  "服务器内部错误"
// @Router       /comments/upload-image [post]
// @Security     BearerAuth
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
// @Summary      获取表情包列表
// @Description  根据分类获取表情包列表
// @Tags         评论管理
// @Accept       json
// @Produce      json
// @Param        category  query     string  false  "表情分类"
// @Success      200       {object}  response.Response{data=[]dto.StickerInfo}  "成功"
// @Failure      500       {object}  response.Response  "服务器内部错误"
// @Router       /sticker/list [get]
// @Security     BearerAuth
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
// @Summary      获取表情分类
// @Description  获取所有表情分类列表
// @Tags         评论管理
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response{data=[]dto.StickerCategoryInfo}  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /sticker/categories [get]
// @Security     BearerAuth
func (ctrl *CommentController) ListStickerCategories(c *gin.Context) {
	categories, err := comment.ListStickerCategories()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, categories)
}

// SearchUsers 搜索用户（用于@功能）
// @Summary      搜索用户
// @Description  搜索用户用于评论@功能
// @Tags         评论管理
// @Accept       json
// @Produce      json
// @Param        keyword    query     string  true   "搜索关键词"
// @Param        page_size  query     int     false  "返回数量"
// @Success      200        {object}  response.Response{data=[]dto.UserSearchInfo}  "成功"
// @Failure      400        {object}  response.Response  "参数错误"
// @Failure      500        {object}  response.Response  "服务器内部错误"
// @Router       /user/search [get]
// @Security     BearerAuth
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
