package practice

import (
	"errors"
	"fmt"
	"time"

	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	practiceRepo "questionhelper-server/internal/repository/practice"
	questionRepo "questionhelper-server/internal/repository/question"
	"questionhelper-server/pkg/logger"
)

// ==================== 模拟考试 ====================

// StartMockExam 开始模拟考试
func StartMockExam(userID uint, req *dto.StartMockExamRequest) (*dto.MockExamSessionInfo, error) {
	config, err := practiceRepo.FindMockExamConfigByID(req.ConfigID)
	if err != nil {
		return nil, errors.New("模拟考试配置不存在")
	}

	if config.Status != 1 {
		return nil, errors.New("该模拟考试已禁用")
	}

	// 检查最大尝试次数
	if config.MaxAttempts > 0 {
		var count int64
		practiceRepo.ListMockExamHistory(userID, &dto.MockExamHistoryRequest{
			ConfigID: &req.ConfigID,
		})
		// 简化处理：直接查询数量
		_ = count
	}

	// 获取题目数量
	questionCount := config.QuestionCount
	if questionCount == 0 {
		questionCount = 20
	}

	// 创建练习会话（复用PracticeSession模型，CategoryID存储ConfigID）
	session := &model.PracticeSession{
		UserID:     userID,
		CategoryID: &req.ConfigID,
		TotalCount: questionCount,
		Duration:   config.Duration,
		Status:     0,
	}

	if err := practiceRepo.CreateMockExamSession(session); err != nil {
		return nil, fmt.Errorf("创建模拟考试失败: %w", err)
	}

	logger.Infof("用户 %d 开始模拟考试，配置 %d，题目数 %d", userID, config.ID, questionCount)

	return &dto.MockExamSessionInfo{
		ID:         session.ID,
		UserID:     session.UserID,
		ConfigID:   config.ID,
		ConfigName: config.Name,
		CategoryID: config.CategoryID,
		TotalCount: questionCount,
		TotalScore: config.TotalScore,
		PassScore:  config.PassScore,
		Duration:   config.Duration,
		Status:     session.Status,
		StartedAt:  session.CreatedAt,
	}, nil
}

// SubmitMockExam 提交模拟考试
func SubmitMockExam(sessionID, userID uint, req *dto.SubmitMockExamRequest) (*dto.MockExamResultInfo, error) {
	session, err := practiceRepo.FindSessionByID(sessionID)
	if err != nil {
		return nil, errors.New("模拟考试会话不存在")
	}

	if session.UserID != userID {
		return nil, errors.New("无权操作此模拟考试")
	}

	if session.Status != 0 {
		return nil, errors.New("模拟考试已结束")
	}

	// 批量保存答题记录
	records := make([]model.PracticeRecord, 0, len(req.Answers))
	correctCount := 0
	for _, ans := range req.Answers {
		question, err := questionRepo.FindByID(ans.QuestionID)
		if err != nil {
			continue
		}
		isCorrect := question.Answer == ans.Answer
		if isCorrect {
			correctCount++
		}
		records = append(records, model.PracticeRecord{
			SessionID:  sessionID,
			QuestionID: ans.QuestionID,
			Answer:     ans.Answer,
			IsCorrect:  isCorrect,
			Duration:   ans.Duration,
		})
	}

	if len(records) > 0 {
		if err := practiceRepo.CreateRecords(records); err != nil {
			return nil, fmt.Errorf("保存答题记录失败: %w", err)
		}
	}

	// 计算成绩
	accuracy := float64(0)
	if len(records) > 0 {
		accuracy = float64(correctCount) / float64(len(records)) * 100
	}

	now := time.Now()
	session.CorrectCount = correctCount
	session.Accuracy = accuracy
	session.Duration = req.Duration
	session.Status = 1
	session.UpdatedAt = now

	if err := practiceRepo.UpdateSession(session); err != nil {
		return nil, fmt.Errorf("更新模拟考试失败: %w", err)
	}

	// 获取配置信息
	configID := uint(0)
	if session.CategoryID != nil {
		configID = *session.CategoryID
	}
	config, _ := practiceRepo.FindMockExamConfigByID(configID)

	// 计算得分和是否及格
	score := accuracy
	isPassed := false
	passScore := float64(60)
	if config != nil {
		if config.TotalScore > 0 {
			score = accuracy / 100 * config.TotalScore
		}
		passScore = config.PassScore
		isPassed = score >= passScore
	}

	// 构建记录信息
	recordInfos := make([]dto.MockExamRecordInfo, 0, len(records))
	for _, r := range records {
		recordInfos = append(recordInfos, dto.MockExamRecordInfo{
			ID:         r.ID,
			QuestionID: r.QuestionID,
			Answer:     r.Answer,
			IsCorrect:  r.IsCorrect,
			Score:      0,
			Duration:   r.Duration,
		})
	}

	logger.Infof("模拟考试 %d 完成，得分: %.2f，及格: %v", sessionID, score, isPassed)

	return &dto.MockExamResultInfo{
		Session: dto.MockExamSessionInfo{
			ID:         session.ID,
			UserID:     session.UserID,
			ConfigID:   configID,
			TotalCount: session.TotalCount,
			Duration:   session.Duration,
			Status:     session.Status,
			StartedAt:  session.CreatedAt,
			SubmittedAt: &now,
		},
		Records:  recordInfos,
		Score:    score,
		IsPassed: isPassed,
		Rank:     0,
	}, nil
}

