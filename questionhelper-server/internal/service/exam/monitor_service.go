package exam

import (
	"fmt"
	"time"

	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	examRepo "questionhelper-server/internal/repository/exam"
	notifRepo "questionhelper-server/internal/repository/notification"
	questionRepo "questionhelper-server/internal/repository/question"
	"questionhelper-server/pkg/logger"
)

// GetExamMonitor 获取考试监控信息
func GetExamMonitor(examID uint) (*dto.ExamMonitorInfo, error) {
	exam, err := examRepo.FindExamByID(examID)
	if err != nil {
		return nil, fmt.Errorf("查询考试失败: %w", err)
	}

	// 获取所有考试记录
	req := &dto.PageRequest{Page: 1, PageSize: 1000}
	records, _, err := examRepo.ListExamRecords(&examID, nil, req)
	if err != nil {
		return nil, fmt.Errorf("查询考试记录失败: %w", err)
	}

	monitor := &dto.ExamMonitorInfo{
		ExamID:        exam.ID,
		ExamTitle:     exam.Title,
		TotalStudents: len(records),
	}

	// 统计在线和已提交
	now := time.Now()
	for _, record := range records {
		if record.Status >= 1 {
			monitor.SubmittedCount++
		}

		// 判断是否在线(5分钟内有活动)
		if record.Status == 0 && now.Sub(record.UpdatedAt) < 5*time.Minute {
			monitor.OnlineCount++

			// 计算答题进度
			answers, _ := examRepo.GetAnswerRecords(record.ID)
			progress := float64(0)
			totalCount := exam.Paper.TotalCount
			if totalCount > 0 {
				progress = float64(len(answers)) / float64(totalCount) * 100
			}

			monitor.OnlineUsers = append(monitor.OnlineUsers, dto.OnlineUserInfo{
				UserID:       record.UserID,
				Progress:     progress,
				DurationUsed: int(now.Sub(record.StartTime).Seconds()),
				LastActive:   record.UpdatedAt,
				IP:           record.IP,
			})
		}
	}

	// 获取异常记录
	warnings, _, err := examRepo.FindWarningsByExamID(examID, 1, 50)
	if err == nil {
		monitor.WarningCount = len(warnings)
		for _, w := range warnings {
			monitor.Warnings = append(monitor.Warnings, dto.WarningInfo{
				ID:        w.ID,
				UserID:    w.UserID,
				Type:      w.Type,
				Detail:    w.Detail,
				CreatedAt: w.CreatedAt,
			})
		}
	}

	return monitor, nil
}

