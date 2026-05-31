package exam

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	examRepo "questionhelper-server/internal/repository/exam"
	"questionhelper-server/pkg/database"
	"questionhelper-server/pkg/logger"
)

// PreviewPaper 试卷预览(按题型分组)
func PreviewPaper(paperID uint) (*dto.PaperPreviewResponse, error) {
	paper, err := examRepo.FindPaperByID(paperID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("试卷不存在")
		}
		return nil, fmt.Errorf("查询试卷失败: %w", err)
	}

	// 获取试卷题目
	paperQuestions, err := examRepo.GetPaperQuestions(paperID)
	if err != nil {
		return nil, fmt.Errorf("查询试卷题目失败: %w", err)
	}

	// 按题型分组
	sectionMap := make(map[int8]*dto.PaperSection)
	typeNames := map[int8]string{
		1: "单选题",
		2: "多选题",
		3: "判断题",
		4: "填空题",
		5: "简答题",
	}

	for _, pq := range paperQuestions {
		// 获取题目信息
		var questionInfo *dto.QuestionInfo
		question, err := FindQuestionByID(pq.QuestionID)
		if err == nil {
			info := toQuestionInfo(question)
			questionInfo = &info
		}

		// 获取题型
		var questionType int8
		if questionInfo != nil {
			questionType = questionInfo.Type
		}

		section, exists := sectionMap[questionType]
		if !exists {
			section = &dto.PaperSection{
				Type:     questionType,
				TypeName: typeNames[questionType],
			}
			sectionMap[questionType] = section
		}

		section.Questions = append(section.Questions, dto.PaperQuestionInfo{
			ID:         pq.ID,
			QuestionID: pq.QuestionID,
			Score:      pq.Score,
			Sort:       pq.Sort,
			Question:   questionInfo,
		})
		section.Count++
		section.Score += pq.Score
	}

	// 转换为切片
	sections := make([]dto.PaperSection, 0, len(sectionMap))
	for _, section := range sectionMap {
		sections = append(sections, *section)
	}

	return &dto.PaperPreviewResponse{
		PaperInfo: toPaperInfo(paper),
		Sections:  sections,
	}, nil
}

// CopyPaper 复制试卷
func CopyPaper(paperID uint, userID uint, title string) (*dto.PaperInfo, error) {
	original, err := examRepo.FindPaperByID(paperID)
	if err != nil {
		return nil, fmt.Errorf("查询试卷失败: %w", err)
	}

	// 创建新试卷
	newPaper := &model.Paper{
		Title:       title,
		Description: original.Description,
		TotalScore:  original.TotalScore,
		TotalCount:  original.TotalCount,
		Type:        original.Type,
		Status:      0, // 草稿
		Visibility:  original.Visibility,
		CreatorID:   userID,
	}

	if newPaper.Title == "" {
		newPaper.Title = original.Title + " (副本)"
	}

	if err := examRepo.CreatePaper(newPaper); err != nil {
		return nil, fmt.Errorf("复制试卷失败: %w", err)
	}

	// 复制题目
	originalQuestions, err := examRepo.GetPaperQuestions(paperID)
	if err == nil {
		for _, pq := range originalQuestions {
			newPQ := model.PaperQuestion{
				PaperID:    newPaper.ID,
				QuestionID: pq.QuestionID,
				Score:      pq.Score,
				Sort:       pq.Sort,
			}
			examRepo.AddPaperQuestions(newPaper.ID, []model.PaperQuestion{newPQ})
		}
	}

	logger.Infof("复制试卷成功: %d -> %d", paperID, newPaper.ID)
	info := toPaperInfo(newPaper)
	return &info, nil
}

// PublishPaper 发布/归档试卷
func PublishPaper(paperID uint, status int8) error {
	paper, err := examRepo.FindPaperByID(paperID)
	if err != nil {
		return fmt.Errorf("查询试卷失败: %w", err)
	}

	paper.Status = status
	if err := examRepo.UpdatePaper(paper); err != nil {
		return fmt.Errorf("更新试卷状态失败: %w", err)
	}

	logger.Infof("试卷 %d 状态更新为 %d", paperID, status)
	return nil
}

