package question

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

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

// ListQuestionsWithPermission 带数据权限过滤的题目列表
func ListQuestionsWithPermission(userID uint, req *dto.QuestionListRequest) ([]dto.QuestionInfo, int64, error) {
	req.UserID = &userID
	return ListQuestions(req)
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

	// 创建版本快照
	version := &model.QuestionVersion{
		QuestionID: q.ID,
		Title:      q.Title,
		Content:    q.Content,
		Answer:     q.Answer,
		Analysis:   q.Analysis,
		Version:    q.Version + 1,
		CreatorID:  q.CreatorID,
	}
	if err := questionRepo.CreateVersion(version); err != nil {
		logger.Errorf("创建版本快照失败: %v", err)
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

// validStatusTransitions 定义合法的状态流转
// 0=草稿, 1=待审核, 2=审核通过, 3=已发布, 4=审核拒绝, 5=已下架
var validStatusTransitions = map[int8][]int8{
	0: {1, 3}, // 草稿 -> 待审核, 已发布
	1: {2, 4}, // 待审核 -> 审核通过, 审核拒绝
	2: {3},    // 审核通过 -> 已发布
	3: {5},    // 已发布 -> 已下架
	4: {0},    // 审核拒绝 -> 草稿
	5: {0},    // 已下架 -> 草稿
}

// isValidStatusTransition 校验状态流转是否合法
func isValidStatusTransition(from, to int8) bool {
	allowed, ok := validStatusTransitions[from]
	if !ok {
		return false
	}
	for _, s := range allowed {
		if s == to {
			return true
		}
	}
	return false
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

	// 状态流转校验
	if !isValidStatusTransition(q.Status, status) {
		return fmt.Errorf("不允许从状态 %d 变更为 %d", q.Status, status)
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
func LikeQuestion(userID, questionID uint) error {
	_, err := questionRepo.FindByID(questionID)
	if err != nil {
		return errors.New("题目不存在")
	}

	// 检查是否已点赞，防止重复
	if questionRepo.IsLiked(userID, questionID) {
		return errors.New("已经点赞过此题目")
	}

	// 记录点赞
	like := &model.QuestionLike{
		UserID:     userID,
		QuestionID: questionID,
	}
	if err := questionRepo.AddLike(like); err != nil {
		return fmt.Errorf("点赞记录失败: %w", err)
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

// ==================== 分类管理 ====================

// CreateCategory 创建分类
func CreateCategory(req *dto.CreateCategoryRequest) error {
	category := &model.Category{
		ParentID: req.ParentID,
		Name:     req.Name,
		Sort:     req.Sort,
	}

	if err := questionRepo.CreateCategory(category); err != nil {
		return fmt.Errorf("创建分类失败: %w", err)
	}

	logger.Infof("创建分类成功: %s", req.Name)
	return nil
}

// UpdateCategory 更新分类
func UpdateCategory(id uint, req *dto.UpdateCategoryRequest) error {
	category, err := questionRepo.FindCategoryByID(id)
	if err != nil {
		return errors.New("分类不存在")
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.ParentID != nil {
		category.ParentID = req.ParentID
	}
	category.Sort = req.Sort

	if err := questionRepo.UpdateCategory(category); err != nil {
		return fmt.Errorf("更新分类失败: %w", err)
	}

	logger.Infof("更新分类 %d 成功", id)
	return nil
}

// DeleteCategory 删除分类
func DeleteCategory(id uint) error {
	if err := questionRepo.DeleteCategoryByID(id); err != nil {
		return fmt.Errorf("删除分类失败: %w", err)
	}

	logger.Infof("删除分类 %d 成功", id)
	return nil
}

// ==================== 知识点管理 ====================

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

// ==================== 内容审核 ====================

// ListPendingReviews 待审核列表
func ListPendingReviews(req *dto.PageRequest) ([]dto.ReviewInfo, int64, error) {
	status := int8(0) // 待审核
	reviews, total, err := questionRepo.ListReviews(&status, req.Page, req.PageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("查询待审核列表失败: %w", err)
	}

	list := make([]dto.ReviewInfo, 0, len(reviews))
	for _, r := range reviews {
		reviewerID := uint(0)
		if r.ReviewerID != nil {
			reviewerID = *r.ReviewerID
		}
		list = append(list, dto.ReviewInfo{
			ID:         r.ID,
			QuestionID: r.QuestionID,
			ReviewerID: reviewerID,
			Status:     r.Status,
			Reason:     r.Opinion,
			CreatedAt:  r.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return list, total, nil
}

// GetReviewDetail 审核详情
func GetReviewDetail(id uint) (*dto.ReviewInfo, error) {
	review, err := questionRepo.FindReviewByID(id)
	if err != nil {
		return nil, errors.New("审核记录不存在")
	}

	reviewerID := uint(0)
	if review.ReviewerID != nil {
		reviewerID = *review.ReviewerID
	}

	return &dto.ReviewInfo{
		ID:         review.ID,
		QuestionID: review.QuestionID,
		ReviewerID: reviewerID,
		Status:     review.Status,
		Reason:     review.Opinion,
		CreatedAt:  review.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// ApproveReview 审核通过
func ApproveReview(id, reviewerID uint) error {
	review, err := questionRepo.FindReviewByID(id)
	if err != nil {
		return errors.New("审核记录不存在")
	}

	review.Status = 1 // 通过
	review.ReviewerID = &reviewerID
	if err := questionRepo.UpdateReview(review); err != nil {
		return fmt.Errorf("更新审核状态失败: %w", err)
	}

	// 更新题目状态为已发布
	q, err := questionRepo.FindByID(review.QuestionID)
	if err == nil {
		q.Status = 3
		questionRepo.Update(q)
	}

	logger.Infof("审核 %d 通过", id)
	return nil
}

// RejectReview 审核拒绝
func RejectReview(id, reviewerID uint, reason string) error {
	review, err := questionRepo.FindReviewByID(id)
	if err != nil {
		return errors.New("审核记录不存在")
	}

	review.Status = 2 // 拒绝
	review.ReviewerID = &reviewerID
	review.Opinion = reason
	if err := questionRepo.UpdateReview(review); err != nil {
		return fmt.Errorf("更新审核状态失败: %w", err)
	}

	// 更新题目状态为草稿
	q, err := questionRepo.FindByID(review.QuestionID)
	if err == nil {
		q.Status = 0
		questionRepo.Update(q)
	}

	logger.Infof("审核 %d 拒绝，原因: %s", id, reason)
	return nil
}

// ==================== 敏感词管理 ====================

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
