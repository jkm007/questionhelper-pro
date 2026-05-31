package statistics

import (
	"fmt"
	"time"

	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	"questionhelper-server/pkg/database"
)

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
		// 获取所有考试的及格分数
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
