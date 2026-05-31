package dto

// QuestionStats 题目统计概览
type QuestionStats struct {
	TotalCount     int64              `json:"total_count"`
	DraftCount     int64              `json:"draft_count"`
	PublishedCount int64              `json:"published_count"`
	ArchivedCount  int64              `json:"archived_count"`
	TodayCount     int64              `json:"today_count"`
	WeekCount      int64              `json:"week_count"`
}

// TypeStats 题型统计
type TypeStats struct {
	Type  int8   `json:"type"`
	Name  string `json:"name"`
	Count int64  `json:"count"`
	Rate  float64 `json:"rate"`
}

// CategoryStats 分类统计
type CategoryStats struct {
	CategoryID   uint   `json:"category_id"`
	CategoryName string `json:"category_name"`
	Count        int64  `json:"count"`
	Rate         float64 `json:"rate"`
}

// DifficultyStats 难度统计
type DifficultyStats struct {
	Difficulty int8   `json:"difficulty"`
	Name       string `json:"name"`
	Count      int64  `json:"count"`
	Rate       float64 `json:"rate"`
}

// TrendStats 创建趋势
type TrendStats struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

// QuestionStatsResponse 题目统计响应
type QuestionStatsResponse struct {
	Overview     QuestionStats    `json:"overview"`
	ByType       []TypeStats      `json:"by_type"`
	ByCategory   []CategoryStats  `json:"by_category"`
	ByDifficulty []DifficultyStats `json:"by_difficulty"`
	Trend        []TrendStats     `json:"trend"`
}
