package content

import (
	"time"

	"gorm.io/gorm"

	"questionhelper-server/internal/model"
	"questionhelper-server/pkg/database"
)

// ==================== 创作者申请 ====================

// CreateApplication 创建创作者申请
func CreateApplication(app *model.RoleApplication) error {
	return database.DB.Create(app).Error
}

// FindApplicationByUserAndRole 根据用户和角色查找申请
func FindApplicationByUserAndRole(userID, roleID uint) (*model.RoleApplication, error) {
	var app model.RoleApplication
	err := database.DB.Where("user_id = ? AND role_id = ? AND status = 0", userID, roleID).First(&app).Error
	return &app, err
}

// ==================== 创作者等级 ====================

// FindCreatorLevel 根据积分查找对应等级
func FindCreatorLevel(points int) (*model.CreatorLevel, error) {
	var level model.CreatorLevel
	err := database.DB.Where("min_points <= ? AND max_points >= ?", points, points).First(&level).Error
	return &level, err
}

// FindCreatorLevelByLevel 根据等级值查找等级
func FindCreatorLevelByLevel(level int) (*model.CreatorLevel, error) {
	var creatorLevel model.CreatorLevel
	err := database.DB.Where("level = ?", level).First(&creatorLevel).Error
	return &creatorLevel, err
}

// ListCreatorLevels 获取所有创作者等级
func ListCreatorLevels() ([]model.CreatorLevel, error) {
	var levels []model.CreatorLevel
	err := database.DB.Order("level ASC").Find(&levels).Error
	return levels, err
}

// ==================== 创作者积分 ====================

// FindCreatorPoint 根据用户ID查找积分
func FindCreatorPoint(userID uint) (*model.CreatorPoint, error) {
	var point model.CreatorPoint
	err := database.DB.Where("user_id = ?", userID).First(&point).Error
	return &point, err
}

// CreateCreatorPoint 创建创作者积分
func CreateCreatorPoint(point *model.CreatorPoint) error {
	return database.DB.Create(point).Error
}

// UpdateCreatorPoint 更新创作者积分
func UpdateCreatorPoint(point *model.CreatorPoint) error {
	return database.DB.Save(point).Error
}

// ListCreatorPointLogs 获取积分记录
func ListCreatorPointLogs(userID uint, changeType string, offset, limit int) ([]model.CreatorPointLog, int64, error) {
	var logs []model.CreatorPointLog
	var total int64

	db := database.DB.Model(&model.CreatorPointLog{}).Where("user_id = ?", userID)
	if changeType != "" {
		db = db.Where("change_type = ?", changeType)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset(offset).Limit(limit).Order("id DESC").Find(&logs).Error
	return logs, total, err
}

// CreateCreatorPointLog 创建积分记录
func CreateCreatorPointLog(log *model.CreatorPointLog) error {
	return database.DB.Create(log).Error
}

// ==================== 创作者协议 ====================

// FindActiveAgreement 获取当前生效的协议
func FindActiveAgreement() (*model.CreatorAgreement, error) {
	var agreement model.CreatorAgreement
	err := database.DB.Where("is_active = ?", true).Order("id DESC").First(&agreement).Error
	return &agreement, err
}

// FindAgreementByID 根据ID获取协议
func FindAgreementByID(id uint) (*model.CreatorAgreement, error) {
	var agreement model.CreatorAgreement
	err := database.DB.First(&agreement, id).Error
	return &agreement, err
}

// FindAgreementSign 查找协议签署记录
func FindAgreementSign(agreementID, userID uint) (*model.CreatorAgreementSign, error) {
	var sign model.CreatorAgreementSign
	err := database.DB.Where("agreement_id = ? AND user_id = ?", agreementID, userID).First(&sign).Error
	return &sign, err
}

// CreateAgreementSign 创建协议签署记录
func CreateAgreementSign(sign *model.CreatorAgreementSign) error {
	return database.DB.Create(sign).Error
}

// ==================== 创作者作品集 ====================

// FindPortfolioByID 根据ID查找作品集
func FindPortfolioByID(id uint) (*model.CreatorPortfolio, error) {
	var portfolio model.CreatorPortfolio
	err := database.DB.First(&portfolio, id).Error
	return &portfolio, err
}

// CreatePortfolio 创建作品集
func CreatePortfolio(portfolio *model.CreatorPortfolio) error {
	return database.DB.Create(portfolio).Error
}

// UpdatePortfolio 更新作品集
func UpdatePortfolio(portfolio *model.CreatorPortfolio) error {
	return database.DB.Save(portfolio).Error
}

// DeletePortfolio 删除作品集
func DeletePortfolio(id uint) error {
	return database.DB.Delete(&model.CreatorPortfolio{}, id).Error
}

// ListPortfolios 获取作品集列表
func ListPortfolios(userID *uint, offset, limit int) ([]model.CreatorPortfolio, int64, error) {
	var portfolios []model.CreatorPortfolio
	var total int64

	db := database.DB.Model(&model.CreatorPortfolio{})
	if userID != nil {
		db = db.Where("user_id = ?", *userID)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset(offset).Limit(limit).Order("id DESC").Find(&portfolios).Error
	return portfolios, total, err
}

// CountPortfolioItems 统计作品集内容数量
func CountPortfolioItems(portfolioID uint) (int, error) {
	var count int64
	err := database.DB.Model(&model.CreatorPortfolioItem{}).Where("portfolio_id = ?", portfolioID).Count(&count).Error
	return int(count), err
}

// ==================== 内容版本管理 ====================

// FindContentVersionByID 根据ID查找版本
func FindContentVersionByID(id uint) (*model.ContentVersion, error) {
	var version model.ContentVersion
	err := database.DB.First(&version, id).Error
	return &version, err
}

// ListContentVersions 获取内容版本列表
func ListContentVersions(contentType string, contentID uint) ([]model.ContentVersion, error) {
	var versions []model.ContentVersion
	err := database.DB.Where("content_type = ? AND content_id = ?", contentType, contentID).
		Order("version DESC").Find(&versions).Error
	return versions, err
}

// CreateContentVersion 创建内容版本
func CreateContentVersion(version *model.ContentVersion) error {
	return database.DB.Create(version).Error
}

// GetLatestVersion 获取最新版本号
func GetLatestVersion(contentType string, contentID uint) (int, error) {
	var version model.ContentVersion
	err := database.DB.Where("content_type = ? AND content_id = ?", contentType, contentID).
		Order("version DESC").First(&version).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, nil
		}
		return 0, err
	}
	return version.Version, nil
}

