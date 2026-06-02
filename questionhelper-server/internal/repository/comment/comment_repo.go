package comment

import (
	"time"

	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	"questionhelper-server/pkg/database"
)

// ==================== Comment ====================

func FindByID(id uint) (*model.Comment, error) {
	var comment model.Comment
	err := database.DB.Preload("User").First(&comment, id).Error
	return &comment, err
}

func Create(comment *model.Comment) error {
	return database.DB.Create(comment).Error
}

func Update(comment *model.Comment) error {
	return database.DB.Save(comment).Error
}

func DeleteByID(id uint) error {
	return database.DB.Delete(&model.Comment{}, id).Error
}

func DeleteByIDs(ids []uint) error {
	return database.DB.Delete(&model.Comment{}, ids).Error
}

func List(req *dto.CommentListRequest) ([]model.Comment, int64, error) {
	var comments []model.Comment
	var total int64

	db := database.DB.Model(&model.Comment{}).
		Where("target_type = ? AND target_id = ?", req.TargetType, req.TargetID)

	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	} else {
		db = db.Where("status = 1")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Preload("User").
		Offset(req.GetOffset()).Limit(req.GetLimit()).
		Order("is_pinned DESC, created_at DESC").Find(&comments).Error

	return comments, total, err
}

// ListByAdmin 管理员评论列表
func ListByAdmin(req *dto.CommentAdminListRequest) ([]model.Comment, int64, error) {
	var comments []model.Comment
	var total int64

	db := database.DB.Model(&model.Comment{})

	if req.TargetType != nil {
		db = db.Where("target_type = ?", *req.TargetType)
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	if req.UserID != nil {
		db = db.Where("user_id = ?", *req.UserID)
	}
	if req.Keyword != "" {
		db = db.Where("content LIKE ?", "%"+req.Keyword+"%")
	}
	if req.StartDate != nil {
		db = db.Where("created_at >= ?", *req.StartDate)
	}
	if req.EndDate != nil {
		db = db.Where("created_at <= ?", *req.EndDate)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Preload("User").
		Offset(req.GetOffset()).Limit(req.GetLimit()).
		Order("created_at DESC").Find(&comments).Error

	return comments, total, err
}

// BatchUpdateStatus 批量更新评论状态
func BatchUpdateStatus(ids []uint, status int8) error {
	return database.DB.Model(&model.Comment{}).Where("id IN ?", ids).
		Update("status", status).Error
}

// CountByStatus 按状态统计评论数
func CountByStatus(status int8) (int64, error) {
	var count int64
	err := database.DB.Model(&model.Comment{}).Where("status = ?", status).Count(&count).Error
	return count, err
}

// CountTotal 统计评论总数
func CountTotal() (int64, error) {
	var count int64
	err := database.DB.Model(&model.Comment{}).Count(&count).Error
	return count, err
}

// CountByDateRange 按日期范围统计评论数
func CountByDateRange(start, end time.Time) (int64, error) {
	var count int64
	err := database.DB.Model(&model.Comment{}).
		Where("created_at >= ? AND created_at < ?", start, end).
		Count(&count).Error
	return count, err
}

// DailyStats 按日统计评论数（最近N天）
func DailyStats(days int) ([]dto.DailyCommentStat, error) {
	var stats []dto.DailyCommentStat
	startDate := time.Now().AddDate(0, 0, -days)

	rows, err := database.DB.Model(&model.Comment{}).
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("created_at >= ?", startDate).
		Group("DATE(created_at)").
		Order("date ASC").
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var stat dto.DailyCommentStat
		if err := rows.Scan(&stat.Date, &stat.Count); err != nil {
			return nil, err
		}
		stats = append(stats, stat)
	}
	return stats, nil
}

// TopTargetStats 评论最多的目标统计
func TopTargetStats(limit int) ([]dto.TargetCommentStat, error) {
	var stats []dto.TargetCommentStat

	rows, err := database.DB.Model(&model.Comment{}).
		Select("target_type, target_id, COUNT(*) as count").
		Group("target_type, target_id").
		Order("count DESC").
		Limit(limit).
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var stat dto.TargetCommentStat
		if err := rows.Scan(&stat.TargetType, &stat.TargetID, &stat.Count); err != nil {
			return nil, err
		}
		stats = append(stats, stat)
	}
	return stats, nil
}

