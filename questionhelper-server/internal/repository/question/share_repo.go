package question

import (
	"time"

	"gorm.io/gorm"

	"questionhelper-server/internal/model"
	"questionhelper-server/pkg/database"
)

// CreateShare 创建分享
func CreateShare(share *model.QuestionShare) error {
	return database.DB.Create(share).Error
}

// FindShareByCode 通过分享码查找
func FindShareByCode(code string) (*model.QuestionShare, error) {
	var share model.QuestionShare
	err := database.DB.Where("share_code = ?", code).First(&share).Error
	return &share, err
}

// FindShareByID 通过ID查找
func FindShareByID(id uint) (*model.QuestionShare, error) {
	var share model.QuestionShare
	err := database.DB.First(&share, id).Error
	return &share, err
}

// UpdateShare 更新分享
func UpdateShare(share *model.QuestionShare) error {
	return database.DB.Save(share).Error
}

// DeleteShare 删除分享
func DeleteShare(id uint) error {
	return database.DB.Delete(&model.QuestionShare{}, id).Error
}

// ListSharesByQuestionID 获取题目分享列表
func ListSharesByQuestionID(questionID uint) ([]model.QuestionShare, error) {
	var shares []model.QuestionShare
	err := database.DB.Where("question_id = ?", questionID).
		Order("id DESC").
		Find(&shares).Error
	return shares, err
}

// ListSharesByUserID 获取用户分享列表
func ListSharesByUserID(userID uint, page, pageSize int) ([]model.QuestionShare, int64, error) {
	var shares []model.QuestionShare
	var total int64

	db := database.DB.Model(&model.QuestionShare{}).Where("user_id = ?", userID)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset((page - 1) * pageSize).
		Limit(pageSize).
		Order("id DESC").
		Find(&shares).Error

	return shares, total, err
}

// IncrementShareViewCount 增加分享查看次数
func IncrementShareViewCount(id uint) error {
	return database.DB.Model(&model.QuestionShare{}).Where("id = ?", id).
		Update("view_count", gorm.Expr("view_count + 1")).Error
}

// DeleteExpiredShares 删除过期分享
func DeleteExpiredShares() (int64, error) {
	result := database.DB.Where("expires_at IS NOT NULL AND expires_at < ?", time.Now()).
		Delete(&model.QuestionShare{})
	return result.RowsAffected, result.Error
}
