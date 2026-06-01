package sensitive

import (
	"errors"
	"fmt"
	"strings"

	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	questionRepo "questionhelper-server/internal/repository/question"
	"questionhelper-server/pkg/logger"
)

// ListSensitiveWords 敏感词列表
func ListSensitiveWords(req *dto.PageRequest) ([]model.SensitiveWord, int64, error) {
	return questionRepo.ListSensitiveWords("", req.Page, req.PageSize)
}

// CreateSensitiveWord 创建敏感词
func CreateSensitiveWord(req *dto.CreateSensitiveWordRequest) error {
	word := &model.SensitiveWord{
		Word: req.Word,
	}

	if err := questionRepo.CreateSensitiveWord(word); err != nil {
		return fmt.Errorf("创建敏感词失败: %w", err)
	}

	logger.Infof("创建敏感词成功: %s", req.Word)
	return nil
}

// DeleteSensitiveWord 删除敏感词
func DeleteSensitiveWord(id uint) error {
	if err := questionRepo.DeleteSensitiveWordByID(id); err != nil {
		return fmt.Errorf("删除敏感词失败: %w", err)
	}

	logger.Infof("删除敏感词 %d 成功", id)
	return nil
}

// ImportSensitiveWords 导入敏感词
func ImportSensitiveWords(data []byte) (int, error) {
	lines := strings.Split(string(data), "\n")
	words := make([]model.SensitiveWord, 0)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			words = append(words, model.SensitiveWord{Word: line})
		}
	}

	if len(words) == 0 {
		return 0, errors.New("没有有效的敏感词")
	}

	if err := questionRepo.BatchCreateSensitiveWords(words); err != nil {
		return 0, fmt.Errorf("批量创建敏感词失败: %w", err)
	}

	logger.Infof("导入 %d 个敏感词成功", len(words))
	return len(words), nil
}

// TestSensitiveWord 测试敏感词
func TestSensitiveWord(content string) bool {
	return questionRepo.HasSensitiveWord(content)
}
