package question

import (
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/service/question"
	"questionhelper-server/internal/service/sensitive"
	"questionhelper-server/pkg/response"
)

type QuestionController struct{}

func NewQuestionController() *QuestionController {
	return &QuestionController{}
}

// ListQuestions 题目列表（带数据权限过滤）
// @Summary      获取题目列表
// @Description  获取题目列表，支持分页和筛选，带数据权限过滤
// @Tags         题目管理
// @Accept       json
// @Produce      json
// @Param        keyword       query     string  false  "搜索关键词"
// @Param        category_id   query     uint    false  "分类ID"
// @Param        type          query     int     false  "题型"
// @Param        difficulty    query     int     false  "难度"
// @Param        visibility    query     int     false  "可见性"
// @Param        status        query     int     false  "状态"
// @Param        page          query     int     false  "页码"
// @Param        page_size     query     int     false  "每页数量"
// @Success      200  {object}  response.Response{data=dto.PageResponse{list=[]dto.QuestionInfo}}  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /questions [get]
// @Security     BearerAuth
func (ctrl *QuestionController) ListQuestions(c *gin.Context) {
	var req dto.QuestionListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	userID := c.GetUint("user_id")
	list, total, err := question.ListQuestionsWithPermission(userID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// GetQuestion 获取题目详情
// @Summary      获取题目详情
// @Description  根据题目ID获取题目详细信息
// @Tags         题目管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "题目ID"
// @Success      200  {object}  response.Response{data=dto.QuestionInfo}  "成功"
// @Failure      400  {object}  response.Response  "无效的题目ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /questions/{id} [get]
// @Security     BearerAuth
func (ctrl *QuestionController) GetQuestion(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的题目ID")
		return
	}

	info, err := question.GetQuestion(uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, info)
}

// SearchQuestions 搜索题目
// @Summary      搜索题目
// @Description  根据关键词搜索题目
// @Tags         题目管理
// @Accept       json
// @Produce      json
// @Param        keyword    query     string  true   "搜索关键词"
// @Param        page       query     int     false  "页码"
// @Param        page_size  query     int     false  "每页数量"
// @Success      200  {object}  response.Response{data=dto.PageResponse{list=[]dto.QuestionInfo}}  "成功"
// @Failure      400  {object}  response.Response  "请输入搜索关键词"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /questions/search [get]
// @Security     BearerAuth
func (ctrl *QuestionController) SearchQuestions(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		response.Error(c, http.StatusBadRequest, "请输入搜索关键词")
		return
	}

	var req dto.PageRequest
	c.ShouldBindQuery(&req)

	list, total, err := question.SearchQuestions(keyword, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// FavoriteQuestion 收藏/取消收藏题目
// @Summary      收藏/取消收藏题目
// @Description  收藏或取消收藏指定题目，通过action参数控制
// @Tags         题目管理
// @Accept       json
// @Produce      json
// @Param        id      path      uint                   true  "题目ID"
// @Param        action  query     string                 false  "操作类型(remove为取消收藏)"
// @Param        req     body      dto.FavoriteRequest    false  "收藏信息"
// @Success      200  {object}  response.Response  "操作成功"
// @Failure      400  {object}  response.Response  "无效的题目ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /questions/{id}/favorite [post]
// @Security     BearerAuth
func (ctrl *QuestionController) FavoriteQuestion(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的题目ID")
		return
	}

	action := c.Query("action")
	if action == "remove" {
		if err := question.UnfavoriteQuestion(userID, uint(id)); err != nil {
			response.Error(c, http.StatusInternalServerError, err.Error())
			return
		}
		response.SuccessWithMessage(c, "取消收藏成功", nil)
	} else {
		var req dto.FavoriteRequest
		c.ShouldBindJSON(&req)
		if err := question.FavoriteQuestion(userID, uint(id), &req); err != nil {
			response.Error(c, http.StatusInternalServerError, err.Error())
			return
		}
		response.SuccessWithMessage(c, "收藏成功", nil)
	}
}

// LikeQuestion 点赞题目
// @Summary      点赞题目
// @Description  对指定题目进行点赞
// @Tags         题目管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "题目ID"
// @Success      200  {object}  response.Response  "点赞成功"
// @Failure      400  {object}  response.Response  "无效的题目ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /questions/{id}/like [post]
// @Security     BearerAuth
func (ctrl *QuestionController) LikeQuestion(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的题目ID")
		return
	}

	if err := question.LikeQuestion(userID, uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "点赞成功", nil)
}

