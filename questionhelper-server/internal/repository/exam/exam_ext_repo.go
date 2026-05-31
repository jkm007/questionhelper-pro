package exam

import (
	"time"

	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	"questionhelper-server/pkg/database"
)

// ==================== PaperShare ====================

func CreatePaperShare(share *model.PaperShare) error {
	return database.DB.Create(share).Error
}

func FindPaperShare(paperID, targetID uint, targetType int8) (*model.PaperShare, error) {
	var share model.PaperShare
	err := database.DB.Where("paper_id = ? AND target_id = ? AND target_type = ? AND status = 1",
		paperID, targetID, targetType).First(&share).Error
	return &share, err
}

// ==================== PaperFavorite ====================

func CreatePaperFavorite(fav *model.PaperFavorite) error {
	return database.DB.Create(fav).Error
}

func FindPaperFavorite(userID, paperID uint) (*model.PaperFavorite, error) {
	var fav model.PaperFavorite
	err := database.DB.Where("user_id = ? AND paper_id = ?", userID, paperID).First(&fav).Error
	return &fav, err
}

func DeletePaperFavorite(userID, paperID uint) error {
	return database.DB.Where("user_id = ? AND paper_id = ?", userID, paperID).
		Delete(&model.PaperFavorite{}).Error
}

// ==================== ExamExtension ====================

func CreateExamExtension(ext *model.ExamExtension) error {
	return database.DB.Create(ext).Error
}

func FindExtensionsByExamID(examID uint) ([]model.ExamExtension, error) {
	var extensions []model.ExamExtension
	err := database.DB.Where("exam_id = ? AND status = 1", examID).
		Order("created_at DESC").Find(&extensions).Error
	return extensions, err
}

// ==================== ExamPause ====================

func CreateExamPause(pause *model.ExamPause) error {
	return database.DB.Create(pause).Error
}

func FindLastPause(examID uint) (*model.ExamPause, error) {
	var pause model.ExamPause
	err := database.DB.Where("exam_id = ?", examID).
		Order("id DESC").First(&pause).Error
	return &pause, err
}

func UpdateExamPause(pause *model.ExamPause) error {
	return database.DB.Save(pause).Error
}

// ==================== ScoreReview ====================

func CreateScoreReview(review *model.ScoreReview) error {
	return database.DB.Create(review).Error
}

func FindScoreReviewByID(id uint) (*model.ScoreReview, error) {
	var review model.ScoreReview
	err := database.DB.First(&review, id).Error
	return &review, err
}

func UpdateScoreReview(review *model.ScoreReview) error {
	return database.DB.Save(review).Error
}

func ListScoreReviews(req *dto.ScoreReviewListRequest) ([]model.ScoreReview, int64, error) {
	var reviews []model.ScoreReview
	var total int64

	db := database.DB.Model(&model.ScoreReview{})

	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset(req.GetOffset()).Limit(req.GetLimit()).
		Order("id DESC").Find(&reviews).Error

	return reviews, total, err
}

func FindScoreReviewByRecordAndUser(recordID, userID uint) (*model.ScoreReview, error) {
	var review model.ScoreReview
	err := database.DB.Where("record_id = ? AND user_id = ? AND status = 0",
		recordID, userID).First(&review).Error
	return &review, err
}

// ==================== ExamNotice ====================

func CreateExamNotice(notice *model.ExamNotice) error {
	return database.DB.Create(notice).Error
}

func ListExamNotices(examID uint, page, pageSize int) ([]model.ExamNotice, int64, error) {
	var notices []model.ExamNotice
	var total int64

	db := database.DB.Model(&model.ExamNotice{}).Where("exam_id = ? AND status = 1", examID)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset((page-1)*pageSize).Limit(pageSize).
		Order("is_pinned DESC, priority DESC, created_at DESC").
		Find(&notices).Error

	return notices, total, err
}

// ==================== ExamRanking ====================

func UpsertExamRanking(ranking *model.ExamRanking) error {
	return database.DB.Where("exam_id = ? AND user_id = ?", ranking.ExamID, ranking.UserID).
		Assign(map[string]interface{}{
			"score":         ranking.Score,
			"obj_score":     ranking.ObjScore,
			"subj_score":    ranking.SubjScore,
			"rank_pos":      ranking.RankPos,
			"duration_used": ranking.DurationUsed,
			"accuracy":      ranking.Accuracy,
			"submit_time":   ranking.SubmitTime,
		}).FirstOrCreate(ranking).Error
}

