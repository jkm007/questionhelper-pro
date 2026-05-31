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
