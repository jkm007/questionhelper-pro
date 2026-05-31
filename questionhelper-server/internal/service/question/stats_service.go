package question

import (
	"fmt"
	"time"

	"questionhelper-server/internal/dto"
	"questionhelper-server/pkg/database"
)

// GetQuestionStats 获取题目统计
func GetQuestionStats() (*dto.QuestionStatsResponse, error) {
	stats := &dto.QuestionStatsResponse{}

	// 概览统计
	if err := getOverviewStats(&stats.Overview); err != nil {
		return nil, fmt.Errorf("获取概览统计失败: %w", err)
	}

	// 按题型统计
	byType, err := getTypeStats()
	if err != nil {
		return nil, fmt.Errorf("获取题型统计失败: %w", err)
	}
	stats.ByType = byType

	// 按分类统计
	byCategory, err := getCategoryStats()
	if err != nil {
		return nil, fmt.Errorf("获取分类统计失败: %w", err)
	}
	stats.ByCategory = byCategory

	// 按难度统计
	byDifficulty, err := getDifficultyStats()
	if err != nil {
		return nil, fmt.Errorf("获取难度统计失败: %w", err)
	}
	stats.ByDifficulty = byDifficulty

	// 创建趋势
	trend, err := getCreateTrend(30)
	if err != nil {
		return nil, fmt.Errorf("获取创建趋势失败: %w", err)
	}
	stats.Trend = trend

	return stats, nil
}

func getOverviewStats(stats *dto.QuestionStats) error {
	db := database.DB.Model(nil)

	// 总数
	db.Table("questions").Where("deleted_at IS NULL").Count(&stats.TotalCount)

	// 各状态数量
	db.Table("questions").Where("status = 0 AND deleted_at IS NULL").Count(&stats.DraftCount)
	db.Table("questions").Where("status = 1 AND deleted_at IS NULL").Count(&stats.PublishedCount)
	db.Table("questions").Where("status = 2 AND deleted_at IS NULL").Count(&stats.ArchivedCount)

	// 今日新增
	today := time.Now().Format("2006-01-02")
	db.Table("questions").Where("DATE(created_at) = ? AND deleted_at IS NULL", today).Count(&stats.TodayCount)

	// 本周新增
	weekAgo := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	db.Table("questions").Where("created_at >= ? AND deleted_at IS NULL", weekAgo).Count(&stats.WeekCount)

	return nil
}

func getTypeStats() ([]dto.TypeStats, error) {
	var results []dto.TypeStats
	typeNames := map[int8]string{
		1: "单选题",
		2: "多选题",
		3: "判断题",
		4: "填空题",
		5: "简答题",
	}

	rows, err := database.DB.Raw(`
		SELECT type, COUNT(*) as count
		FROM questions
		WHERE deleted_at IS NULL
		GROUP BY type
		ORDER BY count DESC
	`).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var total int64
	for rows.Next() {
		var stat dto.TypeStats
		rows.Scan(&stat.Type, &stat.Count)
		stat.Name = typeNames[stat.Type]
		total += stat.Count
		results = append(results, stat)
	}

	// 计算占比
	for i := range results {
		if total > 0 {
			results[i].Rate = float64(results[i].Count) / float64(total) * 100
		}
	}

	return results, nil
}

func getCategoryStats() ([]dto.CategoryStats, error) {
	var results []dto.CategoryStats

	rows, err := database.DB.Raw(`
		SELECT c.id, c.name, COUNT(q.id) as count
		FROM categories c
		LEFT JOIN questions q ON q.category_id = c.id AND q.deleted_at IS NULL
		WHERE c.deleted_at IS NULL
		GROUP BY c.id, c.name
		HAVING count > 0
		ORDER BY count DESC
		LIMIT 20
	`).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var total int64
	for rows.Next() {
		var stat dto.CategoryStats
		rows.Scan(&stat.CategoryID, &stat.CategoryName, &stat.Count)
		total += stat.Count
		results = append(results, stat)
	}

	// 计算占比
	for i := range results {
		if total > 0 {
			results[i].Rate = float64(results[i].Count) / float64(total) * 100
		}
	}

	return results, nil
}

func getDifficultyStats() ([]dto.DifficultyStats, error) {
	var results []dto.DifficultyStats
	diffNames := map[int8]string{
		1: "简单",
		2: "中等",
		3: "困难",
	}

	rows, err := database.DB.Raw(`
		SELECT difficulty, COUNT(*) as count
		FROM questions
		WHERE deleted_at IS NULL
		GROUP BY difficulty
		ORDER BY difficulty
	`).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var total int64
	for rows.Next() {
		var stat dto.DifficultyStats
		rows.Scan(&stat.Difficulty, &stat.Count)
		stat.Name = diffNames[stat.Difficulty]
		total += stat.Count
		results = append(results, stat)
	}

	// 计算占比
	for i := range results {
		if total > 0 {
			results[i].Rate = float64(results[i].Count) / float64(total) * 100
		}
	}

	return results, nil
}

func getCreateTrend(days int) ([]dto.TrendStats, error) {
	var results []dto.TrendStats

	startDate := time.Now().AddDate(0, 0, -days).Format("2006-01-02")

	rows, err := database.DB.Raw(`
		SELECT DATE(created_at) as date, COUNT(*) as count
		FROM questions
		WHERE created_at >= ? AND deleted_at IS NULL
		GROUP BY DATE(created_at)
		ORDER BY date
	`, startDate).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var stat dto.TrendStats
		rows.Scan(&stat.Date, &stat.Count)
		results = append(results, stat)
	}

	return results, nil
}
