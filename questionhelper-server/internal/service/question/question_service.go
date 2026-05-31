package question

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"gorm.io/gorm"

	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	questionRepo "questionhelper-server/internal/repository/question"
	"questionhelper-server/pkg/logger"
)

// ListQuestions 题目列表
func ListQuestions(req *dto.QuestionListRequest) ([]dto.QuestionInfo, int64, error) {
	questions, total, err := questionRepo.List(req)
	if err != nil {
		return nil, 0, fmt.Errorf("查询题目列表失败: %w", err)
	}

	list := make([]dto.QuestionInfo, 0, len(questions))
	for _, q := range questions {
		list = append(list, toQuestionInfo(&q))
	}
	return list, total, nil
}

// GetQuestion 获取题目详情
func GetQuestion(id uint) (*dto.QuestionInfo, error) {
	q, err := questionRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("题目不存在")
		}
		return nil, fmt.Errorf("查询题目失败: %w", err)
	}

	info := toQuestionInfo(q)
	return &info, nil
}

// SearchQuestions 搜索题目
func SearchQuestions(keyword string, req *dto.PageRequest) ([]dto.QuestionInfo, int64, error) {
	searchReq := &dto.QuestionListRequest{
		PageRequest: *req,
		Keyword:     keyword,
	}
	return ListQuestions(searchReq)
}

// CreateQuestion 创建题目
func CreateQuestion(creatorID uint, req *dto.CreateQuestionRequest) error {
	q := &model.Question{
		Title:      req.Title,
		Content:    req.Content,
		Type:       req.Type,
		Difficulty: req.Difficulty,
		Answer:     req.Answer,
		Analysis:   req.Analysis,
		CategoryID: req.CategoryID,
		Visibility: req.Visibility,
		CreatorID:  creatorID,
		Status:     1,
	}

	// 处理知识点IDs
	if len(req.KnowledgeIDs) > 0 {
		ids := make([]string, 0, len(req.KnowledgeIDs))
		for _, id := range req.KnowledgeIDs {
			ids = append(ids, strconv.FormatUint(uint64(id), 10))
		}
		knowledgeJSON, _ := json.Marshal(ids)
		q.KnowledgeIDs = string(knowledgeJSON)
	}

	if err := questionRepo.Create(q); err != nil {
		return fmt.Errorf("创建题目失败: %w", err)
	}

	// 创建选项
	if len(req.Options) > 0 {
		options := make([]model.Option, 0, len(req.Options))
		for i, opt := range req.Options {
			options = append(options, model.Option{
				QuestionID: q.ID,
				Label:      opt.Label,
				Content:    opt.Content,
				IsCorrect:  opt.IsCorrect,
				Sort:       i,
			})
		}
		if err := questionRepo.CreateOptions(options); err != nil {
			return fmt.Errorf("创建选项失败: %w", err)
		}
	}

	logger.Infof("创建题目成功: %d", q.ID)
	return nil
}

// UpdateQuestion 更新题目
func UpdateQuestion(id uint, req *dto.UpdateQuestionRequest) error {
	q, err := questionRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("题目不存在")
		}
		return fmt.Errorf("查询题目失败: %w", err)
	}

	if req.Title != "" {
		q.Title = req.Title
	}
	if req.Content != "" {
		q.Content = req.Content
	}
	if req.Type != 0 {
		q.Type = req.Type
	}
	if req.Difficulty != 0 {
		q.Difficulty = req.Difficulty
	}
	if req.Answer != "" {
		q.Answer = req.Answer
	}
	if req.Analysis != "" {
		q.Analysis = req.Analysis
	}
	if req.CategoryID != nil {
		q.CategoryID = *req.CategoryID
	}
	if req.Visibility != 0 {
		q.Visibility = req.Visibility
	}

	// 处理知识点IDs
	if req.KnowledgeIDs != nil {
		ids := make([]string, 0, len(req.KnowledgeIDs))
		for _, id := range req.KnowledgeIDs {
			ids = append(ids, strconv.FormatUint(uint64(id), 10))
		}
		knowledgeJSON, _ := json.Marshal(ids)
		q.KnowledgeIDs = string(knowledgeJSON)
	}

	if err := questionRepo.Update(q); err != nil {
		return fmt.Errorf("更新题目失败: %w", err)
	}

	// 更新选项
	if req.Options != nil {
		// 删除旧选项
		if err := questionRepo.DeleteOptionsByQuestionID(id); err != nil {
			return fmt.Errorf("删除旧选项失败: %w", err)
		}
		// 创建新选项
		options := make([]model.Option, 0, len(req.Options))
		for i, opt := range req.Options {
			options = append(options, model.Option{
				QuestionID: id,
				Label:      opt.Label,
				Content:    opt.Content,
				IsCorrect:  opt.IsCorrect,
				Sort:       i,
			})
		}
		if err := questionRepo.CreateOptions(options); err != nil {
			return fmt.Errorf("创建新选项失败: %w", err)
		}
	}

	logger.Infof("更新题目 %d 成功", id)
	return nil
}