// ListCategories 分类列表
// @Summary      获取分类列表
// @Description  获取所有题目分类列表
// @Tags         分类管理
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response{data=[]dto.CategoryInfo}  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /categories [get]
// @Security     BearerAuth
func (ctrl *QuestionController) ListCategories(c *gin.Context) {
	list, err := question.ListCategories()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, list)
}

// AdminListQuestions 管理员题目列表
// @Summary      管理员获取题目列表
// @Description  管理员获取题目列表，支持分页和筛选
// @Tags         题目管理
// @Accept       json
// @Produce      json
// @Param        keyword       query     string  false  "搜索关键词"
// @Param        category_id   query     uint    false  "分类ID"
// @Param        type          query     int     false  "题型"
// @Param        difficulty    query     int     false  "难度"
// @Param        visibility    query     int     false  "可见性"
// @Param        status        query     int     false  "状态"
// @Param        page          query     int     false  "页码"
// @Param        page_size     query     int     false  "每页数量"
// @Success      200  {object}  response.Response{data=dto.PageResponse{list=[]dto.QuestionInfo}}  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/questions [get]
// @Security     BearerAuth
func (ctrl *QuestionController) AdminListQuestions(c *gin.Context) {
	var req dto.QuestionListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	list, total, err := question.ListQuestions(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// AdminGetQuestion 管理员获取题目详情
// @Summary      管理员获取题目详情
// @Description  管理员根据题目ID获取题目详细信息
// @Tags         题目管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "题目ID"
// @Success      200  {object}  response.Response{data=dto.QuestionInfo}  "成功"
// @Failure      400  {object}  response.Response  "无效的题目ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/questions/{id} [get]
// @Security     BearerAuth
func (ctrl *QuestionController) AdminGetQuestion(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的题目ID")
		return
	}

	info, err := question.GetQuestion(uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, info)
}

// CreateQuestion 创建题目
// @Summary      创建题目
// @Description  管理员创建新题目
// @Tags         题目管理
// @Accept       json
// @Produce      json
// @Param        req  body      dto.CreateQuestionRequest  true  "题目信息"
// @Success      200  {object}  response.Response  "创建成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/questions [post]
// @Security     BearerAuth
func (ctrl *QuestionController) CreateQuestion(c *gin.Context) {
	creatorID := c.GetUint("user_id")

	var req dto.CreateQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := question.CreateQuestion(creatorID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "创建成功", nil)
}

// UpdateQuestion 更新题目
// @Summary      更新题目
// @Description  管理员更新指定题目信息
// @Tags         题目管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint                       true  "题目ID"
// @Param        req  body      dto.UpdateQuestionRequest   true  "更新信息"
// @Success      200  {object}  response.Response  "更新成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/questions/{id} [put]
// @Security     BearerAuth
func (ctrl *QuestionController) UpdateQuestion(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的题目ID")
		return
	}

	var req dto.UpdateQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := question.UpdateQuestion(uint(id), &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "更新成功", nil)
}

// DeleteQuestion 删除题目
// @Summary      删除题目
// @Description  管理员删除指定题目（软删除）
// @Tags         题目管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "题目ID"
// @Success      200  {object}  response.Response  "删除成功"
// @Failure      400  {object}  response.Response  "无效的题目ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/questions/{id} [delete]
// @Security     BearerAuth
func (ctrl *QuestionController) DeleteQuestion(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的题目ID")
		return
	}

	if err := question.DeleteQuestion(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "删除成功", nil)
}

