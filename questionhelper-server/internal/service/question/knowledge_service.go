package question

import (
	"errors"
	"fmt"

	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	questionRepo "questionhelper-server/internal/repository/question"
	"questionhelper-server/pkg/logger"
)

// ListKnowledgePoints 知识点列表
func ListKnowledgePoints(categoryID *uint) ([]model.Knowledge, error) {
	if categoryID != nil {
		return questionRepo.FindKnowledgeByCategoryID(*categoryID)
	}
	return questionRepo.FindAllKnowledge()
}

// CreateKnowledgePoint 创建知识点
func CreateKnowledgePoint(req *dto.CreateKnowledgeRequest) error {
	knowledge := &model.Knowledge{
		CategoryID: req.CategoryID,
		Name:       req.Name,
	}

	if err := questionRepo.CreateKnowledge(knowledge); err != nil {
		return fmt.Errorf("创建知识点失败: %w", err)
	}

	logger.Infof("创建知识点成功: %s", req.Name)
	return nil
}

// UpdateKnowledgePoint 更新知识点
func UpdateKnowledgePoint(id uint, req *dto.UpdateKnowledgeRequest) error {
	knowledge, err := questionRepo.FindKnowledgeByID(id)
	if err != nil {
		return errors.New("知识点不存在")
	}

	if req.Name != "" {
		knowledge.Name = req.Name
	}
	if req.CategoryID > 0 {
		knowledge.CategoryID = req.CategoryID
	}

	if err := questionRepo.UpdateKnowledge(knowledge); err != nil {
		return fmt.Errorf("更新知识点失败: %w", err)
	}

	logger.Infof("更新知识点 %d 成功", id)
	return nil
}

// DeleteKnowledgePoint 删除知识点
func DeleteKnowledgePoint(id uint) error {
	if err := questionRepo.DeleteKnowledgeByID(id); err != nil {
		return fmt.Errorf("删除知识点失败: %w", err)
	}

	logger.Infof("删除知识点 %d 成功", id)
	return nil
}
