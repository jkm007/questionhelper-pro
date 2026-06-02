package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/service/user"
	"questionhelper-server/pkg/response"
)

type TagController struct{}

func NewTagController() *TagController {
	return &TagController{}
}

// CreateTag 创建标签
// @Summary      创建标签
// @Description  创建新标签（管理员）
// @Tags         标签管理
// @Accept       json
// @Produce      json
// @Param        req  body      dto.CreateTagRequest  true  "标签信息"
// @Success      200  {object}  response.Response  "创建成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/tags [post]
// @Security     BearerAuth
func (ctrl *TagController) CreateTag(c *gin.Context) {
	var req dto.CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := user.CreateTag(&req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "创建成功", nil)
}

// UpdateTag 更新标签
// @Summary      更新标签
// @Description  更新指定标签信息（管理员）
// @Tags         标签管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint                   true  "标签ID"
// @Param        req  body      dto.UpdateTagRequest   true  "标签信息"
// @Success      200  {object}  response.Response  "更新成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/tags/{id} [put]
// @Security     BearerAuth
func (ctrl *TagController) UpdateTag(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的标签ID")
		return
	}

	var req dto.UpdateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := user.UpdateTag(uint(id), &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "更新成功", nil)
}

// DeleteTag 删除标签
// @Summary      删除标签
// @Description  删除指定标签（管理员）
// @Tags         标签管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "标签ID"
// @Success      200  {object}  response.Response  "删除成功"
// @Failure      400  {object}  response.Response  "无效的标签ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/tags/{id} [delete]
// @Security     BearerAuth
func (ctrl *TagController) DeleteTag(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的标签ID")
		return
	}

	if err := user.DeleteTag(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "删除成功", nil)
}

// GetTag 获取标签详情
// @Summary      获取标签详情
// @Description  根据ID获取标签详情（管理员）
// @Tags         标签管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "标签ID"
// @Success      200  {object}  response.Response{data=dto.TagInfo}  "成功"
// @Failure      400  {object}  response.Response  "无效的标签ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/tags/{id} [get]
// @Security     BearerAuth
func (ctrl *TagController) GetTag(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的标签ID")
		return
	}

	info, err := user.GetTag(uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, info)
}

// ListTags 标签列表
// @Summary      标签列表
// @Description  分页查询标签列表（管理员）
// @Tags         标签管理
// @Accept       json
// @Produce      json
// @Param        page       query     int     false  "页码"       default(1)
// @Param        page_size  query     int     false  "每页数量"   default(10)
// @Param        keyword    query     string  false  "搜索关键词"
// @Success      200  {object}  response.Response{data=dto.PageResponse{list=[]dto.TagInfo}}  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/tags [get]
// @Security     BearerAuth
func (ctrl *TagController) ListTags(c *gin.Context) {
	var req dto.TagListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	list, total, err := user.ListTags(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// GetAllTags 获取所有标签
// @Summary      获取所有标签
// @Description  获取所有标签列表（用户端）
// @Tags         标签
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response{data=[]dto.TagInfo}  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /tags [get]
// @Security     BearerAuth
func (ctrl *TagController) GetAllTags(c *gin.Context) {
	list, err := user.GetAllTags()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, list)
}

// AddUserTags 给用户添加标签
// @Summary      给用户添加标签
// @Description  为指定用户添加标签（管理员）
// @Tags         标签管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint                  true  "用户ID"
// @Param        req  body      dto.UserTagRequest    true  "标签ID列表"
// @Success      200  {object}  response.Response  "添加成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/users/{id}/tags [post]
// @Security     BearerAuth
func (ctrl *TagController) AddUserTags(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	var req dto.UserTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := user.AddUserTags(uint(userID), req.TagIDs); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "添加成功", nil)
}

// RemoveUserTags 移除用户标签
// @Summary      移除用户标签
// @Description  移除指定用户的标签（管理员）
// @Tags         标签管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint                  true  "用户ID"
// @Param        req  body      dto.UserTagRequest    true  "标签ID列表"
// @Success      200  {object}  response.Response  "移除成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/users/{id}/tags [delete]
// @Security     BearerAuth
func (ctrl *TagController) RemoveUserTags(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	var req dto.UserTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := user.RemoveUserTags(uint(userID), req.TagIDs); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "移除成功", nil)
}

// GetUserTags 获取用户的标签
// @Summary      获取用户的标签
// @Description  获取指定用户的标签列表（管理员）
// @Tags         标签管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "用户ID"
// @Success      200  {object}  response.Response{data=[]dto.TagInfo}  "成功"
// @Failure      400  {object}  response.Response  "无效的用户ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/users/{id}/tags [get]
// @Security     BearerAuth
func (ctrl *TagController) GetUserTags(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	list, err := user.GetUserTags(uint(userID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, list)
}
