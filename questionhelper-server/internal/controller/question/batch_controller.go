package question

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/service/question"
	"questionhelper-server/pkg/response"
)

type BatchController struct{}

func NewBatchController() *BatchController {
	return &BatchController{}
}

// BatchPublish 批量发布
// @Summary      批量发布题目
// @Description  管理员批量发布多个题目
// @Tags         题目管理
// @Accept       json
// @Produce      json
// @Param        req  body      dto.BatchQuestionRequest  true  "题目ID列表"
// @Success      200  {object}  response.Response{data=dto.BatchResult}  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/questions/batch/publish [post]
// @Security     BearerAuth
func (ctrl *BatchController) BatchPublish(c *gin.Context) {
	var req dto.BatchQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	result, err := question.BatchPublish(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, result)
}

// BatchArchive 批量归档
// @Summary      批量归档题目
// @Description  管理员批量归档多个题目
// @Tags         题目管理
// @Accept       json
// @Produce      json
// @Param        req  body      dto.BatchQuestionRequest  true  "题目ID列表"
// @Success      200  {object}  response.Response{data=dto.BatchResult}  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/questions/batch/archive [post]
// @Security     BearerAuth
func (ctrl *BatchController) BatchArchive(c *gin.Context) {
	var req dto.BatchQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	result, err := question.BatchArchive(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, result)
}

// BatchDelete 批量删除
// @Summary      批量删除题目
// @Description  管理员批量删除多个题目
// @Tags         题目管理
// @Accept       json
// @Produce      json
// @Param        req  body      dto.BatchQuestionRequest  true  "题目ID列表"
// @Success      200  {object}  response.Response{data=dto.BatchResult}  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/questions/batch/delete [post]
// @Security     BearerAuth
func (ctrl *BatchController) BatchDelete(c *gin.Context) {
	var req dto.BatchQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	result, err := question.BatchDelete(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, result)
}

// BatchMoveCategory 批量移动分类
// @Summary      批量移动题目分类
// @Description  管理员批量将题目移动到指定分类
// @Tags         题目管理
// @Accept       json
// @Produce      json
// @Param        req  body      dto.BatchMoveRequest  true  "题目ID列表和目标分类"
// @Success      200  {object}  response.Response{data=dto.BatchResult}  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/questions/batch/move [post]
// @Security     BearerAuth
func (ctrl *BatchController) BatchMoveCategory(c *gin.Context) {
	var req dto.BatchMoveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	result, err := question.BatchMoveCategory(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, result)
}
