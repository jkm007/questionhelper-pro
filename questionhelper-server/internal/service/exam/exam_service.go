package exam

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	examRepo "questionhelper-server/internal/repository/exam"
	questionRepo "questionhelper-server/internal/repository/question"
	"questionhelper-server/pkg/logger"
)

// ==================== Paper ====================

// ListPapers 试卷列表
func ListPapers(req *dto.PageRequest) ([]dto.PaperInfo, int64, error) {
	papers, total, err := examRepo.ListPapers(req)
	if err != nil {
		return nil, 0, fmt.Errorf("查询试卷列表失败: %w", err)
	}

	list := make([]dto.PaperInfo, 0, len(papers))
	for _, p := range papers {
		list = append(list, toPaperInfo(&p))
	}
	return list, total, nil
}

// GetPaper 获取试卷详情
func GetPaper(id uint) (*dto.PaperInfo, error) {
	paper, err := examRepo.FindPaperByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("试卷不存在")
		}
		return nil, fmt.Errorf("查询试卷失败: %w", err)
	}
	info := toPaperInfo(paper)
	return &info, nil
}

// CreatePaper 创建试卷
func CreatePaper(creatorID uint, req *dto.CreatePaperRequest) error {
	paper := &model.Paper{
		Title:       req.Title,
		Description: req.Description,
		Type:        req.Type,
		CreatorID:   creatorID,
	}

	if err := examRepo.CreatePaper(paper); err != nil {
		return fmt.Errorf("创建试卷失败: %w", err)
	}

	logger.Infof("创建试卷成功: %d", paper.ID)
	return nil
}

// UpdatePaper 更新试卷
func UpdatePaper(id uint, req *dto.CreatePaperRequest) error {
	paper, err := examRepo.FindPaperByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("试卷不存在")
		}
		return fmt.Errorf("查询试卷失败: %w", err)
	}

	paper.Title = req.Title
	paper.Description = req.Description

	if err := examRepo.UpdatePaper(paper); err != nil {
		return fmt.Errorf("更新试卷失败: %w", err)
	}

	logger.Infof("更新试卷 %d 成功", id)
	return nil
}

// DeletePaper 删除试卷
func DeletePaper(id uint) error {
	_, err := examRepo.FindPaperByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("试卷不存在")
		}
		return fmt.Errorf("查询试卷失败: %w", err)
	}

	if err := examRepo.DeletePaperByID(id); err != nil {
		return fmt.Errorf("删除试卷失败: %w", err)
	}

	logger.Infof("删除试卷 %d 成功", id)
	return nil
}

// AddPaperQuestions 添加试卷题目
func AddPaperQuestions(paperID uint, questions []struct {
	QuestionID uint    `json:"question_id"`
	Score      float64 `json:"score"`
}) error {
	paper, err := examRepo.FindPaperByID(paperID)
	if err != nil {
		return errors.New("试卷不存在")
	}

	// 删除旧题目
	if err := examRepo.DeletePaperQuestions(paperID); err != nil {
		return fmt.Errorf("删除旧题目失败: %w", err)
	}

	// 添加新题目
	paperQuestions := make([]model.PaperQuestion, 0, len(questions))
	totalScore := float64(0)
	for i, q := range questions {
		paperQuestions = append(paperQuestions, model.PaperQuestion{
			PaperID:    paperID,
			QuestionID: q.QuestionID,
			Score:      q.Score,
			Sort:       i,
		})
		totalScore += q.Score
	}

	if err := examRepo.AddPaperQuestions(paperID, paperQuestions); err != nil {
		return fmt.Errorf("添加题目失败: %w", err)
	}

	// 更新试卷统计
	paper.TotalScore = totalScore
	paper.TotalCount = len(questions)
	if err := examRepo.UpdatePaper(paper); err != nil {
		return fmt.Errorf("更新试卷统计失败: %w", err)
	}

	logger.Infof("试卷 %d 添加 %d 道题目", paperID, len(questions))
	return nil
}

// ==================== Exam ====================

