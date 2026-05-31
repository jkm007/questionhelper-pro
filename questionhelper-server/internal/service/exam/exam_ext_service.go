package exam

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	examRepo "questionhelper-server/internal/repository/exam"
	"questionhelper-server/pkg/database"
	"questionhelper-server/pkg/logger"
)

// ==================== 试卷共享 ====================

// SharePaper 试卷共享
func SharePaper(paperID, sharerID uint, req *dto.SharePaperRequest) error {
	// 检查试卷是否存在
	_, err := examRepo.FindPaperByID(paperID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("试卷不存在")
		}
		return fmt.Errorf("查询试卷失败: %w", err)
	}

	// 检查是否已共享
	existing, err := examRepo.FindPaperShare(paperID, req.TargetID, req.TargetType)
	if err == nil && existing.ID > 0 {
		return errors.New("已共享给该目标")
	}

	share := &model.PaperShare{
		PaperID:    paperID,
		SharerID:   sharerID,
		TargetID:   req.TargetID,
		TargetType: req.TargetType,
		Permission: req.Permission,
		Status:     1,
	}

	if err := examRepo.CreatePaperShare(share); err != nil {
		return fmt.Errorf("共享试卷失败: %w", err)
	}

	logger.Infof("用户 %d 共享试卷 %d 给目标 %d", sharerID, paperID, req.TargetID)
	return nil
}

// ==================== 试卷收藏 ====================

// FavoritePaper 收藏/取消收藏试卷
func FavoritePaper(userID, paperID uint, note string) (string, error) {
	// 检查试卷是否存在
	_, err := examRepo.FindPaperByID(paperID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("试卷不存在")
		}
		return "", fmt.Errorf("查询试卷失败: %w", err)
	}

	// 检查是否已收藏
	existing, err := examRepo.FindPaperFavorite(userID, paperID)
	if err == nil && existing.ID > 0 {
		// 已收藏则取消
		if err := examRepo.DeletePaperFavorite(userID, paperID); err != nil {
			return "", fmt.Errorf("取消收藏失败: %w", err)
		}
		logger.Infof("用户 %d 取消收藏试卷 %d", userID, paperID)
		return "已取消收藏", nil
	}

	// 未收藏则添加
	fav := &model.PaperFavorite{
		UserID:  userID,
		PaperID: paperID,
		Note:    note,
	}

	if err := examRepo.CreatePaperFavorite(fav); err != nil {
		return "", fmt.Errorf("收藏试卷失败: %w", err)
	}

	logger.Infof("用户 %d 收藏试卷 %d", userID, paperID)
	return "收藏成功", nil
}

// ==================== 试卷导入 ====================

// ImportPaper 导入试卷(从导出格式JSON)
func ImportPaper(creatorID uint, paper *dto.ExportPaperResponse) (*dto.PaperInfo, error) {
	if paper == nil {
		return nil, errors.New("导入数据为空")
	}

	newPaper := &model.Paper{
		Title:       paper.Title,
		Description: paper.Description,
		TotalScore:  paper.TotalScore,
		TotalCount:  len(paper.Questions),
		Type:        1,
		Status:      0,
		CreatorID:   creatorID,
	}

	if err := examRepo.CreatePaper(newPaper); err != nil {
		return nil, fmt.Errorf("创建试卷失败: %w", err)
	}

	// 导入题目关联
	if len(paper.Questions) > 0 {
		paperQuestions := make([]model.PaperQuestion, 0, len(paper.Questions))
		for i, q := range paper.Questions {
			paperQuestions = append(paperQuestions, model.PaperQuestion{
				PaperID:    newPaper.ID,
				QuestionID: uint(q.Sort), // 使用Sort作为QuestionID的映射
				Score:      q.Score,
				Sort:       i,
			})
		}

		if err := examRepo.AddPaperQuestions(newPaper.ID, paperQuestions); err != nil {
			logger.Errorf("导入试卷题目失败: %v", err)
		}
	}

	logger.Infof("导入试卷成功: %d, 题目数: %d", newPaper.ID, len(paper.Questions))
	info := toPaperInfo(newPaper)
	return &info, nil
}

// ==================== 模板列表 ====================

// ListTemplates 模板列表
func ListTemplates(req *dto.TemplateListRequest) ([]dto.PaperInfo, int64, error) {
	papers, total, err := examRepo.ListTemplatesPaged(req)
	if err != nil {
		return nil, 0, fmt.Errorf("查询模板列表失败: %w", err)
	}

	list := make([]dto.PaperInfo, 0, len(papers))
	for _, p := range papers {
		list = append(list, toPaperInfo(&p))
	}
	return list, total, nil
}