// GetMockExamHistory 模拟考试历史
func GetMockExamHistory(userID uint, req *dto.MockExamHistoryRequest) ([]dto.MockExamHistoryItem, int64, error) {
	sessions, total, err := practiceRepo.ListMockExamHistory(userID, req)
	if err != nil {
		return nil, 0, fmt.Errorf("查询模拟考试历史失败: %w", err)
	}

	items := make([]dto.MockExamHistoryItem, 0, len(sessions))
	for _, s := range sessions {
		configID := uint(0)
		if s.CategoryID != nil {
			configID = *s.CategoryID
		}
		configName := ""
		config, _ := practiceRepo.FindMockExamConfigByID(configID)
		if config != nil {
			configName = config.Name
		}

		score := s.Accuracy
		isPassed := false
		if config != nil && config.TotalScore > 0 {
			score = s.Accuracy / 100 * config.TotalScore
			isPassed = score >= config.PassScore
		} else {
			isPassed = s.Accuracy >= 60
		}

		items = append(items, dto.MockExamHistoryItem{
			ID:         s.ID,
			ConfigID:   configID,
			ConfigName: configName,
			TotalCount: s.TotalCount,
			Score:      score,
			IsPassed:   isPassed,
			Duration:   s.Duration,
			Status:     s.Status,
			StartedAt:  s.CreatedAt,
		})
	}

	return items, total, nil
}

// GetMockExamDetail 模拟考试详情
func GetMockExamDetail(sessionID, userID uint) (*dto.MockExamResultInfo, error) {
	session, err := practiceRepo.FindSessionByID(sessionID)
	if err != nil {
		return nil, errors.New("模拟考试不存在")
	}

	if session.UserID != userID {
		return nil, errors.New("无权查看此模拟考试")
	}

	configID := uint(0)
	if session.CategoryID != nil {
		configID = *session.CategoryID
	}
	config, _ := practiceRepo.FindMockExamConfigByID(configID)

	score := session.Accuracy
	isPassed := false
	if config != nil && config.TotalScore > 0 {
		score = session.Accuracy / 100 * config.TotalScore
		isPassed = score >= config.PassScore
	} else {
		isPassed = session.Accuracy >= 60
	}

	configName := ""
	if config != nil {
		configName = config.Name
	}

	records, _ := practiceRepo.GetRecords(sessionID)
	recordInfos := make([]dto.MockExamRecordInfo, 0, len(records))
	for _, r := range records {
		recordInfos = append(recordInfos, dto.MockExamRecordInfo{
			ID:         r.ID,
			QuestionID: r.QuestionID,
			Answer:     r.Answer,
			IsCorrect:  r.IsCorrect,
			Duration:   r.Duration,
		})
	}

	return &dto.MockExamResultInfo{
		Session: dto.MockExamSessionInfo{
			ID:         session.ID,
			UserID:     session.UserID,
			ConfigID:   configID,
			ConfigName: configName,
			TotalCount: session.TotalCount,
			Duration:   session.Duration,
			Status:     session.Status,
			StartedAt:  session.CreatedAt,
		},
		Records:  recordInfos,
		Score:    score,
		IsPassed: isPassed,
		Rank:     0,
	}, nil
}