// ListExams 考试列表（管理员）
func ListExams(req *dto.ExamListRequest) ([]dto.ExamInfo, int64, error) {
	exams, total, err := examRepo.ListExams(req)
	if err != nil {
		return nil, 0, fmt.Errorf("查询考试列表失败: %w", err)
	}

	list := make([]dto.ExamInfo, 0, len(exams))
	for _, e := range exams {
		list = append(list, toExamInfo(&e))
	}
	return list, total, nil
}

// ListAvailableExams 可用考试列表（学生）
func ListAvailableExams(userID uint, req *dto.PageRequest) ([]dto.ExamInfo, int64, error) {
	exams, total, err := examRepo.ListAvailableExams(userID, req)
	if err != nil {
		return nil, 0, fmt.Errorf("查询可用考试失败: %w", err)
	}

	list := make([]dto.ExamInfo, 0, len(exams))
	for _, e := range exams {
		list = append(list, toExamInfo(&e))
	}
	return list, total, nil
}

// GetExam 获取考试详情
func GetExam(id uint) (*dto.ExamInfo, error) {
	exam, err := examRepo.FindExamByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("考试不存在")
		}
		return nil, fmt.Errorf("查询考试失败: %w", err)
	}
	info := toExamInfo(exam)
	return &info, nil
}

// CreateExam 创建考试
func CreateExam(creatorID uint, req *dto.CreateExamRequest) error {
	exam := &model.Exam{
		Title:       req.Title,
		Description: req.Description,
		PaperID:     req.PaperID,
		ClassID:     req.ClassID,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		Duration:    req.Duration,
		TotalScore:  req.TotalScore,
		PassScore:   req.PassScore,
		MaxAttempts: req.MaxAttempts,
		Shuffle:     req.Shuffle,
		ShowAnswer:  req.ShowAnswer,
		AntiCheat:   req.AntiCheat,
		CreatorID:   creatorID,
	}

	if exam.MaxAttempts == 0 {
		exam.MaxAttempts = 1
	}

	if err := examRepo.CreateExam(exam); err != nil {
		return fmt.Errorf("创建考试失败: %w", err)
	}

	logger.Infof("创建考试成功: %d", exam.ID)
	return nil
}

// UpdateExam 更新考试
func UpdateExam(id uint, req *dto.UpdateExamRequest) error {
	exam, err := examRepo.FindExamByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("考试不存在")
		}
		return fmt.Errorf("查询考试失败: %w", err)
	}

	if exam.Status != 0 {
		return errors.New("只能修改未发布的考试")
	}

	if req.Title != "" {
		exam.Title = req.Title
	}
	if req.Description != "" {
		exam.Description = req.Description
	}
	if !req.StartTime.IsZero() {
		exam.StartTime = req.StartTime
	}
	if !req.EndTime.IsZero() {
		exam.EndTime = req.EndTime
	}
	if req.Duration > 0 {
		exam.Duration = req.Duration
	}
	if req.TotalScore > 0 {
		exam.TotalScore = req.TotalScore
	}
	if req.PassScore > 0 {
		exam.PassScore = req.PassScore
	}
	if req.MaxAttempts > 0 {
		exam.MaxAttempts = req.MaxAttempts
	}
	if req.Shuffle != nil {
		exam.Shuffle = *req.Shuffle
	}

	if err := examRepo.UpdateExam(exam); err != nil {
		return fmt.Errorf("更新考试失败: %w", err)
	}

	logger.Infof("更新考试 %d 成功", id)
	return nil
}

// DeleteExam 删除考试
func DeleteExam(id uint) error {
	exam, err := examRepo.FindExamByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("考试不存在")
		}
		return fmt.Errorf("查询考试失败: %w", err)
	}

	if exam.Status == 1 {
		return errors.New("不能删除进行中的考试")
	}

	if err := examRepo.DeleteExamByID(id); err != nil {
		return fmt.Errorf("删除考试失败: %w", err)
	}

	logger.Infof("删除考试 %d 成功", id)
	return nil
}

// PublishExam 发布考试
func PublishExam(id uint) error {
	exam, err := examRepo.FindExamByID(id)
	if err != nil {
		return errors.New("考试不存在")
	}

	if exam.Status != 0 {
		return errors.New("只能发布未发布的考试")
	}

	exam.Status = 1
	if err := examRepo.UpdateExam(exam); err != nil {
		return fmt.Errorf("发布考试失败: %w", err)
	}

	logger.Infof("发布考试 %d 成功", id)
	return nil
}

