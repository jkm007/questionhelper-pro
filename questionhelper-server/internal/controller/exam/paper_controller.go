package exam

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/service/exam"
	"questionhelper-server/pkg/response"
)

type PaperController struct{}

func NewPaperController() *PaperController {
	return &PaperController{}
}

// PreviewPaper 试卷预览
func (ctrl *PaperController) PreviewPaper(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的试卷ID")
		return
	}

	result, err := exam.PreviewPaper(uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, result)
}

// CopyPaper 复制试卷
func (ctrl *PaperController) CopyPaper(c *gin.Context) {
	userID := c.GetUint("user_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的试卷ID")
		return
	}

	var req dto.CopyPaperRequest
	c.ShouldBindJSON(&req)

	result, err := exam.CopyPaper(uint(id), userID, req.Title)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, result)
}

// PublishPaper 发布试卷
func (ctrl *PaperController) PublishPaper(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的试卷ID")
		return
	}

	var req dto.PublishPaperRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := exam.PublishPaper(uint(id), req.Status); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "操作成功", nil)
}

// SaveAsTemplate 保存为模板
func (ctrl *PaperController) SaveAsTemplate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的试卷ID")
		return
	}

	var req dto.SaveTemplateRequest
	c.ShouldBindJSON(&req)

	if err := exam.SaveAsTemplate(uint(id), req.Category, req.Tags); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "保存模板成功", nil)
}

// CreateFromTemplate 从模板创建
func (ctrl *PaperController) CreateFromTemplate(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		TemplateID uint   `json:"template_id" binding:"required"`
		Title      string `json:"title"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	result, err := exam.CreateFromTemplate(req.TemplateID, userID, req.Title)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, result)
}

// ListTemplates 模板列表
func (ctrl *PaperController) ListTemplates(c *gin.Context) {
	var req dto.TemplateListRequest
	c.ShouldBindQuery(&req)

	list, total, err := exam.ListTemplates(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// ExportPaper 导出试卷
func (ctrl *PaperController) ExportPaper(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的试卷ID")
		return
	}

	result, err := exam.ExportPaper(uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 设置下载头
	c.Header("Content-Type", "application/json")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=paper_%d.json", id))
	c.JSON(http.StatusOK, result)
}

// GetPaperStats 试卷统计
func (ctrl *PaperController) GetPaperStats(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的试卷ID")
		return
	}

	result, err := exam.GetPaperStats(uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, result)
}