// ==================== 考试延长 ====================

// ExtendExam 延长考试时间
func ExtendExam(examID, operatorID uint, req *dto.ExtendExamRequest) error {
	exam, err := examRepo.FindExamByID(examID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("考试不存在")
		}
		return fmt.Errorf("查询考试失败: %w", err)
	}

	if exam.Status != 1 {
		return errors.New("只能延长进行中的考试")
	}

	if req.Minutes <= 0 {
		return errors.New("延长时间必须大于0")
	}

	// 记录原始结束时间（仅首次）
	if exam.OriginalEndTime == nil {
		originalEndTime := exam.EndTime
		exam.OriginalEndTime = &originalEndTime
	}

	oldEndTime := exam.EndTime
	newEndTime := oldEndTime.Add(time.Duration(req.Minutes) * time.Minute)
	exam.EndTime = newEndTime
	exam.ExtendReason = req.Reason

	// 记录延期记录
	extension := &model.ExamExtension{
		ExamID:        examID,
		OperatorID:    operatorID,
		OldEndTime:    oldEndTime,
		NewEndTime:    newEndTime,
		ExtendMinutes: req.Minutes,
		Reason:        req.Reason,
		Status:        1,
	}

	if err := examRepo.CreateExamExtension(extension); err != nil {
		logger.Errorf("创建延期记录失败: %v", err)
	}

	if err := examRepo.UpdateExam(exam); err != nil {
		return fmt.Errorf("延长考试失败: %w", err)
	}

	logger.Infof("考试 %d 延长 %d 分钟，新结束时间: %v", examID, req.Minutes, newEndTime)
	return nil
}

// ==================== 考试暂停/恢复 ====================

// PauseExam 暂停考试
func PauseExam(examID, operatorID uint, reason string) error {
	exam, err := examRepo.FindExamByID(examID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("考试不存在")
		}
		return fmt.Errorf("查询考试失败: %w", err)
	}

	if exam.Status != 1 {
		return errors.New("只能暂停进行中的考试")
	}

	if exam.StatusPause {
		return errors.New("考试已经处于暂停状态")
	}

	exam.StatusPause = true
	now := time.Now()

	pause := &model.ExamPause{
		ExamID:     examID,
		OperatorID: operatorID,
		Action:     1, // 暂停
		Reason:     reason,
		PausedAt:   &now,
	}

	if err := examRepo.CreateExamPause(pause); err != nil {
		logger.Errorf("创建暂停记录失败: %v", err)
	}

	if err := examRepo.UpdateExam(exam); err != nil {
		return fmt.Errorf("暂停考试失败: %w", err)
	}

	logger.Infof("考试 %d 已暂停，操作人: %d", examID, operatorID)
	return nil
}

// ResumeExam 恢复考试
func ResumeExam(examID, operatorID uint, reason string) error {
	exam, err := examRepo.FindExamByID(examID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("考试不存在")
		}
		return fmt.Errorf("查询考试失败: %w", err)
	}

	if exam.Status != 1 {
		return errors.New("只能恢复进行中的考试")
	}

	if !exam.StatusPause {
		return errors.New("考试未处于暂停状态")
	}

	exam.StatusPause = false
	now := time.Now()

	// 更新最近的暂停记录
	lastPause, err := examRepo.FindLastPause(examID)
	if err == nil && lastPause.Action == 1 && lastPause.ResumedAt == nil {
		lastPause.Action = 2 // 恢复
		lastPause.ResumedAt = &now
		if lastPause.PausedAt != nil {
			lastPause.Duration = int(now.Sub(*lastPause.PausedAt).Seconds())
		}
		examRepo.UpdateExamPause(lastPause)
	} else {
		pause := &model.ExamPause{
			ExamID:     examID,
			OperatorID: operatorID,
			Action:     2, // 恢复
			Reason:     reason,
			ResumedAt:  &now,
		}
		examRepo.CreateExamPause(pause)
	}

	if err := examRepo.UpdateExam(exam); err != nil {
		return fmt.Errorf("恢复考试失败: %w", err)
	}

	logger.Infof("考试 %d 已恢复，操作人: %d", examID, operatorID)
	return nil
}

// ==================== 成绩复核 ====================

