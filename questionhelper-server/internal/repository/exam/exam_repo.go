package exam

import (
	"time"

	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	"questionhelper-server/pkg/database"
)

// ==================== Exam ====================

func FindExamByID(id uint) (*model.Exam, error) {
	var exam model.Exam
	err := database.DB.Preload("Paper").Preload("Creator").First(&exam, id).Error
	return &exam, err
}

func CreateExam(exam *model.Exam) error {
	return database.DB.Create(exam).Error
}

func UpdateExam(exam *model.Exam) error {
	return database.DB.Save(exam).Error
}

func DeleteExamByID(id uint) error {
	return database.DB.Delete(&model.Exam{}, id).Error
}

func ListExams(req *dto.ExamListRequest) ([]model.Exam, int64, error) {
	var exams []model.Exam
	var total int64

	db := database.DB.Model(&model.Exam{})

	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	if req.ClassID != nil {
		db = db.Where("class_id = ?", *req.ClassID)
	}
	if req.Keyword != "" {
		db = db.Where("title LIKE ?", "%"+req.Keyword+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Preload("Paper").Preload("Creator").
		Offset(req.GetOffset()).Limit(req.GetLimit()).
		Order("id DESC").Find(&exams).Error

	return exams, total, err
}

func ListAvailableExams(userID uint, req *dto.PageRequest) ([]model.Exam, int64, error) {
	var exams []model.Exam
	var total int64

	now := time.Now()
	db := database.DB.Model(&model.Exam{}).
		Where("status = ? AND start_time <= ? AND end_time >= ?", 1, now, now)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Preload("Paper").
		Offset(req.GetOffset()).Limit(req.GetLimit()).
		Order("start_time DESC").Find(&exams).Error

	return exams, total, err
}

// ==================== Score Analysis ====================

func GetExamScoreStats(examID uint) (map[string]interface{}, error) {
	var result struct {
		Count int     `json:"count"`
		Avg   float64 `json:"avg"`
		Max   float64 `json:"max"`
		Min   float64 `json:"min"`
	}

	err := database.DB.Model(&model.ExamRecord{}).
		Where("exam_id = ? AND status = ?", examID, 2).
		Select("COUNT(*) as count, AVG(score) as avg, MAX(score) as max, MIN(score) as min").
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	// 统计及格率
	var passCount int64
	exam, err := FindExamByID(examID)
	if err != nil {
		return nil, err
	}
	database.DB.Model(&model.ExamRecord{}).
		Where("exam_id = ? AND status = ? AND score >= ?", examID, 2, exam.PassScore).
		Count(&passCount)

	passRate := float64(0)
	if result.Count > 0 {
		passRate = float64(passCount) / float64(result.Count) * 100
	}

	return map[string]interface{}{
		"total_count": result.Count,
		"avg_score":   result.Avg,
		"max_score":   result.Max,
		"min_score":   result.Min,
		"pass_count":  passCount,
		"pass_rate":   passRate,
	}, nil
}
