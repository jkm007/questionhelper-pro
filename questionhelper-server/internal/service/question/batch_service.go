package question

import (
	"fmt"

	"questionhelper-server/internal/dto"
	questionRepo "questionhelper-server/internal/repository/question"
	"questionhelper-server/pkg/logger"
)

// BatchPublish 批量发布题目
func BatchPublish(req *dto.BatchQuestionRequest) (*dto.BatchResult, error) {
	result := &dto.BatchResult{
		Total: len(req.IDs),
	}

	for _, id := range req.IDs {
		// 状态校验：只有审核通过(status=3,待审核可视为已通过审核流程)的题目才能发布
		question, err := questionRepo.FindByID(id)
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, dto.BatchError{
				ID:     id,
				Reason: fmt.Sprintf("题目不存在: %v", err),
			})
			continue
		}

		// 只有草稿(0)或待审核(3)状态的题目才允许发布
		if question.Status != 0 && question.Status != 3 {
			result.Failed++
			result.Errors = append(result.Errors, dto.BatchError{
				ID:     id,
				Reason: fmt.Sprintf("题目状态(%d)不允许发布，只有草稿或待审核状态的题目才能发布", question.Status),
			})
			continue
		}

		err = UpdateQuestionStatus(id, 1) // 1=已发布
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, dto.BatchError{
				ID:     id,
				Reason: err.Error(),
			})
		} else {
			result.Success++
		}
	}

	logger.Infof("批量发布完成: %d/%d", result.Success, result.Total)
	return result, nil
}

// BatchArchive 批量归档题目
func BatchArchive(req *dto.BatchQuestionRequest) (*dto.BatchResult, error) {
	result := &dto.BatchResult{
		Total: len(req.IDs),
	}

	for _, id := range req.IDs {
		err := UpdateQuestionStatus(id, 2) // 2=已归档
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, dto.BatchError{
				ID:     id,
				Reason: err.Error(),
			})
		} else {
			result.Success++
		}
	}

	logger.Infof("批量归档完成: %d/%d", result.Success, result.Total)
	return result, nil
}

// BatchDelete 批量删除题目
func BatchDelete(req *dto.BatchQuestionRequest) (*dto.BatchResult, error) {
	result := &dto.BatchResult{
		Total: len(req.IDs),
	}

	for _, id := range req.IDs {
		err := questionRepo.DeleteByID(id)
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, dto.BatchError{
				ID:     id,
				Reason: err.Error(),
			})
		} else {
			result.Success++
		}
	}

	logger.Infof("批量删除完成: %d/%d", result.Success, result.Total)
	return result, nil
}

// BatchMoveCategory 批量移动分类
func BatchMoveCategory(req *dto.BatchMoveRequest) (*dto.BatchResult, error) {
	result := &dto.BatchResult{
		Total: len(req.IDs),
	}

	for _, id := range req.IDs {
		question, err := questionRepo.FindByID(id)
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, dto.BatchError{
				ID:     id,
				Reason: fmt.Sprintf("题目不存在: %v", err),
			})
			continue
		}

		question.CategoryID = req.CategoryID
		if err := questionRepo.Update(question); err != nil {
			result.Failed++
			result.Errors = append(result.Errors, dto.BatchError{
				ID:     id,
				Reason: err.Error(),
			})
		} else {
			result.Success++
		}
	}

	logger.Infof("批量移动分类完成: %d/%d", result.Success, result.Total)
	return result, nil
}
