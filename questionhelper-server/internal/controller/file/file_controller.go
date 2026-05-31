package file

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/dto"
	fileService "questionhelper-server/internal/service/file"
	"questionhelper-server/pkg/response"
)

type FileController struct{}

func NewFileController() *FileController {
	return &FileController{}
}

// ==================== Upload ====================

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

// UploadImage 上传图片（压缩+缩略图）
func (ctrl *FileController) UploadImage(c *gin.Context) {
	userID := c.GetUint("user_id")

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "请选择图片文件")
		return
	}
	defer file.Close()

	result, err := fileService.UploadImage(userID, header.Filename, header.Size, header.Header.Get("Content-Type"), file)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// BatchUpload 批量文件上传
func (ctrl *FileController) BatchUpload(c *gin.Context) {
	userID := c.GetUint("user_id")

	form, err := c.MultipartForm()
	if err != nil {
		response.Error(c, http.StatusBadRequest, "请上传文件: "+err.Error())
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		response.Error(c, http.StatusBadRequest, "请选择文件")
		return
	}

	result := fileService.BatchUpload(userID, files)
	response.Success(c, result)
}

// ==================== Download ====================

// DownloadFile 文件下载
func (ctrl *FileController) DownloadFile(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的文件ID")
		return
	}

	file, err := fileService.GetFileForDownload(uint(id), userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}

	// 记录访问日志
	ip := c.ClientIP()
	ua := c.GetHeader("User-Agent")
	go fileService.LogFileAccess(file.ID, userID, "download", ip, ua)

	c.File(file.Path)
}

// ==================== Thumbnail ====================

// GetThumbnail 获取缩略图
func (ctrl *FileController) GetThumbnail(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的文件ID")
		return
	}
	size := c.Param("size")

	// 记录访问日志
	go fileService.LogFileAccess(uint(id), userID, "view", c.ClientIP(), c.GetHeader("User-Agent"))

	filePath, contentType, err := fileService.GetThumbnail(uint(id), size)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}

	c.Header("Content-Type", contentType)
	c.Header("Cache-Control", "public, max-age=86400")
	c.File(filePath)
}

// ==================== File Reference ====================

// AddReference 添加文件引用
func (ctrl *FileController) AddReference(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的文件ID")
		return
	}

	var req dto.AddReferenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := fileService.AddReference(uint(id), userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "添加引用成功", nil)
}

// GetReferences 获取文件引用列表
func (ctrl *FileController) GetReferences(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的文件ID")
		return
	}

	refs, err := fileService.GetReferences(uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, refs)
}

// DeleteReference 删除文件引用
func (ctrl *FileController) DeleteReference(c *gin.Context) {
	userID := c.GetUint("user_id")
	fileID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的文件ID")
		return
	}
	refID, err := strconv.ParseUint(c.Param("refId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的引用ID")
		return
	}

	if err := fileService.DeleteReference(uint(fileID), uint(refID), userID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "删除引用成功", nil)
}

// ==================== Access Log ====================

// GetAccessLogs 获取文件访问日志
func (ctrl *FileController) GetAccessLogs(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的文件ID")
		return
	}

	var req dto.PageRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	logs, total, err := fileService.GetAccessLogs(uint(id), req.Page, req.PageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, logs, total, req.Page, req.PageSize)
}

// ==================== Legacy ====================

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
