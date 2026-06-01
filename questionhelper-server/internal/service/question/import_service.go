package question

import (
	"encoding/json"
	"fmt"

	"questionhelper-server/internal/model"
	questionRepo "questionhelper-server/internal/repository/question"
	"questionhelper-server/pkg/logger"
)

// ImportQuestions 导入题目
func ImportQuestions(creatorID uint, categoryID uint, visibility int8, data []byte) (int, error) {
	// 解析JSON格式的题目数据
	var importData []struct {
		Title      string `json:"title"`
		Content    string `json:"content"`
		Type       int8   `json:"type"`
		Difficulty int8   `json:"difficulty"`
		Answer     string `json:"answer"`
		Analysis   string `json:"analysis"`
		Options    []struct {
			Label     string `json:"label"`
			Content   string `json:"content"`
			IsCorrect bool   `json:"is_correct"`
		} `json:"options"`
	}

	if err := json.Unmarshal(data, &importData); err != nil {
		return 0, fmt.Errorf("解析导入数据失败: %w", err)
	}

	questions := make([]model.Question, 0, len(importData))
	for _, item := range importData {
		q := model.Question{
			Title:      item.Title,
			Content:    item.Content,
			Type:       item.Type,
			Difficulty: item.Difficulty,
			Answer:     item.Answer,
			Analysis:   item.Analysis,
			CategoryID: categoryID,
			Visibility: visibility,
			CreatorID:  creatorID,
			Status:     1,
		}
		questions = append(questions, q)
	}

	if err := questionRepo.BatchCreate(questions); err != nil {
		return 0, fmt.Errorf("批量创建题目失败: %w", err)
	}

	logger.Infof("导入 %d 道题目成功", len(questions))
	return len(questions), nil
}
