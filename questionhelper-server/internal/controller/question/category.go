package question

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/service/question"
	"questionhelper-server/pkg/response"
)

// GetCategoryTree 分类树
// @Summary      获取分类树
// @Description  获取题目分类的树形结构
// @Tags         分类管理
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response{data=[]dto.CategoryInfo}  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /categories/tree [get]
// @Security     BearerAuth
func (ctrl *QuestionController) GetCategoryTree(c *gin.Context) {
	tree, err := question.GetCategoryTree()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, tree)
}

// CreateCategory 创建分类
// @Summary      创建分类
// @Description  管理员创建新的题目分类
// @Tags         分类管理
// @Accept       json
// @Produce      json
// @Param        req  body      dto.CreateCategoryRequest  true  "分类信息"
// @Success      200  {object}  response.Response  "创建分类成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/categories [post]
// @Security     BearerAuth
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
// @Summary      更新分类
// @Description  管理员更新指定分类信息
// @Tags         分类管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint                        true  "分类ID"
// @Param        req  body      dto.UpdateCategoryRequest    true  "更新信息"
// @Success      200  {object}  response.Response  "更新分类成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/categories/{id} [put]
// @Security     BearerAuth
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
// @Summary      删除分类
// @Description  管理员删除指定分类
// @Tags         分类管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "分类ID"
// @Success      200  {object}  response.Response  "删除分类成功"
// @Failure      400  {object}  response.Response  "无效的分类ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/categories/{id} [delete]
// @Security     BearerAuth
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
