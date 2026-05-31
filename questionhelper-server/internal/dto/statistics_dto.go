package dto

// UserStatistics 用户统计
type UserStatistics struct {
	TotalUsers    int64 `json:"total_users"`
	ActiveUsers   int64 `json:"active_users"`
	NewUsersToday int64 `json:"new_users_today"`
	NewUsersWeek  int64 `json:"new_users_week"`
}

// PracticeStatistics 练习统计
type PracticeStatistics struct {
	TotalSessions  int64   `json:"total_sessions"`
	TotalQuestions int64   `json:"total_questions"`
	AvgAccuracy    float64 `json:"avg_accuracy"`
	TotalDuration  int64   `json:"total_duration"`
}

// ExamStatistics 考试统计
type ExamStatistics struct {
	TotalExams    int64   `json:"total_exams"`
	TotalRecords  int64   `json:"total_records"`
	AvgScore      float64 `json:"avg_score"`
	PassRate      float64 `json:"pass_rate"`
}

// ClassStatistics 班级统计
type ClassStatistics struct {
	TotalClasses  int64 `json:"total_classes"`
	TotalMembers  int64 `json:"total_members"`
	AvgMemberCount float64 `json:"avg_member_count"`
}

// RankInfo 排行榜信息
type RankInfo struct {
	Rank     int    `json:"rank"`
	UserID   uint   `json:"user_id"`
	UserName string `json:"user_name"`
	Avatar   string `json:"avatar"`
	Score    float64 `json:"score"`
	Count    int    `json:"count"`
}

// RankingRequest 排行榜请求
type RankingRequest struct {
	PageRequest
	Type int `form:"type" binding:"required,oneof=1 2 3"` // 1:练习 2:考试 3:积分
}

// DashboardInfo 仪表盘信息
type DashboardInfo struct {
	UserStats     UserStatistics     `json:"user_stats"`
	PracticeStats PracticeStatistics `json:"practice_stats"`
	ExamStats     ExamStatistics     `json:"exam_stats"`
	ClassStats    ClassStatistics    `json:"class_stats"`
}
