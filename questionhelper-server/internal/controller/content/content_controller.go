package content

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/service/content"
	"questionhelper-server/pkg/response"
)

type ContentController struct{}

func NewContentController() *ContentController {
	return &ContentController{}
}

// ==================== 创作者申请 ====================

// ApplyCreator 申请成为创作者
// @Summary      申请成为创作者
// @Description  用户提交创作者申请
// @Tags         内容管理
// @Accept       json
// @Produce      json
// @Param        req  body      dto.ApplyCreatorRequest  true  "申请数据"
// @Success      200  {object}  response.Response  "申请提交成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /user/apply/creator [post]
// @Security     BearerAuth
func (ctrl *ContentController) ApplyCreator(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.ApplyCreatorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := content.ApplyCreator(userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "申请提交成功", nil)
}

// ==================== 创作者信息 ====================

// GetCreatorProfile 获取创作者资料
// @Summary      获取创作者资料
// @Description  获取当前用户的创作者资料
// @Tags         内容管理
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /creator/profile [get]
// @Security     BearerAuth
func (ctrl *ContentController) GetCreatorProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	info, err := content.GetCreatorProfile(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, info)
}

// GetCreatorLevel 获取创作者等级
// @Summary      获取创作者等级
// @Description  获取当前用户的创作者等级信息
// @Tags         内容管理
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /creator/level [get]
// @Security     BearerAuth
func (ctrl *ContentController) GetCreatorLevel(c *gin.Context) {
	userID := c.GetUint("user_id")

	levels, err := content.GetCreatorLevel(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, levels)
}

// GetCreatorPoints 获取积分信息
// @Summary      获取创作者积分信息
// @Description  获取当前用户的创作者积分信息
// @Tags         内容管理
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /creator/points [get]
// @Security     BearerAuth
func (ctrl *ContentController) GetCreatorPoints(c *gin.Context) {
	userID := c.GetUint("user_id")

	info, err := content.GetCreatorPoints(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, info)
}

// GetCreatorPointLogs 获取积分记录
// @Summary      获取创作者积分记录
// @Description  获取当前用户的创作者积分变动记录，支持分页
// @Tags         内容管理
// @Accept       json
// @Produce      json
// @Param        page       query     int  false  "页码"
// @Param        page_size  query     int  false  "每页数量"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /creator/points/logs [get]
// @Security     BearerAuth
func (ctrl *ContentController) GetCreatorPointLogs(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.CreatorPointLogsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	logs, total, err := content.GetCreatorPointLogs(userID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, logs, total, req.Page, req.PageSize)
}

// ==================== 创作者协议 ====================

// GetCreatorAgreement 获取协议
// @Summary      获取创作者协议
// @Description  获取当前用户的创作者协议信息
// @Tags         内容管理
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /creator/agreement [get]
// @Security     BearerAuth
func (ctrl *ContentController) GetCreatorAgreement(c *gin.Context) {
	userID := c.GetUint("user_id")

	info, err := content.GetCreatorAgreement(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, info)
}

// SignAgreement 签署协议
// @Summary      签署创作者协议
// @Description  当前用户签署创作者协议
// @Tags         内容管理
// @Accept       json
// @Produce      json
// @Param        req  body      dto.SignAgreementRequest  true  "签署数据"
// @Success      200  {object}  response.Response  "签署成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /creator/agreement/sign [post]
// @Security     BearerAuth
func (ctrl *ContentController) SignAgreement(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.SignAgreementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := content.SignAgreement(userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "签署成功", nil)
}

// ==================== 创作者作品集 ====================

// GetPortfolio 获取作品集详情
// @Summary      获取作品集详情
// @Description  根据ID获取作品集详情
// @Tags         内容管理
// @Accept       json
// @Produce      json
// @Param        id  path      uint  true  "作品集ID"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "无效的作品集ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /creator/portfolios/{id} [get]
// @Security     BearerAuth
func (ctrl *ContentController) GetPortfolio(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的作品集ID")
		return
	}

	info, err := content.GetPortfolio(uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, info)
}