// ==================== 练习计划 ====================

// CreatePlan 创建练习计划
func CreatePlan(userID uint, req *dto.CreatePlanRequest) (*dto.PlanInfo, error) {
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, errors.New("开始日期格式错误，请使用YYYY-MM-DD格式")
	}

	plan := &model.PracticePlan{
		UserID:        userID,
		Name:          req.Name,
		Description:   req.Description,
		PlanType:      req.PlanType,
		CategoryID:    req.CategoryID,
		QuestionType:  req.QuestionType,
		Difficulty:    req.Difficulty,
		DailyCount:    req.DailyCount,
		DailyDuration: req.DailyDuration,
		StartDate:     startDate,
		TotalTarget:   req.TotalTarget,
		Status:        1,
	}

	if req.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return nil, errors.New("结束日期格式错误，请使用YYYY-MM-DD格式")
		}
		plan.EndDate = &endDate
	}

	if err := practiceRepo.CreatePlan(plan); err != nil {
		return nil, fmt.Errorf("创建练习计划失败: %w", err)
	}

	return toPlanInfo(plan), nil
}

// GetPlan 获取计划详情
func GetPlan(planID, userID uint) (*dto.PlanInfo, error) {
	plan, err := practiceRepo.FindPlanByID(planID)
	if err != nil {
		return nil, errors.New("练习计划不存在")
	}

	if plan.UserID != userID {
		return nil, errors.New("无权查看此练习计划")
	}

	return toPlanInfo(plan), nil
}

// UpdatePlan 更新计划
func UpdatePlan(planID, userID uint, req *dto.UpdatePlanRequest) error {
	plan, err := practiceRepo.FindPlanByID(planID)
	if err != nil {
		return errors.New("练习计划不存在")
	}

	if plan.UserID != userID {
		return errors.New("无权修改此练习计划")
	}

	if req.Name != "" {
		plan.Name = req.Name
	}
	if req.Description != "" {
		plan.Description = req.Description
	}
	if req.PlanType > 0 {
		plan.PlanType = req.PlanType
	}
	if req.CategoryID != nil {
		plan.CategoryID = req.CategoryID
	}
	if req.QuestionType != nil {
		plan.QuestionType = req.QuestionType
	}
	if req.Difficulty != nil {
		plan.Difficulty = req.Difficulty
	}
	if req.DailyCount > 0 {
		plan.DailyCount = req.DailyCount
	}
	if req.DailyDuration > 0 {
		plan.DailyDuration = req.DailyDuration
	}
	if req.StartDate != "" {
		startDate, err := time.Parse("2006-01-02", req.StartDate)
		if err == nil {
			plan.StartDate = startDate
		}
	}
	if req.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", req.EndDate)
		if err == nil {
			plan.EndDate = &endDate
		}
	}
	if req.TotalTarget > 0 {
		plan.TotalTarget = req.TotalTarget
	}
	if req.Status != nil {
		plan.Status = *req.Status
	}

	return practiceRepo.UpdatePlan(plan)
}