// UpdateQuestionStatus 更新题目状态
// @Summary      更新题目状态
// @Description  管理员更新指定题目的状态（草稿/发布/归档）
// @Tags         题目管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint              true  "题目ID"
// @Param        req  body      dto.StatusRequest  true  "状态信息"
// @Success      200  {object}  response.Response  "状态更新成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/questions/{id}/status [put]
// @Security     BearerAuth
func (ctrl *QuestionController) UpdateQuestionStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的题目ID")
		return
	}

	var req dto.StatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := question.UpdateQuestionStatus(uint(id), req.Status); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "状态更新成功", nil)
}

// CreateFavoriteFolder 创建收藏夹
// @Summary      创建收藏夹
// @Description  创建新的收藏夹
// @Tags         收藏夹
// @Accept       json
// @Produce      json
// @Param        req  body      dto.CreateFolderRequest  true  "收藏夹信息"
// @Success      200  {object}  response.Response  "创建收藏夹成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /favorites/folders [post]
// @Security     BearerAuth
func (ctrl *QuestionController) CreateFavoriteFolder(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.CreateFolderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := question.CreateFolder(userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "创建收藏夹成功", nil)
}

// UpdateFavoriteFolder 更新收藏夹
// @Summary      更新收藏夹
// @Description  更新指定收藏夹信息
// @Tags         收藏夹
// @Accept       json
// @Produce      json
// @Param        id   path      uint                      true  "收藏夹ID"
// @Param        req  body      dto.UpdateFolderRequest    true  "更新信息"
// @Success      200  {object}  response.Response  "更新收藏夹成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /favorites/folders/{id} [put]
// @Security     BearerAuth
func (ctrl *QuestionController) UpdateFavoriteFolder(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的收藏夹ID")
		return
	}

	var req dto.UpdateFolderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := question.UpdateFolder(uint(id), userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "更新收藏夹成功", nil)
}

// DeleteFavoriteFolder 删除收藏夹
// @Summary      删除收藏夹
// @Description  删除指定收藏夹
// @Tags         收藏夹
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "收藏夹ID"
// @Success      200  {object}  response.Response  "删除收藏夹成功"
// @Failure      400  {object}  response.Response  "无效的收藏夹ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /favorites/folders/{id} [delete]
// @Security     BearerAuth
func (ctrl *QuestionController) DeleteFavoriteFolder(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的收藏夹ID")
		return
	}

	if err := question.DeleteFolder(uint(id), userID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "删除收藏夹成功", nil)
}

// ListFavoriteFolders 获取收藏夹列表
// @Summary      获取收藏夹列表
// @Description  获取当前用户的收藏夹列表
// @Tags         收藏夹
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response{data=[]dto.FolderInfo}  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /favorites/folders [get]
// @Security     BearerAuth
func (ctrl *QuestionController) ListFavoriteFolders(c *gin.Context) {
	userID := c.GetUint("user_id")

	list, err := question.ListFolders(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, list)
}

// ==================== 内容审核 ====================

// ListPendingReviews 待审核列表
// @Summary      获取待审核列表
// @Description  管理员获取待审核的题目列表
// @Tags         题目审核
// @Accept       json
// @Produce      json
// @Param        page        query     int  false  "页码"
// @Param        page_size   query     int  false  "每页数量"
// @Success      200  {object}  response.Response{data=dto.PageResponse{list=[]dto.ReviewInfo}}  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/reviews [get]
// @Security     BearerAuth
func (ctrl *QuestionController) ListPendingReviews(c *gin.Context) {
	var req dto.PageRequest
	c.ShouldBindQuery(&req)

	list, total, err := question.ListPendingReviews(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// GetReviewDetail 审核详情
// @Summary      获取审核详情
// @Description  管理员获取指定审核记录的详细信息
// @Tags         题目审核
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "审核ID"
// @Success      200  {object}  response.Response{data=dto.ReviewInfo}  "成功"
// @Failure      400  {object}  response.Response  "无效的审核ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/reviews/{id} [get]
// @Security     BearerAuth
func (ctrl *QuestionController) GetReviewDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的审核ID")
		return
	}

	info, err := question.GetReviewDetail(uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, info)
}

