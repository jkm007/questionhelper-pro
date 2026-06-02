package wrong

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/dto"
	fileService "questionhelper-server/internal/service/file"
	wrongService "questionhelper-server/internal/service/wrong"
	"questionhelper-server/pkg/response"
)

type WrongController struct{}

func NewWrongController() *WrongController {
	return &WrongController{}
}

// ListWrongQuestions 错题列表
// @Summary      获取错题列表
// @Description  分页获取当前用户的错题列表
// @Tags         错题本
// @Accept       json
// @Produce      json
// @Param        page       query     int  false  "页码"
// @Param        page_size  query     int  false  "每页数量"
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /wrong [get]
// @Security     BearerAuth
func (ctrl *WrongController) ListWrongQuestions(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.WrongListRequest
	c.ShouldBindQuery(&req)

	list, total, err := wrongService.ListWrongQuestions(userID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// GetWrongQuestion 获取错题详情
// @Summary      获取错题详情
// @Description  根据ID获取错题详细信息
// @Tags         错题本
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "错题ID"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "无效的错题ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /wrong/{id} [get]
// @Security     BearerAuth
func (ctrl *WrongController) GetWrongQuestion(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的错题ID")
		return
	}

	info, err := wrongService.GetWrongQuestion(uint(id), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, info)
}

// ReviewWrongQuestion 复习错题
// @Summary      复习错题
// @Description  提交错题复习答案并判断对错
// @Tags         错题本
// @Accept       json
// @Produce      json
// @Param        id   path      uint                        true  "错题ID"
// @Param        req  body      dto.ReviewWrongRequest       true  "复习答案"
// @Success      200  {object}  response.Response  "复习结果"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /wrong/{id}/review [post]
// @Security     BearerAuth
func (ctrl *WrongController) ReviewWrongQuestion(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的错题ID")
		return
	}

	var req dto.ReviewWrongRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	isCorrect, err := wrongService.ReviewWrongQuestion(uint(id), userID, req.Answer)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	if isCorrect {
		response.SuccessWithMessage(c, "回答正确！", gin.H{"is_correct": true})
	} else {
		response.SuccessWithMessage(c, "回答错误", gin.H{"is_correct": false})
	}
}

// RemoveWrongQuestion 移除错题
// @Summary      移除错题
// @Description  根据ID从错题本中移除错题
// @Tags         错题本
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "错题ID"
// @Success      200  {object}  response.Response  "移除成功"
// @Failure      400  {object}  response.Response  "无效的错题ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /wrong/{id} [delete]
// @Security     BearerAuth
func (ctrl *WrongController) RemoveWrongQuestion(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的错题ID")
		return
	}

	if err := wrongService.RemoveWrongQuestion(uint(id), userID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "移除成功", nil)
}

// GetWrongAnalysis 错题分析
// @Summary      获取错题分析
// @Description  获取当前用户的错题统计分析数据
// @Tags         错题本
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /wrong/analysis [get]
// @Security     BearerAuth
func (ctrl *WrongController) GetWrongAnalysis(c *gin.Context) {
	userID := c.GetUint("user_id")

	analysis, err := wrongService.GetWrongAnalysis(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, analysis)
}

// ==================== 错题搜索 ====================

