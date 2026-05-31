package exam

import (
	"time"

	"questionhelper-server/internal/model"
	"questionhelper-server/pkg/database"
)

// CreateWarning 创建考试异常记录
func CreateWarning(warning *model.ExamWarning) error {
	return database.DB.Create(warning).Error
}

// FindWarningsByRecordID 获取考试记录的异常
func FindWarningsByRecordID(recordID uint) ([]model.ExamWarning, error) {
	var warnings []model.ExamWarning
	err := database.DB.Where("record_id = ?", recordID).
		Order("created_at DESC").
		Find(&warnings).Error
	return warnings, err
}

// FindWarningsByExamID 获取考试的异常记录
func FindWarningsByExamID(examID uint, page, pageSize int) ([]model.ExamWarning, int64, error) {
	var warnings []model.ExamWarning
	var total int64

	db := database.DB.Model(&model.ExamWarning{}).Where("exam_id = ?", examID)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset((page - 1) * pageSize).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&warnings).Error

	return warnings, total, err
}

// CountWarningsByExamID 统计考试异常数量
func CountWarningsByExamID(examID uint) (int64, error) {
	var count int64
	err := database.DB.Model(&model.ExamWarning{}).Where("exam_id = ?", examID).Count(&count).Error
	return count, err
}

// FindWarningsBefore 删除过期异常记录
func DeleteWarningsBefore(before time.Time) (int64, error) {
	result := database.DB.Where("created_at < ?", before).Delete(&model.ExamWarning{})
	return result.RowsAffected, result.Error
}
