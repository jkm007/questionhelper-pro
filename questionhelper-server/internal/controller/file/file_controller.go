package file

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	fileService "questionhelper-server/internal/service/file"
	"questionhelper-server/pkg/response"
)

type FileController struct{}

func NewFileController() *FileController {
	return &FileController{}
}

// UploadFile 上传文件
func (ctrl *FileController) UploadFile(c *gin.Context) {
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

	response.Success(c, gin.H{
		"id":       result.ID,
		"name":     result.Name,
		"original": result.Original,
		"url":      result.URL,
		"size":     result.Size,
	})
}

// DeleteFile 删除文件
func (ctrl *FileController) DeleteFile(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的文件ID")
		return
	}

	if err := fileService.DeleteFile(uint(id), userID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "删除成功", nil)
}
