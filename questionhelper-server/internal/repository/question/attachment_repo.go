package question

import (
	"questionhelper-server/internal/model"
	"questionhelper-server/pkg/database"
)

// CreateAttachment 创建附件
func CreateAttachment(attachment *model.QuestionAttachment) error {
	return database.DB.Create(attachment).Error
}

// DeleteAttachment 删除附件
func DeleteAttachment(id uint) error {
	return database.DB.Delete(&model.QuestionAttachment{}, id).Error
}

// FindAttachmentsByQuestionID 获取题目附件列表
func FindAttachmentsByQuestionID(questionID uint) ([]model.QuestionAttachment, error) {
	var attachments []model.QuestionAttachment
	err := database.DB.Where("question_id = ?", questionID).
		Order("sort ASC, id ASC").
		Find(&attachments).Error
	return attachments, err
}

// FindAttachmentByID 获取附件详情
func FindAttachmentByID(id uint) (*model.QuestionAttachment, error) {
	var attachment model.QuestionAttachment
	err := database.DB.First(&attachment, id).Error
	return &attachment, err
}

// DeleteAttachmentsByQuestionID 删除题目所有附件
func DeleteAttachmentsByQuestionID(questionID uint) error {
	return database.DB.Where("question_id = ?", questionID).
		Delete(&model.QuestionAttachment{}).Error
}
