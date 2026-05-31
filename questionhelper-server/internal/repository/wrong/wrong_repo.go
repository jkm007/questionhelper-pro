package wrong

import (
	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	"questionhelper-server/pkg/database"
)

func FindByID(id uint) (*model.WrongQuestion, error) {
	var wrong model.WrongQuestion
	err := database.DB.Preload("Question").Preload("Question.Options").First(&wrong, id).Error
	return &wrong, err
}

func FindByUserAndQuestion(userID, questionID uint) (*model.WrongQuestion, error) {
	var wrong model.WrongQuestion
	err := database.DB.Where("user_id = ? AND question_id = ?", userID, questionID).First(&wrong).Error
	return &wrong, err
}

func Create(wrong *model.WrongQuestion) error {
	return database.DB.Create(wrong).Error
}

func Update(wrong *model.WrongQuestion) error {
	return database.DB.Save(wrong).Error
}

func DeleteByID(id uint) error {
	return database.DB.Delete(&model.WrongQuestion{}, id).Error
}

func List(userID uint, req *dto.WrongListRequest) ([]model.WrongQuestion, int64, error) {
	var wrongs []model.WrongQuestion
	var total int64

	db := database.DB.Model(&model.WrongQuestion{}).Where("user_id = ?", userID)

	if req.Mastered != nil {
		db = db.Where("mastered = ?", *req.Mastered)
	}
	if req.Source != nil {
		db = db.Where("source = ?", *req.Source)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Preload("Question").Preload("Question.Options").
		Offset(req.GetOffset()).Limit(req.GetLimit()).
		Order("updated_at DESC").Find(&wrongs).Error

	return wrongs, total, err
}

// GetAnalysis 错题分析
func GetAnalysis(userID uint) (map[string]interface{}, error) {
	var totalCount int64
	var masteredCount int64

	database.DB.Model(&model.WrongQuestion{}).Where("user_id = ?", userID).Count(&totalCount)
	database.DB.Model(&model.WrongQuestion{}).Where("user_id = ? AND mastered = ?", userID, true).Count(&masteredCount)

	// 按分类统计
	var byCategory []struct {
		CategoryID   uint   `json:"category_id"`
		CategoryName string `json:"category_name"`
		Count        int    `json:"count"`
	}
	database.DB.Model(&model.WrongQuestion{}).
		Select("questions.category_id, categories.name as category_name, COUNT(*) as count").
		Joins("JOIN questions ON questions.id = wrong_questions.question_id").
		Joins("LEFT JOIN categories ON categories.id = questions.category_id").
		Where("wrong_questions.user_id = ?", userID).
		Group("questions.category_id, categories.name").
		Scan(&byCategory)

	// 按题型统计
	var byType []struct {
		Type  int8 `json:"type"`
		Count int  `json:"count"`
	}
	database.DB.Model(&model.WrongQuestion{}).
		Select("questions.type, COUNT(*) as count").
		Joins("JOIN questions ON questions.id = wrong_questions.question_id").
		Where("wrong_questions.user_id = ?", userID).
		Group("questions.type").
		Scan(&byType)

	// 按难度统计
	var byDifficulty []struct {
		Difficulty int8 `json:"difficulty"`
		Count      int  `json:"count"`
	}
	database.DB.Model(&model.WrongQuestion{}).
		Select("questions.difficulty, COUNT(*) as count").
		Joins("JOIN questions ON questions.id = wrong_questions.question_id").
		Where("wrong_questions.user_id = ?", userID).
		Group("questions.difficulty").
		Scan(&byDifficulty)

	return map[string]interface{}{
		"total_count":    totalCount,
		"mastered_count": masteredCount,
		"by_category":    byCategory,
		"by_type":        byType,
		"by_difficulty":  byDifficulty,
	}, nil
}
