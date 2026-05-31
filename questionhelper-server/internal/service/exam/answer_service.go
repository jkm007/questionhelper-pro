package exam

import (
	"errors"
	"fmt"

	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	examRepo "questionhelper-server/internal/repository/exam"
	questionRepo "questionhelper-server/internal/repository/question"
	"questionhelper-server/pkg/logger"
)

// SaveAnswer 保存答案(单题)
func SaveAnswer(recordID uint, req *dto.SaveAnswerRequest) error {
	record, err := examRepo.FindExamRecord(recordID)
	if err != nil {
		return fmt.Errorf("查询考试记录失败: %w", err)
	}

	if record.Status != 0 {
		return errors.New("考试已结束，无法保存答案")
	}

	// 查找或创建答案记录
	answers, _ := examRepo.GetAnswerRecords(recordID)
	var existing *model.AnswerRecord
	for i, a := range answers {
		if a.QuestionID == req.QuestionID {
			existing = &answers[i]
			break
		}
	}

	if existing != nil {
		existing.Answer = req.Answer
		if req.IsMarked != nil {
			existing.IsMarked = *req.IsMarked
		}
		return examRepo.UpdateAnswerRecord(existing)
	}

	// 创建新答案
	answer := &model.AnswerRecord{
		RecordID:   recordID,
		QuestionID: req.QuestionID,
		Answer:     req.Answer,
	}
	if req.IsMarked != nil {
		answer.IsMarked = *req.IsMarked
	}

	return examRepo.CreateSingleAnswerRecord(answer)
}

// SaveAnswers 批量保存答案
func SaveAnswers(recordID uint, req *dto.SaveAnswerBatchRequest) error {
	for _, ans := range req.Answers {
		if err := SaveAnswer(recordID, &ans); err != nil {
			logger.Errorf("保存答案失败: %v", err)
		}
	}
	return nil
}

// GetStandardAnswers 获取标准答案
func GetStandardAnswers(examID, userID uint) (*dto.StandardAnswersResponse, error) {
	// 获取考试记录
	record, err := examRepo.FindExamRecordByUser(examID, userID)
	if err != nil {
		return nil, errors.New("未找到考试记录")
	}

	// 获取考试信息
	exam, err := examRepo.FindExamByID(examID)
	if err != nil {
		return nil, fmt.Errorf("查询考试失败: %w", err)
	}

	// 检查是否允许查看答案
	if exam.ShowAnswer == 0 {
		return nil, errors.New("该考试不允许查看答案")
	}
	if exam.ShowAnswer == 1 && record.Status < 1 {
		return nil, errors.New("交卷后才能查看答案")
	}

	// 获取试卷题目
	paperQuestions, err := examRepo.GetPaperQuestions(exam.PaperID)
	if err != nil {
		return nil, fmt.Errorf("查询试卷题目失败: %w", err)
	}

	// 获取用户答案
	userAnswers, _ := examRepo.GetAnswerRecords(record.ID)
	answerMap := make(map[uint]*model.AnswerRecord)
	for i, a := range userAnswers {
		answerMap[a.QuestionID] = &userAnswers[i]
	}

	// 构建标准答案
	response := &dto.StandardAnswersResponse{
		ExamID: examID,
	}

	for _, pq := range paperQuestions {
		// 获取题目信息
		question, err := FindQuestionByID(pq.QuestionID)
		if err != nil {
			continue
		}

		sa := dto.StandardAnswer{
			QuestionID: pq.QuestionID,
			Title:      question.Title,
			Type:       question.Type,
			Score:      pq.Score,
			Answer:     question.Answer,
			Analysis:   question.Analysis,
		}

		// 添加选项
		if len(question.Options) > 0 {
			sa.Options = make([]dto.OptionInfo, 0, len(question.Options))
			for _, opt := range question.Options {
				sa.Options = append(sa.Options, dto.OptionInfo{
					Label:     opt.Label,
					Content:   opt.Content,
					IsCorrect: opt.IsCorrect,
				})
			}
		}

		// 添加用户答案
		if ua, exists := answerMap[pq.QuestionID]; exists {
			sa.UserAnswer = ua.Answer
			sa.UserScore = ua.Score
			sa.IsCorrect = ua.IsCorrect
		}

		response.Questions = append(response.Questions, sa)
	}

	return response, nil
}

