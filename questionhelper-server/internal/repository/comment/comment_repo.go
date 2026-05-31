package comment

import (
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

func List(req *dto.CommentListRequest) ([]model.Comment, int64, error) {
	var comments []model.Comment
	var total int64

	db := database.DB.Model(&model.Comment{}).
		Where("target_type = ? AND target_id = ?", req.TargetType, req.TargetID)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Preload("User").
		Offset(req.GetOffset()).Limit(req.GetLimit()).
		Order("created_at DESC").Find(&comments).Error

	return comments, total, err
}

// ListByAdmin 管理员评论列表
func ListByAdmin(req *dto.CommentListRequest) ([]model.Comment, int64, error) {
	var comments []model.Comment
	var total int64

	db := database.DB.Model(&model.Comment{}).
		Where("target_type = ? AND target_id = ?", req.TargetType, req.TargetID)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Preload("User").
		Offset(req.GetOffset()).Limit(req.GetLimit()).
		Order("created_at DESC").Find(&comments).Error

	return comments, total, err
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

func FindReport(id uint) (*model.CommentReport, error) {
	var report model.CommentReport
	err := database.DB.First(&report, id).Error
	return &report, err
}

func UpdateReport(report *model.CommentReport) error {
	return database.DB.Save(report).Error
}

func ListReports(req *dto.PageRequest) ([]model.CommentReport, int64, error) {
	var reports []model.CommentReport
	var total int64

	db := database.DB.Model(&model.CommentReport{})

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset(req.GetOffset()).Limit(req.GetLimit()).
		Order("created_at DESC").Find(&reports).Error

	return reports, total, err
}
