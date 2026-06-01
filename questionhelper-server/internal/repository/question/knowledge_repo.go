package question

import (
	"questionhelper-server/internal/model"
	"questionhelper-server/pkg/database"
)

// ==================== Knowledge ====================

func FindKnowledgeByID(id uint) (*model.Knowledge, error) {
	var k model.Knowledge
	err := database.DB.First(&k, id).Error
	return &k, err
}

func FindKnowledgeByCategoryID(categoryID uint) ([]model.Knowledge, error) {
	var knowledge []model.Knowledge
	err := database.DB.Where("category_id = ?", categoryID).
		Order("sort ASC, id ASC").Find(&knowledge).Error
	return knowledge, err
}

func FindAllKnowledge() ([]model.Knowledge, error) {
	var knowledge []model.Knowledge
	err := database.DB.Order("category_id ASC, sort ASC").Find(&knowledge).Error
	return knowledge, err
}

func CreateKnowledge(k *model.Knowledge) error {
	return database.DB.Create(k).Error
}

func UpdateKnowledge(k *model.Knowledge) error {
	return database.DB.Save(k).Error
}

func DeleteKnowledgeByID(id uint) error {
	return database.DB.Delete(&model.Knowledge{}, id).Error
}