// SaveAsTemplate 保存为模板
func SaveAsTemplate(paperID uint, category, tags string) error {
	paper, err := examRepo.FindPaperByID(paperID)
	if err != nil {
		return fmt.Errorf("查询试卷失败: %w", err)
	}

	paper.IsTemplate = true
	paper.Category = category
	paper.Tags = tags

	if err := examRepo.UpdatePaper(paper); err != nil {
		return fmt.Errorf("保存模板失败: %w", err)
	}

	logger.Infof("试卷 %d 保存为模板成功", paperID)
	return nil
}

// CreateFromTemplate 从模板创建试卷
func CreateFromTemplate(templateID uint, userID uint, title string) (*dto.PaperInfo, error) {
	template, err := examRepo.FindPaperByID(templateID)
	if err != nil {
		return nil, fmt.Errorf("查询模板失败: %w", err)
	}

	if !template.IsTemplate {
		return nil, errors.New("该试卷不是模板")
	}

	return CopyPaper(templateID, userID, title)
}

// ListTemplates 模板列表 - moved to exam_ext_service.go

// ExportPaper 导出试卷
func ExportPaper(paperID uint) (*dto.ExportPaperResponse, error) {
	paper, err := examRepo.FindPaperByID(paperID)
	if err != nil {
		return nil, fmt.Errorf("查询试卷失败: %w", err)
	}

	paperQuestions, err := examRepo.GetPaperQuestions(paperID)
	if err != nil {
		return nil, fmt.Errorf("查询试卷题目失败: %w", err)
	}

	export := &dto.ExportPaperResponse{
		Title:       paper.Title,
		Description: paper.Description,
		TotalScore:  paper.TotalScore,
	}

	for i, pq := range paperQuestions {
		eq := dto.ExportQuestion{
			Sort:  i + 1,
			Score: pq.Score,
		}

		question, err := FindQuestionByID(pq.QuestionID)
		if err == nil {
			eq.Title = question.Title
			eq.Type = question.Type
			eq.Difficulty = question.Difficulty
			eq.Answer = question.Answer

			if len(question.Options) > 0 {
				eq.Options = make([]dto.OptionInfo, 0, len(question.Options))
				for _, opt := range question.Options {
					eq.Options = append(eq.Options, dto.OptionInfo{
						Label:     opt.Label,
						Content:   opt.Content,
						IsCorrect: opt.IsCorrect,
					})
				}
			}
		}

		export.Questions = append(export.Questions, eq)
	}

	return export, nil
}

// GetPaperStats 试卷统计
func GetPaperStats(paperID uint) (*dto.PaperStatsResponse, error) {
	// 查询使用该试卷的考试记录
	var exams []model.Exam
	database.DB.Where("paper_id = ?", paperID).Find(&exams)

	if len(exams) == 0 {
		return &dto.PaperStatsResponse{}, nil
	}

	// 收集所有考试记录的分数（每个考试使用自身的及格线）
	var scores []float64
	type examScore struct {
		score     float64
		passScore float64
	}
	var examScores []examScore
	for _, e := range exams {
		var records []model.ExamRecord
		database.DB.Where("exam_id = ? AND status >= 1", e.ID).Find(&records)
		for _, r := range records {
			scores = append(scores, r.Score)
			examScores = append(examScores, examScore{score: r.Score, passScore: e.PassScore})
		}
	}

	if len(scores) == 0 {
		return &dto.PaperStatsResponse{
			UsageCount: len(exams),
		}, nil
	}

	// 计算统计
	stats := &dto.PaperStatsResponse{
		UsageCount: len(exams),
	}

	var total, max, min, passCount float64
	min = scores[0]
	for i, s := range scores {
		total += s
		if s > max {
			max = s
		}
		if s < min {
			min = s
		}
		// 使用考试自身的及格线
		if s >= examScores[i].passScore {
			passCount++
		}
	}

	stats.AvgScore = total / float64(len(scores))
	stats.MaxScore = max
	stats.MinScore = min
	stats.PassRate = passCount / float64(len(scores)) * 100

	// 分数分布
	distMap := make(map[string]int)
	for _, s := range scores {
		key := fmt.Sprintf("%d-%d", int(s)/10*10, int(s)/10*10+10)
		distMap[key]++
	}

	for k, v := range distMap {
		stats.ScoreDistribution = append(stats.ScoreDistribution, dto.ScoreDist{
			Range: k,
			Count: v,
		})
	}

	return stats, nil
}

// toPaperInfo 转换为 PaperInfo DTO (使用exam_service.go中的版本)

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