// SearchWrongQuestions 搜索错题
// @Summary      搜索错题
// @Description  根据关键词搜索错题
// @Tags         错题本
// @Accept       json
// @Produce      json
// @Param        req  query     dto.WrongSearchRequest  true  "搜索参数"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /wrong/search [get]
// @Security     BearerAuth
func (ctrl *WrongController) SearchWrongQuestions(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.WrongSearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	list, total, err := wrongService.SearchWrongQuestions(userID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// ==================== 批量操作 ====================

// BatchDeleteWrongQuestions 批量删除错题
// @Summary      批量删除错题
// @Description  根据ID列表批量删除错题
// @Tags         错题本
// @Accept       json
// @Produce      json
// @Param        req  body      dto.WrongBatchDeleteRequest  true  "删除的错题ID列表"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /wrong/batch/delete [post]
// @Security     BearerAuth
func (ctrl *WrongController) BatchDeleteWrongQuestions(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.WrongBatchDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	result, err := wrongService.BatchDeleteWrongQuestions(req.IDs, userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, result)
}

// BatchExportWrongQuestions 批量导出错题
// @Summary      批量导出错题
// @Description  根据ID列表批量导出错题为CSV文件
// @Tags         错题本
// @Accept       json
// @Produce      text/csv
// @Param        req  body      dto.WrongBatchDeleteRequest  true  "导出的错题ID列表"
// @Success      200  {string}  string  "CSV文件内容"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /wrong/batch/export [post]
// @Security     BearerAuth
func (ctrl *WrongController) BatchExportWrongQuestions(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.WrongBatchDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	buf, err := wrongService.BatchExportWrongQuestions(req.IDs, userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=wrong_questions.csv")
	c.Data(http.StatusOK, "text/csv; charset=utf-8", buf.Bytes())
}

// BatchTagWrongQuestions 批量打标签
// @Summary      批量打标签
// @Description  为多个错题批量添加标签
// @Tags         错题本
// @Accept       json
// @Produce      json
// @Param        req  body      object{ids=[]uint,tag_ids=[]uint}  true  "错题ID和标签ID列表"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /wrong/batch/tag [post]
// @Security     BearerAuth
func (ctrl *WrongController) BatchTagWrongQuestions(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		IDs   []uint `json:"ids" binding:"required,min=1"`
		TagIDs []uint `json:"tag_ids" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	result, err := wrongService.BatchTagWrongQuestions(req.IDs, req.TagIDs, userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, result)
}

// ==================== 错题标签管理 ====================

// ListWrongTags 获取用户标签列表
// @Summary      获取用户标签列表
// @Description  获取当前用户的所有错题标签
// @Tags         错题标签
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /wrong/tags [get]
// @Security     BearerAuth
func (ctrl *WrongController) ListWrongTags(c *gin.Context) {
	userID := c.GetUint("user_id")

	list, err := wrongService.ListWrongTags(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, list)
}

// CreateWrongTag 创建标签
// @Summary      创建错题标签
// @Description  创建一个新的错题标签
// @Tags         错题标签
// @Accept       json
// @Produce      json
// @Param        req  body      dto.CreateWrongTagRequest  true  "标签信息"
// @Success      200  {object}  response.Response  "创建成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /wrong/tags [post]
// @Security     BearerAuth
func (ctrl *WrongController) CreateWrongTag(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.CreateWrongTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := wrongService.CreateWrongTag(userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "创建成功", nil)
}

// UpdateWrongTag 更新标签
// @Summary      更新错题标签
// @Description  根据ID更新错题标签信息
// @Tags         错题标签
// @Accept       json
// @Produce      json
// @Param        id   path      uint                        true  "标签ID"
// @Param        req  body      dto.UpdateWrongTagRequest    true  "更新信息"
// @Success      200  {object}  response.Response  "更新成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /wrong/tags/{id} [put]
// @Security     BearerAuth
func (ctrl *WrongController) UpdateWrongTag(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的标签ID")
		return
	}

	var req dto.UpdateWrongTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := wrongService.UpdateWrongTag(uint(id), userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "更新成功", nil)
}

// DeleteWrongTag 删除标签
// @Summary      删除错题标签
// @Description  根据ID删除错题标签
// @Tags         错题标签
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "标签ID"
// @Success      200  {object}  response.Response  "删除成功"
// @Failure      400  {object}  response.Response  "无效的标签ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /wrong/tags/{id} [delete]
// @Security     BearerAuth
func (ctrl *WrongController) DeleteWrongTag(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的标签ID")
		return
	}

	if err := wrongService.DeleteWrongTag(uint(id), userID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "删除成功", nil)
}

// AddTagToWrongQuestion 为错题添加标签
// @Summary      为错题添加标签
// @Description  为指定错题添加标签
// @Tags         错题标签
// @Accept       json
// @Produce      json
// @Param        id   path      uint                       true  "错题ID"
// @Param        req  body      dto.WrongTagIDsRequest      true  "标签ID列表"
// @Success      200  {object}  response.Response  "添加标签成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /wrong/{id}/tags [post]
// @Security     BearerAuth
func (ctrl *WrongController) AddTagToWrongQuestion(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的错题ID")
		return
	}

	var req dto.WrongTagIDsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := wrongService.AddTagToWrongQuestion(uint(id), userID, req.TagIDs); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "添加标签成功", nil)
}

// RemoveTagFromWrongQuestion 移除错题标签
// @Summary      移除错题标签
// @Description  从指定错题中移除标签
// @Tags         错题标签
// @Accept       json
// @Produce      json
// @Param        id     path      uint  true  "错题ID"
// @Param        tagId  path      uint  true  "标签ID"
// @Success      200  {object}  response.Response  "移除标签成功"
// @Failure      400  {object}  response.Response  "无效的ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /wrong/{id}/tags/{tagId} [delete]
// @Security     BearerAuth
func (ctrl *WrongController) RemoveTagFromWrongQuestion(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的错题ID")
		return
	}

	tagID, err := strconv.ParseUint(c.Param("tagId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的标签ID")
		return
	}

	if err := wrongService.RemoveTagFromWrongQuestion(uint(id), uint(tagID), userID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "移除标签成功", nil)
}

// ==================== 错题备注 ====================

// ListWrongNotes 获取错题备注
// @Summary      获取错题备注列表
// @Description  获取指定错题的所有备注
// @Tags         错题本
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "错题ID"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "无效的错题ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /wrong/{id}/notes [get]
// @Security     BearerAuth
func (ctrl *WrongController) ListWrongNotes(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的错题ID")
		return
	}

	list, err := wrongService.ListWrongNotes(uint(id), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, list)
}

// CreateWrongNote 创建备注
// @Summary      创建错题备注
// @Description  为指定错题创建备注
// @Tags         错题本
// @Accept       json
// @Produce      json
// @Param        id   path      uint                         true  "错题ID"
// @Param        req  body      dto.CreateWrongNoteRequest    true  "备注信息"
// @Success      200  {object}  response.Response  "创建成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /wrong/{id}/notes [post]
// @Security     BearerAuth
func (ctrl *WrongController) CreateWrongNote(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的错题ID")
		return
	}

	var req dto.CreateWrongNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := wrongService.CreateWrongNote(uint(id), userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "创建成功", nil)
}

// UpdateWrongNote 更新备注
// @Summary      更新错题备注
// @Description  根据ID更新错题备注
// @Tags         错题本
// @Accept       json
// @Produce      json
// @Param        id      path      uint                         true  "错题ID"
// @Param        noteId  path      uint                         true  "备注ID"
// @Param        req     body      dto.UpdateWrongNoteRequest    true  "更新信息"
// @Success      200  {object}  response.Response  "更新成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /wrong/{id}/notes/{noteId} [put]
// @Security     BearerAuth
func (ctrl *WrongController) UpdateWrongNote(c *gin.Context) {
	userID := c.GetUint("user_id")
	wrongID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的错题ID")
		return
	}

	noteID, err := strconv.ParseUint(c.Param("noteId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的备注ID")
		return
	}

	var req dto.UpdateWrongNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := wrongService.UpdateWrongNote(uint(wrongID), uint(noteID), userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "更新成功", nil)
}

// DeleteWrongNote 删除备注
// @Summary      删除错题备注
// @Description  根据ID删除错题备注
// @Tags         错题本
// @Accept       json
// @Produce      json
// @Param        id      path      uint  true  "错题ID"
// @Param        noteId  path      uint  true  "备注ID"
// @Success      200  {object}  response.Response  "删除成功"
// @Failure      400  {object}  response.Response  "无效的ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /wrong/{id}/notes/{noteId} [delete]
// @Security     BearerAuth
func (ctrl *WrongController) DeleteWrongNote(c *gin.Context) {
	userID := c.GetUint("user_id")
	wrongID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的错题ID")
		return
	}

	noteID, err := strconv.ParseUint(c.Param("noteId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的备注ID")
		return
	}

	if err := wrongService.DeleteWrongNote(uint(wrongID), uint(noteID), userID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "删除成功", nil)
}

// ==================== 错题附件 ====================

// UploadWrongAttachment 上传附件
// @Summary      上传错题附件
// @Description  为指定错题上传附件文件
// @Tags         错题本
// @Accept       multipart/form-data
// @Produce      json
// @Param        id    path      uint   true  "错题ID"
// @Param        file  formData  file   true  "附件文件"
// @Success      200  {object}  response.Response  "上传成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /wrong/{id}/attachments [post]
// @Security     BearerAuth
func (ctrl *WrongController) UploadWrongAttachment(c *gin.Context) {
	userID := c.GetUint("user_id")
	wrongID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的错题ID")
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "请选择文件")
		return
	}
	defer file.Close()

	// 上传文件
	uploadedFile, err := fileService.UploadFile(userID, header.Filename, header.Size, header.Header.Get("Content-Type"), file)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 确定文件类型
	ext := strings.ToLower(filepath.Ext(header.Filename))
	fileType := "document"
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp":
		fileType = "image"
	case ".mp3":
		fileType = "audio"
	case ".mp4":
		fileType = "video"
	}

	// 创建错题附件记录
	if err := wrongService.CreateWrongAttachment(uint(wrongID), userID, header.Filename, fileType, uploadedFile.URL, header.Size, uploadedFile.ID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{
		"id":        uploadedFile.ID,
		"file_name": header.Filename,
		"file_type": fileType,
		"file_url":  uploadedFile.URL,
		"file_size": header.Size,
	})
}

// ListWrongAttachments 获取附件列表
// @Summary      获取错题附件列表
// @Description  获取指定错题的所有附件
// @Tags         错题本
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "错题ID"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "无效的错题ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /wrong/{id}/attachments [get]
// @Security     BearerAuth
func (ctrl *WrongController) ListWrongAttachments(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的错题ID")
		return
	}

	list, err := wrongService.ListWrongAttachments(uint(id), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, list)
}

// DeleteWrongAttachment 删除附件
// @Summary      删除错题附件
// @Description  根据ID删除错题附件
// @Tags         错题本
// @Accept       json
// @Produce      json
// @Param        id           path      uint  true  "错题ID"
// @Param        attachmentId  path      uint  true  "附件ID"
// @Success      200  {object}  response.Response  "删除成功"
// @Failure      400  {object}  response.Response  "无效的ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /wrong/{id}/attachments/{attachmentId} [delete]
// @Security     BearerAuth
func (ctrl *WrongController) DeleteWrongAttachment(c *gin.Context) {
	userID := c.GetUint("user_id")
	wrongID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的错题ID")
		return
	}

	attachmentID, err := strconv.ParseUint(c.Param("attachmentId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的附件ID")
		return
	}

	if err := wrongService.DeleteWrongAttachment(uint(wrongID), uint(attachmentID), userID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "删除成功", nil)
}

// ==================== 错题收藏 ====================

// FavoriteWrongQuestion 收藏错题
// @Summary      收藏错题
// @Description  收藏指定错题
// @Tags         错题本
// @Accept       json
// @Produce      json
// @Param        id   path      uint                        true  "错题ID"
// @Param        req  body      dto.FavoriteWrongRequest     true  "收藏信息"
// @Success      200  {object}  response.Response  "收藏成功"
// @Failure      400  {object}  response.Response  "无效的错题ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /wrong/{id}/favorite [post]
// @Security     BearerAuth
func (ctrl *WrongController) FavoriteWrongQuestion(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的错题ID")
		return
	}

	var req dto.FavoriteWrongRequest
	c.ShouldBindJSON(&req)

	if err := wrongService.FavoriteWrongQuestion(uint(id), userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "收藏成功", nil)
}

// UnfavoriteWrongQuestion 取消收藏
// @Summary      取消收藏错题
// @Description  取消收藏指定错题
// @Tags         错题本
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "错题ID"
// @Success      200  {object}  response.Response  "取消收藏成功"
// @Failure      400  {object}  response.Response  "无效的错题ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /wrong/{id}/favorite [delete]
// @Security     BearerAuth
func (ctrl *WrongController) UnfavoriteWrongQuestion(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的错题ID")
		return
	}

	if err := wrongService.UnfavoriteWrongQuestion(uint(id), userID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "取消收藏成功", nil)
}

// ListWrongFavorites 收藏列表
// @Summary      获取收藏列表
// @Description  分页获取当前用户的收藏错题列表
// @Tags         错题本
// @Accept       json
// @Produce      json
// @Param        page       query     int  false  "页码"
// @Param        page_size  query     int  false  "每页数量"
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /wrong/favorites [get]
// @Security     BearerAuth
func (ctrl *WrongController) ListWrongFavorites(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.WrongFavoriteListRequest
	c.ShouldBindQuery(&req)

	list, total, err := wrongService.ListWrongFavorites(userID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// ==================== 错题复习 ====================

// GetTodayReviewQuestions 今日待复习
// @Summary      获取今日待复习题目
// @Description  获取当前用户今日需要复习的错题列表
// @Tags         错题复习
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /wrong/review/today [get]
// @Security     BearerAuth
func (ctrl *WrongController) GetTodayReviewQuestions(c *gin.Context) {
	userID := c.GetUint("user_id")

	list, err := wrongService.GetTodayReviewQuestions(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, list)
}

// RecordReviewResult 记录复习结果
// @Summary      记录复习结果
// @Description  记录错题复习结果（进阶复习模式）
// @Tags         错题复习
// @Accept       json
// @Produce      json
// @Param        id   path      uint                              true  "错题ID"
// @Param        req  body      dto.ReviewWrongAdvancedRequest     true  "复习结果"
// @Success      200  {object}  response.Response  "复习记录已保存"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /wrong/{id}/review/record [post]
// @Security     BearerAuth
func (ctrl *WrongController) RecordReviewResult(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的错题ID")
		return
	}

	var req dto.ReviewWrongAdvancedRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := wrongService.ReviewWrongQuestionAdvanced(uint(id), userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "复习记录已保存", nil)
}

// GetReviewHistory 复习历史
// @Summary      获取复习历史
// @Description  分页获取当前用户的错题复习历史记录
// @Tags         错题复习
// @Accept       json
// @Produce      json
// @Param        page       query     int  false  "页码"
// @Param        page_size  query     int  false  "每页数量"
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /wrong/review/history [get]
// @Security     BearerAuth
func (ctrl *WrongController) GetReviewHistory(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.WrongReviewListRequest
	c.ShouldBindQuery(&req)

	list, total, err := wrongService.GetReviewHistory(userID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// ==================== 错题导出 ====================

// ExportWrongQuestions 导出错题
// @Summary      导出错题
// @Description  根据条件导出错题为CSV文件
// @Tags         错题本
// @Accept       json
// @Produce      text/csv
// @Param        req  body      dto.WrongExportRequest  true  "导出条件"
// @Success      200  {string}  string  "CSV文件内容"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /wrong/export [post]
// @Security     BearerAuth
func (ctrl *WrongController) ExportWrongQuestions(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.WrongExportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	buf, err := wrongService.ExportWrongQuestions(userID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	filename := fmt.Sprintf("wrong_questions_%s.csv", time.Now().Format("20060102150405"))
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Data(http.StatusOK, "text/csv; charset=utf-8", buf.Bytes())
}

// ==================== 错题分析扩展 ====================

// GetWrongTrendAnalysis 错题趋势分析
// @Summary      获取错题趋势分析
// @Description  获取错题趋势分析数据
// @Tags         错题本
// @Accept       json
// @Produce      json
// @Param        req  query     dto.WrongTrendRequest  true  "查询参数"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /wrong/analysis/trend [get]
// @Security     BearerAuth
func (ctrl *WrongController) GetWrongTrendAnalysis(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.WrongTrendRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	data, err := wrongService.GetWrongTrendAnalysis(userID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, data)
}

// GetWrongCategoryAnalysis 按分类统计
// @Summary      获取错题分类统计
// @Description  按分类统计当前用户的错题分布
// @Tags         错题本
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /wrong/analysis/category [get]
// @Security     BearerAuth
func (ctrl *WrongController) GetWrongCategoryAnalysis(c *gin.Context) {
	userID := c.GetUint("user_id")

	data, err := wrongService.GetWrongCategoryAnalysis(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, data)
}

// GetWrongAccuracyAnalysis 正确率分析
// @Summary      获取正确率分析
// @Description  获取当前用户的错题正确率分析数据
// @Tags         错题本
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /wrong/analysis/accuracy [get]
// @Security     BearerAuth
func (ctrl *WrongController) GetWrongAccuracyAnalysis(c *gin.Context) {
	userID := c.GetUint("user_id")

	data, err := wrongService.GetWrongAccuracyAnalysis(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, data)
}

