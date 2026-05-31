package question

import (
	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	"questionhelper-server/pkg/database"
)

// ==================== QuestionNote ====================

func FindNoteByID(id uint) (*model.QuestionNote, error) {
	var note model.QuestionNote
	err := database.DB.First(&note, id).Error
	return &note, err
}

func CreateNote(note *model.QuestionNote) error {
	return database.DB.Create(note).Error
}

func UpdateNote(note *model.QuestionNote) error {
	return database.DB.Save(note).Error
}

func DeleteNote(id uint) error {
	return database.DB.Delete(&model.QuestionNote{}, id).Error
}

func ListNotes(questionID uint, req *dto.NoteListRequest) ([]model.QuestionNote, int64, error) {
	var notes []model.QuestionNote
	var total int64

	db := database.DB.Model(&model.QuestionNote{}).Where("question_id = ?", questionID)

	if req.IsPublic != nil {
		db = db.Where("is_public = ?", *req.IsPublic)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset(req.GetOffset()).Limit(req.GetLimit()).
		Order("created_at DESC").Find(&notes).Error

	return notes, total, err
}

func FindNoteByUserAndQuestion(userID, questionID uint) (*model.QuestionNote, error) {
	var note model.QuestionNote
	err := database.DB.Where("user_id = ? AND question_id = ?", userID, questionID).First(&note).Error
	return &note, err
}

// ==================== QuestionDifficultyRating ====================

func FindDifficultyRating(userID, questionID uint) (*model.QuestionDifficultyRating, error) {
	var rating model.QuestionDifficultyRating
	err := database.DB.Where("user_id = ? AND question_id = ?", userID, questionID).First(&rating).Error
	return &rating, err
}

func CreateDifficultyRating(rating *model.QuestionDifficultyRating) error {
	return database.DB.Create(rating).Error
}

func UpdateDifficultyRating(rating *model.QuestionDifficultyRating) error {
	return database.DB.Save(rating).Error
}

func GetDifficultyRatingSummary(questionID uint) (float64, int, error) {
	var avg float64
	var count int64

	err := database.DB.Model(&model.QuestionDifficultyRating{}).
		Where("question_id = ?", questionID).
		Select("COALESCE(AVG(rating), 0)").Scan(&avg).Error
	if err != nil {
		return 0, 0, err
	}

	database.DB.Model(&model.QuestionDifficultyRating{}).
		Where("question_id = ?", questionID).Count(&count)

	return avg, int(count), nil
}

// ==================== QuestionQualityRating ====================

func FindQualityRating(userID, questionID uint) (*model.QuestionQualityRating, error) {
	var rating model.QuestionQualityRating
	err := database.DB.Where("user_id = ? AND question_id = ?", userID, questionID).First(&rating).Error
	return &rating, err
}

func CreateQualityRating(rating *model.QuestionQualityRating) error {
	return database.DB.Create(rating).Error
}

func UpdateQualityRating(rating *model.QuestionQualityRating) error {
	return database.DB.Save(rating).Error
}

func GetQualityRatingSummary(questionID uint) (float64, int, error) {
	var avg float64
	var count int64

	err := database.DB.Model(&model.QuestionQualityRating{}).
		Where("question_id = ?", questionID).
		Select("COALESCE(AVG(score), 0)").Scan(&avg).Error
	if err != nil {
		return 0, 0, err
	}

	database.DB.Model(&model.QuestionQualityRating{}).
		Where("question_id = ?", questionID).Count(&count)

	return avg, int(count), nil
}

// ==================== QuestionCorrection ====================

func FindCorrectionByID(id uint) (*model.QuestionCorrection, error) {
	var correction model.QuestionCorrection
	err := database.DB.First(&correction, id).Error
	return &correction, err
}

func CreateCorrection(correction *model.QuestionCorrection) error {
	return database.DB.Create(correction).Error
}

func UpdateCorrection(correction *model.QuestionCorrection) error {
	return database.DB.Save(correction).Error
}

func ListCorrections(questionID uint, req *dto.CorrectionListRequest) ([]model.QuestionCorrection, int64, error) {
	var corrections []model.QuestionCorrection
	var total int64

	db := database.DB.Model(&model.QuestionCorrection{}).Where("question_id = ?", questionID)

	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset(req.GetOffset()).Limit(req.GetLimit()).
		Order("created_at DESC").Find(&corrections).Error

	return corrections, total, err
}
