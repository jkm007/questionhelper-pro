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
