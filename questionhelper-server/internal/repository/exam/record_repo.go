package exam

import (
	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	"questionhelper-server/pkg/database"
)

// ==================== ExamRecord ====================

func FindExamRecord(id uint) (*model.ExamRecord, error) {
	var record model.ExamRecord
	err := database.DB.First(&record, id).Error
	return &record, err
}

func FindExamRecordByUser(examID, userID uint) (*model.ExamRecord, error) {
	var record model.ExamRecord
	err := database.DB.Where("exam_id = ? AND user_id = ?", examID, userID).
		Order("id DESC").First(&record).Error
	return &record, err
}

func CreateExamRecord(record *model.ExamRecord) error {
	return database.DB.Create(record).Error
}

func UpdateExamRecord(record *model.ExamRecord) error {
	return database.DB.Save(record).Error
}

func CountUserAttempts(examID, userID uint) (int64, error) {
	var count int64
	err := database.DB.Model(&model.ExamRecord{}).
		Where("exam_id = ? AND user_id = ?", examID, userID).
		Count(&count).Error
	return count, err
}

func ListExamRecords(examID *uint, userID *uint, req *dto.PageRequest) ([]model.ExamRecord, int64, error) {
	var records []model.ExamRecord
	var total int64

	db := database.DB.Model(&model.ExamRecord{})
	if examID != nil {
		db = db.Where("exam_id = ?", *examID)
	}
	if userID != nil {
		db = db.Where("user_id = ?", *userID)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset(req.GetOffset()).Limit(req.GetLimit()).
		Order("id DESC").Find(&records).Error

	return records, total, err
}

// ==================== AnswerRecord ====================

func CreateAnswerRecords(records []model.AnswerRecord) error {
	if len(records) == 0 {
		return nil
	}
	return database.DB.Create(&records).Error
}

func CreateSingleAnswerRecord(record *model.AnswerRecord) error {
	return database.DB.Create(record).Error
}

func UpdateAnswerRecord(record *model.AnswerRecord) error {
	return database.DB.Save(record).Error
}

func FindAnswerRecordByID(id uint) (*model.AnswerRecord, error) {
	var record model.AnswerRecord
	err := database.DB.First(&record, id).Error
	return &record, err
}

func GetAnswerRecords(recordID uint) ([]model.AnswerRecord, error) {
	var records []model.AnswerRecord
	err := database.DB.Where("record_id = ?", recordID).Find(&records).Error
	return records, err
}

// DeleteAnswerRecordsByRecordID 删除指定考试记录的所有答题记录
func DeleteAnswerRecordsByRecordID(recordID uint) error {
	return database.DB.Where("record_id = ?", recordID).Delete(&model.AnswerRecord{}).Error
}

// CountOngoingByUserID 统计用户进行中的考试记录数量
func CountOngoingByUserID(userID uint) (int64, error) {
	var count int64
	err := database.DB.Model(&model.ExamRecord{}).
		Where("user_id = ? AND status = 0", userID).
		Count(&count).Error
	return count, err
}