// ==================== 内容标签 ====================

// FindContentTagByID 根据ID查找标签
func FindContentTagByID(id uint) (*model.ContentTag, error) {
	var tag model.ContentTag
	err := database.DB.First(&tag, id).Error
	return &tag, err
}

// FindContentTagByName 根据名称查找标签
func FindContentTagByName(name string) (*model.ContentTag, error) {
	var tag model.ContentTag
	err := database.DB.Where("name = ?", name).First(&tag).Error
	return &tag, err
}

// FindContentTagBySlug 根据Slug查找标签
func FindContentTagBySlug(slug string) (*model.ContentTag, error) {
	var tag model.ContentTag
	err := database.DB.Where("slug = ?", slug).First(&tag).Error
	return &tag, err
}

// CreateContentTag 创建标签
func CreateContentTag(tag *model.ContentTag) error {
	return database.DB.Create(tag).Error
}

// UpdateContentTag 更新标签
func UpdateContentTag(tag *model.ContentTag) error {
	return database.DB.Save(tag).Error
}

// ListContentTags 获取内容标签列表
func ListContentTags(contentType string, contentID uint) ([]model.ContentTag, error) {
	var tags []model.ContentTag
	err := database.DB.Joins("JOIN content_tag_relations ON content_tag_relations.tag_id = content_tags.id").
		Where("content_tag_relations.content_type = ? AND content_tag_relations.content_id = ?", contentType, contentID).
		Find(&tags).Error
	return tags, err
}

// FindContentTagRelation 查找内容标签关联
func FindContentTagRelation(tagID uint, contentType string, contentID uint) (*model.ContentTagRelation, error) {
	var relation model.ContentTagRelation
	err := database.DB.Where("tag_id = ? AND content_type = ? AND content_id = ?", tagID, contentType, contentID).
		First(&relation).Error
	return &relation, err
}

// CreateContentTagRelation 创建内容标签关联
func CreateContentTagRelation(relation *model.ContentTagRelation) error {
	return database.DB.Create(relation).Error
}

// DeleteContentTagRelation 删除内容标签关联
func DeleteContentTagRelation(tagID uint, contentType string, contentID uint) error {
	return database.DB.Where("tag_id = ? AND content_type = ? AND content_id = ?", tagID, contentType, contentID).
		Delete(&model.ContentTagRelation{}).Error
}

// ListHotTags 获取热门标签
func ListHotTags(limit int) ([]model.ContentTag, error) {
	var tags []model.ContentTag
	err := database.DB.Where("is_active = ?", true).
		Order("usage_count DESC").Limit(limit).Find(&tags).Error
	return tags, err
}

// ==================== 内容收藏 ====================

// FindContentFavorite 查找内容收藏
func FindContentFavorite(userID uint, contentType string, contentID uint) (*model.ContentFavorite, error) {
	var favorite model.ContentFavorite
	err := database.DB.Where("user_id = ? AND content_type = ? AND content_id = ?", userID, contentType, contentID).
		First(&favorite).Error
	return &favorite, err
}

