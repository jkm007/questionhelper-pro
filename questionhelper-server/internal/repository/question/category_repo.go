package question

import (
	"gorm.io/gorm"

	"questionhelper-server/internal/model"
	"questionhelper-server/pkg/database"
)

// ==================== Category ====================

func FindCategoryByID(id uint) (*model.Category, error) {
	var cat model.Category
	err := database.DB.First(&cat, id).Error
	return &cat, err
}

func FindCategoryTree() ([]model.Category, error) {
	var categories []model.Category
	err := database.DB.Where("parent_id IS NULL").
		Order("sort ASC, id ASC").
		Preload("Children", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort ASC, id ASC")
		}).
		Find(&categories).Error
	return categories, err
}

func FindAllCategories() ([]model.Category, error) {
	var categories []model.Category
	err := database.DB.Order("sort ASC, id ASC").Find(&categories).Error
	return categories, err
}

func CreateCategory(cat *model.Category) error {
	return database.DB.Create(cat).Error
}

func UpdateCategory(cat *model.Category) error {
	return database.DB.Save(cat).Error
}

func DeleteCategoryByID(id uint) error {
	return database.DB.Delete(&model.Category{}, id).Error
}
