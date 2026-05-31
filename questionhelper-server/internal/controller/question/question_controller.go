package question

import (
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/service/question"
	"questionhelper-server/pkg/response"
)

type QuestionController struct{}

func NewQuestionController() *QuestionController {
	return &QuestionController{}
}

// ListQuestions 题目列表（带数据权限过滤）
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
func (ctrl *QuestionController) ListCategories(c *gin.Context) {
	list, err := question.ListCategories()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, list)
}

// GetCategoryTree 分类树
func (ctrl *QuestionController) GetCategoryTree(c *gin.Context) {
	tree, err := question.GetCategoryTree()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, tree)
}

// ListKnowledgePoints 知识点列表
func (ctrl *QuestionController) ListKnowledgePoints(c *gin.Context) {
	var categoryID *uint
	if cid := c.Query("category_id"); cid != "" {
		id, err := strconv.ParseUint(cid, 10, 32)
		if err == nil {
			uid := uint(id)
			categoryID = &uid
		}
	}

	list, err := question.ListKnowledgePoints(categoryID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, list)
}

// AdminListQuestions 管理员题目列表
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

// ImportQuestions 导入题目
func (ctrl *QuestionController) ImportQuestions(c *gin.Context) {
	creatorID := c.GetUint("user_id")

	// 获取分类ID
	categoryIDStr := c.PostForm("category_id")
	categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的分类ID")
		return
	}

	// 获取可见性
	visibilityStr := c.PostForm("visibility")
	visibility, err := strconv.ParseInt(visibilityStr, 10, 8)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的可见性")
		return
	}

	// 获取上传的文件
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

	count, err := question.ImportQuestions(creatorID, uint(categoryID), int8(visibility), data)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithMessage(c, "导入成功", gin.H{"count": count})
}

// CreateFavoriteFolder 创建收藏夹
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
func (ctrl *QuestionController) ListFavoriteFolders(c *gin.Context) {
	userID := c.GetUint("user_id")

	list, err := question.ListFolders(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, list)
}

// ExportQuestions 导出题目
func (ctrl *QuestionController) ExportQuestions(c *gin.Context) {
	var categoryID *uint
	if cid := c.Query("category_id"); cid != "" {
		id, err := strconv.ParseUint(cid, 10, 32)
		if err == nil {
			uid := uint(id)
			categoryID = &uid
		}
	}

	list, err := question.ExportQuestions(categoryID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Header("Content-Type", "application/json")
	c.Header("Content-Disposition", "attachment; filename=questions.json")
	c.JSON(http.StatusOK, list)
}

// ==================== 分类管理 ====================

// CreateCategory 创建分类
func (ctrl *QuestionController) CreateCategory(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := question.CreateCategory(&req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "创建分类成功", nil)
}

// UpdateCategory 更新分类
func (ctrl *QuestionController) UpdateCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的分类ID")
		return
	}

	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := question.UpdateCategory(uint(id), &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "更新分类成功", nil)
}

// DeleteCategory 删除分类
func (ctrl *QuestionController) DeleteCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的分类ID")
		return
	}

	if err := question.DeleteCategory(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "删除分类成功", nil)
}

// ==================== 知识点管理 ====================

// CreateKnowledgePoint 创建知识点
func (ctrl *QuestionController) CreateKnowledgePoint(c *gin.Context) {
	var req dto.CreateKnowledgeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := question.CreateKnowledgePoint(&req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "创建知识点成功", nil)
}

// UpdateKnowledgePoint 更新知识点
func (ctrl *QuestionController) UpdateKnowledgePoint(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的知识点ID")
		return
	}

	var req dto.UpdateKnowledgeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := question.UpdateKnowledgePoint(uint(id), &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "更新知识点成功", nil)
}

// DeleteKnowledgePoint 删除知识点
func (ctrl *QuestionController) DeleteKnowledgePoint(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的知识点ID")
		return
	}

	if err := question.DeleteKnowledgePoint(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "删除知识点成功", nil)
}

// ==================== 内容审核 ====================

// ListPendingReviews 待审核列表
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
func (ctrl *QuestionController) ListSensitiveWords(c *gin.Context) {
	var req dto.PageRequest
	c.ShouldBindQuery(&req)

	list, total, err := question.ListSensitiveWords(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// CreateSensitiveWord 创建敏感词
func (ctrl *QuestionController) CreateSensitiveWord(c *gin.Context) {
	var req dto.CreateSensitiveWordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := question.CreateSensitiveWord(&req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "创建敏感词成功", nil)
}

// DeleteSensitiveWord 删除敏感词
func (ctrl *QuestionController) DeleteSensitiveWord(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的敏感词ID")
		return
	}

	if err := question.DeleteSensitiveWord(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "删除敏感词成功", nil)
}

// ImportSensitiveWords 导入敏感词
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

	count, err := question.ImportSensitiveWords(data)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithMessage(c, "导入成功", gin.H{"count": count})
}

// TestSensitiveWord 测试敏感词
func (ctrl *QuestionController) TestSensitiveWord(c *gin.Context) {
	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	result := question.TestSensitiveWord(req.Content)
	response.Success(c, gin.H{"has_sensitive": result})
}

// ==================== 题目笔记 ====================

// GetQuestionNotes 获取题目笔记
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
