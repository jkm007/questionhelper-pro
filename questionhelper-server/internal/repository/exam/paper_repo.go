package exam

import (
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

func UpdatePaperQuestion(pq *model.PaperQuestion) error {
	return database.DB.Save(pq).Error
}