// SubmitScoreReview 申请成绩复核
func SubmitScoreReview(userID uint, examID uint, req *dto.ScoreReviewRequest) error {
	// 获取考试记录
	record, err := examRepo.FindExamRecordByUser(examID, userID)
	if err != nil {
		return errors.New("未找到考试记录")
	}

	if record.Status < 1 {
		return errors.New("考试尚未提交，无法申请复核")
	}

	// 检查是否已有待处理的复核
	existing, err := examRepo.FindScoreReviewByRecordAndUser(record.ID, userID)
	if err == nil && existing.ID > 0 {
		return errors.New("已存在待处理的复核申请")
	}

	review := &model.ScoreReview{
		RecordID: record.ID,
		UserID:   userID,
		Reason:   req.Reason,
		OldScore: record.Score,
		Status:   0, // 待复核
	}

	if err := examRepo.CreateScoreReview(review); err != nil {
		return fmt.Errorf("提交复核申请失败: %w", err)
	}

	logger.Infof("用户 %d 提交考试 %d 的成绩复核申请", userID, examID)
	return nil
}

// ListScoreReviews 复核申请列表
func ListScoreReviews(req *dto.ScoreReviewListRequest) ([]dto.ScoreReviewInfo, int64, error) {
	reviews, total, err := examRepo.ListScoreReviews(req)
	if err != nil {
		return nil, 0, fmt.Errorf("查询复核列表失败: %w", err)
	}

	list := make([]dto.ScoreReviewInfo, 0, len(reviews))
	for _, r := range reviews {
		info := dto.ScoreReviewInfo{
			ID:         r.ID,
			RecordID:   r.RecordID,
			UserID:     r.UserID,
			Reason:     r.Reason,
			OldScore:   r.OldScore,
			NewScore:   r.NewScore,
			ReviewerID: r.ReviewerID,
			ReviewNote: r.ReviewNote,
			Status:     r.Status,
			ReviewedAt: r.ReviewedAt,
			CreatedAt:  r.CreatedAt,
		}

		// 获取用户信息
		var user model.User
		if err := database.DB.First(&user, r.UserID).Error; err == nil {
			info.Username = user.Username
			info.Nickname = user.Nickname
		}

		// 获取考试信息
		var record model.ExamRecord
		if err := database.DB.First(&record, r.RecordID).Error; err == nil {
			info.ExamID = record.ExamID
			var examModel model.Exam
			if err := database.DB.First(&examModel, record.ExamID).Error; err == nil {
				info.ExamTitle = examModel.Title
			}
		}

		list = append(list, info)
	}

	return list, total, nil
}

// HandleScoreReview 处理复核
func HandleScoreReview(reviewID, reviewerID uint, req *dto.HandleScoreReviewRequest) error {
	review, err := examRepo.FindScoreReviewByID(reviewID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("复核申请不存在")
		}
		return fmt.Errorf("查询复核失败: %w", err)
	}

	if review.Status != 0 {
		return errors.New("该复核已处理")
	}

	review.Status = req.Status
	review.ReviewerID = &reviewerID
	review.ReviewNote = req.ReviewNote
	now := time.Now()
	review.ReviewedAt = &now

	if req.Status == 1 { // 已复核（通过）
		review.NewScore = req.NewScore

		// 更新考试记录分数
		record, err := examRepo.FindExamRecord(review.RecordID)
		if err == nil {
			record.Score = req.NewScore
			examRepo.UpdateExamRecord(record)

			// 更新排名
			examRepo.BuildExamRankings(record.ExamID)
		}
	}

	if err := examRepo.UpdateScoreReview(review); err != nil {
		return fmt.Errorf("处理复核失败: %w", err)
	}

	logger.Infof("复核 %d 已处理，状态: %d，处理人: %d", reviewID, req.Status, reviewerID)
	return nil
}

// ==================== 考试公告 ====================

// CreateExamNotice 创建考试公告
func CreateExamNotice(examID, creatorID uint, req *dto.CreateExamNoticeRequest) error {
	_, err := examRepo.FindExamByID(examID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("考试不存在")
		}
		return fmt.Errorf("查询考试失败: %w", err)
	}

	notice := &model.ExamNotice{
		ExamID:    examID,
		Title:     req.Title,
		Content:   req.Content,
		Priority:  req.Priority,
		IsPinned:  req.IsPinned,
		CreatorID: creatorID,
		Status:    1,
	}

	if err := examRepo.CreateExamNotice(notice); err != nil {
		return fmt.Errorf("创建公告失败: %w", err)
	}

	logger.Infof("创建考试 %d 公告成功: %s", examID, req.Title)
	return nil
}

