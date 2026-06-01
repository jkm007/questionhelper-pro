package question

import (
	"gorm.io/gorm"

	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	"questionhelper-server/pkg/database"
)

// ==================== Question ====================

func FindByID(id uint) (*model.Question, error) {
	var q model.Question
	err := database.DB.Preload("Options").Preload("Category").Preload("Creator").First(&q, id).Error
	return &q, err
}

func Create(q *model.Question) error {
	return database.DB.Create(q).Error
}

func Update(q *model.Question) error {
	return database.DB.Save(q).Error
}

func DeleteByID(id uint) error {
	return database.DB.Delete(&model.Question{}, id).Error
}

func List(req *dto.QuestionListRequest) ([]model.Question, int64, error) {
	var questions []model.Question
	var total int64

	db := database.DB.Model(&model.Question{})

	if req.Keyword != "" {
		db = db.Where("title LIKE ? OR content LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}
	if req.CategoryID != nil {
		db = db.Where("category_id = ?", *req.CategoryID)
	}
	if req.Type != nil {
		db = db.Where("type = ?", *req.Type)
	}
	if req.Difficulty != nil {
		db = db.Where("difficulty = ?", *req.Difficulty)
	}
	if req.Visibility != nil {
		db = db.Where("visibility = ?", *req.Visibility)
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}

	// 数据权限过滤：普通用户只能看到公开题目和自己创建的题目
	if req.UserID != nil {
		db = db.Where("(visibility = 1 OR creator_id = ?)", *req.UserID)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Preload("Options").Preload("Category").Preload("Creator").
		Offset(req.GetOffset()).Limit(req.GetLimit()).
		Order("id DESC").Find(&questions).Error

	return questions, total, err
}

func ListByCreator(creatorID uint, req *dto.QuestionListRequest) ([]model.Question, int64, error) {
	var questions []model.Question
	var total int64

	db := database.DB.Model(&model.Question{}).Where("creator_id = ?", creatorID)

	if req.Keyword != "" {
		db = db.Where("title LIKE ?", "%"+req.Keyword+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Preload("Options").Preload("Category").
		Offset(req.GetOffset()).Limit(req.GetLimit()).
		Order("id DESC").Find(&questions).Error

	return questions, total, err
}

// ==================== Option ====================

func DeleteOptionsByQuestionID(questionID uint) error {
	return database.DB.Where("question_id = ?", questionID).Delete(&model.Option{}).Error
}

func CreateOptions(options []model.Option) error {
	if len(options) == 0 {
		return nil
	}
	return database.DB.Create(&options).Error
}

// ==================== Like ====================

func IncrementLikeCount(questionID uint) error {
	return database.DB.Model(&model.Question{}).Where("id = ?", questionID).
		UpdateColumn("like_count", gorm.Expr("like_count + 1")).Error
}

func DecrementLikeCount(questionID uint) error {
	return database.DB.Model(&model.Question{}).Where("id = ?", questionID).
		UpdateColumn("like_count", gorm.Expr("GREATEST(like_count - 1, 0)")).Error
}

// IncrementFavoriteCount 原子递增收藏数
func IncrementFavoriteCount(questionID uint) error {
	return database.DB.Model(&model.Question{}).Where("id = ?", questionID).
		UpdateColumn("favorite_count", gorm.Expr("favorite_count + 1")).Error
}

// DecrementFavoriteCount 原子递减收藏数
func DecrementFavoriteCount(questionID uint) error {
	return database.DB.Model(&model.Question{}).Where("id = ?", questionID).
		UpdateColumn("favorite_count", gorm.Expr("GREATEST(favorite_count - 1, 0)")).Error
}

// ==================== Import/Export ====================

func BatchCreate(questions []model.Question) error {
	if len(questions) == 0 {
		return nil
	}
	return database.DB.Create(&questions).Error
}

func FindAllForExport(categoryID *uint) ([]model.Question, error) {
	var questions []model.Question
	db := database.DB.Preload("Options").Preload("Category")
	if categoryID != nil {
		db = db.Where("category_id = ?", *categoryID)
	}
	err := db.Order("id ASC").Find(&questions).Error
	return questions, err
}
