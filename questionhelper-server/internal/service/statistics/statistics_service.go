package statistics

import (
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"

	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	"questionhelper-server/pkg/database"
	"questionhelper-server/pkg/logger"
)

// ==================== 基础统计 ====================

// GetUserStatistics 用户统计
func GetUserStatistics() (*dto.UserStatistics, error) {
	var stats dto.UserStatistics

	database.DB.Model(&model.User{}).Count(&stats.TotalUsers)
	database.DB.Model(&model.User{}).Where("status = ?", 1).Count(&stats.ActiveUsers)

	today := time.Now().Format("2006-01-02")
	database.DB.Model(&model.User{}).Where("DATE(created_at) = ?", today).Count(&stats.NewUsersToday)

	weekAgo := time.Now().AddDate(0, 0, -7)
	database.DB.Model(&model.User{}).Where("created_at >= ?", weekAgo).Count(&stats.NewUsersWeek)

	return &stats, nil
}

// GetPracticeStatistics 练习统计
func GetPracticeStatistics() (*dto.PracticeStatistics, error) {
	var stats dto.PracticeStatistics

	database.DB.Model(&model.PracticeSession{}).Count(&stats.TotalSessions)

	var totalQuestions int64
	database.DB.Model(&model.PracticeRecord{}).Count(&totalQuestions)
	stats.TotalQuestions = totalQuestions

	var avgAccuracy float64
	database.DB.Model(&model.PracticeSession{}).
		Where("status = ?", 1).
		Select("COALESCE(AVG(accuracy), 0)").
		Scan(&avgAccuracy)
	stats.AvgAccuracy = avgAccuracy

	var totalDuration int64
	database.DB.Model(&model.PracticeSession{}).
		Select("COALESCE(SUM(duration), 0)").
		Scan(&totalDuration)
	stats.TotalDuration = totalDuration

	return &stats, nil
}

// GetExamStatistics 考试统计
func GetExamStatistics() (*dto.ExamStatistics, error) {
	var stats dto.ExamStatistics

	database.DB.Model(&model.Exam{}).Count(&stats.TotalExams)
	database.DB.Model(&model.ExamRecord{}).Count(&stats.TotalRecords)

	var avgScore float64
	database.DB.Model(&model.ExamRecord{}).
		Where("status = ?", 2).
		Select("COALESCE(AVG(score), 0)").
		Scan(&avgScore)
	stats.AvgScore = avgScore

	var totalRecords int64
	var passRecords int64
	database.DB.Model(&model.ExamRecord{}).Where("status = ?", 2).Count(&totalRecords)
	if totalRecords > 0 {
		database.DB.Model(&model.ExamRecord{}).
			Joins("JOIN exams ON exams.id = exam_records.exam_id").
			Where("exam_records.status = ? AND exam_records.score >= exams.pass_score", 2).
			Count(&passRecords)
		stats.PassRate = float64(passRecords) / float64(totalRecords) * 100
	}

	return &stats, nil
}

// GetClassStatistics 班级统计
func GetClassStatistics() (*dto.ClassStatistics, error) {
	var stats dto.ClassStatistics

	database.DB.Model(&model.Class{}).Count(&stats.TotalClasses)
	database.DB.Model(&model.ClassMember{}).Count(&stats.TotalMembers)

	var avgMemberCount float64
	database.DB.Model(&model.Class{}).
		Select("COALESCE(AVG(member_count), 0)").
		Scan(&avgMemberCount)
	stats.AvgMemberCount = avgMemberCount

	return &stats, nil
}

// GetDashboard 仪表盘数据
func GetDashboard() (*dto.DashboardInfo, error) {
	userStats, err := GetUserStatistics()
	if err != nil {
		return nil, fmt.Errorf("获取用户统计失败: %w", err)
	}

	practiceStats, err := GetPracticeStatistics()
	if err != nil {
		return nil, fmt.Errorf("获取练习统计失败: %w", err)
	}

	examStats, err := GetExamStatistics()
	if err != nil {
		return nil, fmt.Errorf("获取考试统计失败: %w", err)
	}

	classStats, err := GetClassStatistics()
	if err != nil {
		return nil, fmt.Errorf("获取班级统计失败: %w", err)
	}

	return &dto.DashboardInfo{
		UserStats:     *userStats,
		PracticeStats: *practiceStats,
		ExamStats:     *examStats,
		ClassStats:    *classStats,
	}, nil
}

// GetRanking 排行榜
func GetRanking(req *dto.RankingRequest) ([]dto.RankInfo, int64, error) {
	var ranks []dto.RankInfo
	var total int64

	switch req.Type {
	case 1: // 练习排行榜
		var results []struct {
			UserID       uint    `json:"user_id"`
			UserName     string  `json:"user_name"`
			Avatar       string  `json:"avatar"`
			TotalCount   int     `json:"total_count"`
			CorrectCount int     `json:"correct_count"`
		}

		database.DB.Model(&model.PracticeSession{}).
			Select("practice_sessions.user_id, users.nickname as user_name, users.avatar, SUM(practice_sessions.total_count) as total_count, SUM(practice_sessions.correct_count) as correct_count").
			Joins("JOIN users ON users.id = practice_sessions.user_id").
			Where("practice_sessions.status = ?", 1).
			Group("practice_sessions.user_id, users.nickname, users.avatar").
			Order("correct_count DESC").
			Offset(req.GetOffset()).Limit(req.GetLimit()).
			Scan(&results)

		database.DB.Model(&model.PracticeSession{}).
			Where("status = ?", 1).
			Distinct("user_id").Count(&total)

		for i, r := range results {
			accuracy := float64(0)
			if r.TotalCount > 0 {
				accuracy = float64(r.CorrectCount) / float64(r.TotalCount) * 100
			}
			ranks = append(ranks, dto.RankInfo{
				Rank:     i + 1 + req.GetOffset(),
				UserID:   r.UserID,
				UserName: r.UserName,
				Avatar:   r.Avatar,
				Score:    accuracy,
				Count:    r.TotalCount,
			})
		}

	case 2: // 考试排行榜
		var results []struct {
			UserID   uint    `json:"user_id"`
			UserName string  `json:"user_name"`
			Avatar   string  `json:"avatar"`
			AvgScore float64 `json:"avg_score"`
			Count    int     `json:"count"`
		}

		database.DB.Model(&model.ExamRecord{}).
			Select("exam_records.user_id, users.nickname as user_name, users.avatar, AVG(exam_records.score) as avg_score, COUNT(*) as count").
			Joins("JOIN users ON users.id = exam_records.user_id").
			Where("exam_records.status = ?", 2).
			Group("exam_records.user_id, users.nickname, users.avatar").
			Order("avg_score DESC").
			Offset(req.GetOffset()).Limit(req.GetLimit()).
			Scan(&results)

		database.DB.Model(&model.ExamRecord{}).
			Where("status = ?", 2).
			Distinct("user_id").Count(&total)

		for i, r := range results {
			ranks = append(ranks, dto.RankInfo{
				Rank:     i + 1 + req.GetOffset(),
				UserID:   r.UserID,
				UserName: r.UserName,
				Avatar:   r.Avatar,
				Score:    r.AvgScore,
				Count:    r.Count,
			})
		}

	default:
		return nil, 0, fmt.Errorf("不支持的排行榜类型: %d", req.Type)
	}

	return ranks, total, nil
}

// ==================== 用户留存分析 ====================

// GetRetention 获取用户留存统计
func GetRetention(req *dto.RetentionRequest) ([]dto.RetentionItem, error) {
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("无效的开始日期: %w", err)
	}
	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("无效的结束日期: %w", err)
	}

	var items []dto.RetentionItem

	// 从数据库查询已有的留存记录
	var records []model.UserRetention
	query := database.DB.Where("period = ?", req.Period).
		Where("date >= ? AND date <= ?", startDate, endDate).
		Order("date ASC")

	if err := query.Find(&records).Error; err != nil {
		return nil, fmt.Errorf("查询留存数据失败: %w", err)
	}

	for _, r := range records {
		items = append(items, dto.RetentionItem{
			Date:          r.Date.Format("2006-01-02"),
			NewUsers:      r.NewUsers,
			RetainedUsers: r.RetainedUsers,
			RetentionRate: r.RetentionRate,
			Period:        r.Period,
		})
	}

	// 如果没有预计算的数据，尝试实时计算
	if len(items) == 0 {
		items, err = calculateRetention(req.Period, startDate, endDate)
		if err != nil {
			return nil, err
		}
	}

	return items, nil
}