// MarkQuestion 标记/取消标记题目
func MarkQuestion(recordID, questionID uint, isMarked bool) error {
	answers, _ := examRepo.GetAnswerRecords(recordID)
	for i, a := range answers {
		if a.QuestionID == questionID {
			answers[i].IsMarked = isMarked
			return examRepo.UpdateAnswerRecord(&answers[i])
		}
	}

	// 如果没有答案记录，创建一个
	answer := &model.AnswerRecord{
		RecordID:   recordID,
		QuestionID: questionID,
		IsMarked:   isMarked,
	}
	return examRepo.CreateSingleAnswerRecord(answer)
}

// GetMarkedQuestions 获取标记题目
func GetMarkedQuestions(recordID uint) ([]dto.SaveAnswerRequest, error) {
	answers, err := examRepo.GetAnswerRecords(recordID)
	if err != nil {
		return nil, err
	}

	var marked []dto.SaveAnswerRequest
	for _, a := range answers {
		if a.IsMarked {
			marked = append(marked, dto.SaveAnswerRequest{
				QuestionID: a.QuestionID,
				Answer:     a.Answer,
				IsMarked:   &a.IsMarked,
			})
		}
	}

	return marked, nil
}

// GetExamGuide 获取答题指引
func GetExamGuide(examID uint) (*dto.ExamGuideResponse, error) {
	exam, err := examRepo.FindExamByID(examID)
	if err != nil {
		return nil, fmt.Errorf("查询考试失败: %w", err)
	}

	return &dto.ExamGuideResponse{
		ExamID:      exam.ID,
		Title:       exam.Title,
		Description: exam.Description,
		Duration:    exam.Duration,
		TotalScore:  exam.TotalScore,
		PassScore:   exam.PassScore,
		Rules: []string{
			"请在规定时间内完成答题",
			"答案将自动保存",
			"交卷后无法修改答案",
		},
		Tips: []string{
			"先做会做的题目，再回头做难题",
			"注意题目要求，单选/多选不要混淆",
			"合理分配时间，不要在一道题上花费太多时间",
		},
	}, nil
}

// SubmitFeedback 提交考后反馈
func SubmitFeedback(examID, userID uint, req *dto.FeedbackRequest) error {
	// 简化实现：记录到日志
	logger.Infof("用户 %d 对考试 %d 提交反馈: 评分=%d, 内容=%s", userID, examID, req.Rating, req.Content)
	return nil
}

// AutoSaveAnswer 自动保存答案(Redis)
func AutoSaveAnswer(recordID uint, questionID uint, answer string) error {
	// 简化实现：直接保存到数据库
	return SaveAnswer(recordID, &dto.SaveAnswerRequest{
		QuestionID: questionID,
		Answer:     answer,
	})
}

// GetAutoSavedAnswers 获取自动保存的答案
func GetAutoSavedAnswers(recordID uint) (map[uint]string, error) {
	answers, err := examRepo.GetAnswerRecords(recordID)
	if err != nil {
		return nil, err
	}

	result := make(map[uint]string)
	for _, a := range answers {
		result[a.QuestionID] = a.Answer
	}

	return result, nil
}

// RecordWarning 记录考试异常
func RecordWarning(recordID, examID, userID uint, warningType, detail, ip, userAgent string) error {
	warning := &model.ExamWarning{
		RecordID:  recordID,
		ExamID:    examID,
		UserID:    userID,
		Type:      warningType,
		Detail:    detail,
		IP:        ip,
		UserAgent: userAgent,
	}

	return examRepo.CreateWarning(warning)
}

// FindRecordByID 查找考试记录
func FindRecordByID(recordID uint) (*model.ExamRecord, error) {
	return examRepo.FindExamRecord(recordID)
}

// FindQuestionByID 查找题目(通过question_repo)
func FindQuestionByID(id uint) (*model.Question, error) {
	return questionRepo.FindByID(id)
}