// ListForExport 导出评论列表
func ListForExport(req *dto.CommentExportRequest) ([]model.Comment, error) {
	var comments []model.Comment

	db := database.DB.Model(&model.Comment{})
	if req.TargetType != nil {
		db = db.Where("target_type = ?", *req.TargetType)
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	if req.StartDate != nil {
		db = db.Where("created_at >= ?", *req.StartDate)
	}
	if req.EndDate != nil {
		db = db.Where("created_at <= ?", *req.EndDate)
	}

	err := db.Preload("User").Order("created_at DESC").Find(&comments).Error
	return comments, err
}

// ==================== CommentLike ====================

func FindLike(commentID, userID uint) (*model.CommentLike, error) {
	var like model.CommentLike
	err := database.DB.Where("comment_id = ? AND user_id = ?", commentID, userID).First(&like).Error
	return &like, err
}

func CreateLike(like *model.CommentLike) error {
	return database.DB.Create(like).Error
}

func DeleteLike(commentID, userID uint) error {
	return database.DB.Where("comment_id = ? AND user_id = ?", commentID, userID).Delete(&model.CommentLike{}).Error
}

func IncrementLikeCount(commentID uint) error {
	return database.DB.Model(&model.Comment{}).Where("id = ?", commentID).
		UpdateColumn("like_count", database.DB.Raw("like_count + 1")).Error
}

func DecrementLikeCount(commentID uint) error {
	return database.DB.Model(&model.Comment{}).Where("id = ?", commentID).
		UpdateColumn("like_count", database.DB.Raw("GREATEST(like_count - 1, 0)")).Error
}

// ==================== CommentReport ====================

func CreateReport(report *model.CommentReport) error {
	return database.DB.Create(report).Error
}

func FindReportByUserAndComment(userID, commentID uint) (*model.CommentReport, error) {
	var report model.CommentReport
	err := database.DB.Where("user_id = ? AND comment_id = ?", userID, commentID).First(&report).Error
	return &report, err
}

func FindReport(id uint) (*model.CommentReport, error) {
	var report model.CommentReport
	err := database.DB.First(&report, id).Error
	return &report, err
}

func UpdateReport(report *model.CommentReport) error {
	return database.DB.Save(report).Error
}

func ListReports(req *dto.ReportListRequest) ([]model.CommentReport, int64, error) {
	var reports []model.CommentReport
	var total int64

	db := database.DB.Model(&model.CommentReport{})

	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	if req.ReasonType != "" {
		db = db.Where("reason_type = ?", req.ReasonType)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset(req.GetOffset()).Limit(req.GetLimit()).
		Order("created_at DESC").Find(&reports).Error

	return reports, total, err
}

// CountPendingReports 统计待处理举报数
func CountPendingReports() (int64, error) {
	var count int64
	err := database.DB.Model(&model.CommentReport{}).Where("status = 0").Count(&count).Error
	return count, err
}

// ==================== CommentBlacklist ====================

func CreateBlacklist(bl *model.CommentBlacklist) error {
	return database.DB.Create(bl).Error
}

func FindBlacklistByID(id uint) (*model.CommentBlacklist, error) {
	var bl model.CommentBlacklist
	err := database.DB.Preload("User").First(&bl, id).Error
	return &bl, err
}

func DeleteBlacklist(id uint) error {
	return database.DB.Delete(&model.CommentBlacklist{}, id).Error
}

func FindBlacklistByUserAndTarget(userID uint, targetType int8, targetID uint) (*model.CommentBlacklist, error) {
	var bl model.CommentBlacklist
	err := database.DB.Where("user_id = ? AND target_type = ? AND target_id = ?", userID, targetType, targetID).
		First(&bl).Error
	return &bl, err
}

func ListBlacklists(req *dto.BlacklistListRequest) ([]model.CommentBlacklist, int64, error) {
	var list []model.CommentBlacklist
	var total int64

	db := database.DB.Model(&model.CommentBlacklist{})
	if req.UserID != nil {
		db = db.Where("user_id = ?", *req.UserID)
	}
	if req.TargetType != nil {
		db = db.Where("target_type = ?", *req.TargetType)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Preload("User").
		Offset(req.GetOffset()).Limit(req.GetLimit()).
		Order("created_at DESC").Find(&list).Error

	return list, total, err
}

// ==================== Sticker ====================

func ListStickers(category string) ([]model.Sticker, error) {
	var stickers []model.Sticker
	db := database.DB.Model(&model.Sticker{}).Where("status = 1")
	if category != "" {
		db = db.Where("category = ?", category)
	}
	err := db.Order("sort_order DESC, id ASC").Find(&stickers).Error
	return stickers, err
}

func ListStickerCategories() ([]dto.StickerCategoryInfo, error) {
	var categories []dto.StickerCategoryInfo

	rows, err := database.DB.Model(&model.Sticker{}).
		Where("status = 1").
		Select("category, COUNT(*) as count").
		Group("category").
		Order("count DESC").
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var cat dto.StickerCategoryInfo
		if err := rows.Scan(&cat.Name, &cat.Count); err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}
	return categories, nil
}

// ==================== CommentAuditRule ====================

func CreateAuditRule(rule *model.CommentAuditRule) error {
	return database.DB.Create(rule).Error
}

func FindAuditRuleByID(id uint) (*model.CommentAuditRule, error) {
	var rule model.CommentAuditRule
	err := database.DB.First(&rule, id).Error
	return &rule, err
}

func UpdateAuditRule(rule *model.CommentAuditRule) error {
	return database.DB.Save(rule).Error
}

func DeleteAuditRule(id uint) error {
	return database.DB.Delete(&model.CommentAuditRule{}, id).Error
}

func ListAuditRules() ([]model.CommentAuditRule, error) {
	var rules []model.CommentAuditRule
	err := database.DB.Order("priority DESC, id ASC").Find(&rules).Error
	return rules, err
}

// FindEnabledAuditRules 获取所有启用的审核规则
func FindEnabledAuditRules() ([]model.CommentAuditRule, error) {
	var rules []model.CommentAuditRule
	err := database.DB.Where("status = 1").Order("priority DESC, id ASC").Find(&rules).Error
	return rules, err
}

// ==================== User Search ====================

func SearchUsers(keyword string, limit int) ([]model.User, error) {
	var users []model.User
	err := database.DB.Model(&model.User{}).
		Where("status = 1 AND (username LIKE ? OR nickname LIKE ?)", "%"+keyword+"%", "%"+keyword+"%").
		Limit(limit).
		Order("id ASC").
		Find(&users).Error
	return users, err
}