// ListExamNotices 考试公告列表
func ListExamNotices(examID uint, page, pageSize int) ([]dto.ExamNoticeInfo, int64, error) {
	notices, total, err := examRepo.ListExamNotices(examID, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("查询公告列表失败: %w", err)
	}

	list := make([]dto.ExamNoticeInfo, 0, len(notices))
	for _, n := range notices {
		list = append(list, dto.ExamNoticeInfo{
			ID:        n.ID,
			ExamID:    n.ExamID,
			Title:     n.Title,
			Content:   n.Content,
			Priority:  n.Priority,
			IsPinned:  n.IsPinned,
			CreatorID: n.CreatorID,
			Status:    n.Status,
			CreatedAt: n.CreatedAt,
		})
	}

	return list, total, nil
}

// ==================== 考试防作弊 ====================

// ReportSwitchScreen 上报切屏
func ReportSwitchScreen(recordID, examID, userID uint, detail, ip, userAgent string) error {
	// 更新切屏次数
	record, err := examRepo.FindExamRecord(recordID)
	if err != nil {
		return errors.New("考试记录不存在")
	}

	record.SwitchCount++
	examRepo.UpdateExamRecord(record)

	// 记录异常
	warning := &model.ExamWarning{
		RecordID:  recordID,
		ExamID:    examID,
		UserID:    userID,
		Type:      "switch_screen",
		Detail:    fmt.Sprintf("切屏次数: %d, %s", record.SwitchCount, detail),
		IP:        ip,
		UserAgent: userAgent,
	}

	if err := examRepo.CreateWarning(warning); err != nil {
		logger.Errorf("记录切屏异常失败: %v", err)
	}

	logger.Infof("用户 %d 切屏上报，考试 %d，累计 %d 次", userID, examID, record.SwitchCount)
	return nil
}

// ResumeExamForStudent 断线续考
func ResumeExamForStudent(examID, userID uint, ip string) (*dto.ExamRecordInfo, error) {
	exam, err := examRepo.FindExamByID(examID)
	if err != nil {
		return nil, errors.New("考试不存在")
	}

	if exam.Status != 1 {
		return nil, errors.New("考试未开放")
	}

	if exam.StatusPause {
		return nil, errors.New("考试已暂停，请等待恢复")
	}

	// 查找进行中的记录
	record, err := examRepo.FindExamRecordByUser(examID, userID)
	if err != nil {
		return nil, errors.New("未找到考试记录，无法续考")
	}

	if record.Status != 0 {
		return nil, errors.New("考试已提交，无法续考")
	}

	// 更新IP
	record.IP = ip
	examRepo.UpdateExamRecord(record)

	info := &dto.ExamRecordInfo{
		ID:         record.ID,
		ExamID:     record.ExamID,
		UserID:     record.UserID,
		Score:      record.Score,
		Status:     record.Status,
		StartTime:  record.StartTime,
		SubmitTime: record.SubmitTime,
	}

	logger.Infof("用户 %d 断线续考 %d", userID, examID)
	return info, nil
}

// ==================== 考试查询增强 ====================

// ListUpcomingExams 即将开始的考试
func ListUpcomingExams(userID uint, req *dto.PageRequest) ([]dto.ExamUpcomingInfo, int64, error) {
	exams, total, err := examRepo.ListUpcomingExams(userID, req)
	if err != nil {
		return nil, 0, fmt.Errorf("查询即将开始的考试失败: %w", err)
	}

	now := time.Now()
	list := make([]dto.ExamUpcomingInfo, 0, len(exams))
	for _, e := range exams {
		list = append(list, dto.ExamUpcomingInfo{
			ID:          e.ID,
			Title:       e.Title,
			Description: e.Description,
			StartTime:   e.StartTime,
			EndTime:     e.EndTime,
			Duration:    e.Duration,
			TotalScore:  e.TotalScore,
			PassScore:   e.PassScore,
			Status:      e.Status,
			TimeLeft:    int(e.StartTime.Sub(now).Seconds()),
		})
	}

	return list, total, nil
}

// GetExamRankings 成绩排名
func GetExamRankings(examID uint, page, pageSize int) ([]dto.ExamRankingInfo, int64, error) {
	// 先检查是否有排名数据，没有则构建
	rankings, total, err := examRepo.ListExamRankings(examID, page, pageSize)
	if err != nil || total == 0 {
		// 尝试构建排名
		if buildErr := examRepo.BuildExamRankings(examID); buildErr != nil {
			return nil, 0, fmt.Errorf("构建排名失败: %w", buildErr)
		}
		rankings, total, err = examRepo.ListExamRankings(examID, page, pageSize)
		if err != nil {
			return nil, 0, fmt.Errorf("查询排名失败: %w", err)
		}
	}

	list := make([]dto.ExamRankingInfo, 0, len(rankings))
	for _, r := range rankings {
		info := dto.ExamRankingInfo{
			ID:           r.ID,
			ExamID:       r.ExamID,
			UserID:       r.UserID,
			Score:        r.Score,
			ObjScore:     r.ObjScore,
			SubjScore:    r.SubjScore,
			RankPos:      r.RankPos,
			DurationUsed: r.DurationUsed,
			Accuracy:     r.Accuracy,
			SubmitTime:   r.SubmitTime.Format("2006-01-02 15:04:05"),
		}

		// 获取用户信息
		var user model.User
		if err := database.DB.First(&user, r.UserID).Error; err == nil {
			info.Username = user.Username
			info.Nickname = user.Nickname
		}

		list = append(list, info)
	}

	return list, total, nil
}