// DeletePlan 删除计划
func DeletePlan(planID, userID uint) error {
	plan, err := practiceRepo.FindPlanByID(planID)
	if err != nil {
		return errors.New("练习计划不存在")
	}

	if plan.UserID != userID {
		return errors.New("无权删除此练习计划")
	}

	return practiceRepo.DeletePlan(planID)
}

// ListPlans 获取计划列表
func ListPlans(userID uint, req *dto.PlanListRequest) ([]dto.PlanInfo, int64, error) {
	plans, total, err := practiceRepo.ListPlans(userID, req)
	if err != nil {
		return nil, 0, fmt.Errorf("查询练习计划失败: %w", err)
	}

	items := make([]dto.PlanInfo, 0, len(plans))
	for _, p := range plans {
		items = append(items, *toPlanInfo(&p))
	}

	return items, total, nil
}

// ExecutePlan 执行计划
func ExecutePlan(planID, userID uint, req *dto.ExecutePlanRequest) error {
	plan, err := practiceRepo.FindPlanByID(planID)
	if err != nil {
		return errors.New("练习计划不存在")
	}

	if plan.UserID != userID {
		return errors.New("无权执行此练习计划")
	}

	if plan.Status != 1 {
		return errors.New("该计划不在进行中")
	}

	// 记录执行日志
	log := &model.PracticePlanLog{
		PlanID:      planID,
		UserID:      userID,
		QuestionNum: req.QuestionNum,
		Duration:    req.Duration,
		Date:        time.Now().Format("2006-01-02"),
	}

	if err := practiceRepo.CreatePlanLog(log); err != nil {
		return fmt.Errorf("记录计划执行失败: %w", err)
	}

	// 更新计划进度
	plan.TotalDone += req.QuestionNum
	if plan.TotalTarget > 0 {
		plan.Progress = float64(plan.TotalDone) / float64(plan.TotalTarget) * 100
		if plan.Progress > 100 {
			plan.Progress = 100
		}
	}
	if plan.TotalTarget > 0 && plan.TotalDone >= plan.TotalTarget {
		plan.Status = 2 // 已完成
	}

	return practiceRepo.UpdatePlan(plan)
}

func toPlanInfo(plan *model.PracticePlan) *dto.PlanInfo {
	info := &dto.PlanInfo{
		ID:            plan.ID,
		UserID:        plan.UserID,
		Name:          plan.Name,
		Description:   plan.Description,
		PlanType:      plan.PlanType,
		CategoryID:    plan.CategoryID,
		QuestionType:  plan.QuestionType,
		Difficulty:    plan.Difficulty,
		DailyCount:    plan.DailyCount,
		DailyDuration: plan.DailyDuration,
		StartDate:     plan.StartDate.Format("2006-01-02"),
		TotalTarget:   plan.TotalTarget,
		TotalDone:     plan.TotalDone,
		Progress:      plan.Progress,
		Status:        plan.Status,
		CreatedAt:     plan.CreatedAt,
		UpdatedAt:     plan.UpdatedAt,
	}
	if plan.EndDate != nil {
		endDate := plan.EndDate.Format("2006-01-02")
		info.EndDate = &endDate
	}
	return info
}

// ==================== 每日练习 ====================

// GetTodayPractice 获取今日练习
func GetTodayPractice(userID uint) (*dto.DailyPracticeInfo, error) {
	today := time.Now().Format("2006-01-02")
	dp, err := practiceRepo.FindDailyPractice(userID, today)
	if err != nil {
		// 今日无记录，返回空
		return &dto.DailyPracticeInfo{
			UserID:      userID,
			Date:        today,
			IsCompleted: false,
			Target:      10, // 默认目标
		}, nil
	}

	return &dto.DailyPracticeInfo{
		ID:            dp.ID,
		UserID:        dp.UserID,
		Date:          dp.Date,
		TotalCount:    dp.TotalCount,
		CorrectCount:  dp.CorrectCount,
		Accuracy:      dp.Accuracy,
		Duration:      dp.Duration,
		SessionCount:  dp.SessionCount,
		WrongCount:    dp.WrongCount,
		NewQuestion:   dp.NewQuestion,
		CategoryStats: dp.CategoryStats,
		IsCompleted:   dp.TotalCount >= 10, // 默认目标10题
		Target:        10,
	}, nil
}