// CloseExam 结束考试
func CloseExam(id uint) error {
	exam, err := examRepo.FindExamByID(id)
	if err != nil {
		return errors.New("考试不存在")
	}

	if exam.Status != 1 {
		return errors.New("只能结束进行中的考试")
	}

	exam.Status = 2
	if err := examRepo.UpdateExam(exam); err != nil {
		return fmt.Errorf("结束考试失败: %w", err)
	}

	logger.Infof("结束考试 %d 成功", id)
	return nil
}

// ==================== Exam Taking ====================

// StartExam 开始考试
func StartExam(examID, userID uint, ip string) (*dto.ExamRecordInfo, error) {
	exam, err := examRepo.FindExamByID(examID)
	if err != nil {
		return nil, errors.New("考试不存在")
	}

	if exam.Status != 1 {
		return nil, errors.New("考试未开放")
	}

	now := time.Now()
	if now.Before(exam.StartTime) {
		return nil, errors.New("考试尚未开始")
	}
	if now.After(exam.EndTime) {
		return nil, errors.New("考试已结束")
	}

	// 检查考试次数
	attempts, err := examRepo.CountUserAttempts(examID, userID)
	if err != nil {
		return nil, fmt.Errorf("查询考试次数失败: %w", err)
	}
	if int(attempts) >= exam.MaxAttempts {
		return nil, errors.New("已达最大考试次数")
	}

	// 创建考试记录
	record := &model.ExamRecord{
		ExamID:    examID,
		UserID:    userID,
		Status:    0,
		StartTime: now,
		IP:        ip,
	}

	if err := examRepo.CreateExamRecord(record); err != nil {
		return nil, fmt.Errorf("创建考试记录失败: %w", err)
	}

	logger.Infof("用户 %d 开始考试 %d", userID, examID)

	info := &dto.ExamRecordInfo{
		ID:        record.ID,
		ExamID:    record.ExamID,
		UserID:    record.UserID,
		Status:    record.Status,
		StartTime: record.StartTime,
	}
	return info, nil
}

// SubmitExam 提交考试
func SubmitExam(recordID uint, req *dto.SubmitExamRequest) error {
	record, err := examRepo.FindExamRecord(recordID)
	if err != nil {
		return errors.New("考试记录不存在")
	}

	if record.Status != 0 {
		return errors.New("考试已提交")
	}

	// 获取试卷题目
	paperQuestions, err := examRepo.GetPaperQuestions(record.ExamID)
	if err != nil {
		return fmt.Errorf("获取试卷题目失败: %w", err)
	}

	// 构建题目分数映射
	scoreMap := make(map[uint]float64)
	for _, pq := range paperQuestions {
		scoreMap[pq.QuestionID] = pq.Score
	}

	// 计算得分
	totalScore := float64(0)
	answerRecords := make([]model.AnswerRecord, 0, len(req.Answers))
	for _, ans := range req.Answers {
		// 获取题目正确答案
		question, err := questionRepo.FindByID(ans.QuestionID)
		if err != nil {
			continue
		}

		isCorrect := question.Answer == ans.Answer
		score := float64(0)
		if isCorrect {
			score = scoreMap[ans.QuestionID]
			totalScore += score
		}

		answerRecords = append(answerRecords, model.AnswerRecord{
			RecordID:   recordID,
			QuestionID: ans.QuestionID,
			Answer:     ans.Answer,
			Score:      score,
			IsCorrect:  isCorrect,
		})
	}

	// 保存答题记录
	if err := examRepo.CreateAnswerRecords(answerRecords); err != nil {
		return fmt.Errorf("保存答题记录失败: %w", err)
	}

	// 更新考试记录
	now := time.Now()
	record.Score = totalScore
	record.Status = 1
	record.SubmitTime = &now

	if err := examRepo.UpdateExamRecord(record); err != nil {
		return fmt.Errorf("更新考试记录失败: %w", err)
	}

	logger.Infof("用户 %d 提交考试 %d，得分: %.2f", record.UserID, record.ExamID, totalScore)
	return nil
}

