package question

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/service/question"
	"questionhelper-server/pkg/response"
)

// ExportQuestions 导出题目
// @Summary      导出题目
// @Description  管理员导出题目为JSON文件，支持按分类筛选
// @Tags         题目管理
// @Accept       json
// @Produce      application/json
// @Param        category_id  query     uint  false  "分类ID"
// @Success      200  {array}   dto.QuestionInfo  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/questions/export [get]
// @Security     BearerAuth
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
