package dto

import "time"

// ExamMonitorInfo 考试监控信息
type ExamMonitorInfo struct {
	ExamID        uint              `json:"exam_id"`
	ExamTitle     string            `json:"exam_title"`
	TotalStudents int               `json:"total_students"`
	OnlineCount   int               `json:"online_count"`
	SubmittedCount int              `json:"submitted_count"`
	WarningCount  int               `json:"warning_count"`
	OnlineUsers   []OnlineUserInfo  `json:"online_users"`
	Warnings      []WarningInfo     `json:"warnings"`
}

// OnlineUserInfo 在线考生信息
type OnlineUserInfo struct {
	UserID       uint      `json:"user_id"`
	Username     string    `json:"username"`
	Nickname     string    `json:"nickname"`
	Progress     float64   `json:"progress"`      // 答题进度(百分比)
	CurrentQID   uint      `json:"current_qid"`   // 当前题目
	DurationUsed int       `json:"duration_used"`  // 已用时间(秒)
	LastActive   time.Time `json:"last_active"`
	IP           string    `json:"ip"`
}

// WarningInfo 异常信息
type WarningInfo struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"user_id"`
	Username   string    `json:"username"`
	Nickname   string    `json:"nickname"`
	Type       string    `json:"type"`
	Detail     string    `json:"detail"`
	CreatedAt  time.Time `json:"created_at"`
}

// ReviewRequest 阅卷请求
type ReviewRequest struct {
	RecordID uint           `json:"record_id" binding:"required"`
	Answers  []ReviewAnswer `json:"answers" binding:"required"`
	Comment  string         `json:"comment"` // 总评
}

// ReviewAnswer 阅卷答案
type ReviewAnswer struct {
	AnswerID uint    `json:"answer_id" binding:"required"`
	Score    float64 `json:"score" binding:"required,min=0"`
	Note     string  `json:"note"` // 阅卷备注
}

// BatchReviewRequest 批量阅卷请求
type BatchReviewRequest struct {
	Records []ReviewRequest `json:"records" binding:"required"`
}

// ExamAnalysisResponse 考试分析响应
type ExamAnalysisResponse struct {
	BasicInfo    ExamBasicInfo    `json:"basic_info"`
	ScoreStats   ScoreStatistics  `json:"score_stats"`
	QuestionStats []QuestionStat `json:"question_stats"`
	TopStudents  []StudentScore  `json:"top_students"`
}

// ExamBasicInfo 考试基本信息
type ExamBasicInfo struct {
	ExamID        uint    `json:"exam_id"`
	Title         string  `json:"title"`
	TotalStudents int     `json:"total_students"`
	SubmitCount   int     `json:"submit_count"`
	AvgDuration   int     `json:"avg_duration"` // 平均用时(秒)
}

// ScoreStatistics 分数统计
type ScoreStatistics struct {
	AvgScore  float64     `json:"avg_score"`
	MaxScore  float64     `json:"max_score"`
	MinScore  float64     `json:"min_score"`
	Median    float64     `json:"median"`
	PassRate  float64     `json:"pass_rate"`
	ExcellentRate float64 `json:"excellent_rate"` // 优秀率(>=90%)
	Distribution []ScoreDist `json:"distribution"`
}

// QuestionStat 题目统计
type QuestionStat struct {
	QuestionID uint    `json:"question_id"`
	Title      string  `json:"title"`
	Type       int8    `json:"type"`
	CorrectRate float64 `json:"correct_rate"`
	AvgScore   float64 `json:"avg_score"`
	FullScore  int     `json:"full_score"`
}

// StudentScore 学生成绩
type StudentScore struct {
	UserID   uint    `json:"user_id"`
	Username string  `json:"username"`
	Nickname string  `json:"nickname"`
	Score    float64 `json:"score"`
	Duration int     `json:"duration"`
	Rank     int     `json:"rank"`
}