// GetExamResult 获取考试结果
func GetExamResult(examID, userID uint) (interface{}, error) {
	record, err := examRepo.FindExamRecordByUser(examID, userID)
	if err != nil {
		return nil, errors.New("未找到考试记录")
	}

	exam, err := examRepo.FindExamByID(examID)
	if err != nil {
		return nil, errors.New("考试不存在")
	}

	result := map[string]interface{}{
		"record": dto.ExamRecordInfo{
			ID:         record.ID,
			ExamID:     record.ExamID,
			UserID:     record.UserID,
			Score:      record.Score,
			Status:     record.Status,
			StartTime:  record.StartTime,
			SubmitTime: record.SubmitTime,
		},
		"total_score": exam.TotalScore,
		"pass_score":  exam.PassScore,
		"is_pass":     record.Score >= exam.PassScore,
	}

	// 如果需要显示答案
	if exam.ShowAnswer == 1 && record.Status >= 1 {
		answers, _ := examRepo.GetAnswerRecords(record.ID)
		result["answers"] = answers
	}

	return result, nil
}

// GetExamHistory 考试历史
func GetExamHistory(userID uint, req *dto.PageRequest) ([]dto.ExamRecordInfo, int64, error) {
	records, total, err := examRepo.ListExamRecords(nil, &userID, req)
	if err != nil {
		return nil, 0, fmt.Errorf("查询考试历史失败: %w", err)
	}

	list := make([]dto.ExamRecordInfo, 0, len(records))
	for _, r := range records {
		list = append(list, dto.ExamRecordInfo{
			ID:         r.ID,
			ExamID:     r.ExamID,
			UserID:     r.UserID,
			Score:      r.Score,
			Status:     r.Status,
			StartTime:  r.StartTime,
			SubmitTime: r.SubmitTime,
		})
	}
	return list, total, nil
}

// ==================== Score ====================

// ListScores 成绩列表
func ListScores(examID *uint, req *dto.PageRequest) ([]dto.ExamRecordInfo, int64, error) {
	records, total, err := examRepo.ListExamRecords(examID, nil, req)
	if err != nil {
		return nil, 0, fmt.Errorf("查询成绩列表失败: %w", err)
	}

	list := make([]dto.ExamRecordInfo, 0, len(records))
	for _, r := range records {
		list = append(list, dto.ExamRecordInfo{
			ID:         r.ID,
			ExamID:     r.ExamID,
			UserID:     r.UserID,
			Score:      r.Score,
			Status:     r.Status,
			StartTime:  r.StartTime,
			SubmitTime: r.SubmitTime,
		})
	}
	return list, total, nil
}

// GetScore 获取成绩详情
func GetScore(recordID uint) (interface{}, error) {
	record, err := examRepo.FindExamRecord(recordID)
	if err != nil {
		return nil, errors.New("考试记录不存在")
	}

	answers, _ := examRepo.GetAnswerRecords(recordID)

	return map[string]interface{}{
		"record":  record,
		"answers": answers,
	}, nil
}

// GetScoreAnalysis 成绩分析
func GetScoreAnalysis(examID uint) (map[string]interface{}, error) {
	return examRepo.GetExamScoreStats(examID)
}

// ==================== Helpers ====================

func toPaperInfo(p *model.Paper) dto.PaperInfo {
	return dto.PaperInfo{
		ID:          p.ID,
		Title:       p.Title,
		Description: p.Description,
		TotalScore:  p.TotalScore,
		TotalCount:  p.TotalCount,
		Type:        p.Type,
	}
}

func toExamInfo(e *model.Exam) dto.ExamInfo {
	return dto.ExamInfo{
		ID:          e.ID,
		Title:       e.Title,
		Description: e.Description,
		PaperID:     e.PaperID,
		ClassID:     e.ClassID,
		StartTime:   e.StartTime,
		EndTime:     e.EndTime,
		Duration:    e.Duration,
		TotalScore:  e.TotalScore,
		PassScore:   e.PassScore,
		MaxAttempts: e.MaxAttempts,
		Shuffle:     e.Shuffle,
		ShowAnswer:  e.ShowAnswer,
		AntiCheat:   e.AntiCheat,
		Status:      e.Status,
		CreatorID:   e.CreatorID,
	}
}