// ListPortfolios 获取作品集列表
// @Summary      获取作品集列表
// @Description  获取作品集列表，支持分页
// @Tags         内容管理
// @Accept       json
// @Produce      json
// @Param        page       query     int  false  "页码"
// @Param        page_size  query     int  false  "每页数量"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /creator/portfolios [get]
// @Security     BearerAuth
func (ctrl *ContentController) ListPortfolios(c *gin.Context) {
	var req dto.PortfolioListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	list, total, err := content.ListPortfolios(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// CreatePortfolio 创建作品集
// @Summary      创建作品集
// @Description  创建一个新的作品集
// @Tags         内容管理
// @Accept       json
// @Produce      json
// @Param        req  body      dto.CreatePortfolioRequest  true  "作品集数据"
// @Success      200  {object}  response.Response  "创建成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /creator/portfolios [post]
// @Security     BearerAuth
func (ctrl *ContentController) CreatePortfolio(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.CreatePortfolioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := content.CreatePortfolio(userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "创建成功", nil)
}

// UpdatePortfolio 更新作品集
// @Summary      更新作品集
// @Description  根据ID更新作品集信息
// @Tags         内容管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint                        true  "作品集ID"
// @Param        req  body      dto.UpdatePortfolioRequest  true  "作品集数据"
// @Success      200  {object}  response.Response  "更新成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /creator/portfolios/{id} [put]
// @Security     BearerAuth
func (ctrl *ContentController) UpdatePortfolio(c *gin.Context) {
	userID := c.GetUint("user_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的作品集ID")
		return
	}

	var req dto.UpdatePortfolioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := content.UpdatePortfolio(userID, uint(id), &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "更新成功", nil)
}

// DeletePortfolio 删除作品集
// @Summary      删除作品集
// @Description  根据ID删除作品集
// @Tags         内容管理
// @Accept       json
// @Produce      json
// @Param        id  path      uint  true  "作品集ID"
// @Success      200  {object}  response.Response  "删除成功"
// @Failure      400  {object}  response.Response  "无效的作品集ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /creator/portfolios/{id} [delete]
// @Security     BearerAuth
func (ctrl *ContentController) DeletePortfolio(c *gin.Context) {
	userID := c.GetUint("user_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的作品集ID")
		return
	}

	if err := content.DeletePortfolio(userID, uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "删除成功", nil)
}

// ==================== 内容版本管理 ====================

// ListContentVersions 获取版本列表
// @Summary      获取内容版本列表
// @Description  获取指定内容的版本历史列表
// @Tags         内容管理
// @Accept       json
// @Produce      json
// @Param        type  path      string  true  "内容类型"
// @Param        id    path      uint    true  "内容ID"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "无效的内容ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /content/{type}/{id}/versions [get]
// @Security     BearerAuth
func (ctrl *ContentController) ListContentVersions(c *gin.Context) {
	contentType := c.Param("type")
	contentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的内容ID")
		return
	}

	versions, err := content.ListContentVersions(contentType, uint(contentID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, versions)
}

// GetContentVersion 获取版本详情
// @Summary      获取内容版本详情
// @Description  获取指定内容的特定版本详情
// @Tags         内容管理
// @Accept       json
// @Produce      json
// @Param        type       path      string  true  "内容类型"
// @Param        id         path      uint    true  "内容ID"
// @Param        versionId  path      uint    true  "版本ID"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /content/{type}/{id}/versions/{versionId} [get]
// @Security     BearerAuth
func (ctrl *ContentController) GetContentVersion(c *gin.Context) {
	contentType := c.Param("type")
	contentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的内容ID")
		return
	}

	versionID, err := strconv.ParseUint(c.Param("versionId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的版本ID")
		return
	}

	info, err := content.GetContentVersion(contentType, uint(contentID), uint(versionID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, info)
}