// CompleteDailyPractice 完成今日练习
func CompleteDailyPractice(userID uint, req *dto.CompleteDailyRequest) (*dto.DailyPracticeInfo, error) {
	today := time.Now().Format("2006-01-02")

	dp, err := practiceRepo.FindDailyPractice(userID, today)
	if err != nil {
		// 创建新记录
		dp = &model.DailyPractice{
			UserID:       userID,
			Date:         today,
			TotalCount:   req.QuestionCount,
			CorrectCount: req.CorrectCount,
			Duration:     req.Duration,
			SessionCount: 1,
		}
		if req.QuestionCount > 0 {
			dp.Accuracy = float64(req.CorrectCount) / float64(req.QuestionCount) * 100
		}
		if err := practiceRepo.CreateDailyPractice(dp); err != nil {
			return nil, fmt.Errorf("创建每日练习记录失败: %w", err)
		}
	} else {
		// 更新已有记录
		dp.TotalCount += req.QuestionCount
		dp.CorrectCount += req.CorrectCount
		dp.Duration += req.Duration
		dp.SessionCount++
		if dp.TotalCount > 0 {
			dp.Accuracy = float64(dp.CorrectCount) / float64(dp.TotalCount) * 100
		}
		if err := practiceRepo.UpdateDailyPractice(dp); err != nil {
			return nil, fmt.Errorf("更新每日练习记录失败: %w", err)
		}
	}

	return &dto.DailyPracticeInfo{
		ID:            dp.ID,
		UserID:        dp.UserID,
		Date:          dp.Date,
		TotalCount:    dp.TotalCount,
		CorrectCount:  dp.CorrectCount,
		Accuracy:      dp.Accuracy,
		Duration:      dp.Duration,
		SessionCount:  dp.SessionCount,
		WrongCount:    dp.WrongCount,
		NewQuestion:   dp.NewQuestion,
		CategoryStats: dp.CategoryStats,
		IsCompleted:   dp.TotalCount >= 10,
		Target:        10,
	}, nil
}

// ==================== 练习打卡 ====================

// Checkin 打卡
func Checkin(userID uint, req *dto.CheckinRequest) (*dto.CheckinInfo, error) {
	today := time.Now().Format("2006-01-02")

	checkin, err := practiceRepo.FindCheckin(userID, today)
	if err == nil {
		// 今日已打卡
		return &dto.CheckinInfo{
			ID:            checkin.ID,
			UserID:        checkin.UserID,
			Date:          checkin.Date,
			IsCheckin:     checkin.IsCheckin,
			QuestionCount: checkin.QuestionCount,
			Duration:      checkin.Duration,
			Streak:        checkin.Streak,
			Reward:        checkin.Reward,
		}, nil
	}

	// 计算连续打卡天数
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	streak := 1
	prevCheckin, err := practiceRepo.FindCheckin(userID, yesterday)
	if err == nil && prevCheckin.IsCheckin {
		streak = prevCheckin.Streak + 1
	}

	// 计算奖励积分
	reward := 10
	if streak >= 7 {
		reward = 50
	} else if streak >= 3 {
		reward = 20
	}

	checkin = &model.PracticeCheckin{
		UserID:        userID,
		Date:          today,
		IsCheckin:     true,
		QuestionCount: req.QuestionCount,
		Duration:      req.Duration,
		Streak:        streak,
		Reward:        reward,
	}

	if err := practiceRepo.CreateCheckin(checkin); err != nil {
		return nil, fmt.Errorf("打卡失败: %w", err)
	}

	logger.Infof("用户 %d 打卡成功，连续 %d 天，奖励 %d 积分", userID, streak, reward)

	return &dto.CheckinInfo{
		ID:            checkin.ID,
		UserID:        checkin.UserID,
		Date:          checkin.Date,
		IsCheckin:     checkin.IsCheckin,
		QuestionCount: checkin.QuestionCount,
		Duration:      checkin.Duration,
		Streak:        checkin.Streak,
		Reward:        checkin.Reward,
	}, nil
}