// ApproveReview 审核通过
// @Summary      审核通过
// @Description  管理员审核通过指定题目
// @Tags         题目审核
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "审核ID"
// @Success      200  {object}  response.Response  "审核通过"
// @Failure      400  {object}  response.Response  "无效的审核ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/reviews/{id}/approve [post]
// @Security     BearerAuth
func (ctrl *QuestionController) ApproveReview(c *gin.Context) {
	reviewerID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的审核ID")
		return
	}

	if err := question.ApproveReview(uint(id), reviewerID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "审核通过", nil)
}

// RejectReview 审核拒绝
// @Summary      审核拒绝
// @Description  管理员审核拒绝指定题目，需提供拒绝原因
// @Tags         题目审核
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "审核ID"
// @Param        req  body      object{reason=string}  true  "拒绝原因"
// @Success      200  {object}  response.Response  "审核拒绝"
// @Failure      400  {object}  response.Response  "无效的审核ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/reviews/{id}/reject [post]
// @Security     BearerAuth
func (ctrl *QuestionController) RejectReview(c *gin.Context) {
	reviewerID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的审核ID")
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	c.ShouldBindJSON(&req)

	if err := question.RejectReview(uint(id), reviewerID, req.Reason); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "审核拒绝", nil)
}

// ==================== 敏感词管理 ====================

// ListSensitiveWords 敏感词列表
// @Summary      获取敏感词列表
// @Description  管理员获取敏感词列表，支持分页
// @Tags         敏感词管理
// @Accept       json
// @Produce      json
// @Param        page        query     int  false  "页码"
// @Param        page_size   query     int  false  "每页数量"
// @Success      200  {object}  response.Response{data=dto.PageResponse}  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/sensitive-words [get]
// @Security     BearerAuth
func (ctrl *QuestionController) ListSensitiveWords(c *gin.Context) {
	var req dto.PageRequest
	c.ShouldBindQuery(&req)

	list, total, err := sensitive.ListSensitiveWords(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// CreateSensitiveWord 创建敏感词
// @Summary      创建敏感词
// @Description  管理员创建新的敏感词
// @Tags         敏感词管理
// @Accept       json
// @Produce      json
// @Param        req  body      dto.CreateSensitiveWordRequest  true  "敏感词信息"
// @Success      200  {object}  response.Response  "创建敏感词成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/sensitive-words [post]
// @Security     BearerAuth
func (ctrl *QuestionController) CreateSensitiveWord(c *gin.Context) {
	var req dto.CreateSensitiveWordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := sensitive.CreateSensitiveWord(&req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "创建敏感词成功", nil)
}

// DeleteSensitiveWord 删除敏感词
// @Summary      删除敏感词
// @Description  管理员删除指定敏感词
// @Tags         敏感词管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "敏感词ID"
// @Success      200  {object}  response.Response  "删除敏感词成功"
// @Failure      400  {object}  response.Response  "无效的敏感词ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/sensitive-words/{id} [delete]
// @Security     BearerAuth
func (ctrl *QuestionController) DeleteSensitiveWord(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的敏感词ID")
		return
	}

	if err := sensitive.DeleteSensitiveWord(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "删除敏感词成功", nil)
}

// ImportSensitiveWords 导入敏感词
// @Summary      导入敏感词
// @Description  管理员通过文件批量导入敏感词
// @Tags         敏感词管理
// @Accept       multipart/form-data
// @Produce      json
// @Param        file  formData  file  true  "敏感词文件"
// @Success      200  {object}  response.Response{data=object{count=int}}  "导入成功"
// @Failure      400  {object}  response.Response  "请上传文件"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/sensitive-words/import [post]
// @Security     BearerAuth
func (ctrl *QuestionController) ImportSensitiveWords(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "请上传文件")
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "读取文件失败")
		return
	}

	count, err := sensitive.ImportSensitiveWords(data)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithMessage(c, "导入成功", gin.H{"count": count})
}

