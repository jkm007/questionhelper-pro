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
// @Summary      试卷预览
// @Description  预览指定试卷的内容和结构
// @Tags         试卷管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "试卷ID"
// @Success      200  {object}  response.Response{data=dto.PaperPreviewInfo}  "成功"
// @Failure      400  {object}  response.Response  "无效的试卷ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/papers/{id}/preview [get]
// @Security     BearerAuth
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
// @Summary      复制试卷
// @Description  复制指定试卷为新试卷
// @Tags         试卷管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint                true  "试卷ID"
// @Param        req  body      dto.CopyPaperRequest  true  "复制信息"
// @Success      200  {object}  response.Response{data=dto.PaperInfo}  "成功"
// @Failure      400  {object}  response.Response  "无效的试卷ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/papers/{id}/copy [post]
// @Security     BearerAuth
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
// @Summary      发布试卷
// @Description  发布或下架指定试卷
// @Tags         试卷管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint                      true  "试卷ID"
// @Param        req  body      dto.PublishPaperRequest    true  "发布信息"
// @Success      200  {object}  response.Response  "操作成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/papers/{id}/publish [put]
// @Security     BearerAuth
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
// @Summary      保存为模板
// @Description  将指定试卷保存为模板
// @Tags         试卷管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint                      true  "试卷ID"
// @Param        req  body      dto.SaveTemplateRequest    true  "模板信息"
// @Success      200  {object}  response.Response  "保存模板成功"
// @Failure      400  {object}  response.Response  "无效的试卷ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/papers/{id}/save-template [post]
// @Security     BearerAuth
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
// @Summary      从模板创建试卷
// @Description  根据模板创建新试卷
// @Tags         试卷管理
// @Accept       json
// @Produce      json
// @Param        req  body      object{template_id=uint,title=string}  true  "模板创建信息"
// @Success      200  {object}  response.Response{data=dto.PaperInfo}  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/templates/create [post]
// @Security     BearerAuth
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
// @Summary      获取模板列表
// @Description  获取试卷模板列表，支持分页
// @Tags         试卷管理
// @Accept       json
// @Produce      json
// @Param        page       query     int  false  "页码"
// @Param        page_size  query     int  false  "每页数量"
// @Success      200  {object}  response.PageResponse{data=[]dto.TemplateInfo}  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/templates [get]
// @Security     BearerAuth
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
// @Summary      导出试卷
// @Description  导出指定试卷为JSON格式
// @Tags         试卷管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "试卷ID"
// @Success      200  {object}  dto.ExportPaperResponse  "成功"
// @Failure      400  {object}  response.Response  "无效的试卷ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/papers/{id}/export [get]
// @Security     BearerAuth
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
// @Summary      获取试卷统计
// @Description  获取指定试卷的统计数据
// @Tags         试卷管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "试卷ID"
// @Success      200  {object}  response.Response{data=dto.PaperStatsInfo}  "成功"
// @Failure      400  {object}  response.Response  "无效的试卷ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/papers/{id}/stats [get]
// @Security     BearerAuth
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