// GetCheckinCalendar 打卡日历
func GetCheckinCalendar(userID uint, req *dto.CheckinCalendarRequest) ([]dto.CheckinCalendarItem, error) {
	checkins, err := practiceRepo.ListCheckins(userID, req.Year, req.Month)
	if err != nil {
		return nil, fmt.Errorf("查询打卡日历失败: %w", err)
	}

	items := make([]dto.CheckinCalendarItem, 0, len(checkins))
	for _, c := range checkins {
		items = append(items, dto.CheckinCalendarItem{
			Date:          c.Date,
			IsCheckin:     c.IsCheckin,
			QuestionCount: c.QuestionCount,
			Streak:        c.Streak,
		})
	}

	return items, nil
}

// ==================== 排行榜 ====================

// GetLeaderboard 获取排行榜
func GetLeaderboard(req *dto.LeaderboardRequest) ([]dto.LeaderboardItem, error) {
	limit := req.Limit
	if limit <= 0 {
		limit = 50
	}

	items, err := practiceRepo.GetLeaderboard(req.RankType, limit)
	if err != nil {
		return nil, fmt.Errorf("获取排行榜失败: %w", err)
	}

	result := make([]dto.LeaderboardItem, 0, len(items))
	for _, item := range items {
		result = append(result, dto.LeaderboardItem{
			UserID:     item.UserID,
			Nickname:   item.Nickname,
			Avatar:     item.Avatar,
			RankPos:    item.RankPos,
			Score:      item.Score,
			Accuracy:   item.Accuracy,
			TotalCount: item.TotalCount,
			Duration:   item.Duration,
		})
	}

	return result, nil
}

// ==================== 闯关模式 ====================

// GetChallengeLevels 获取关卡列表
func GetChallengeLevels(userID uint) ([]dto.ChallengeLevelInfo, error) {
	levels, err := practiceRepo.ListChallengeLevels()
	if err != nil {
		return nil, fmt.Errorf("获取关卡列表失败: %w", err)
	}

	// 获取用户进度
	progressMap := make(map[uint]*model.UserChallengeProgress)
	progressList, _ := practiceRepo.ListChallengeProgress(userID)
	for i := range progressList {
		progressMap[progressList[i].ChallengeLevel] = &progressList[i]
	}

	items := make([]dto.ChallengeLevelInfo, 0, len(levels))
	for _, level := range levels {
		info := dto.ChallengeLevelInfo{
			ID:            level.ID,
			Name:          level.Name,
			Description:   level.Description,
			Level:         level.Level,
			CategoryID:    level.CategoryID,
			QuestionCount: level.QuestionCount,
			PassAccuracy:  level.PassAccuracy,
			PassScore:     level.PassScore,
			TimeLimit:     level.TimeLimit,
			Difficulty:    level.Difficulty,
			Icon:          level.Icon,
			Badge:         level.Badge,
			PreLevel:      level.PreLevel,
			Status:        level.Status,
			Sort:          level.Sort,
		}

		// 检查是否锁定
		isLocked := false
		if level.PreLevel > 0 {
			preProgress, ok := progressMap[0] // 简化：通过level查找
			if !ok || preProgress.Status != 2 {
				// 检查前置关卡是否通关
				preLevel, _ := practiceRepo.FindChallengeLevelByLevel(level.PreLevel)
				if preLevel != nil {
					preProg, ok := progressMap[preLevel.ID]
					if !ok || preProg.Status != 2 {
						isLocked = true
					}
				}
			}
		}
		info.IsLocked = isLocked

		// 填充用户进度
		if prog, ok := progressMap[level.ID]; ok {
			info.UserStatus = prog.Status
			info.BestAccuracy = prog.BestAccuracy
			info.Attempts = prog.Attempts
		}

		items = append(items, info)
	}

	return items, nil
}