// calculateRetention 实时计算留存率
func calculateRetention(period string, startDate, endDate time.Time) ([]dto.RetentionItem, error) {
	var items []dto.RetentionItem

	switch period {
	case "day":
		// 按日计算
		for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
			dateStr := d.Format("2006-01-02")
			nextDay := d.AddDate(0, 0, 1)

			// 当日新增用户
			var newUsers int64
			database.DB.Model(&model.User{}).
				Where("DATE(created_at) = ?", dateStr).
				Count(&newUsers)

			// 次日仍活跃的用户
			var retainedUsers int64
			if newUsers > 0 {
				database.DB.Model(&model.User{}).
					Where("DATE(created_at) = ? AND last_login_at >= ?", dateStr, nextDay).
					Count(&retainedUsers)
			}

			rate := float64(0)
			if newUsers > 0 {
				rate = float64(retainedUsers) / float64(newUsers) * 100
			}

			items = append(items, dto.RetentionItem{
				Date:          dateStr,
				NewUsers:      int(newUsers),
				RetainedUsers: int(retainedUsers),
				RetentionRate: rate,
				Period:        "day",
			})
		}

	case "week":
		// 按周计算
		for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 7) {
			weekEnd := d.AddDate(0, 0, 7)
			dateStr := d.Format("2006-01-02")

			var newUsers int64
			database.DB.Model(&model.User{}).
				Where("created_at >= ? AND created_at < ?", d, weekEnd).
				Count(&newUsers)

			var retainedUsers int64
			if newUsers > 0 {
				database.DB.Model(&model.User{}).
					Where("created_at >= ? AND created_at < ? AND last_login_at >= ?", d, weekEnd, weekEnd).
					Count(&retainedUsers)
			}

			rate := float64(0)
			if newUsers > 0 {
				rate = float64(retainedUsers) / float64(newUsers) * 100
			}

			items = append(items, dto.RetentionItem{
				Date:          dateStr,
				NewUsers:      int(newUsers),
				RetainedUsers: int(retainedUsers),
				RetentionRate: rate,
				Period:        "week",
			})
		}

	case "month":
		// 按月计算
		for d := startDate; d.Before(endDate) || d.Equal(endDate); {
			monthEnd := d.AddDate(0, 1, 0)
			dateStr := d.Format("2006-01-02")

			var newUsers int64
			database.DB.Model(&model.User{}).
				Where("created_at >= ? AND created_at < ?", d, monthEnd).
				Count(&newUsers)

			var retainedUsers int64
			if newUsers > 0 {
				database.DB.Model(&model.User{}).
					Where("created_at >= ? AND created_at < ? AND last_login_at >= ?", d, monthEnd, monthEnd).
					Count(&retainedUsers)
			}

			rate := float64(0)
			if newUsers > 0 {
				rate = float64(retainedUsers) / float64(newUsers) * 100
			}

			items = append(items, dto.RetentionItem{
				Date:          dateStr,
				NewUsers:      int(newUsers),
				RetainedUsers: int(retainedUsers),
				RetentionRate: rate,
				Period:        "month",
			})

			d = monthEnd
		}
	}

	return items, nil
}

// ==================== 用户流失分析 ====================

// GetChurn 获取用户流失统计
func GetChurn(req *dto.ChurnRequest) ([]dto.ChurnItem, error) {
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("无效的开始日期: %w", err)
	}
	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("无效的结束日期: %w", err)
	}

	var items []dto.ChurnItem

	// 查询已有的流失记录
	var records []model.UserChurn
	query := database.DB.Where("period = ?", req.Period).
		Where("date >= ? AND date <= ?", startDate, endDate).
		Order("date ASC")

	if err := query.Find(&records).Error; err != nil {
		return nil, fmt.Errorf("查询流失数据失败: %w", err)
	}

	for _, r := range records {
		items = append(items, dto.ChurnItem{
			Date:         r.Date.Format("2006-01-02"),
			ChurnedUsers: r.ChurnedUsers,
			ChurnRate:    r.ChurnRate,
			ChurnReasons: r.ChurnReasons,
			Period:       r.Period,
		})
	}

	// 如果没有预计算数据，实时计算
	if len(items) == 0 {
		items, err = calculateChurn(req.Period, startDate, endDate)
		if err != nil {
			return nil, err
		}
	}

	return items, nil
}

// calculateChurn 实时计算流失
func calculateChurn(period string, startDate, endDate time.Time) ([]dto.ChurnItem, error) {
	var items []dto.ChurnItem

	// 流失定义：超过30天未登录的用户
	churnDays := 30

	switch period {
	case "day":
		for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
			dateStr := d.Format("2006-01-02")
			churnThreshold := d.AddDate(0, 0, -churnDays)

			var totalUsers int64
			database.DB.Model(&model.User{}).Where("created_at < ?", d).Count(&totalUsers)

			var churnedUsers int64
			database.DB.Model(&model.User{}).
				Where("created_at < ? AND (last_login_at IS NULL OR last_login_at < ?) AND status = 1", d, churnThreshold).
				Count(&churnedUsers)

			rate := float64(0)
			if totalUsers > 0 {
				rate = float64(churnedUsers) / float64(totalUsers) * 100
			}

			items = append(items, dto.ChurnItem{
				Date:         dateStr,
				ChurnedUsers: int(churnedUsers),
				ChurnRate:    rate,
				Period:       "day",
			})
		}

	case "week":
		for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 7) {
			weekEnd := d.AddDate(0, 0, 7)
			dateStr := d.Format("2006-01-02")
			churnThreshold := d.AddDate(0, 0, -churnDays)

			var totalUsers int64
			database.DB.Model(&model.User{}).Where("created_at < ?", weekEnd).Count(&totalUsers)

			var churnedUsers int64
			database.DB.Model(&model.User{}).
				Where("created_at < ? AND (last_login_at IS NULL OR last_login_at < ?) AND status = 1", weekEnd, churnThreshold).
				Count(&churnedUsers)

			rate := float64(0)
			if totalUsers > 0 {
				rate = float64(churnedUsers) / float64(totalUsers) * 100
			}

			items = append(items, dto.ChurnItem{
				Date:         dateStr,
				ChurnedUsers: int(churnedUsers),
				ChurnRate:    rate,
				Period:       "week",
			})
		}

	case "month":
		for d := startDate; d.Before(endDate) || d.Equal(endDate); {
			monthEnd := d.AddDate(0, 1, 0)
			dateStr := d.Format("2006-01-02")
			churnThreshold := d.AddDate(0, 0, -churnDays)

			var totalUsers int64
			database.DB.Model(&model.User{}).Where("created_at < ?", monthEnd).Count(&totalUsers)

			var churnedUsers int64
			database.DB.Model(&model.User{}).
				Where("created_at < ? AND (last_login_at IS NULL OR last_login_at < ?) AND status = 1", monthEnd, churnThreshold).
				Count(&churnedUsers)

			rate := float64(0)
			if totalUsers > 0 {
				rate = float64(churnedUsers) / float64(totalUsers) * 100
			}

			items = append(items, dto.ChurnItem{
				Date:         dateStr,
				ChurnedUsers: int(churnedUsers),
				ChurnRate:    rate,
				Period:       "month",
			})

			d = monthEnd
		}
	}

	return items, nil
}

// ==================== 用户行为事件 ====================

// CreateEvent 创建用户行为事件
func CreateEvent(userID uint, ip string, req *dto.CreateEventRequest) error {
	event := &model.UserEvent{
		UserID:     userID,
		EventType:  req.EventType,
		EventName:  req.EventName,
		Page:       req.Page,
		Element:    req.Element,
		ExtraData:  req.ExtraData,
		SessionID:  req.SessionID,
		DeviceType: req.DeviceType,
		IP:         ip,
	}

	if err := database.DB.Create(event).Error; err != nil {
		return fmt.Errorf("创建行为事件失败: %w", err)
	}
	return nil
}