// RollbackContentVersion 回滚版本
// @Summary      回滚内容版本
// @Description  将指定内容回滚到特定版本
// @Tags         内容管理
// @Accept       json
// @Produce      json
// @Param        type  path      string                       true  "内容类型"
// @Param        id    path      uint                         true  "内容ID"
// @Param        req   body      dto.RollbackVersionRequest   true  "回滚参数"
// @Success      200  {object}  response.Response  "回滚成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /content/{type}/{id}/versions/rollback [post]
// @Security     BearerAuth
func (ctrl *ContentController) RollbackContentVersion(c *gin.Context) {
	userID := c.GetUint("user_id")
	contentType := c.Param("type")
	contentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的内容ID")
		return
	}

	var req dto.RollbackVersionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := content.RollbackContentVersion(userID, contentType, uint(contentID), &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "回滚成功", nil)
}

// ==================== 内容标签 ====================

// GetContentTags 获取内容标签
// @Summary      获取内容标签
// @Description  获取指定内容的标签列表
// @Tags         内容管理
// @Accept       json
// @Produce      json
// @Param        type  path      string  true  "内容类型"
// @Param        id    path      uint    true  "内容ID"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "无效的内容ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /content/{type}/{id}/tags [get]
// @Security     BearerAuth
func (ctrl *ContentController) GetContentTags(c *gin.Context) {
	contentType := c.Param("type")
	contentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的内容ID")
		return
	}

	tags, err := content.GetContentTags(contentType, uint(contentID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, tags)
}

// AddContentTag 添加内容标签
// @Summary      添加内容标签
// @Description  为指定内容添加标签
// @Tags         内容管理
// @Accept       json
// @Produce      json
// @Param        type  path      string                   true  "内容类型"
// @Param        id    path      uint                     true  "内容ID"
// @Param        req   body      dto.AddContentTagRequest  true  "标签数据"
// @Success      200  {object}  response.Response  "添加成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /content/{type}/{id}/tags [post]
// @Security     BearerAuth
func (ctrl *ContentController) AddContentTag(c *gin.Context) {
	contentType := c.Param("type")
	contentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的内容ID")
		return
	}

	var req dto.AddContentTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := content.AddContentTag(contentType, uint(contentID), &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "添加成功", nil)
}

// RemoveContentTag 移除内容标签
// @Summary      移除内容标签
// @Description  移除指定内容的标签
// @Tags         内容管理
// @Accept       json
// @Produce      json
// @Param        type   path      string  true  "内容类型"
// @Param        id     path      uint    true  "内容ID"
// @Param        tagId  path      uint    true  "标签ID"
// @Success      200  {object}  response.Response  "移除成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /content/{type}/{id}/tags/{tagId} [delete]
// @Security     BearerAuth
func (ctrl *ContentController) RemoveContentTag(c *gin.Context) {
	contentType := c.Param("type")
	contentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的内容ID")
		return
	}

	tagID, err := strconv.ParseUint(c.Param("tagId"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的标签ID")
		return
	}

	if err := content.RemoveContentTag(contentType, uint(contentID), uint(tagID)); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "移除成功", nil)
}

// GetHotTags 获取热门标签
// @Summary      获取热门标签
// @Description  获取热门标签列表
// @Tags         内容管理
// @Accept       json
// @Produce      json
// @Param        limit  query     int  false  "返回数量"
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /content/tags/hot [get]
// @Security     BearerAuth
func (ctrl *ContentController) GetHotTags(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	tags, err := content.GetHotTags(limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, tags)
}

// ==================== 内容收藏 ====================

// CheckContentFavorite 检查是否收藏
// @Summary      检查是否已收藏
// @Description  检查当前用户是否已收藏指定内容
// @Tags         内容管理
// @Accept       json
// @Produce      json
// @Param        type  path      string  true  "内容类型"
// @Param        id    path      uint    true  "内容ID"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "无效的内容ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /content/{type}/{id}/favorites [get]
// @Security     BearerAuth
func (ctrl *ContentController) CheckContentFavorite(c *gin.Context) {
	userID := c.GetUint("user_id")
	contentType := c.Param("type")
	contentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的内容ID")
		return
	}

	isFavorite, err := content.CheckContentFavorite(userID, contentType, uint(contentID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, gin.H{"is_favorite": isFavorite})
}