// GetChallengeLevel 获取关卡详情
func GetChallengeLevel(levelID, userID uint) (*dto.ChallengeLevelInfo, error) {
	level, err := practiceRepo.FindChallengeLevelByID(levelID)
	if err != nil {
		return nil, errors.New("关卡不存在")
	}

	info := &dto.ChallengeLevelInfo{
		ID:            level.ID,
		Name:          level.Name,
		Description:   level.Description,
		Level:         level.Level,
		CategoryID:    level.CategoryID,
		QuestionCount: level.QuestionCount,
		PassAccuracy:  level.PassAccuracy,
		PassScore:     level.PassScore,
		TimeLimit:     level.TimeLimit,
		Difficulty:    level.Difficulty,
		Icon:          level.Icon,
		Badge:         level.Badge,
		PreLevel:      level.PreLevel,
		Status:        level.Status,
		Sort:          level.Sort,
	}

	// 获取用户进度
	progress, err := practiceRepo.FindChallengeProgress(userID, levelID)
	if err == nil {
		info.UserStatus = progress.Status
		info.BestAccuracy = progress.BestAccuracy
		info.Attempts = progress.Attempts
	}

	return info, nil
}

// StartChallenge 开始闯关
func StartChallenge(userID, levelID uint) (*dto.ChallengeSessionInfo, error) {
	level, err := practiceRepo.FindChallengeLevelByID(levelID)
	if err != nil {
		return nil, errors.New("关卡不存在")
	}

	if level.Status != 1 {
		return nil, errors.New("该关卡已禁用")
	}

	// 检查前置关卡
	if level.PreLevel > 0 {
		preLevel, _ := practiceRepo.FindChallengeLevelByLevel(level.PreLevel)
		if preLevel != nil {
			preProgress, err := practiceRepo.FindChallengeProgress(userID, preLevel.ID)
			if err != nil || preProgress.Status != 2 {
				return nil, errors.New("请先通过前置关卡")
			}
		}
	}

	// 创建练习会话
	session := &model.PracticeSession{
		UserID:     userID,
		CategoryID: &levelID,
		TotalCount: level.QuestionCount,
		Duration:   level.TimeLimit,
		Status:     0,
	}

	if err := practiceRepo.CreateSession(session); err != nil {
		return nil, fmt.Errorf("创建闯关会话失败: %w", err)
	}

	// 更新或创建进度
	progress, err := practiceRepo.FindChallengeProgress(userID, levelID)
	if err != nil {
		progress = &model.UserChallengeProgress{
			UserID:         userID,
			ChallengeLevel: levelID,
			Status:         1, // 进行中
			Attempts:       1,
		}
		practiceRepo.CreateChallengeProgress(progress)
	} else {
		progress.Status = 1
		progress.Attempts++
		practiceRepo.UpdateChallengeProgress(progress)
	}

	return &dto.ChallengeSessionInfo{
		ID:           session.ID,
		UserID:       session.UserID,
		LevelID:      levelID,
		LevelName:    level.Name,
		TotalCount:   level.QuestionCount,
		PassAccuracy: level.PassAccuracy,
		TimeLimit:    level.TimeLimit,
		Status:       session.Status,
		StartedAt:    session.CreatedAt,
	}, nil
}

