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
// @Summary      上传文件
// @Description  上传单个文件
// @Tags         文件管理
// @Accept       multipart/form-data
// @Produce      json
// @Param        file  formData  file  true  "文件"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "请选择文件"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /files [post]
// @Security     BearerAuth
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
// @Summary      上传图片
// @Description  上传图片文件，自动压缩和生成缩略图
// @Tags         文件管理
// @Accept       multipart/form-data
// @Produce      json
// @Param        file  formData  file  true  "图片文件"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "请选择图片文件"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /file/upload/image [post]
// @Security     BearerAuth
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
// @Summary      批量上传文件
// @Description  批量上传多个文件
// @Tags         文件管理
// @Accept       multipart/form-data
// @Produce      json
// @Param        files  formData  file  true  "文件列表"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "请选择文件"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /file/upload/batch [post]
// @Security     BearerAuth
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
// @Summary      下载文件
// @Description  根据文件ID下载文件
// @Tags         文件管理
// @Accept       json
// @Produce      application/octet-stream
// @Param        id  path      uint  true  "文件ID"
// @Success      200  {string}  string  "文件内容"
// @Failure      400  {object}  response.Response  "无效的文件ID"
// @Failure      404  {object}  response.Response  "文件未找到"
// @Router       /file/{id}/download [get]
// @Security     BearerAuth
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
// @Summary      获取缩略图
// @Description  根据文件ID和尺寸获取缩略图
// @Tags         文件管理
// @Accept       json
// @Produce      image/*
// @Param        id    path      uint    true  "文件ID"
// @Param        size  path      string  true  "缩略图尺寸"
// @Success      200  {string}  string  "缩略图内容"
// @Failure      400  {object}  response.Response  "无效的文件ID"
// @Failure      404  {object}  response.Response  "缩略图未找到"
// @Router       /file/{id}/thumbnail/{size} [get]
// @Security     BearerAuth
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
// @Summary      添加文件引用
// @Description  为文件添加引用关联
// @Tags         文件管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint                      true  "文件ID"
// @Param        req  body      dto.AddReferenceRequest   true  "引用数据"
// @Success      200  {object}  response.Response  "添加引用成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /file/{id}/reference [post]
// @Security     BearerAuth
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
// @Summary      获取文件引用列表
// @Description  获取指定文件的引用列表
// @Tags         文件管理
// @Accept       json
// @Produce      json
// @Param        id  path      uint  true  "文件ID"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "无效的文件ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /file/{id}/references [get]
// @Security     BearerAuth
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
// @Summary      删除文件引用
// @Description  删除指定文件的引用关联
// @Tags         文件管理
// @Accept       json
// @Produce      json
// @Param        id     path      uint  true  "文件ID"
// @Param        refId  path      uint  true  "引用ID"
// @Success      200  {object}  response.Response  "删除引用成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /file/{id}/reference/{refId} [delete]
// @Security     BearerAuth
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
// @Summary      获取文件访问日志
// @Description  获取指定文件的访问日志，支持分页
// @Tags         文件管理
// @Accept       json
// @Produce      json
// @Param        id         path      uint  true   "文件ID"
// @Param        page       query     int   false  "页码"
// @Param        page_size  query     int   false  "每页数量"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /file/{id}/access-logs [get]
// @Security     BearerAuth
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

// ==================== Query ====================

// GetFile 获取文件信息
// @Summary      获取文件信息
// @Description  根据文件ID获取文件详细信息
// @Tags         文件管理
// @Accept       json
// @Produce      json
// @Param        id  path      uint  true  "文件ID"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "无效的文件ID"
// @Failure      404  {object}  response.Response  "文件不存在"
// @Router       /file/{id} [get]
// @Security     BearerAuth
func (ctrl *FileController) GetFile(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的文件ID")
		return
	}
	file, err := fileService.GetFile(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "文件不存在")
		return
	}
	response.Success(c, file)
}

// ListFiles 文件列表
// @Summary      获取文件列表
// @Description  获取文件列表，支持按类型、关键词过滤和排序
// @Tags         文件管理
// @Accept       json
// @Produce      json
// @Param        page        query     int     false  "页码"           default(1)
// @Param        page_size   query     int     false  "每页数量"       default(20)
// @Param        file_type   query     string  false  "文件类型(扩展名)"
// @Param        keyword     query     string  false  "文件名关键词"
// @Param        sort_by     query     string  false  "排序字段"       Enums(created_at, updated_at, size, name)  default(created_at)
// @Param        sort_order  query     string  false  "排序方式"       Enums(asc, desc)  default(desc)
// @Success      200  {object}  response.PageResponse  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /file/list [get]
// @Security     BearerAuth
func (ctrl *FileController) ListFiles(c *gin.Context) {
	var req dto.FileListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	files, total, err := fileService.ListFilesWithFilter(req.Page, req.PageSize, req.FileType, req.Keyword, req.SortBy, req.SortOrder)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.Page(c, files, total, req.Page, req.PageSize)
}

// ==================== Legacy ====================

// DeleteFile 删除文件
// @Summary      删除文件
// @Description  根据文件ID删除文件
// @Tags         文件管理
// @Accept       json
// @Produce      json
// @Param        id  path      uint  true  "文件ID"
// @Success      200  {object}  response.Response  "删除成功"
// @Failure      400  {object}  response.Response  "无效的文件ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /files/{id} [delete]
// @Security     BearerAuth
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