// AddContentFavorite 添加收藏
// @Summary      添加内容收藏
// @Description  收藏指定内容
// @Tags         内容管理
// @Accept       json
// @Produce      json
// @Param        type  path      string                          true  "内容类型"
// @Param        id    path      uint                            true  "内容ID"
// @Param        req   body      dto.AddContentFavoriteRequest   false  "收藏数据"
// @Success      200  {object}  response.Response  "收藏成功"
// @Failure      400  {object}  response.Response  "无效的内容ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /content/{type}/{id}/favorites [post]
// @Security     BearerAuth
func (ctrl *ContentController) AddContentFavorite(c *gin.Context) {
	userID := c.GetUint("user_id")
	contentType := c.Param("type")
	contentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的内容ID")
		return
	}

	var req dto.AddContentFavoriteRequest
	c.ShouldBindJSON(&req)

	if err := content.AddContentFavorite(userID, contentType, uint(contentID), &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "收藏成功", nil)
}

// RemoveContentFavorite 取消收藏
// @Summary      取消内容收藏
// @Description  取消收藏指定内容
// @Tags         内容管理
// @Accept       json
// @Produce      json
// @Param        type  path      string  true  "内容类型"
// @Param        id    path      uint    true  "内容ID"
// @Success      200  {object}  response.Response  "取消收藏成功"
// @Failure      400  {object}  response.Response  "无效的内容ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /content/{type}/{id}/favorites [delete]
// @Security     BearerAuth
func (ctrl *ContentController) RemoveContentFavorite(c *gin.Context) {
	userID := c.GetUint("user_id")
	contentType := c.Param("type")
	contentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的内容ID")
		return
	}

	if err := content.RemoveContentFavorite(userID, contentType, uint(contentID)); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "取消收藏成功", nil)
}

