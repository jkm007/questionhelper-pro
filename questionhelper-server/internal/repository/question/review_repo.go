package question

import (
	"time"

	"questionhelper-server/internal/model"
	"questionhelper-server/pkg/database"
)

// CreateReview 创建审核记录
func CreateReview(review *model.QuestionReview) error {
	return database.DB.Create(review).Error
}

// UpdateReview 更新审核记录
func UpdateReview(review *model.QuestionReview) error {
	return database.DB.Save(review).Error
}

// FindReviewByID 获取审核记录
func FindReviewByID(id uint) (*model.QuestionReview, error) {
	var review model.QuestionReview
	err := database.DB.First(&review, id).Error
	return &review, err
}

// FindPendingReviewByQuestionID 获取题目待审核记录
func FindPendingReviewByQuestionID(questionID uint) (*model.QuestionReview, error) {
	var review model.QuestionReview
	err := database.DB.Where("question_id = ? AND status = 0", questionID).
		Order("id DESC").First(&review).Error
	return &review, err
}

// ListReviews 审核列表
func ListReviews(status *int8, page, pageSize int) ([]model.QuestionReview, int64, error) {
	var reviews []model.QuestionReview
	var total int64

	db := database.DB.Model(&model.QuestionReview{})
	if status != nil {
		db = db.Where("status = ?", *status)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset((page - 1) * pageSize).
		Limit(pageSize).
		Order("id DESC").
		Find(&reviews).Error

	return reviews, total, err
}

// ListReviewsByReviewerID 获取审核人审核列表
func ListReviewsByReviewerID(reviewerID uint, page, pageSize int) ([]model.QuestionReview, int64, error) {
	var reviews []model.QuestionReview
	var total int64

	db := database.DB.Model(&model.QuestionReview{}).Where("reviewer_id = ?", reviewerID)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset((page - 1) * pageSize).
		Limit(pageSize).
		Order("id DESC").
		Find(&reviews).Error

	return reviews, total, err
}

// BatchCreateReviews 批量创建审核记录
func BatchCreateReviews(reviews []model.QuestionReview) error {
	return database.DB.CreateInBatches(reviews, 100).Error
}

// FindTimeoutReviews 查找超时审核
func FindTimeoutReviews(before time.Time) ([]model.QuestionReview, error) {
	var reviews []model.QuestionReview
	err := database.DB.Where("status = 0 AND submitted_at < ?", before).
		Find(&reviews).Error
	return reviews, err
}

// CountReviewsByStatus 统计各状态审核数量
func CountReviewsByStatus() (map[int8]int64, error) {
	var results []struct {
		Status int8
		Count  int64
	}
	err := database.DB.Model(&model.QuestionReview{}).
		Select("status, count(*) as count").
		Group("status").
		Find(&results).Error
	if err != nil {
		return nil, err
	}

	counts := make(map[int8]int64)
	for _, r := range results {
		counts[r.Status] = r.Count
	}
	return counts, nil
}