// SubmitChallenge 提交闯关
func SubmitChallenge(sessionID, userID uint, req *dto.SubmitChallengeRequest) (*dto.ChallengeResultInfo, error) {
	session, err := practiceRepo.FindSessionByID(sessionID)
	if err != nil {
		return nil, errors.New("闯关会话不存在")
	}

	if session.UserID != userID {
		return nil, errors.New("无权操作此闯关会话")
	}

	if session.Status != 0 {
		return nil, errors.New("闯关已结束")
	}

	levelID := uint(0)
	if session.CategoryID != nil {
		levelID = *session.CategoryID
	}

	level, err := practiceRepo.FindChallengeLevelByID(levelID)
	if err != nil {
		return nil, errors.New("关卡配置不存在")
	}

	// 批量保存答题记录
	records := make([]model.PracticeRecord, 0, len(req.Answers))
	correctCount := 0
	for _, ans := range req.Answers {
		question, err := questionRepo.FindByID(ans.QuestionID)
		if err != nil {
			continue
		}
		isCorrect := question.Answer == ans.Answer
		if isCorrect {
			correctCount++
		}
		records = append(records, model.PracticeRecord{
			SessionID:  sessionID,
			QuestionID: ans.QuestionID,
			Answer:     ans.Answer,
			IsCorrect:  isCorrect,
			Duration:   ans.Duration,
		})
	}

	if len(records) > 0 {
		if err := practiceRepo.CreateRecords(records); err != nil {
			return nil, fmt.Errorf("保存答题记录失败: %w", err)
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
	session.Duration = req.Duration
	session.Status = 1
	practiceRepo.UpdateSession(session)

	// 判断是否通关
	isPassed := accuracy >= level.PassAccuracy

	// 更新用户进度
	progress, _ := practiceRepo.FindChallengeProgress(userID, levelID)
	if progress != nil {
		if accuracy > progress.BestAccuracy {
			progress.BestAccuracy = accuracy
		}
		if req.Duration > 0 && (progress.BestDuration == 0 || req.Duration < progress.BestDuration) {
			progress.BestDuration = req.Duration
		}
		if isPassed {
			progress.Status = 2 // 已通关
			now := time.Now()
			progress.PassedAt = &now
			progress.Score = level.PassScore
		} else {
			progress.Status = 3 // 失败
		}
		practiceRepo.UpdateChallengeProgress(progress)
	}

	score := 0
	if isPassed {
		score = level.PassScore
	}

	logger.Infof("闯关 %d 提交，正确率: %.2f%%，通关: %v", levelID, accuracy, isPassed)

	return &dto.ChallengeResultInfo{
		LevelID:      levelID,
		LevelName:    level.Name,
		IsPassed:     isPassed,
		Accuracy:     accuracy,
		Duration:     req.Duration,
		Score:        score,
		BestAccuracy: accuracy,
		Attempts:     1,
	}, nil
}

// GetChallengeProgress 获取闯关进度
func GetChallengeProgress(userID uint) (*dto.ChallengeProgressInfo, error) {
	levels, err := practiceRepo.ListChallengeLevels()
	if err != nil {
		return nil, fmt.Errorf("获取关卡列表失败: %w", err)
	}

	progressList, _ := practiceRepo.ListChallengeProgress(userID)
	progressMap := make(map[uint]*model.UserChallengeProgress)
	for i := range progressList {
		progressMap[progressList[i].ChallengeLevel] = &progressList[i]
	}

	passedLevels := 0
	totalScore := 0
	currentLevel := 1
	levelProgress := make([]dto.ChallengeLevelProgress, 0, len(levels))

	for _, level := range levels {
		lp := dto.ChallengeLevelProgress{
			LevelID:   level.ID,
			LevelName: level.Name,
		}

		if prog, ok := progressMap[level.ID]; ok {
			lp.Status = prog.Status
			lp.BestAccuracy = prog.BestAccuracy
			lp.Attempts = prog.Attempts
			lp.Score = prog.Score
			if prog.Status == 2 {
				passedLevels++
				totalScore += prog.Score
				currentLevel = level.Level + 1
			}
		}

		levelProgress = append(levelProgress, lp)
	}

	return &dto.ChallengeProgressInfo{
		TotalLevels:  len(levels),
		PassedLevels: passedLevels,
		TotalScore:   totalScore,
		CurrentLevel: currentLevel,
		ProgressList: levelProgress,
	}, nil
}