// TestSensitiveWord 测试敏感词
// @Summary      测试敏感词
// @Description  测试文本是否包含敏感词
// @Tags         敏感词管理
// @Accept       json
// @Produce      json
// @Param        req  body      object{content=string}  true  "测试内容"
// @Success      200  {object}  response.Response{data=object{has_sensitive=bool}}  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/sensitive-words/test [post]
// @Security     BearerAuth
func (ctrl *QuestionController) TestSensitiveWord(c *gin.Context) {
	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	result := sensitive.TestSensitiveWord(req.Content)
	response.Success(c, gin.H{"has_sensitive": result})
}

// ==================== 题目笔记 ====================

// GetQuestionNotes 获取题目笔记
// @Summary      获取题目笔记
// @Description  获取指定题目的笔记列表
// @Tags         题目管理
// @Accept       json
// @Produce      json
// @Param        id         path      uint   true   "题目ID"
// @Param        is_public  query     bool   false  "是否公开"
// @Param        page       query     int    false  "页码"
// @Param        page_size  query     int    false  "每页数量"
// @Success      200  {object}  response.Response{data=dto.PageResponse{list=[]dto.NoteInfo}}  "成功"
// @Failure      400  {object}  response.Response  "无效的题目ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /questions/{id}/notes [get]
// @Security     BearerAuth
func (ctrl *QuestionController) GetQuestionNotes(c *gin.Context) {
	questionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的题目ID")
		return
	}

	var req dto.NoteListRequest
	c.ShouldBindQuery(&req)

	// 如果是普通用户，只看公开笔记和自己的笔记
	userID := c.GetUint("user_id")
	isPublic := true
	if req.IsPublic == nil {
		req.IsPublic = &isPublic
	}

	list, total, err := question.GetQuestionNotes(uint(questionID), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	_ = userID // 用于后续扩展权限过滤
	response.Page(c, list, total, req.Page, req.PageSize)
}

// CreateNote 创建笔记
// @Summary      创建笔记
// @Description  为指定题目创建笔记
// @Tags         题目管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint                    true  "题目ID"
// @Param        req  body      dto.CreateNoteRequest    true  "笔记信息"
// @Success      200  {object}  response.Response{data=dto.NoteInfo}  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /questions/{id}/notes [post]
// @Security     BearerAuth
func (ctrl *QuestionController) CreateNote(c *gin.Context) {
	userID := c.GetUint("user_id")
	questionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的题目ID")
		return
	}

	var req dto.CreateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	note, err := question.CreateNote(userID, uint(questionID), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, note)
}

// UpdateNote 更新笔记
// @Summary      更新笔记
// @Description  更新指定题目的笔记
// @Tags         题目管理
// @Accept       json
// @Produce      json
// @Param        id       path      uint                    true  "题目ID"
// @Param        noteId   path      uint                    true  "笔记ID"
// @Param        req      body      dto.UpdateNoteRequest    true  "更新信息"
// @Success      200  {object}  response.Response  "更新成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /questions/{id}/notes/{noteId} [put]
// @Security     BearerAuth
func (ctrl *QuestionController) UpdateNote(c *gin.Context) {
	userID := c.GetUint("user_id")
	questionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的题目ID")
		return
	}

	noteID, err := strconv.ParseUint(c.Param("noteId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的笔记ID")
		return
	}

	var req dto.UpdateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := question.UpdateNote(userID, uint(questionID), uint(noteID), &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "更新成功", nil)
}