// ExportScores 导出成绩
func ExportScores(examID uint) ([]dto.ExamRecordInfo, error) {
	req := &dto.PageRequest{Page: 1, PageSize: 10000}
	records, _, err := examRepo.ListExamRecords(&examID, nil, req)
	if err != nil {
		return nil, fmt.Errorf("查询成绩失败: %w", err)
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

	return list, nil
}

// ReviewExam 阅卷
func ReviewExam(examID uint, req *dto.ReviewRequest) error {
	// 获取考试记录
	record, err := examRepo.FindExamRecord(req.RecordID)
	if err != nil {
		return fmt.Errorf("查询考试记录失败: %w", err)
	}

	// 获取考试信息
	exam, err := examRepo.FindExamByID(examID)
	if err != nil {
		return fmt.Errorf("查询考试失败: %w", err)
	}

	// 更新每道题的分数
	for _, ans := range req.Answers {
		answerRecord, err := examRepo.FindAnswerRecordByID(ans.AnswerID)
		if err != nil {
			continue
		}

		answerRecord.Score = ans.Score
		answerRecord.ReviewNote = ans.Note
		answerRecord.IsReviewed = true
		now := time.Now()
		answerRecord.ReviewedAt = &now
		// 设置阅卷人（从examID关联的考试创建者获取，或使用reviewerID参数）
		// 这里使用考试的CreatorID作为阅卷人
		reviewerID := exam.CreatorID
		answerRecord.ReviewedBy = &reviewerID

		examRepo.UpdateAnswerRecord(answerRecord)
	}

	// 重新计算总分
	answers, _ := examRepo.GetAnswerRecords(record.ID)
	var totalScore float64
	var objScore, subjScore float64
	allReviewed := true
	for _, a := range answers {
		totalScore += a.Score
		if a.IsReviewed {
			// 已阅卷的主观题分数计入主观题得分
			subjScore += a.Score
		} else {
			allReviewed = false
		}
		// 客观题得分从原始记录中累计
		if a.MaxScore > 0 && a.IsCorrect {
			objScore += a.Score
		}
	}

	record.Score = totalScore
	record.ObjScore = objScore
	record.SubjScore = subjScore
	if allReviewed {
		record.Status = 2 // 已阅卷
	}

	examRepo.UpdateExamRecord(record)

	// 发送通知给学生
	notification := &model.Notification{
		UserID:     record.UserID,
		Type:       2, // 考试通知
		Title:      fmt.Sprintf("考试「%s」阅卷完成", exam.Title),
		Content:    fmt.Sprintf("您的考试已阅卷完成，总分: %.2f。%s", totalScore, req.Comment),
		TargetType: "exam",
		TargetID:   examID,
	}

	if err := notifRepo.Create(notification); err != nil {
		logger.Errorf("发送阅卷通知失败: %v", err)
	}

	return nil
}

// GetExamAnalysis 获取考试分析
func GetExamAnalysis(examID uint) (*dto.ExamAnalysisResponse, error) {
	exam, err := examRepo.FindExamByID(examID)
	if err != nil {
		return nil, fmt.Errorf("查询考试失败: %w", err)
	}

	req := &dto.PageRequest{Page: 1, PageSize: 10000}
	records, _, err := examRepo.ListExamRecords(&examID, nil, req)
	if err != nil {
		return nil, fmt.Errorf("查询考试记录失败: %w", err)
	}

	analysis := &dto.ExamAnalysisResponse{
		BasicInfo: dto.ExamBasicInfo{
			ExamID:        exam.ID,
			Title:         exam.Title,
			TotalStudents: len(records),
		},
	}

	if len(records) == 0 {
		return analysis, nil
	}

	// 统计分数
	var scores []float64
	var totalDuration int
	submitCount := 0

	for _, r := range records {
		if r.Status >= 1 {
			scores = append(scores, r.Score)
			submitCount++
			totalDuration += r.DurationUsed
		}
	}

	analysis.BasicInfo.SubmitCount = submitCount
	if submitCount > 0 {
		analysis.BasicInfo.AvgDuration = totalDuration / submitCount
	}

	if len(scores) > 0 {
		// 计算统计指标
		stats := dto.ScoreStatistics{}
		var total, max, min, passCount, excellentCount float64
		min = scores[0]

		for _, s := range scores {
			total += s
			if s > max {
				max = s
			}
			if s < min {
				min = s
			}
			if s >= exam.PassScore {
				passCount++
			}
			if s >= exam.TotalScore*0.9 {
				excellentCount++
			}
		}

		stats.AvgScore = total / float64(len(scores))
		stats.MaxScore = max
		stats.MinScore = min
		stats.PassRate = passCount / float64(len(scores)) * 100
		stats.ExcellentRate = excellentCount / float64(len(scores)) * 100

		// 分数分布
		distMap := make(map[string]int)
		for _, s := range scores {
			key := fmt.Sprintf("%d-%d", int(s)/10*10, int(s)/10*10+10)
			distMap[key]++
		}
		for k, v := range distMap {
			stats.Distribution = append(stats.Distribution, dto.ScoreDist{
				Range: k,
				Count: v,
			})
		}

		analysis.ScoreStats = stats
	}

	// 题目级别分析
	paperQuestions, err := examRepo.GetPaperQuestions(exam.PaperID)
	if err == nil {
		questionStats := make([]dto.QuestionStat, 0, len(paperQuestions))
		for _, pq := range paperQuestions {
			var correctCount, totalAnswered int
			var totalQuestionScore float64

			for _, r := range records {
				if r.Status < 1 {
					continue
				}
				answers, _ := examRepo.GetAnswerRecords(r.ID)
				for _, a := range answers {
					if a.QuestionID == pq.QuestionID {
						totalAnswered++
						totalQuestionScore += a.Score
						if a.IsCorrect {
							correctCount++
						}
						break
					}
				}
			}

			qs := dto.QuestionStat{
				QuestionID: pq.QuestionID,
				Score:      pq.Score,
			}

			if totalAnswered > 0 {
				qs.CorrectRate = float64(correctCount) / float64(totalAnswered) * 100
				qs.AvgScore = totalQuestionScore / float64(totalAnswered)
			}

			// 获取题目信息用于标题和类型
			question, err := questionRepo.FindByID(pq.QuestionID)
			if err == nil {
				qs.Title = question.Title
				qs.Type = question.Type
				qs.Difficulty = question.Difficulty
			}

			// 计算区分度（高分组正确率 - 低分组正确率）
			if totalAnswered >= 4 {
				// 按总分排序，取前27%为高分组，后27%为低分组
				highGroupCorrect := 0
				lowGroupCorrect := 0
				highGroupTotal := 0
				lowGroupTotal := 0
				cutoff := int(float64(totalAnswered) * 0.27)

				// 收集每人的该题得分和总分
				type studentResult struct {
					isCorrect bool
					totalScore float64
				}
				var results []studentResult
				for _, r := range records {
					if r.Status < 1 {
						continue
					}
					answers, _ := examRepo.GetAnswerRecords(r.ID)
					for _, a := range answers {
						if a.QuestionID == pq.QuestionID {
							results = append(results, studentResult{
								isCorrect: a.IsCorrect,
								totalScore: r.Score,
							})
							break
						}
					}
				}

				// 按总分排序
				for i := 0; i < len(results); i++ {
					for j := i + 1; j < len(results); j++ {
						if results[j].totalScore > results[i].totalScore {
							results[i], results[j] = results[j], results[i]
						}
					}
				}

				for i, r := range results {
					if i < cutoff {
						highGroupTotal++
						if r.isCorrect {
							highGroupCorrect++
						}
					}
					if i >= len(results)-cutoff {
						lowGroupTotal++
						if r.isCorrect {
							lowGroupCorrect++
						}
					}
				}

				if highGroupTotal > 0 && lowGroupTotal > 0 {
					highRate := float64(highGroupCorrect) / float64(highGroupTotal)
					lowRate := float64(lowGroupCorrect) / float64(lowGroupTotal)
					qs.Discrimination = highRate - lowRate
				}
			}

			questionStats = append(questionStats, qs)
		}
		analysis.QuestionStats = questionStats
	}

	return analysis, nil
}
