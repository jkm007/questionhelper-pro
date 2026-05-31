package wrong

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	wrongRepo "questionhelper-server/internal/repository/wrong"
	"questionhelper-server/pkg/logger"
)

// ListWrongQuestions 错题列表
func ListWrongQuestions(userID uint, req *dto.WrongListRequest) ([]dto.WrongQuestionInfo, int64, error) {
	wrongs, total, err := wrongRepo.List(userID, req)
	if err != nil {
		return nil, 0, fmt.Errorf("查询错题列表失败: %w", err)
	}

	list := make([]dto.WrongQuestionInfo, 0, len(wrongs))
	for _, w := range wrongs {
		list = append(list, toWrongQuestionInfo(&w))
	}
	return list, total, nil
}

// GetWrongQuestion 获取错题详情
func GetWrongQuestion(id, userID uint) (*dto.WrongQuestionInfo, error) {
	wrong, err := wrongRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("错题不存在")
		}
		return nil, fmt.Errorf("查询错题失败: %w", err)
	}

	if wrong.UserID != userID {
		return nil, errors.New("无权查看此错题")
	}

	info := toWrongQuestionInfo(wrong)
	return &info, nil
}

// ReviewWrongQuestion 复习错题
func ReviewWrongQuestion(id, userID uint, answer string) (bool, error) {
	wrong, err := wrongRepo.FindByID(id)
	if err != nil {
		return false, errors.New("错题不存在")
	}

	if wrong.UserID != userID {
		return false, errors.New("无权操作此错题")
	}

	// 检查答案是否正确
	isCorrect := wrong.Question.Answer == answer

	if isCorrect {
		// 答对了，标记为已掌握
		wrong.Mastered = true
		wrong.LastAnswer = answer
		if err := wrongRepo.Update(wrong); err != nil {
			return false, fmt.Errorf("更新错题失败: %w", err)
		}
		logger.Infof("用户 %d 复习错题 %d 成功", userID, id)
		return true, nil
	}

	// 答错了，更新错误次数
	wrong.WrongCount++
	wrong.LastAnswer = answer
	if err := wrongRepo.Update(wrong); err != nil {
		return false, fmt.Errorf("更新错题失败: %w", err)
	}

	return false, nil
}

// RemoveWrongQuestion 移除错题
func RemoveWrongQuestion(id, userID uint) error {
	wrong, err := wrongRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("错题不存在")
		}
		return fmt.Errorf("查询错题失败: %w", err)
	}

	if wrong.UserID != userID {
		return errors.New("无权删除此错题")
	}

	if err := wrongRepo.DeleteByID(id); err != nil {
		return fmt.Errorf("删除错题失败: %w", err)
	}

	logger.Infof("用户 %d 删除错题 %d", userID, id)
	return nil
}

// GetWrongAnalysis 错题分析
func GetWrongAnalysis(userID uint) (map[string]interface{}, error) {
	return wrongRepo.GetAnalysis(userID)
}

// toWrongQuestionInfo 转换为 WrongQuestionInfo DTO
func toWrongQuestionInfo(w *model.WrongQuestion) dto.WrongQuestionInfo {
	info := dto.WrongQuestionInfo{
		ID:         w.ID,
		UserID:     w.UserID,
		QuestionID: w.QuestionID,
		Source:     w.Source,
		SourceID:   w.SourceID,
		WrongCount: w.WrongCount,
		LastAnswer: w.LastAnswer,
		Mastered:   w.Mastered,
		CreatedAt:  w.CreatedAt,
	}

	if w.Question.ID > 0 {
		info.Question = dto.QuestionInfo{
			ID:         w.Question.ID,
			Title:      w.Question.Title,
			Content:    w.Question.Content,
			Type:       w.Question.Type,
			Difficulty: w.Question.Difficulty,
			Answer:     w.Question.Answer,
			Analysis:   w.Question.Analysis,
		}

		if len(w.Question.Options) > 0 {
			info.Question.Options = make([]dto.OptionInfo, 0, len(w.Question.Options))
			for _, opt := range w.Question.Options {
				info.Question.Options = append(info.Question.Options, dto.OptionInfo{
					ID:        opt.ID,
					Label:     opt.Label,
					Content:   opt.Content,
					IsCorrect: opt.IsCorrect,
				})
			}
		}
	}

	return info
}