func ListExamRankings(examID uint, page, pageSize int) ([]model.ExamRanking, int64, error) {
	var rankings []model.ExamRanking
	var total int64

	db := database.DB.Model(&model.ExamRanking{}).Where("exam_id = ?", examID)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset((page-1)*pageSize).Limit(pageSize).
		Order("rank_pos ASC").Find(&rankings).Error

	return rankings, total, err
}

func FindExamRankingByUser(examID, userID uint) (*model.ExamRanking, error) {
	var ranking model.ExamRanking
	err := database.DB.Where("exam_id = ? AND user_id = ?", examID, userID).
		First(&ranking).Error
	return &ranking, err
}

func BuildExamRankings(examID uint) error {
	// 获取所有已提交的记录
	var records []model.ExamRecord
	err := database.DB.Where("exam_id = ? AND status >= 1", examID).
		Order("score DESC, duration_used ASC").Find(&records).Error
	if err != nil {
		return err
	}

	// 删除旧排名
	database.DB.Where("exam_id = ?", examID).Delete(&model.ExamRanking{})

	// 获取考试信息用于计算正确率
	examInfo, err := FindExamByID(examID)
	if err != nil {
		return err
	}

	// 构建新排名
	for i, record := range records {
		// 计算正确率
		answers, _ := GetAnswerRecords(record.ID)
		correctCount := 0
		for _, a := range answers {
			if a.IsCorrect {
				correctCount++
			}
		}
		accuracy := float64(0)
		if len(answers) > 0 {
			accuracy = float64(correctCount) / float64(len(answers)) * 100
		}

		ranking := &model.ExamRanking{
			ExamID:       examID,
			UserID:       record.UserID,
			Score:        record.Score,
			ObjScore:     record.ObjScore,
			SubjScore:    record.SubjScore,
			RankPos:      i + 1,
			DurationUsed: record.DurationUsed,
			Accuracy:     accuracy,
			SubmitTime:   getTimeOrDefault(record.SubmitTime),
		}
		database.DB.Create(ranking)
	}

	_ = examInfo // used for potential future scoring logic
	return nil
}

func getTimeOrDefault(t *time.Time) time.Time {
	if t != nil {
		return *t
	}
	return time.Now()
}

// ==================== ExamFeedback ====================

func CreateExamFeedback(feedback *model.ExamFeedback) error {
	return database.DB.Create(feedback).Error
}

func FindExamFeedbackByUser(examID, userID uint) (*model.ExamFeedback, error) {
	var feedback model.ExamFeedback
	err := database.DB.Where("exam_id = ? AND user_id = ?", examID, userID).
		First(&feedback).Error
	return &feedback, err
}

func ListExamFeedbacks(examID uint, page, pageSize int) ([]model.ExamFeedback, int64, error) {
	var feedbacks []model.ExamFeedback
	var total int64

	db := database.DB.Model(&model.ExamFeedback{}).Where("exam_id = ?", examID)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset((page-1)*pageSize).Limit(pageSize).
		Order("created_at DESC").Find(&feedbacks).Error

	return feedbacks, total, err
}

// ==================== Exam (extended queries) ====================

func ListUpcomingExams(userID uint, req *dto.PageRequest) ([]model.Exam, int64, error) {
	var exams []model.Exam
	var total int64

	now := time.Now()
	db := database.DB.Model(&model.Exam{}).
		Where("status = ? AND start_time > ?", 1, now)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Preload("Paper").
		Offset(req.GetOffset()).Limit(req.GetLimit()).
		Order("start_time ASC").Find(&exams).Error

	return exams, total, err
}

func ListTemplatesPaged(req *dto.TemplateListRequest) ([]model.Paper, int64, error) {
	var papers []model.Paper
	var total int64

	db := database.DB.Model(&model.Paper{}).Where("is_template = ?", true)

	if req.Category != "" {
		db = db.Where("category = ?", req.Category)
	}
	if req.Keyword != "" {
		db = db.Where("title LIKE ?", "%"+req.Keyword+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset(req.GetOffset()).Limit(req.GetLimit()).
		Order("id DESC").Find(&papers).Error

	return papers, total, err
}
