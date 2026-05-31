package question

import (
	"encoding/json"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	questionRepo "questionhelper-server/internal/repository/question"
	"questionhelper-server/pkg/logger"
)

// CreateVersion 创建版本记录
func CreateVersion(question *model.Question, changeLog string, creatorID uint) error {
	// 序列化选项
	optionsJSON := "[]"
	if len(question.Options) > 0 {
		data, _ := json.Marshal(question.Options)
		optionsJSON = string(data)
	}

	version := &model.QuestionVersion{
		QuestionID: question.ID,
		Version:    question.Version,
		Title:      question.Title,
		Content:    question.Content,
		Type:       question.Type,
		Difficulty: question.Difficulty,
		Answer:     question.Answer,
		Analysis:   question.Analysis,
		CategoryID: question.CategoryID,
		Options:    optionsJSON,
		ChangeLog:  changeLog,
		CreatorID:  creatorID,
	}

	if err := questionRepo.CreateVersion(version); err != nil {
		return fmt.Errorf("创建版本记录失败: %w", err)
	}

	logger.Infof("题目 %d 创建版本 %d 成功", question.ID, question.Version)
	return nil
}

// ListVersions 获取题目版本列表
func ListVersions(questionID uint) ([]dto.VersionInfo, error) {
	versions, err := questionRepo.FindVersionsByQuestionID(questionID)
	if err != nil {
		return nil, fmt.Errorf("查询版本列表失败: %w", err)
	}

	list := make([]dto.VersionInfo, 0, len(versions))
	for _, v := range versions {
		list = append(list, dto.VersionInfo{
			ID:         v.ID,
			QuestionID: v.QuestionID,
			Version:    v.Version,
			Title:      v.Title,
			ChangeLog:  v.ChangeLog,
			CreatorID:  v.CreatorID,
			CreatedAt:  v.CreatedAt,
		})
	}

	return list, nil
}

// GetVersionDetail 获取版本详情
func GetVersionDetail(id uint) (*dto.VersionDetail, error) {
	v, err := questionRepo.FindVersionByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("版本不存在")
		}
		return nil, fmt.Errorf("查询版本失败: %w", err)
	}

	// 解析选项
	var options []dto.OptionInfo
	if v.Options != "" {
		json.Unmarshal([]byte(v.Options), &options)
	}

	return &dto.VersionDetail{
		VersionInfo: dto.VersionInfo{
			ID:         v.ID,
			QuestionID: v.QuestionID,
			Version:    v.Version,
			Title:      v.Title,
			ChangeLog:  v.ChangeLog,
			CreatorID:  v.CreatorID,
			CreatedAt:  v.CreatedAt,
		},
		Content:    v.Content,
		Type:       v.Type,
		Difficulty: v.Difficulty,
		Answer:     v.Answer,
		Analysis:   v.Analysis,
		CategoryID: v.CategoryID,
		Options:    options,
	}, nil
}

// RollbackVersion 回滚到指定版本
func RollbackVersion(questionID, version int, userID uint) error {
	// 获取指定版本
	v, err := questionRepo.FindVersionByNumber(questionID, version)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("版本不存在")
		}
		return fmt.Errorf("查询版本失败: %w", err)
	}

	// 获取当前题目
	question, err := questionRepo.FindByID(uint(questionID))
	if err != nil {
		return fmt.Errorf("查询题目失败: %w", err)
	}

	// 保存当前版本
	if err := CreateVersion(question, fmt.Sprintf("回滚前自动保存(版本%d)", question.Version), userID); err != nil {
		return fmt.Errorf("保存当前版本失败: %w", err)
	}

	// 解析选项
	var options []model.Option
	if v.Options != "" {
		json.Unmarshal([]byte(v.Options), &options)
	}

	// 更新题目
	question.Title = v.Title
	question.Content = v.Content
	question.Type = v.Type
	question.Difficulty = v.Difficulty
	question.Answer = v.Answer
	question.Analysis = v.Analysis
	question.CategoryID = v.CategoryID
	question.Version = question.Version + 1

	if err := questionRepo.Update(question); err != nil {
		return fmt.Errorf("回滚题目失败: %w", err)
	}

	// 更新选项
	if len(options) > 0 {
		questionRepo.DeleteOptionsByQuestionID(question.ID)
		for i := range options {
			options[i].QuestionID = question.ID
		}
		questionRepo.CreateOptions(options)
	}

	logger.Infof("题目 %d 回滚到版本 %d 成功", questionID, version)
	return nil
}
