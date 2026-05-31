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

// ListQuestions 题目列表
func (ctrl *QuestionController) ListQuestions(c *gin.Context) {
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
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的题目ID")
		return
	}

	if err := question.LikeQuestion(uint(id)); err != nil {
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
