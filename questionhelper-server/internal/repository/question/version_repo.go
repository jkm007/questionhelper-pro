package question

import (
	"questionhelper-server/internal/model"
	"questionhelper-server/pkg/database"
)

// CreateVersion 创建版本记录
func CreateVersion(version *model.QuestionVersion) error {
	return database.DB.Create(version).Error
}

// FindVersionsByQuestionID 获取题目版本列表
func FindVersionsByQuestionID(questionID uint) ([]model.QuestionVersion, error) {
	var versions []model.QuestionVersion
	err := database.DB.Where("question_id = ?", questionID).
		Order("version DESC").
		Find(&versions).Error
	return versions, err
}

// FindVersionByID 获取版本详情
func FindVersionByID(id uint) (*model.QuestionVersion, error) {
	var version model.QuestionVersion
	err := database.DB.First(&version, id).Error
	return &version, err
}

// FindVersionByNumber 获取指定版本
func FindVersionByNumber(questionID, version int) (*model.QuestionVersion, error) {
	var v model.QuestionVersion
	err := database.DB.Where("question_id = ? AND version = ?", questionID, version).
		First(&v).Error
	return &v, err
}

// GetLatestVersion 获取最新版本号
func GetLatestVersion(questionID uint) (int, error) {
	var version model.QuestionVersion
	err := database.DB.Where("question_id = ?", questionID).
		Order("version DESC").
		First(&version).Error
	if err != nil {
		return 0, err
	}
	return version.Version, nil
}
