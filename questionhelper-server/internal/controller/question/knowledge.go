package question

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/service/question"
	"questionhelper-server/pkg/response"
)

// ListKnowledgePoints 知识点列表
// @Summary      获取知识点列表
// @Description  获取知识点列表，支持按分类筛选
// @Tags         知识点管理
// @Accept       json
// @Produce      json
// @Param        category_id  query     uint  false  "分类ID"
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/knowledge-points [get]
// @Security     BearerAuth
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

// CreateKnowledgePoint 创建知识点
// @Summary      创建知识点
// @Description  管理员创建新的知识点
// @Tags         知识点管理
// @Accept       json
// @Produce      json
// @Param        req  body      dto.CreateKnowledgeRequest  true  "知识点信息"
// @Success      200  {object}  response.Response  "创建知识点成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/knowledge-points [post]
// @Security     BearerAuth
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
// @Summary      更新知识点
// @Description  管理员更新指定知识点信息
// @Tags         知识点管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint                         true  "知识点ID"
// @Param        req  body      dto.UpdateKnowledgeRequest    true  "更新信息"
// @Success      200  {object}  response.Response  "更新知识点成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/knowledge-points/{id} [put]
// @Security     BearerAuth
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
// @Summary      删除知识点
// @Description  管理员删除指定知识点
// @Tags         知识点管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "知识点ID"
// @Success      200  {object}  response.Response  "删除知识点成功"
// @Failure      400  {object}  response.Response  "无效的知识点ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/knowledge-points/{id} [delete]
// @Security     BearerAuth
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