// CreateContentFavorite 创建内容收藏
func CreateContentFavorite(favorite *model.ContentFavorite) error {
	return database.DB.Create(favorite).Error
}

// DeleteContentFavorite 删除内容收藏
func DeleteContentFavorite(userID uint, contentType string, contentID uint) error {
	return database.DB.Where("user_id = ? AND content_type = ? AND content_id = ?", userID, contentType, contentID).
		Delete(&model.ContentFavorite{}).Error
}

// ListContentFavorites 获取收藏列表
func ListContentFavorites(userID uint, contentType, folderName string, offset, limit int) ([]model.ContentFavorite, int64, error) {
	var favorites []model.ContentFavorite
	var total int64

	db := database.DB.Model(&model.ContentFavorite{}).Where("user_id = ?", userID)
	if contentType != "" {
		db = db.Where("content_type = ?", contentType)
	}
	if folderName != "" {
		db = db.Where("folder_name = ?", folderName)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset(offset).Limit(limit).Order("id DESC").Find(&favorites).Error
	return favorites, total, err
}

// ==================== 内容搜索 ====================

// CreateSearchLog 创建搜索记录
func CreateSearchLog(log *model.ContentSearchLog) error {
	return database.DB.Create(log).Error
}

// ListSearchHistory 获取搜索历史
func ListSearchHistory(userID uint, limit int) ([]model.ContentSearchLog, error) {
	var logs []model.ContentSearchLog
	err := database.DB.Where("user_id = ?", userID).
		Order("searched_at DESC").Limit(limit).Find(&logs).Error
	return logs, err
}

// DeleteSearchHistory 删除搜索历史
func DeleteSearchHistory(userID uint) error {
	return database.DB.Where("user_id = ?", userID).Delete(&model.ContentSearchLog{}).Error
}

// ListHotSearches 获取热门搜索
func ListHotSearches(limit int) ([]struct {
	Keyword string
	Count   int
}, error) {
	var results []struct {
		Keyword string
		Count   int
	}
	err := database.DB.Model(&model.ContentSearchLog{}).
		Select("keyword, COUNT(*) as count").
		Where("created_at > ?", time.Now().AddDate(0, 0, -7)).
		Group("keyword").
		Order("count DESC").
		Limit(limit).
		Find(&results).Error
	return results, err
}

// SearchContentSuggestions 获取搜索建议
func SearchContentSuggestions(keyword string, limit int) ([]string, error) {
	var keywords []string
	err := database.DB.Model(&model.ContentSearchLog{}).
		Select("DISTINCT keyword").
		Where("keyword LIKE ?", "%"+keyword+"%").
		Limit(limit).
		Pluck("keyword", &keywords).Error
	return keywords, err
}

// ==================== 审核流程 ====================

// FindReviewInstanceByID 根据ID查找审核实例
func FindReviewInstanceByID(id uint) (*model.ReviewInstance, error) {
	var instance model.ReviewInstance
	err := database.DB.First(&instance, id).Error
	return &instance, err
}

// ListReviewInstances 获取审核实例列表
func ListReviewInstances(status int8, contentType string, offset, limit int) ([]model.ReviewInstance, int64, error) {
	var instances []model.ReviewInstance
	var total int64

	db := database.DB.Model(&model.ReviewInstance{})
	if status >= 0 {
		db = db.Where("status = ?", status)
	}
	if contentType != "" {
		db = db.Where("content_type = ?", contentType)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset(offset).Limit(limit).Order("id DESC").Find(&instances).Error
	return instances, total, err
}

// UpdateReviewInstance 更新审核实例
func UpdateReviewInstance(instance *model.ReviewInstance) error {
	return database.DB.Save(instance).Error
}

// ListReviewStepLogs 获取审核步骤记录
func ListReviewStepLogs(instanceID uint) ([]model.ReviewStepLog, error) {
	var logs []model.ReviewStepLog
	err := database.DB.Where("instance_id = ?", instanceID).Order("step_index ASC").Find(&logs).Error
	return logs, err
}

// CreateReviewStepLog 创建审核步骤记录
func CreateReviewStepLog(log *model.ReviewStepLog) error {
	return database.DB.Create(log).Error
}

// CreateReviewReply 创建审核回复
func CreateReviewReply(reply *model.ReviewReply) error {
	return database.DB.Create(reply).Error
}

// ListReviewReplies 获取审核回复列表
func ListReviewReplies(instanceID uint) ([]model.ReviewReply, error) {
	var replies []model.ReviewReply
	err := database.DB.Where("instance_id = ?", instanceID).Order("reply_at ASC").Find(&replies).Error
	return replies, err
}

// CreateReviewNotification 创建审核通知
func CreateReviewNotification(notification *model.ReviewNotification) error {
	return database.DB.Create(notification).Error
}