// GetExamScores 考试成绩列表
func GetExamScores(examID uint, req *dto.PageRequest) ([]dto.ExamRecordInfo, int64, error) {
	records, total, err := examRepo.ListExamRecords(&examID, nil, req)
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

// ExportExamScores 导出考试成绩
func ExportExamScores(examID uint) ([]dto.ExamRecordInfo, error) {
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

// GetExamStatistics 考试统计
func GetExamStatistics(examID uint) (*dto.ExamStatisticsResponse, error) {
	examModel, err := examRepo.FindExamByID(examID)
	if err != nil {
		return nil, fmt.Errorf("查询考试失败: %w", err)
	}

	req := &dto.PageRequest{Page: 1, PageSize: 10000}
	records, _, err := examRepo.ListExamRecords(&examID, nil, req)
	if err != nil {
		return nil, fmt.Errorf("查询考试记录失败: %w", err)
	}

	stats := &dto.ExamStatisticsResponse{
		BasicInfo: dto.ExamBasicInfo{
			ExamID:        examModel.ID,
			Title:         examModel.Title,
			TotalStudents: len(records),
		},
	}

	if len(records) == 0 {
		return stats, nil
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

	stats.BasicInfo.SubmitCount = submitCount
	if submitCount > 0 {
		stats.BasicInfo.AvgDuration = totalDuration / submitCount
	}

	if len(scores) > 0 {
		scoreStats := dto.ScoreStatistics{}
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
			if s >= examModel.PassScore {
				passCount++
			}
			if s >= examModel.TotalScore*0.9 {
				excellentCount++
			}
		}

		scoreStats.AvgScore = total / float64(len(scores))
		scoreStats.MaxScore = max
		scoreStats.MinScore = min
		scoreStats.PassRate = passCount / float64(len(scores)) * 100
		scoreStats.ExcellentRate = excellentCount / float64(len(scores)) * 100

		// 分数分布
		distMap := make(map[string]int)
		for _, s := range scores {
			key := fmt.Sprintf("%d-%d", int(s)/10*10, int(s)/10*10+10)
			distMap[key]++
		}
		for k, v := range distMap {
			scoreStats.Distribution = append(scoreStats.Distribution, dto.ScoreDist{
				Range: k,
				Count: v,
			})
		}

		stats.ScoreStats = scoreStats
	}

	// Top 学生
	pager := &dto.PageRequest{Page: 1, PageSize: 10}
	topRecords, _, _ := examRepo.ListExamRecords(&examID, nil, pager)
	for i, r := range topRecords {
		if r.Status >= 1 {
			student := dto.StudentScore{
				UserID:   r.UserID,
				Score:    r.Score,
				Duration: r.DurationUsed,
				Rank:     i + 1,
			}
			var user model.User
			if err := database.DB.First(&user, r.UserID).Error; err == nil {
				student.Username = user.Username
				student.Nickname = user.Nickname
			}
			stats.TopStudents = append(stats.TopStudents, student)
		}
	}

	// 异常统计
	warningCount, _ := examRepo.CountWarningsByExamID(examID)
	stats.WarningCount = int(warningCount)

	return stats, nil
}

// SubmitExamFeedback 提交考试反馈
func SubmitExamFeedback(examID, userID uint, req *dto.FeedbackRequest) error {
	_, err := examRepo.FindExamByID(examID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("考试不存在")
		}
		return fmt.Errorf("查询考试失败: %w", err)
	}

	feedback := &model.ExamFeedback{
		ExamID:       examID,
		UserID:       userID,
		FeedbackType: 1,
		Content:      req.Content,
		Status:       0,
	}

	if err := examRepo.CreateExamFeedback(feedback); err != nil {
		return fmt.Errorf("提交反馈失败: %w", err)
	}

	logger.Infof("用户 %d 提交考试 %d 反馈", userID, examID)
	return nil
}

// FindRecordByUser 查找用户的考试记录
func FindRecordByUser(examID, userID uint) (*model.ExamRecord, error) {
	return examRepo.FindExamRecordByUser(examID, userID)
}