// DeleteQuestion 删除题目
func DeleteQuestion(id uint) error {
	_, err := questionRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("题目不存在")
		}
		return fmt.Errorf("查询题目失败: %w", err)
	}

	if err := questionRepo.DeleteByID(id); err != nil {
		return fmt.Errorf("删除题目失败: %w", err)
	}

	logger.Infof("删除题目 %d 成功", id)
	return nil
}

// UpdateQuestionStatus 更新题目状态
func UpdateQuestionStatus(id uint, status int8) error {
	q, err := questionRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("题目不存在")
		}
		return fmt.Errorf("查询题目失败: %w", err)
	}

	q.Status = status
	if err := questionRepo.Update(q); err != nil {
		return fmt.Errorf("更新题目状态失败: %w", err)
	}

	logger.Infof("更新题目 %d 状态为 %d", id, status)
	return nil
}

// ListCategories 分类列表
func ListCategories() ([]dto.CategoryInfo, error) {
	categories, err := questionRepo.FindAllCategories()
	if err != nil {
		return nil, fmt.Errorf("查询分类失败: %w", err)
	}

	list := make([]dto.CategoryInfo, 0, len(categories))
	for _, c := range categories {
		list = append(list, toCategoryInfo(&c))
	}
	return list, nil
}

// GetCategoryTree 分类树
func GetCategoryTree() ([]dto.CategoryInfo, error) {
	categories, err := questionRepo.FindCategoryTree()
	if err != nil {
		return nil, fmt.Errorf("查询分类树失败: %w", err)
	}

	tree := make([]dto.CategoryInfo, 0, len(categories))
	for _, c := range categories {
		tree = append(tree, toCategoryInfo(&c))
	}
	return tree, nil
}

// ListKnowledgePoints 知识点列表
func ListKnowledgePoints(categoryID *uint) ([]model.Knowledge, error) {
	if categoryID != nil {
		return questionRepo.FindKnowledgeByCategoryID(*categoryID)
	}
	return questionRepo.FindAllKnowledge()
}

// LikeQuestion 点赞题目
func LikeQuestion(questionID uint) error {
	_, err := questionRepo.FindByID(questionID)
	if err != nil {
		return errors.New("题目不存在")
	}
	return questionRepo.IncrementLikeCount(questionID)
}

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

// ExportQuestions 导出题目
func ExportQuestions(categoryID *uint) ([]dto.QuestionInfo, error) {
	questions, err := questionRepo.FindAllForExport(categoryID)
	if err != nil {
		return nil, fmt.Errorf("查询导出题目失败: %w", err)
	}

	list := make([]dto.QuestionInfo, 0, len(questions))
	for _, q := range questions {
		list = append(list, toQuestionInfo(&q))
	}
	return list, nil
}

// toQuestionInfo 转换为 QuestionInfo DTO
func toQuestionInfo(q *model.Question) dto.QuestionInfo {
	info := dto.QuestionInfo{
		ID:         q.ID,
		Title:      q.Title,
		Content:    q.Content,
		Type:       q.Type,
		Difficulty: q.Difficulty,
		Answer:     q.Answer,
		Analysis:   q.Analysis,
		CategoryID: q.CategoryID,
		Visibility: q.Visibility,
		CreatorID:  q.CreatorID,
		Status:     q.Status,
		ViewCount:  q.ViewCount,
		LikeCount:  q.LikeCount,
	}

	if q.Category.ID > 0 {
		info.Category = dto.CategoryInfo{
			ID:   q.Category.ID,
			Name: q.Category.Name,
		}
	}

	if q.Creator.ID > 0 {
		info.CreatorName = q.Creator.Nickname
	}

	if len(q.Options) > 0 {
		info.Options = make([]dto.OptionInfo, 0, len(q.Options))
		for _, opt := range q.Options {
			info.Options = append(info.Options, dto.OptionInfo{
				ID:        opt.ID,
				Label:     opt.Label,
				Content:   opt.Content,
				IsCorrect: opt.IsCorrect,
			})
		}
	}

	return info
}

// toCategoryInfo 转换为 CategoryInfo DTO
func toCategoryInfo(c *model.Category) dto.CategoryInfo {
	info := dto.CategoryInfo{
		ID:       c.ID,
		ParentID: c.ParentID,
		Name:     c.Name,
	}
	if len(c.Children) > 0 {
		info.Children = make([]dto.CategoryInfo, 0, len(c.Children))
		for _, child := range c.Children {
			info.Children = append(info.Children, toCategoryInfo(&child))
		}
	}
	return info
}
