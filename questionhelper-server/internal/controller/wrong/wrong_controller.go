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
func (ctrl *WrongController) GetWrongAccuracyAnalysis(c *gin.Context) {
	userID := c.GetUint("user_id")

	data, err := wrongService.GetWrongAccuracyAnalysis(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, data)
}