// AnalyzeEvents 行为事件分析
func AnalyzeEvents(req *dto.EventAnalysisRequest) (*dto.EventSummary, error) {
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("无效的开始日期: %w", err)
	}
	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("无效的结束日期: %w", err)
	}
	endDate = endDate.AddDate(0, 0, 1) // 包含结束日期

	db := database.DB.Model(&model.UserEvent{}).
		Where("created_at >= ? AND created_at < ?", startDate, endDate)

	if req.EventType != "" {
		db = db.Where("event_type = ?", req.EventType)
	}

	// 总数统计
	var totalEvents int64
	var totalUsers int64
	db.Count(&totalEvents)
	db.Distinct("user_id").Count(&totalUsers)

	// 分组统计
	var items []dto.EventAnalysisItem

	groupField := "DATE(created_at)"
	dimensionField := "DATE(created_at) as dimension"
	switch req.GroupBy {
	case "event_type":
		groupField = "event_type"
		dimensionField = "event_type as dimension"
	case "page":
		groupField = "page"
		dimensionField = "page as dimension"
	case "device_type":
		groupField = "device_type"
		dimensionField = "device_type as dimension"
	default:
		// date
	}

	var results []struct {
		Dimension string `json:"dimension"`
		Count     int64  `json:"count"`
		Users     int64  `json:"users"`
	}

	db.Select(dimensionField + ", COUNT(*) as count, COUNT(DISTINCT user_id) as users").
		Group(groupField).
		Order("count DESC").
		Scan(&results)

	for _, r := range results {
		items = append(items, dto.EventAnalysisItem{
			Dimension: r.Dimension,
			Count:     r.Count,
			Users:     r.Users,
		})
	}

	return &dto.EventSummary{
		TotalEvents: totalEvents,
		TotalUsers:  totalUsers,
		Items:       items,
	}, nil
}

// ==================== 用户分群 ====================

