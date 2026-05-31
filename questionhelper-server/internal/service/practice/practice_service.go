package practice

import (
	"errors"
	"fmt"

	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	questionRepo "questionhelper-server/internal/repository/question"
	practiceRepo "questionhelper-server/internal/repository/practice"
	wrongRepo "questionhelper-server/internal/repository/wrong"
	"questionhelper-server/pkg/logger"
)

// StartPractice 开始练习
func StartPractice(userID uint, req *dto.StartPracticeRequest) (*dto.PracticeSessionInfo, error) {
	count := req.Count
	if count == 0 {
		count = 10
	}

	// 获取题目
	questionReq := &dto.QuestionListRequest{
		PageRequest: dto.PageRequest{Page: 1, PageSize: count},
		CategoryID:  req.CategoryID,
		Type:        req.Type,
		Difficulty:  req.Difficulty,
	}

	questions, _, err := questionRepo.List(questionReq)
	if err != nil {
		return nil, fmt.Errorf("获取题目失败: %w", err)
	}

	if len(questions) == 0 {
		return nil, errors.New("没有符合条件的题目")
	}

	// 创建练习会话
	session := &model.PracticeSession{
		UserID:     userID,
		CategoryID: req.CategoryID,
		TotalCount: len(questions),
		Status:     0,
	}

	if err := practiceRepo.CreateSession(session); err != nil {
		return nil, fmt.Errorf("创建练习会话失败: %w", err)
	}

	logger.Infof("用户 %d 开始练习，会话 %d，题目数 %d", userID, session.ID, len(questions))

	info := &dto.PracticeSessionInfo{
		ID:         session.ID,
		UserID:     session.UserID,
		CategoryID: session.CategoryID,
		TotalCount: session.TotalCount,
		Status:     session.Status,
		CreatedAt:  session.CreatedAt,
	}
	return info, nil
}

// SubmitPractice 提交练习答案
func SubmitPractice(sessionID, userID uint, req *dto.SubmitPracticeAnswerRequest) error {
	session, err := practiceRepo.FindSessionByID(sessionID)
	if err != nil {
		return errors.New("练习会话不存在")
	}

	if session.UserID != userID {
		return errors.New("无权操作此练习会话")
	}

	if session.Status != 0 {
		return errors.New("练习已结束")
	}

	// 获取题目正确答案
	question, err := questionRepo.FindByID(req.QuestionID)
	if err != nil {
		return errors.New("题目不存在")
	}

	isCorrect := question.Answer == req.Answer

	// 保存答题记录
	record := &model.PracticeRecord{
		SessionID:  sessionID,
		QuestionID: req.QuestionID,
		Answer:     req.Answer,
		IsCorrect:  isCorrect,
		Duration:   req.Duration,
	}

	if err := practiceRepo.CreateRecords([]model.PracticeRecord{*record}); err != nil {
		return fmt.Errorf("保存答题记录失败: %w", err)
	}

	// 更新错题本
	if !isCorrect {
		updateWrongQuestion(session.UserID, req.QuestionID, 1, sessionID, req.Answer)
	}

	return nil
}

// FinishPractice 完成练习
func FinishPractice(sessionID uint) error {
	session, err := practiceRepo.FindSessionByID(sessionID)
	if err != nil {
		return errors.New("练习会话不存在")
	}

	if session.Status != 0 {
		return errors.New("练习已结束")
	}

	// 获取答题记录
	records, err := practiceRepo.GetRecords(sessionID)
	if err != nil {
		return fmt.Errorf("获取答题记录失败: %w", err)
	}

	// 统计正确数
	correctCount := 0
	for _, r := range records {
		if r.IsCorrect {
			correctCount++
		}
	}

	// 计算正确率
	accuracy := float64(0)
	if len(records) > 0 {
		accuracy = float64(correctCount) / float64(len(records)) * 100
	}

	// 更新会话
	session.CorrectCount = correctCount
	session.Accuracy = accuracy
	session.Status = 1

	if err := practiceRepo.UpdateSession(session); err != nil {
		return fmt.Errorf("更新练习会话失败: %w", err)
	}

	logger.Infof("练习 %d 完成，正确率: %.2f%%", sessionID, accuracy)
	return nil
}

// GetPracticeResult 获取练习结果
func GetPracticeResult(sessionID, userID uint) (interface{}, error) {
	session, err := practiceRepo.FindSessionByID(sessionID)
	if err != nil {
		return nil, errors.New("练习会话不存在")
	}

	if session.UserID != userID {
		return nil, errors.New("无权查看此练习")
	}

	records, err := practiceRepo.GetRecords(sessionID)
	if err != nil {
		return nil, fmt.Errorf("获取答题记录失败: %w", err)
	}

	recordInfos := make([]dto.PracticeRecordInfo, 0, len(records))
	for _, r := range records {
		recordInfos = append(recordInfos, dto.PracticeRecordInfo{
			ID:         r.ID,
			SessionID:  r.SessionID,
			QuestionID: r.QuestionID,
			Answer:     r.Answer,
			IsCorrect:  r.IsCorrect,
			Duration:   r.Duration,
		})
	}

	return map[string]interface{}{
		"session": dto.PracticeSessionInfo{
			ID:           session.ID,
			UserID:       session.UserID,
			CategoryID:   session.CategoryID,
			TotalCount:   session.TotalCount,
			CorrectCount: session.CorrectCount,
			Accuracy:     session.Accuracy,
			Duration:     session.Duration,
			Status:       session.Status,
			CreatedAt:    session.CreatedAt,
		},
		"records": recordInfos,
	}, nil
}

// GetPracticeHistory 练习历史
func GetPracticeHistory(userID uint, req *dto.PracticeListRequest) ([]dto.PracticeSessionInfo, int64, error) {
	sessions, total, err := practiceRepo.ListSessions(userID, req)
	if err != nil {
		return nil, 0, fmt.Errorf("查询练习历史失败: %w", err)
	}

	list := make([]dto.PracticeSessionInfo, 0, len(sessions))
	for _, s := range sessions {
		list = append(list, dto.PracticeSessionInfo{
			ID:           s.ID,
			UserID:       s.UserID,
			CategoryID:   s.CategoryID,
			TotalCount:   s.TotalCount,
			CorrectCount: s.CorrectCount,
			Accuracy:     s.Accuracy,
			Duration:     s.Duration,
			Status:       s.Status,
			CreatedAt:    s.CreatedAt,
		})
	}
	return list, total, nil
}

// GetPracticeStats 练习统计
func GetPracticeStats(userID uint) (map[string]interface{}, error) {
	return practiceRepo.GetUserPracticeStats(userID)
}

// updateWrongQuestion 更新错题本
func updateWrongQuestion(userID, questionID, source, sourceID uint, answer string) {
	existing, err := wrongRepo.FindByUserAndQuestion(userID, questionID)
	if err == nil {
		// 已存在，更新错误次数
		existing.WrongCount++
		existing.LastAnswer = answer
		existing.Mastered = false
		wrongRepo.Update(existing)
		return
	}

	// 不存在，创建新的
	wrong := &model.WrongQuestion{
		UserID:     userID,
		QuestionID: questionID,
		Source:     int8(source),
		SourceID:   sourceID,
		WrongCount: 1,
		LastAnswer: answer,
		Mastered:   false,
	}
	wrongRepo.Create(wrong)
}