// ListContentFavorites 获取收藏列表
// @Summary      获取收藏列表
// @Description  获取当前用户的收藏列表，支持分页
// @Tags         内容管理
// @Accept       json
// @Produce      json
// @Param        page       query     int  false  "页码"
// @Param        page_size  query     int  false  "每页数量"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /favorites [get]
// @Security     BearerAuth
func (ctrl *ContentController) ListContentFavorites(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.FavoriteListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	list, total, err := content.ListContentFavorites(userID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// ==================== 内容预览 ====================

// GetContentPreview 获取内容预览
// @Summary      获取内容预览
// @Description  获取指定内容的预览数据
// @Tags         内容管理
// @Accept       json
// @Produce      json
// @Param        type  path      string  true  "内容类型"
// @Param        id    path      uint    true  "内容ID"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "无效的内容ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /content/{type}/{id}/preview [get]
// @Security     BearerAuth
func (ctrl *ContentController) GetContentPreview(c *gin.Context) {
	contentType := c.Param("type")
	contentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的内容ID")
		return
	}

	info, err := content.GetContentPreview(contentType, uint(contentID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, info)
}

// ==================== 综合搜索 ====================

// Search 综合搜索
// @Summary      综合搜索
// @Description  根据关键词进行综合搜索，支持分页
// @Tags         内容管理
// @Accept       json
// @Produce      json
// @Param        keyword    query     string  true   "搜索关键词"
// @Param        page       query     int     false  "页码"
// @Param        page_size  query     int     false  "每页数量"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /search [get]
// @Security     BearerAuth
func (ctrl *ContentController) Search(c *gin.Context) {
	var req dto.SearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	list, total, err := content.Search(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// GetSearchSuggestions 获取搜索建议
// @Summary      获取搜索建议
// @Description  根据关键词获取搜索建议
// @Tags         内容管理
// @Accept       json
// @Produce      json
// @Param        keyword  query     string  true   "搜索关键词"
// @Param        limit    query     int     false  "返回数量"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "请输入搜索关键词"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /search/suggestions [get]
// @Security     BearerAuth
func (ctrl *ContentController) GetSearchSuggestions(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		response.Error(c, http.StatusBadRequest, "请输入搜索关键词")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	suggestions, err := content.GetSearchSuggestions(keyword, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, suggestions)
}

// GetHotSearches 获取热门搜索
// @Summary      获取热门搜索
// @Description  获取热门搜索关键词列表
// @Tags         内容管理
// @Accept       json
// @Produce      json
// @Param        limit  query     int  false  "返回数量"
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /search/hot [get]
// @Security     BearerAuth
func (ctrl *ContentController) GetHotSearches(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	hotSearches, err := content.GetHotSearches(limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, hotSearches)
}

// GetSearchHistory 获取搜索历史
// @Summary      获取搜索历史
// @Description  获取当前用户的搜索历史
// @Tags         内容管理
// @Accept       json
// @Produce      json
// @Param        limit  query     int  false  "返回数量"
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /search/history [get]
// @Security     BearerAuth
func (ctrl *ContentController) GetSearchHistory(c *gin.Context) {
	userID := c.GetUint("user_id")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	history, err := content.GetSearchHistory(userID, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, history)
}

// ==================== 审核流程 ====================

// ListPendingReviews 获取待审核列表
// @Summary      获取待审核列表
// @Description  获取待审核内容列表，支持分页
// @Tags         内容管理
// @Accept       json
// @Produce      json
// @Param        page       query     int  false  "页码"
// @Param        page_size  query     int  false  "每页数量"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /reviews/pending [get]
// @Security     BearerAuth
func (ctrl *ContentController) ListPendingReviews(c *gin.Context) {
	var req dto.ReviewListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	list, total, err := content.ListPendingReviews(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// GetReviewDetail 获取审核详情
// @Summary      获取审核详情
// @Description  根据ID获取审核详情
// @Tags         内容管理
// @Accept       json
// @Produce      json
// @Param        id  path      uint  true  "审核ID"
// @Success      200  {object}  response.Response  "成功"
// @Failure      400  {object}  response.Response  "无效的审核ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /reviews/{id} [get]
// @Security     BearerAuth
func (ctrl *ContentController) GetReviewDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的审核ID")
		return
	}

	info, err := content.GetReviewDetail(uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, info)
}

// ApproveReview 通过审核
// @Summary      通过审核
// @Description  审核通过指定内容
// @Tags         内容管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint                       true  "审核ID"
// @Param        req  body      dto.ApproveReviewRequest   false  "审核意见"
// @Success      200  {object}  response.Response  "审核通过"
// @Failure      400  {object}  response.Response  "无效的审核ID"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /reviews/{id}/approve [post]
// @Security     BearerAuth
func (ctrl *ContentController) ApproveReview(c *gin.Context) {
	reviewerID := c.GetUint("user_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的审核ID")
		return
	}

	var req dto.ApproveReviewRequest
	c.ShouldBindJSON(&req)

	if err := content.ApproveReview(reviewerID, uint(id), &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "审核通过", nil)
}

// RejectReview 拒绝审核
// @Summary      拒绝审核
// @Description  审核拒绝指定内容
// @Tags         内容管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint                      true  "审核ID"
// @Param        req  body      dto.RejectReviewRequest   true  "拒绝原因"
// @Success      200  {object}  response.Response  "审核拒绝"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /reviews/{id}/reject [post]
// @Security     BearerAuth
func (ctrl *ContentController) RejectReview(c *gin.Context) {
	reviewerID := c.GetUint("user_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的审核ID")
		return
	}

	var req dto.RejectReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := content.RejectReview(reviewerID, uint(id), &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "审核拒绝", nil)
}

// AddReviewComment 添加审核意见
// @Summary      添加审核意见
// @Description  为审核添加评论意见
// @Tags         内容管理
// @Accept       json
// @Produce      json
// @Param        id   path      uint                         true  "审核ID"
// @Param        req  body      dto.AddReviewCommentRequest  true  "评论数据"
// @Success      200  {object}  response.Response  "添加意见成功"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /reviews/{id}/comment [post]
// @Security     BearerAuth
func (ctrl *ContentController) AddReviewComment(c *gin.Context) {
	userID := c.GetUint("user_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的审核ID")
		return
	}

	var req dto.AddReviewCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := content.AddReviewComment(userID, uint(id), &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithMessage(c, "添加意见成功", nil)
}