// ListSegments 分群列表
func ListSegments(req *dto.SegmentListRequest) ([]dto.SegmentInfo, int64, error) {
	var segments []model.UserSegment
	var total int64

	db := database.DB.Model(&model.UserSegment{})

	if req.IsActive != nil {
		db = db.Where("is_active = ?", *req.IsActive)
	}
	if req.Keyword != "" {
		db = db.Where("name LIKE ?", "%"+req.Keyword+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("统计分群数量失败: %w", err)
	}

	if err := db.Order("created_at DESC").
		Offset(req.GetOffset()).Limit(req.GetLimit()).
		Find(&segments).Error; err != nil {
		return nil, 0, fmt.Errorf("查询分群列表失败: %w", err)
	}

	var list []dto.SegmentInfo
	for _, s := range segments {
		list = append(list, dto.SegmentInfo{
			ID:          s.ID,
			Name:        s.Name,
			Description: s.Description,
			Rules:       s.Rules,
			UserCount:   s.UserCount,
			IsActive:    s.IsActive,
			CreatorID:   s.CreatorID,
			CreatedAt:   s.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   s.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return list, total, nil
}

// CreateSegment 创建分群
func CreateSegment(creatorID uint, req *dto.CreateSegmentRequest) error {
	segment := &model.UserSegment{
		Name:        req.Name,
		Description: req.Description,
		Rules:       req.Rules,
		IsActive:    true,
		CreatorID:   creatorID,
	}

	if err := database.DB.Create(segment).Error; err != nil {
		return fmt.Errorf("创建分群失败: %w", err)
	}
	return nil
}

// GetSegment 获取分群详情
func GetSegment(id uint) (*dto.SegmentDetail, error) {
	var segment model.UserSegment
	if err := database.DB.First(&segment, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("分群不存在")
		}
		return nil, fmt.Errorf("查询分群失败: %w", err)
	}

	// 查询分群成员
	var members []model.UserSegmentMember
	database.DB.Where("segment_id = ?", id).Find(&members)

	var memberInfos []dto.SegmentMemberInfo
	for _, m := range members {
		var user model.User
		database.DB.Select("id, nickname, avatar").First(&user, m.UserID)
		memberInfos = append(memberInfos, dto.SegmentMemberInfo{
			ID:       m.ID,
			UserID:   m.UserID,
			UserName: user.Nickname,
			Avatar:   user.Avatar,
			JoinedAt: m.JoinedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &dto.SegmentDetail{
		SegmentInfo: dto.SegmentInfo{
			ID:          segment.ID,
			Name:        segment.Name,
			Description: segment.Description,
			Rules:       segment.Rules,
			UserCount:   segment.UserCount,
			IsActive:    segment.IsActive,
			CreatorID:   segment.CreatorID,
			CreatedAt:   segment.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   segment.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
		Members: memberInfos,
	}, nil
}

// UpdateSegment 更新分群
func UpdateSegment(id uint, req *dto.UpdateSegmentRequest) error {
	var segment model.UserSegment
	if err := database.DB.First(&segment, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("分群不存在")
		}
		return fmt.Errorf("查询分群失败: %w", err)
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Rules != "" {
		updates["rules"] = req.Rules
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if len(updates) > 0 {
		if err := database.DB.Model(&segment).Updates(updates).Error; err != nil {
			return fmt.Errorf("更新分群失败: %w", err)
		}
	}
	return nil
}

// DeleteSegment 删除分群
func DeleteSegment(id uint) error {
	var segment model.UserSegment
	if err := database.DB.First(&segment, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("分群不存在")
		}
		return fmt.Errorf("查询分群失败: %w", err)
	}

	// 删除分群成员
	database.DB.Where("segment_id = ?", id).Delete(&model.UserSegmentMember{})

	// 软删除分群
	if err := database.DB.Delete(&segment).Error; err != nil {
		return fmt.Errorf("删除分群失败: %w", err)
	}
	return nil
}

// ==================== 用户路径分析 ====================

// GetPathAnalysis 获取路径分析
func GetPathAnalysis(req *dto.PathAnalysisRequest) (*dto.PathAnalysisResult, error) {
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("无效的开始日期: %w", err)
	}
	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("无效的结束日期: %w", err)
	}
	endDate = endDate.AddDate(0, 0, 1)

	db := database.DB.Model(&model.UserPageView{}).
		Where("created_at >= ? AND created_at < ?", startDate, endDate)

	if req.DeviceType != "" {
		db = db.Where("device_type = ?", req.DeviceType)
	}

	limit := req.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	// 页面访问统计
	var pages []dto.PathItem
	var pageResults []struct {
		Page       string  `json:"page"`
		Count      int64   `json:"count"`
		Users      int64   `json:"users"`
		AvgTime    float64 `json:"avg_time"`
	}

	db.Select("page, COUNT(*) as count, COUNT(DISTINCT user_id) as users, COALESCE(AVG(duration), 0) as avg_time").
		Group("page").
		Order("count DESC").
		Limit(limit).
		Scan(&pageResults)

	for _, p := range pageResults {
		pages = append(pages, dto.PathItem{
			Page:    p.Page,
			Count:   p.Count,
			Users:   p.Users,
			AvgTime: p.AvgTime,
		})
	}

	// 页面流转统计
	var transitions []dto.PathTransition
	var transResults []struct {
		From  string `json:"from_page"`
		To    string `json:"to_page"`
		Count int64  `json:"count"`
		Users int64  `json:"users"`
	}

	database.DB.Raw(`
		SELECT a.page as from_page, b.page as to_page, COUNT(*) as count, COUNT(DISTINCT a.user_id) as users
		FROM user_page_views a
		JOIN user_page_views b ON a.user_id = b.user_id AND a.session_id = b.session_id
		WHERE a.created_at >= ? AND a.created_at < ?
		AND b.created_at > a.created_at
		AND b.created_at = (
			SELECT MIN(c.created_at) FROM user_page_views c
			WHERE c.user_id = a.user_id AND c.session_id = a.session_id AND c.created_at > a.created_at
		)
		GROUP BY a.page, b.page
		ORDER BY count DESC
		LIMIT ?
	`, startDate, endDate, limit).Scan(&transResults)

	for _, t := range transResults {
		transitions = append(transitions, dto.PathTransition{
			From:  t.From,
			To:    t.To,
			Count: t.Count,
			Users: t.Users,
		})
	}

	return &dto.PathAnalysisResult{
		Pages:       pages,
		Transitions: transitions,
	}, nil
}

// ==================== 转化漏斗 ====================

// ListFunnels 漏斗列表
func ListFunnels(page, pageSize int) ([]dto.FunnelInfo, int64, error) {
	var funnels []model.ConversionFunnel
	var total int64

	if err := database.DB.Model(&model.ConversionFunnel{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("统计漏斗数量失败: %w", err)
	}

	offset := (page - 1) * pageSize
	if offset < 0 {
		offset = 0
	}

	if err := database.DB.Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&funnels).Error; err != nil {
		return nil, 0, fmt.Errorf("查询漏斗列表失败: %w", err)
	}

	var list []dto.FunnelInfo
	for _, f := range funnels {
		list = append(list, dto.FunnelInfo{
			ID:          f.ID,
			Name:        f.Name,
			Description: f.Description,
			Steps:       f.Steps,
			IsActive:    f.IsActive,
			CreatorID:   f.CreatorID,
			CreatedAt:   f.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   f.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return list, total, nil
}

// CreateFunnel 创建漏斗
func CreateFunnel(creatorID uint, req *dto.CreateFunnelRequest) error {
	// 验证步骤JSON格式
	var steps []interface{}
	if err := json.Unmarshal([]byte(req.Steps), &steps); err != nil {
		return fmt.Errorf("无效的漏斗步骤JSON: %w", err)
	}

	funnel := &model.ConversionFunnel{
		Name:        req.Name,
		Description: req.Description,
		Steps:       req.Steps,
		IsActive:    true,
		CreatorID:   creatorID,
	}

	if err := database.DB.Create(funnel).Error; err != nil {
		return fmt.Errorf("创建漏斗失败: %w", err)
	}
	return nil
}

// GetFunnelStats 获取漏斗统计
func GetFunnelStats(id uint, req *dto.FunnelStatsRequest) (*dto.FunnelStatsResult, error) {
	var funnel model.ConversionFunnel
	if err := database.DB.First(&funnel, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("漏斗不存在")
		}
		return nil, fmt.Errorf("查询漏斗失败: %w", err)
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("无效的开始日期: %w", err)
	}
	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("无效的结束日期: %w", err)
	}

	// 查询已有的漏斗统计数据
	var stats []model.ConversionFunnelStat
	database.DB.Where("funnel_id = ? AND date >= ? AND date <= ?", id, startDate, endDate).
		Order("step_index ASC").
		Find(&stats)

	// 按步骤聚合
	stepMap := make(map[int]*dto.FunnelStepStat)
	for _, s := range stats {
		if existing, ok := stepMap[s.StepIndex]; ok {
			existing.UserCount += s.UserCount
		} else {
			stepMap[s.StepIndex] = &dto.FunnelStepStat{
				StepIndex: s.StepIndex,
				StepName:  s.StepName,
				UserCount: s.UserCount,
			}
		}
	}

	// 解析步骤JSON获取步骤名
	var steps []struct {
		Name string `json:"name"`
	}
	json.Unmarshal([]byte(funnel.Steps), &steps)

	var stepStats []dto.FunnelStepStat
	var prevCount int
	for i, step := range steps {
		stat, ok := stepMap[i]
		if !ok {
			stat = &dto.FunnelStepStat{
				StepIndex: i,
				StepName:  step.Name,
			}
		}
		stat.StepName = step.Name

		// 计算转化率
		if i == 0 {
			stat.TotalRate = 100
			stat.ConversionRate = 100
			prevCount = stat.UserCount
		} else {
			if prevCount > 0 {
				stat.ConversionRate = float64(stat.UserCount) / float64(prevCount) * 100
			}
			if stepMap[0] != nil && stepMap[0].UserCount > 0 {
				stat.TotalRate = float64(stat.UserCount) / float64(stepMap[0].UserCount) * 100
			}
			prevCount = stat.UserCount
		}

		stepStats = append(stepStats, *stat)
	}

	return &dto.FunnelStatsResult{
		FunnelID:   funnel.ID,
		FunnelName: funnel.Name,
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
		Steps:      stepStats,
	}, nil
}

// ==================== 数据预警 ====================

// ListAlertRules 预警规则列表
func ListAlertRules(req *dto.AlertRuleListRequest) ([]dto.AlertRuleInfo, int64, error) {
	var rules []model.AlertRule
	var total int64

	db := database.DB.Model(&model.AlertRule{})

	if req.IsActive != nil {
		db = db.Where("is_active = ?", *req.IsActive)
	}
	if req.Keyword != "" {
		db = db.Where("name LIKE ?", "%"+req.Keyword+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("统计预警规则数量失败: %w", err)
	}

	if err := db.Order("created_at DESC").
		Offset(req.GetOffset()).Limit(req.GetLimit()).
		Find(&rules).Error; err != nil {
		return nil, 0, fmt.Errorf("查询预警规则列表失败: %w", err)
	}

	var list []dto.AlertRuleInfo
	for _, r := range rules {
		list = append(list, dto.AlertRuleInfo{
			ID:         r.ID,
			Name:       r.Name,
			MetricName: r.MetricName,
			Condition:  r.Condition,
			Threshold:  r.Threshold,
			Duration:   r.Duration,
			NotifyType: r.NotifyType,
			IsActive:   r.IsActive,
			CreatorID:  r.CreatorID,
			CreatedAt:  r.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:  r.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return list, total, nil
}

// CreateAlertRule 创建预警规则
func CreateAlertRule(creatorID uint, req *dto.CreateAlertRuleRequest) error {
	rule := &model.AlertRule{
		Name:       req.Name,
		MetricName: req.MetricName,
		Condition:  req.Condition,
		Threshold:  req.Threshold,
		Duration:   req.Duration,
		NotifyType: req.NotifyType,
		IsActive:   true,
		CreatorID:  creatorID,
	}

	if err := database.DB.Create(rule).Error; err != nil {
		return fmt.Errorf("创建预警规则失败: %w", err)
	}
	return nil
}

// ListAlertRecords 预警记录列表
func ListAlertRecords(req *dto.AlertRecordListRequest) ([]dto.AlertRecordInfo, int64, error) {
	var records []model.AlertRecord
	var total int64

	db := database.DB.Model(&model.AlertRecord{})

	if req.Level != "" {
		db = db.Where("level = ?", req.Level)
	}
	if req.RuleID > 0 {
		db = db.Where("rule_id = ?", req.RuleID)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("统计预警记录数量失败: %w", err)
	}

	if err := db.Order("created_at DESC").
		Offset(req.GetOffset()).Limit(req.GetLimit()).
		Find(&records).Error; err != nil {
		return nil, 0, fmt.Errorf("查询预警记录列表失败: %w", err)
	}

	// 批量获取规则名称
	ruleIDs := make(map[uint]bool)
	for _, r := range records {
		ruleIDs[r.RuleID] = true
	}

	ruleNames := make(map[uint]string)
	for ruleID := range ruleIDs {
		var rule model.AlertRule
		if err := database.DB.Select("id, name").First(&rule, ruleID).Error; err == nil {
			ruleNames[ruleID] = rule.Name
		}
	}

	var list []dto.AlertRecordInfo
	for _, r := range records {
		info := dto.AlertRecordInfo{
			ID:          r.ID,
			RuleID:      r.RuleID,
			RuleName:    ruleNames[r.RuleID],
			MetricValue: r.MetricValue,
			Threshold:   r.Threshold,
			Level:       r.Level,
			Message:     r.Message,
			CreatedAt:   r.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		if r.HandledAt != nil {
			info.HandledAt = r.HandledAt.Format("2006-01-02 15:04:05")
		}
		if r.HandlerID != nil {
			info.HandlerID = r.HandlerID
		}
		list = append(list, info)
	}

	return list, total, nil
}

// ==================== 数据订阅 ====================

// ListSubscriptions 订阅列表
func ListSubscriptions(userID uint, page, pageSize int) ([]dto.SubscriptionInfo, int64, error) {
	var subs []model.DataSubscription
	var total int64

	db := database.DB.Model(&model.DataSubscription{}).Where("user_id = ?", userID)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("统计订阅数量失败: %w", err)
	}

	offset := (page - 1) * pageSize
	if offset < 0 {
		offset = 0
	}

	if err := db.Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&subs).Error; err != nil {
		return nil, 0, fmt.Errorf("查询订阅列表失败: %w", err)
	}

	var list []dto.SubscriptionInfo
	for _, s := range subs {
		info := dto.SubscriptionInfo{
			ID:         s.ID,
			UserID:     s.UserID,
			ReportType: s.ReportType,
			Frequency:  s.Frequency,
			Channels:   s.Channels,
			IsActive:   s.IsActive,
			CreatedAt:  s.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		if s.LastSentAt != nil {
			info.LastSentAt = s.LastSentAt.Format("2006-01-02 15:04:05")
		}
		list = append(list, info)
	}

	return list, total, nil
}

// CreateSubscription 创建订阅
func CreateSubscription(userID uint, req *dto.CreateSubscriptionRequest) error {
	// 验证渠道JSON格式
	var channels []interface{}
	if err := json.Unmarshal([]byte(req.Channels), &channels); err != nil {
		return fmt.Errorf("无效的通知渠道JSON: %w", err)
	}

	sub := &model.DataSubscription{
		UserID:     userID,
		ReportType: req.ReportType,
		Frequency:  req.Frequency,
		Channels:   req.Channels,
		IsActive:   true,
	}

	if err := database.DB.Create(sub).Error; err != nil {
		return fmt.Errorf("创建订阅失败: %w", err)
	}
	return nil
}

// DeleteSubscription 删除订阅
func DeleteSubscription(userID, id uint) error {
	var sub model.DataSubscription
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&sub).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("订阅不存在")
		}
		return fmt.Errorf("查询订阅失败: %w", err)
	}

	if err := database.DB.Delete(&sub).Error; err != nil {
		return fmt.Errorf("删除订阅失败: %w", err)
	}
	return nil
}

// ==================== 数据导出 ====================

// ExportData 导出统计数据
func ExportData(userID uint, req *dto.ExportRequest) (*dto.ExportInfo, error) {
	record := &model.ExportRecord{
		UserID:     userID,
		ExportType: req.ExportType,
		FileFormat: req.FileFormat,
		Status:     "pending",
		Filters:    req.Filters,
	}

	if err := database.DB.Create(record).Error; err != nil {
		return nil, fmt.Errorf("创建导出记录失败: %w", err)
	}

	// TODO: 异步处理导出任务
	logger.Infof("创建导出任务: id=%d, type=%s, format=%s", record.ID, req.ExportType, req.FileFormat)

	return &dto.ExportInfo{
		ID:         record.ID,
		UserID:     record.UserID,
		ExportType: record.ExportType,
		FileFormat: record.FileFormat,
		Status:     record.Status,
		CreatedAt:  record.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// ==================== 数据对比 ====================

// CompareData 数据对比
func CompareData(req *dto.CompareRequest) (*dto.CompareResult, error) {
	p1Start, err := time.Parse("2006-01-02", req.Period1Start)
	if err != nil {
		return nil, fmt.Errorf("无效的第一段开始日期: %w", err)
	}
	p1End, err := time.Parse("2006-01-02", req.Period1End)
	if err != nil {
		return nil, fmt.Errorf("无效的第一段结束日期: %w", err)
	}
	p2Start, err := time.Parse("2006-01-02", req.Period2Start)
	if err != nil {
		return nil, fmt.Errorf("无效的第二段开始日期: %w", err)
	}
	p2End, err := time.Parse("2006-01-02", req.Period2End)
	if err != nil {
		return nil, fmt.Errorf("无效的第二段结束日期: %w", err)
	}

	// 根据快照类型获取数据
	data1, err := getSnapshotData(req.SnapshotType, p1Start, p1End)
	if err != nil {
		return nil, fmt.Errorf("获取第一段数据失败: %w", err)
	}

	data2, err := getSnapshotData(req.SnapshotType, p2Start, p2End)
	if err != nil {
		return nil, fmt.Errorf("获取第二段数据失败: %w", err)
	}

	// 计算差异
	diffValues := make(map[string]float64)
	diffRates := make(map[string]float64)

	for key, v1 := range data1 {
		if v2, ok := data2[key]; ok {
			val1 := toFloat64(v1)
			val2 := toFloat64(v2)
			diffValues[key] = val1 - val2
			if val2 != 0 {
				diffRates[key] = (val1 - val2) / val2 * 100
			}
		}
	}

	return &dto.CompareResult{
		Period1: dto.ComparePeriod{
			StartDate: req.Period1Start,
			EndDate:   req.Period1End,
			Data:      data1,
		},
		Period2: dto.ComparePeriod{
			StartDate: req.Period2Start,
			EndDate:   req.Period2End,
			Data:      data2,
		},
		Diff: dto.CompareDiff{
			Values: diffValues,
			Rates:  diffRates,
		},
	}, nil
}

// getSnapshotData 获取快照数据
func getSnapshotData(snapshotType string, start, end time.Time) (map[string]interface{}, error) {
	data := make(map[string]interface{})

	switch snapshotType {
	case "user":
		var totalUsers, newUsers int64
		database.DB.Model(&model.User{}).Where("created_at <= ?", end).Count(&totalUsers)
		database.DB.Model(&model.User{}).Where("created_at >= ? AND created_at <= ?", start, end).Count(&newUsers)
		data["total_users"] = totalUsers
		data["new_users"] = newUsers

	case "practice":
		var sessions, questions int64
		var avgAccuracy float64
		database.DB.Model(&model.PracticeSession{}).Where("created_at >= ? AND created_at <= ?", start, end).Count(&sessions)
		database.DB.Model(&model.PracticeRecord{}).
			Where("created_at >= ? AND created_at <= ?", start, end).Count(&questions)
		database.DB.Model(&model.PracticeSession{}).
			Where("status = 1 AND created_at >= ? AND created_at <= ?", start, end).
			Select("COALESCE(AVG(accuracy), 0)").Scan(&avgAccuracy)
		data["sessions"] = sessions
		data["questions"] = questions
		data["avg_accuracy"] = avgAccuracy

	case "exam":
		var exams, records int64
		var avgScore float64
		database.DB.Model(&model.Exam{}).Where("created_at >= ? AND created_at <= ?", start, end).Count(&exams)
		database.DB.Model(&model.ExamRecord{}).Where("created_at >= ? AND created_at <= ?", start, end).Count(&records)
		database.DB.Model(&model.ExamRecord{}).
			Where("status = 2 AND created_at >= ? AND created_at <= ?", start, end).
			Select("COALESCE(AVG(score), 0)").Scan(&avgScore)
		data["exams"] = exams
		data["records"] = records
		data["avg_score"] = avgScore

	default:
		// 尝试从快照表查询
		var snapshot model.StatisticsSnapshot
		if err := database.DB.Where("snapshot_type = ? AND date >= ? AND date <= ?", snapshotType, start, end).
			Order("date DESC").First(&snapshot).Error; err == nil {
			json.Unmarshal([]byte(snapshot.Data), &data)
		}
	}

	return data, nil
}

// toFloat64 转换为float64
func toFloat64(v interface{}) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case int64:
		return float64(val)
	case int:
		return float64(val)
	case float32:
		return float64(val)
	default:
		return 0
	}
}

// ==================== 题目分析 ====================

// GetQuestionDifficulty 题目难度分析
func GetQuestionDifficulty(req *dto.QuestionDifficultyRequest) ([]dto.QuestionDifficultyItem, error) {
	var items []dto.QuestionDifficultyItem

	db := database.DB.Model(&model.Question{}).
		Where("status = 1 AND answer_count > 0")

	if req.CategoryID > 0 {
		db = db.Where("category_id = ?", req.CategoryID)
	}

	var questions []model.Question
	if err := db.Select("id, title, type, difficulty, correct_rate, answer_count").
		Order("answer_count DESC").
		Limit(100).
		Find(&questions).Error; err != nil {
		return nil, fmt.Errorf("查询题目失败: %w", err)
	}

	for _, q := range questions {
		// 计算平均用时
		var avgDuration float64
		database.DB.Model(&model.PracticeRecord{}).
			Where("question_id = ?", q.ID).
			Select("COALESCE(AVG(duration), 0)").
			Scan(&avgDuration)

		items = append(items, dto.QuestionDifficultyItem{
			QuestionID:  q.ID,
			Title:       q.Title,
			Type:        q.Type,
			Difficulty:  q.Difficulty,
			CorrectRate: q.CorrectRate,
			AnswerCount: q.AnswerCount,
			AvgDuration: avgDuration,
		})
	}

	return items, nil
}

// GetQuestionDiscrimination 题目区分度分析
func GetQuestionDiscrimination(req *dto.QuestionDiscriminationRequest) ([]dto.QuestionDiscriminationItem, error) {
	var items []dto.QuestionDiscriminationItem

	db := database.DB.Model(&model.Question{}).
		Where("status = 1 AND answer_count >= 10")

	if req.CategoryID > 0 {
		db = db.Where("category_id = ?", req.CategoryID)
	}

	var questions []model.Question
	if err := db.Select("id, title, type, difficulty, correct_rate").
		Order("answer_count DESC").
		Limit(100).
		Find(&questions).Error; err != nil {
		return nil, fmt.Errorf("查询题目失败: %w", err)
	}

	for _, q := range questions {
		// 计算区分度: (高分组正确率 - 低分组正确率)
		// 高分组：总分排名前27%的用户
		// 低分组：总分排名后27%的用户
		highGroupRate, lowGroupRate := calculateDiscrimination(q.ID)
		discrimination := highGroupRate - lowGroupRate

		items = append(items, dto.QuestionDiscriminationItem{
			QuestionID:     q.ID,
			Title:          q.Title,
			Type:           q.Type,
			Difficulty:     q.Difficulty,
			Discrimination: discrimination,
			HighGroupRate:  highGroupRate,
			LowGroupRate:   lowGroupRate,
		})
	}

	return items, nil
}

// calculateDiscrimination 计算区分度
func calculateDiscrimination(questionID uint) (float64, float64) {
	type groupResult struct {
		Total   int64 `json:"total"`
		Correct int64 `json:"correct"`
	}

	// 高分组
	var highResult groupResult
	database.DB.Raw(`
		SELECT
			COUNT(*) as total,
			SUM(CASE WHEN pr.is_correct = 1 THEN 1 ELSE 0 END) as correct
		FROM practice_records pr
		JOIN practice_sessions ps ON pr.session_id = ps.id
		WHERE pr.question_id = ?
		AND ps.user_id IN (
			SELECT user_id FROM practice_sessions
			WHERE status = 1
			GROUP BY user_id
			ORDER BY AVG(accuracy) DESC
			LIMIT (SELECT COUNT(DISTINCT user_id) * 27 / 100 FROM practice_sessions WHERE status = 1)
		)
	`, questionID).Scan(&highResult)

	// 低分组
	var lowResult groupResult
	database.DB.Raw(`
		SELECT
			COUNT(*) as total,
			SUM(CASE WHEN pr.is_correct = 1 THEN 1 ELSE 0 END) as correct
		FROM practice_records pr
		JOIN practice_sessions ps ON pr.session_id = ps.id
		WHERE pr.question_id = ?
		AND ps.user_id IN (
			SELECT user_id FROM practice_sessions
			WHERE status = 1
			GROUP BY user_id
			ORDER BY AVG(accuracy) ASC
			LIMIT (SELECT COUNT(DISTINCT user_id) * 27 / 100 FROM practice_sessions WHERE status = 1)
		)
	`, questionID).Scan(&lowResult)

	highRate := float64(0)
	if highResult.Total > 0 {
		highRate = float64(highResult.Correct) / float64(highResult.Total) * 100
	}

	lowRate := float64(0)
	if lowResult.Total > 0 {
		lowRate = float64(lowResult.Correct) / float64(lowResult.Total) * 100
	}

	return highRate, lowRate
}

// ==================== 成绩预测与预警 ====================

// GetScorePrediction 成绩预测
func GetScorePrediction(userID uint, req *dto.ScorePredictionRequest) (*dto.ScorePrediction, error) {
	targetUserID := req.UserID
	if targetUserID == 0 {
		targetUserID = userID
	}

	var user model.User
	if err := database.DB.Select("id, nickname").First(&user, targetUserID).Error; err != nil {
		return nil, fmt.Errorf("用户不存在")
	}

	// 获取历史考试记录
	var records []model.ExamRecord
	database.DB.Where("user_id = ? AND status = 2", targetUserID).
		Order("created_at DESC").
		Limit(10).
		Find(&records)

	if len(records) == 0 {
		return &dto.ScorePrediction{
			UserID:   targetUserID,
			UserName: user.Nickname,
			ExamID:   req.ExamID,
			Trend:    "stable",
		}, nil
	}

	// 计算历史平均分和最高分
	var totalScore, bestScore float64
	for _, r := range records {
		totalScore += r.Score
		if r.Score > bestScore {
			bestScore = r.Score
		}
	}
	avgScore := totalScore / float64(len(records))

	// 简单线性预测：基于最近成绩的趋势
	predicted := avgScore
	trend := "stable"

	if len(records) >= 3 {
		recentAvg := (records[0].Score + records[1].Score + records[2].Score) / 3
		if recentAvg > avgScore*1.05 {
			trend = "up"
			predicted = recentAvg * 1.02
		} else if recentAvg < avgScore*0.95 {
			trend = "down"
			predicted = recentAvg * 0.98
		}
	}

	// 置信度基于记录数量
	confidence := float64(len(records)) / 10.0 * 100
	if confidence > 95 {
		confidence = 95
	}

	// 获取考试标题
	examTitle := ""
	if req.ExamID > 0 {
		var exam model.Exam
		if err := database.DB.Select("title").First(&exam, req.ExamID).Error; err == nil {
			examTitle = exam.Title
		}
	}

	return &dto.ScorePrediction{
		UserID:          targetUserID,
		UserName:        user.Nickname,
		ExamID:          req.ExamID,
		ExamTitle:       examTitle,
		PredictedScore:  predicted,
		ConfidenceLevel: confidence,
		HistoryAvg:      avgScore,
		HistoryBest:     bestScore,
		Trend:           trend,
	}, nil
}

// GetScoreAlert 成绩预警
func GetScoreAlert(userID uint, req *dto.ScoreAlertRequest) ([]dto.ScoreAlertItem, error) {
	var items []dto.ScoreAlertItem

	targetUserID := req.UserID
	if targetUserID == 0 {
		targetUserID = userID
	}

	// 获取用户最近的考试记录
	var records []model.ExamRecord
	database.DB.Where("user_id = ? AND status = 2", targetUserID).
		Order("created_at DESC").
		Limit(20).
		Find(&records)

	if len(records) == 0 {
		return items, nil
	}

	var user model.User
	database.DB.Select("id, nickname").First(&user, targetUserID)

	// 计算平均分
	var totalScore float64
	for _, r := range records {
		totalScore += r.Score
	}
	avgScore := totalScore / float64(len(records))

	// 检查最近的成绩是否有异常下降
	if len(records) >= 2 {
		recent := records[0]
		prev := records[1]

		var exam model.Exam
		database.DB.Select("id, title").First(&exam, recent.ExamID)

		// 如果最近成绩比平均分低10%以上，发出预警
		if recent.Score < avgScore*0.9 {
			alertLevel := "warning"
			if recent.Score < avgScore*0.7 {
				alertLevel = "critical"
			}

			trend := "down"
			if recent.Score < prev.Score {
				trend = "down"
			}

			items = append(items, dto.ScoreAlertItem{
				UserID:       targetUserID,
				UserName:     user.Nickname,
				ExamID:       recent.ExamID,
				ExamTitle:    exam.Title,
				LatestScore:  recent.Score,
				AvgScore:     avgScore,
				Trend:        trend,
				AlertLevel:   alertLevel,
				AlertMessage: fmt.Sprintf("成绩%.1f分，低于历史平均%.1f分", recent.Score, avgScore),
			})
		}
	}

	return items, nil
}

// ==================== 班级统计(教师视角) ====================

// GetClassOverview 班级概览
func GetClassOverview(classID uint) (*dto.ClassOverview, error) {
	var class model.Class
	if err := database.DB.First(&class, classID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("班级不存在")
		}
		return nil, fmt.Errorf("查询班级失败: %w", err)
	}

	// 获取班级学生数量
	var totalStudents int64
	database.DB.Model(&model.ClassMember{}).Where("class_id = ? AND role = 1", classID).Count(&totalStudents)

	// 获取活跃学生数量（最近30天有活动）
	var activeStudents int64
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	database.DB.Model(&model.ClassMember{}).
		Where("class_id = ? AND role = 1", classID).
		Where("user_id IN (?)",
			database.DB.Model(&model.PracticeSession{}).
				Select("DISTINCT user_id").
				Where("created_at >= ?", thirtyDaysAgo),
		).
		Count(&activeStudents)

	// 获取班级学生ID列表
	var studentIDs []uint
	database.DB.Model(&model.ClassMember{}).
		Where("class_id = ? AND role = 1", classID).
		Pluck("user_id", &studentIDs)

	var avgScore float64
	var passRate float64
	var practiceCount, examCount int64

	if len(studentIDs) > 0 {
		// 计算平均分
		database.DB.Model(&model.ExamRecord{}).
			Where("user_id IN ? AND status = 2", studentIDs).
			Select("COALESCE(AVG(score), 0)").
			Scan(&avgScore)

		// 计算及格率
		var totalRecords, passRecords int64
		database.DB.Model(&model.ExamRecord{}).
			Where("user_id IN ? AND status = 2", studentIDs).
			Count(&totalRecords)
		if totalRecords > 0 {
			database.DB.Model(&model.ExamRecord{}).
				Joins("JOIN exams ON exams.id = exam_records.exam_id").
				Where("exam_records.user_id IN ? AND exam_records.status = 2 AND exam_records.score >= exams.pass_score", studentIDs).
				Count(&passRecords)
			passRate = float64(passRecords) / float64(totalRecords) * 100
		}

		// 练习次数
		database.DB.Model(&model.PracticeSession{}).
			Where("user_id IN ?", studentIDs).
			Count(&practiceCount)

		// 考试次数
		database.DB.Model(&model.ExamRecord{}).
			Where("user_id IN ?", studentIDs).
			Count(&examCount)
	}

	return &dto.ClassOverview{
		ClassID:        class.ID,
		ClassName:      class.Name,
		TotalStudents:  int(totalStudents),
		ActiveStudents: int(activeStudents),
		AvgScore:       avgScore,
		PassRate:       passRate,
		PracticeCount:  practiceCount,
		ExamCount:      examCount,
	}, nil
}

// GetClassStudents 班级学生成绩列表
func GetClassStudents(classID uint, req *dto.ClassStudentListRequest) ([]dto.ClassStudentItem, int64, error) {
	var total int64
	database.DB.Model(&model.ClassMember{}).Where("class_id = ? AND role = 1", classID).Count(&total)

	// 获取班级学生
	var members []model.ClassMember
	database.DB.Where("class_id = ? AND role = 1", classID).
		Offset(req.GetOffset()).Limit(req.GetLimit()).
		Find(&members)

	var items []dto.ClassStudentItem
	for _, m := range members {
		var user model.User
		database.DB.Select("id, nickname, avatar").First(&user, m.UserID)

		// 计算平均分
		var avgScore, bestScore float64
		database.DB.Model(&model.ExamRecord{}).
			Where("user_id = ? AND status = 2", m.UserID).
			Select("COALESCE(AVG(score), 0)").Scan(&avgScore)
		database.DB.Model(&model.ExamRecord{}).
			Where("user_id = ? AND status = 2", m.UserID).
			Select("COALESCE(MAX(score), 0)").Scan(&bestScore)

		// 练习和考试次数
		var practiceCount, examCount int64
		database.DB.Model(&model.PracticeSession{}).Where("user_id = ?", m.UserID).Count(&practiceCount)
		database.DB.Model(&model.ExamRecord{}).Where("user_id = ?", m.UserID).Count(&examCount)

		// 最后活跃时间
		var lastActive string
		var lastSession model.PracticeSession
		if err := database.DB.Where("user_id = ?", m.UserID).Order("created_at DESC").First(&lastSession).Error; err == nil {
			lastActive = lastSession.CreatedAt.Format("2006-01-02 15:04:05")
		}

		items = append(items, dto.ClassStudentItem{
			UserID:        m.UserID,
			UserName:      user.Nickname,
			Avatar:        user.Avatar,
			AvgScore:      avgScore,
			BestScore:     bestScore,
			PracticeCount: int(practiceCount),
			ExamCount:     int(examCount),
			LastActiveAt:  lastActive,
		})
	}

	return items, total, nil
}

// GetClassPracticeStats 班级练习统计
func GetClassPracticeStats(classID uint, req *dto.ClassPracticeRequest) (*dto.ClassPracticeStats, error) {
	// 获取班级学生ID
	var studentIDs []uint
	database.DB.Model(&model.ClassMember{}).
		Where("class_id = ? AND role = 1", classID).
		Pluck("user_id", &studentIDs)

	if len(studentIDs) == 0 {
		return &dto.ClassPracticeStats{}, nil
	}

	db := database.DB.Model(&model.PracticeSession{}).Where("user_id IN ?", studentIDs)

	if req.StartDate != "" {
		db = db.Where("created_at >= ?", req.StartDate)
	}
	if req.EndDate != "" {
		db = db.Where("created_at <= ?", req.EndDate+" 23:59:59")
	}

	var stats dto.ClassPracticeStats
	db.Count(&stats.TotalSessions)

	var totalQuestions int64
	database.DB.Model(&model.PracticeRecord{}).
		Where("session_id IN (?)",
			database.DB.Model(&model.PracticeSession{}).Select("id").Where("user_id IN ?", studentIDs),
		).
		Count(&totalQuestions)
	stats.TotalQuestions = totalQuestions

	db.Where("status = 1").Select("COALESCE(AVG(accuracy), 0)").Scan(&stats.AvgAccuracy)
	db.Select("COALESCE(SUM(duration), 0)").Scan(&stats.TotalDuration)

	// 每日统计
	start := time.Now().AddDate(0, 0, -30)
	if req.StartDate != "" {
		if t, err := time.Parse("2006-01-02", req.StartDate); err == nil {
			start = t
		}
	}

	var dailyResults []struct {
		Date      string  `json:"date"`
		Sessions  int     `json:"sessions"`
		Questions int     `json:"questions"`
		Accuracy  float64 `json:"accuracy"`
	}

	database.DB.Model(&model.PracticeSession{}).
		Select("DATE(created_at) as date, COUNT(*) as sessions, SUM(total_count) as questions, COALESCE(AVG(accuracy), 0) as accuracy").
		Where("user_id IN ? AND created_at >= ?", studentIDs, start).
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&dailyResults)

	for _, d := range dailyResults {
		stats.DailyStats = append(stats.DailyStats, dto.PracticeDayStat{
			Date:      d.Date,
			Sessions:  d.Sessions,
			Questions: d.Questions,
			Accuracy:  d.Accuracy,
		})
	}

	return &stats, nil
}

// GetClassExamStats 班级考试统计
func GetClassExamStats(classID uint, req *dto.ClassExamRequest) (*dto.ClassExamStats, error) {
	var studentIDs []uint
	database.DB.Model(&model.ClassMember{}).
		Where("class_id = ? AND role = 1", classID).
		Pluck("user_id", &studentIDs)

	if len(studentIDs) == 0 {
		return &dto.ClassExamStats{}, nil
	}

	// 获取班级关联的考试
	var examIDs []uint
	database.DB.Model(&model.Exam{}).
		Where("class_id = ?", classID).
		Pluck("id", &examIDs)

	var stats dto.ClassExamStats
	stats.TotalExams = int64(len(examIDs))

	db := database.DB.Model(&model.ExamRecord{}).
		Where("user_id IN ?", studentIDs)
	if len(examIDs) > 0 {
		db = db.Where("exam_id IN ?", examIDs)
	}
	if req.StartDate != "" {
		db = db.Where("created_at >= ?", req.StartDate)
	}
	if req.EndDate != "" {
		db = db.Where("created_at <= ?", req.EndDate+" 23:59:59")
	}

	db.Count(&stats.TotalRecords)
	db.Where("status = 2").Select("COALESCE(AVG(score), 0)").Scan(&stats.AvgScore)

	// 及格率
	var totalGraded, passGraded int64
	db.Where("status = 2").Count(&totalGraded)
	if totalGraded > 0 {
		database.DB.Model(&model.ExamRecord{}).
			Joins("JOIN exams ON exams.id = exam_records.exam_id").
			Where("exam_records.user_id IN ? AND exam_records.status = 2 AND exam_records.score >= exams.pass_score", studentIDs).
			Count(&passGraded)
		stats.PassRate = float64(passGraded) / float64(totalGraded) * 100
	}

	// 考试列表
	for _, examID := range examIDs {
		var exam model.Exam
		database.DB.Select("id, title, start_time").First(&exam, examID)

		var examAvg float64
		var examPass int64
		var examTotal int64
		database.DB.Model(&model.ExamRecord{}).
			Where("exam_id = ? AND user_id IN ? AND status = 2", examID, studentIDs).
			Select("COALESCE(AVG(score), 0)").Scan(&examAvg)
		database.DB.Model(&model.ExamRecord{}).
			Where("exam_id = ? AND user_id IN ? AND status = 2", examID, studentIDs).
			Count(&examTotal)
		if examTotal > 0 {
			database.DB.Model(&model.ExamRecord{}).
				Joins("JOIN exams ON exams.id = exam_records.exam_id").
				Where("exam_records.exam_id = ? AND exam_records.user_id IN ? AND exam_records.status = 2 AND exam_records.score >= exams.pass_score", examID, studentIDs).
				Count(&examPass)
		}

		examPassRate := float64(0)
		if examTotal > 0 {
			examPassRate = float64(examPass) / float64(examTotal) * 100
		}

		stats.ExamList = append(stats.ExamList, dto.ClassExamItem{
			ExamID:    exam.ID,
			Title:     exam.Title,
			StartTime: exam.StartTime.Format("2006-01-02 15:04"),
			AvgScore:  examAvg,
			PassRate:  examPassRate,
			SubCount:  int(examTotal),
		})
	}

	return &stats, nil
}

// GetClassQuestionStats 班级题目统计
func GetClassQuestionStats(classID uint) (*dto.ClassQuestionStats, error) {
	var studentIDs []uint
	database.DB.Model(&model.ClassMember{}).
		Where("class_id = ? AND role = 1", classID).
		Pluck("user_id", &studentIDs)

	if len(studentIDs) == 0 {
		return &dto.ClassQuestionStats{}, nil
	}

	var stats dto.ClassQuestionStats

	// 获取练习过的题目总数
	database.DB.Model(&model.PracticeRecord{}).
		Where("session_id IN (?)",
			database.DB.Model(&model.PracticeSession{}).Select("id").Where("user_id IN ?", studentIDs),
		).
		Distinct("question_id").Count(&stats.TotalQuestions)

	// 总正确率
	var totalRecords, correctRecords int64
	database.DB.Model(&model.PracticeRecord{}).
		Where("session_id IN (?)",
			database.DB.Model(&model.PracticeSession{}).Select("id").Where("user_id IN ?", studentIDs),
		).
		Count(&totalRecords)
	if totalRecords > 0 {
		database.DB.Model(&model.PracticeRecord{}).
			Where("session_id IN (?) AND is_correct = 1",
				database.DB.Model(&model.PracticeSession{}).Select("id").Where("user_id IN ?", studentIDs),
			).
			Count(&correctRecords)
		stats.CorrectRate = float64(correctRecords) / float64(totalRecords) * 100
	}

	// 按题型统计
	var typeResults []dto.QuestionTypeStat
	database.DB.Raw(`
		SELECT q.type, COUNT(DISTINCT pr.question_id) as count,
		COALESCE(AVG(CASE WHEN pr.is_correct = 1 THEN 100.0 ELSE 0 END), 0) as avg_rate
		FROM practice_records pr
		JOIN questions q ON pr.question_id = q.id
		JOIN practice_sessions ps ON pr.session_id = ps.id
		WHERE ps.user_id IN ?
		GROUP BY q.type
	`, studentIDs).Scan(&typeResults)
	stats.TypeStats = typeResults

	// 按难度统计
	var diffResults []dto.QuestionDiffStat
	database.DB.Raw(`
		SELECT q.difficulty, COUNT(DISTINCT pr.question_id) as count,
		COALESCE(AVG(CASE WHEN pr.is_correct = 1 THEN 100.0 ELSE 0 END), 0) as avg_rate
		FROM practice_records pr
		JOIN questions q ON pr.question_id = q.id
		JOIN practice_sessions ps ON pr.session_id = ps.id
		WHERE ps.user_id IN ?
		GROUP BY q.difficulty
	`, studentIDs).Scan(&diffResults)
	stats.DifficultyStats = diffResults

	return &stats, nil
}

// ==================== 移动端统计 ====================

// GetMobileOverview 移动端个人概览
func GetMobileOverview(userID uint) (*dto.MobileOverview, error) {
	var overview dto.MobileOverview

	// 总练习次数
	var totalPractice int64
	database.DB.Model(&model.PracticeSession{}).Where("user_id = ?", userID).Count(&totalPractice)
	overview.TotalPractice = int(totalPractice)

	// 总练习题目数
	var totalQuestions int64
	database.DB.Model(&model.PracticeRecord{}).
		Where("session_id IN (?)",
			database.DB.Model(&model.PracticeSession{}).Select("id").Where("user_id = ?", userID),
		).
		Count(&totalQuestions)
	overview.TotalQuestions = int(totalQuestions)

	// 平均正确率
	database.DB.Model(&model.PracticeSession{}).
		Where("user_id = ? AND status = 1", userID).
		Select("COALESCE(AVG(accuracy), 0)").Scan(&overview.AvgAccuracy)

	// 考试次数
	var examCount int64
	database.DB.Model(&model.ExamRecord{}).Where("user_id = ?", userID).Count(&examCount)
	overview.TotalExams = int(examCount)

	// 考试平均分
	database.DB.Model(&model.ExamRecord{}).
		Where("user_id = ? AND status = 2", userID).
		Select("COALESCE(AVG(score), 0)").Scan(&overview.AvgScore)

	// 学习天数
	var studyDays int64
	database.DB.Model(&model.PracticeSession{}).
		Where("user_id = ?", userID).
		Distinct("DATE(created_at)").Count(&studyDays)
	overview.StudyDays = int(studyDays)

	// 总学习时长
	database.DB.Model(&model.PracticeSession{}).
		Where("user_id = ?", userID).
		Select("COALESCE(SUM(duration), 0)").Scan(&overview.StudyDuration)

	// 排名
	var rank int64
	database.DB.Model(&model.UserAbilityProfile{}).
		Where("rank_score > (?)",
			database.DB.Model(&model.UserAbilityProfile{}).Select("rank_score").Where("user_id = ?", userID),
		).
		Count(&rank)
	overview.Rank = int(rank) + 1

	return &overview, nil
}

// GetMobilePracticeStats 移动端练习统计
func GetMobilePracticeStats(userID uint, req *dto.MobilePracticeRequest) (*dto.MobilePracticeStats, error) {
	days := req.Days
	if days <= 0 || days > 90 {
		days = 30
	}

	var stats dto.MobilePracticeStats

	today := time.Now().Format("2006-01-02")
	weekStart := time.Now().AddDate(0, 0, -int(time.Now().Weekday())).Format("2006-01-02")
	monthStart := time.Now().AddDate(0, 0, -days).Format("2006-01-02")

	// 今日统计
	var todayCount, weekCount, monthCount int64
	database.DB.Model(&model.PracticeSession{}).
		Where("user_id = ? AND DATE(created_at) = ?", userID, today).
		Count(&todayCount)
	stats.TodayCount = int(todayCount)
	database.DB.Model(&model.PracticeSession{}).
		Where("user_id = ? AND DATE(created_at) = ? AND status = 1", userID, today).
		Select("COALESCE(AVG(accuracy), 0)").Scan(&stats.TodayAccuracy)

	// 本周统计
	database.DB.Model(&model.PracticeSession{}).
		Where("user_id = ? AND DATE(created_at) >= ?", userID, weekStart).
		Count(&weekCount)
	stats.WeekCount = int(weekCount)
	database.DB.Model(&model.PracticeSession{}).
		Where("user_id = ? AND DATE(created_at) >= ? AND status = 1", userID, weekStart).
		Select("COALESCE(AVG(accuracy), 0)").Scan(&stats.WeekAccuracy)

	// 本月统计
	database.DB.Model(&model.PracticeSession{}).
		Where("user_id = ? AND DATE(created_at) >= ?", userID, monthStart).
		Count(&monthCount)
	stats.MonthCount = int(monthCount)

	// 每日统计
	startDate := time.Now().AddDate(0, 0, -days)
	var dailyResults []dto.PracticeDayStat
	database.DB.Model(&model.PracticeSession{}).
		Select("DATE(created_at) as date, COUNT(*) as sessions, SUM(total_count) as questions, COALESCE(AVG(accuracy), 0) as accuracy").
		Where("user_id = ? AND created_at >= ?", userID, startDate).
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&dailyResults)
	stats.DailyStats = dailyResults

	return &stats, nil
}

// GetMobileWrongStats 移动端错题统计
func GetMobileWrongStats(userID uint) (*dto.MobileWrongStats, error) {
	var stats dto.MobileWrongStats

	// 总错题数
	database.DB.Model(&model.WrongQuestion{}).Where("user_id = ?", userID).Count(&stats.TotalWrong)

	// 未掌握错题数
	database.DB.Model(&model.WrongQuestion{}).Where("user_id = ? AND mastered = 0", userID).Count(&stats.UndoWrong)

	// 已掌握错题数
	stats.DoneWrong = stats.TotalWrong - stats.UndoWrong

	// 按题型统计
	var typeResults []dto.WrongTypeStat
	database.DB.Raw(`
		SELECT q.type, COUNT(*) as count
		FROM wrong_questions wq
		JOIN questions q ON wq.question_id = q.id
		WHERE wq.user_id = ?
		GROUP BY q.type
	`, userID).Scan(&typeResults)
	stats.TypeStats = typeResults

	// 最近30天每日新增错题
	startDate := time.Now().AddDate(0, 0, -30)
	var dailyResults []dto.WrongDayStat
	database.DB.Model(&model.WrongQuestion{}).
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("user_id = ? AND created_at >= ?", userID, startDate).
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&dailyResults)
	stats.DailyStats = dailyResults

	return &stats, nil
}

// GetMobileTrend 移动端学习趋势
func GetMobileTrend(userID uint, req *dto.MobileTrendRequest) ([]dto.MobileTrendItem, error) {
	days := req.Days
	if days <= 0 || days > 90 {
		days = 30
	}

	startDate := time.Now().AddDate(0, 0, -days)

	var items []dto.MobileTrendItem
	var results []struct {
		Date          string  `json:"date"`
		PracticeCount int     `json:"practice_count"`
		QuestionCount int     `json:"question_count"`
		Accuracy      float64 `json:"accuracy"`
		Duration      int     `json:"duration"`
	}

	database.DB.Model(&model.PracticeSession{}).
		Select("DATE(created_at) as date, COUNT(*) as practice_count, SUM(total_count) as question_count, COALESCE(AVG(accuracy), 0) as accuracy, COALESCE(SUM(duration), 0) as duration").
		Where("user_id = ? AND created_at >= ?", userID, startDate).
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&results)

	for _, r := range results {
		items = append(items, dto.MobileTrendItem{
			Date:          r.Date,
			PracticeCount: r.PracticeCount,
			QuestionCount: r.QuestionCount,
			Accuracy:      r.Accuracy,
			Duration:      r.Duration,
		})
	}

	return items, nil
}
