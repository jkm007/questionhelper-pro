package question

import (
	"questionhelper-server/internal/model"
	"questionhelper-server/pkg/database"
)

// CreateSensitiveWord 创建敏感词
func CreateSensitiveWord(word *model.SensitiveWord) error {
	return database.DB.Create(word).Error
}

// UpdateSensitiveWord 更新敏感词
func UpdateSensitiveWord(word *model.SensitiveWord) error {
	return database.DB.Save(word).Error
}

// DeleteSensitiveWord 删除敏感词
func DeleteSensitiveWord(id uint) error {
	return database.DB.Delete(&model.SensitiveWord{}, id).Error
}

// FindAllSensitiveWords 获取所有启用的敏感词
func FindAllSensitiveWords() ([]model.SensitiveWord, error) {
	var words []model.SensitiveWord
	err := database.DB.Where("status = 1").Find(&words).Error
	return words, err
}

// ListSensitiveWords 敏感词列表
func ListSensitiveWords(category string, page, pageSize int) ([]model.SensitiveWord, int64, error) {
	var words []model.SensitiveWord
	var total int64

	db := database.DB.Model(&model.SensitiveWord{})
	if category != "" {
		db = db.Where("category = ?", category)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset((page - 1) * pageSize).
		Limit(pageSize).
		Order("id ASC").
		Find(&words).Error

	return words, total, err
}

// BatchCreateSensitiveWords 批量创建敏感词
func BatchCreateSensitiveWords(words []model.SensitiveWord) error {
	return database.DB.CreateInBatches(words, 100).Error
}

// DeleteSensitiveWordByID 删除敏感词
func DeleteSensitiveWordByID(id uint) error {
	return database.DB.Delete(&model.SensitiveWord{}, id).Error
}

// HasSensitiveWord 检查内容是否包含敏感词
func HasSensitiveWord(content string) bool {
	var count int64
	database.DB.Model(&model.SensitiveWord{}).Where("status = 1 AND ? LIKE CONCAT('%', word, '%')", content).Count(&count)
	return count > 0
}
