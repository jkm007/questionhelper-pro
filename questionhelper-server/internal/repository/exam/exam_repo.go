package exam

import (
	"time"

	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	"questionhelper-server/pkg/database"
)

// ==================== Paper ====================

func FindPaperByID(id uint) (*model.Paper, error) {
	var paper model.Paper
	err := database.DB.First(&paper, id).Error
	return &paper, err
}

func CreatePaper(paper *model.Paper) error {
	return database.DB.Create(paper).Error
}

func UpdatePaper(paper *model.Paper) error {
	return database.DB.Save(paper).Error
}

func DeletePaperByID(id uint) error {
	return database.DB.Delete(&model.Paper{}, id).Error
}

func ListPapers(req *dto.PageRequest) ([]model.Paper, int64, error) {
	var papers []model.Paper
	var total int64

	db := database.DB.Model(&model.Paper{})

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset(req.GetOffset()).Limit(req.GetLimit()).
		Order("id DESC").Find(&papers).Error

	return papers, total, err
}

// ==================== PaperQuestion ====================

func AddPaperQuestions(paperID uint, questions []model.PaperQuestion) error {
	if len(questions) == 0 {
		return nil
	}
	return database.DB.Create(&questions).Error
}

func GetPaperQuestions(paperID uint) ([]model.PaperQuestion, error) {
	var questions []model.PaperQuestion
	err := database.DB.Where("paper_id = ?", paperID).
		Order("sort ASC").Find(&questions).Error
	return questions, err
}

func DeletePaperQuestions(paperID uint) error {
	return database.DB.Where("paper_id = ?", paperID).Delete(&model.PaperQuestion{}).Error
}

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
