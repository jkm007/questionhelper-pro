package question

import (
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/service/question"
	"questionhelper-server/pkg/response"
)

// ImportQuestions 导入题目
// @Summary      导入题目
// @Description  管理员通过文件批量导入题目
// @Tags         题目管理
// @Accept       multipart/form-data
// @Produce      json
// @Param        file         formData  file    true   "题目文件"
// @Param        category_id  formData  uint    true   "分类ID"
// @Param        visibility   formData  int     true   "可见性(1=公开 2=私有 3=班级)"
// @Success      200  {object}  response.Response{data=object{count=int}}  "导入成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/questions/import [post]
// @Security     BearerAuth
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