// DeleteNote 删除笔记
// @Summary      删除笔记
// @Description  删除指定题目的笔记
// @Tags         题目管理
// @Accept       json
// @Produce      json
// @Param        id       path      uint  true  "题目ID"
// @Param        noteId   path      uint  true  "笔记ID"
// @Success      200  {object}  response.Response  "删除成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /questions/{id}/notes/{noteId} [delete]
// @Security     BearerAuth
func (ctrl *QuestionController) DeleteNote(c *gin.Context) {
	userID := c.GetUint("user_id")
	questionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的题目ID")
		return
	}

	noteID, err := strconv.ParseUint(c.Param("noteId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的笔记ID")
		return
	}

	if err := question.DeleteNote(userID, uint(questionID), uint(noteID)); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "删除成功", nil)
}

// ==================== 题目评价 ====================

// RateDifficulty 评价难度
// @Summary      评价题目难度
// @Description  对指定题目进行难度评价
// @Tags         题目管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint                           true  "题目ID"
// @Param        req  body      dto.DifficultyRatingRequest     true  "评价信息"
// @Success      200  {object}  response.Response  "评价成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /questions/{id}/difficulty-rating [post]
// @Security     BearerAuth
func (ctrl *QuestionController) RateDifficulty(c *gin.Context) {
	userID := c.GetUint("user_id")
	questionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的题目ID")
		return
	}

	var req dto.DifficultyRatingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := question.RateDifficulty(userID, uint(questionID), &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "评价成功", nil)
}

// RateQuality 评价质量
// @Summary      评价题目质量
// @Description  对指定题目进行质量评价
// @Tags         题目管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint                       true  "题目ID"
// @Param        req  body      dto.QualityRatingRequest    true  "评价信息"
// @Success      200  {object}  response.Response  "评价成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /questions/{id}/quality-rating [post]
// @Security     BearerAuth
func (ctrl *QuestionController) RateQuality(c *gin.Context) {
	userID := c.GetUint("user_id")
	questionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的题目ID")
		return
	}

	var req dto.QualityRatingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := question.RateQuality(userID, uint(questionID), &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "评价成功", nil)
}

// ==================== 题目纠错 ====================

// CreateCorrection 提交纠错
// @Summary      提交纠错
// @Description  为指定题目提交纠错信息
// @Tags         题目管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint                          true  "题目ID"
// @Param        req  body      dto.CreateCorrectionRequest    true  "纠错信息"
// @Success      200  {object}  response.Response{data=dto.CorrectionInfo}  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /questions/{id}/corrections [post]
// @Security     BearerAuth
func (ctrl *QuestionController) CreateCorrection(c *gin.Context) {
	userID := c.GetUint("user_id")
	questionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的题目ID")
		return
	}

	var req dto.CreateCorrectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	correction, err := question.CreateCorrection(userID, uint(questionID), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, correction)
}

// GetCorrections 获取纠错列表
// @Summary      获取纠错列表
// @Description  获取指定题目的纠错列表
// @Tags         题目管理
// @Accept       json
// @Produce      json
// @Param        id         path      uint  true   "题目ID"
// @Param        status     query     int   false  "状态筛选"
// @Param        page       query     int   false  "页码"
// @Param        page_size  query     int   false  "每页数量"
// @Success      200  {object}  response.Response{data=dto.PageResponse{list=[]dto.CorrectionInfo}}  "成功"
// @Failure      400  {object}  response.Response  "无效的题目ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /questions/{id}/corrections [get]
// @Security     BearerAuth
func (ctrl *QuestionController) GetCorrections(c *gin.Context) {
	questionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的题目ID")
		return
	}

	var req dto.CorrectionListRequest
	c.ShouldBindQuery(&req)

	list, total, err := question.GetCorrections(uint(questionID), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}
