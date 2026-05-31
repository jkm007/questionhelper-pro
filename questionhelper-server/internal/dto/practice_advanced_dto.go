package dto

import "time"

// ==================== 模拟考试 ====================

// StartMockExamRequest 开始模拟考试请求
type StartMockExamRequest struct {
	ConfigID uint `json:"config_id" binding:"required"`
}

// MockExamSessionInfo 模拟考试会话信息
type MockExamSessionInfo struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"user_id"`
	ConfigID     uint      `json:"config_id"`
	ConfigName   string    `json:"config_name"`
	CategoryID   *uint     `json:"category_id"`
	TotalCount   int       `json:"total_count"`
	TotalScore   float64   `json:"total_score"`
	PassScore    float64   `json:"pass_score"`
	Duration     int       `json:"duration"`
	Status       int8      `json:"status"`
	StartedAt    time.Time `json:"started_at"`
	SubmittedAt  *time.Time `json:"submitted_at"`
}

// SubmitMockExamRequest 提交模拟考试请求
type SubmitMockExamRequest struct {
	Answers  []MockExamAnswer `json:"answers" binding:"required"`
	Duration int              `json:"duration"`
}

// MockExamAnswer 模拟考试答题
type MockExamAnswer struct {
	QuestionID uint   `json:"question_id" binding:"required"`
	Answer     string `json:"answer" binding:"required"`
	Duration   int    `json:"duration"`
}

// MockExamResultInfo 模拟考试结果
type MockExamResultInfo struct {
	Session  MockExamSessionInfo  `json:"session"`
	Records  []MockExamRecordInfo `json:"records"`
	Score    float64              `json:"score"`
	IsPassed bool                 `json:"is_passed"`
	Rank     int                  `json:"rank"`
}

// MockExamRecordInfo 模拟考试记录
type MockExamRecordInfo struct {
	ID         uint    `json:"id"`
	QuestionID uint    `json:"question_id"`
	Answer     string  `json:"answer"`
	Correct    string  `json:"correct"`
	IsCorrect  bool    `json:"is_correct"`
	Score      float64 `json:"score"`
	Duration   int     `json:"duration"`
}

// MockExamHistoryRequest 模拟考试历史请求
type MockExamHistoryRequest struct {
	PageRequest
	ConfigID *uint `form:"config_id"`
	Status   *int8 `form:"status"`
}

// MockExamHistoryItem 模拟考试历史项
type MockExamHistoryItem struct {
	ID         uint       `json:"id"`
	ConfigID   uint       `json:"config_id"`
	ConfigName string     `json:"config_name"`
	TotalCount int        `json:"total_count"`
	Score      float64    `json:"score"`
	IsPassed   bool       `json:"is_passed"`
	Duration   int        `json:"duration"`
	Status     int8       `json:"status"`
	StartedAt  time.Time  `json:"started_at"`
	SubmittedAt *time.Time `json:"submitted_at"`
}

// ==================== 练习计划 ====================

// CreatePlanRequest 创建练习计划请求
type CreatePlanRequest struct {
	Name          string `json:"name" binding:"required,max=200"`
	Description   string `json:"description" binding:"omitempty,max=500"`
	PlanType      int8   `json:"plan_type" binding:"required,oneof=1 2 3"`
	CategoryID    *uint  `json:"category_id"`
	QuestionType  *int8  `json:"question_type" binding:"omitempty,oneof=1 2 3 4 5"`
	Difficulty    *int8  `json:"difficulty" binding:"omitempty,oneof=1 2 3"`
	DailyCount    int    `json:"daily_count" binding:"omitempty,min=1,max=200"`
	DailyDuration int    `json:"daily_duration" binding:"omitempty,min=1,max=480"`
	StartDate     string `json:"start_date" binding:"required"`
	EndDate       string `json:"end_date"`
	TotalTarget   int    `json:"total_target" binding:"omitempty,min=1"`
}

// UpdatePlanRequest 更新练习计划请求
type UpdatePlanRequest struct {
	Name          string `json:"name" binding:"omitempty,max=200"`
	Description   string `json:"description" binding:"omitempty,max=500"`
	PlanType      int8   `json:"plan_type" binding:"omitempty,oneof=1 2 3"`
	CategoryID    *uint  `json:"category_id"`
	QuestionType  *int8  `json:"question_type" binding:"omitempty,oneof=1 2 3 4 5"`
	Difficulty    *int8  `json:"difficulty" binding:"omitempty,oneof=1 2 3"`
	DailyCount    int    `json:"daily_count" binding:"omitempty,min=1,max=200"`
	DailyDuration int    `json:"daily_duration" binding:"omitempty,min=1,max=480"`
	StartDate     string `json:"start_date"`
	EndDate       string `json:"end_date"`
	TotalTarget   int    `json:"total_target" binding:"omitempty,min=1"`
	Status        *int8  `json:"status" binding:"omitempty,oneof=0 1 2 3"`
}

// PlanInfo 练习计划信息
type PlanInfo struct {
	ID            uint       `json:"id"`
	UserID        uint       `json:"user_id"`
	Name          string     `json:"name"`
	Description   string     `json:"description"`
	PlanType      int8       `json:"plan_type"`
	CategoryID    *uint      `json:"category_id"`
	QuestionType  *int8      `json:"question_type"`
	Difficulty    *int8      `json:"difficulty"`
	DailyCount    int        `json:"daily_count"`
	DailyDuration int        `json:"daily_duration"`
	StartDate     string     `json:"start_date"`
	EndDate       *string    `json:"end_date"`
	TotalTarget   int        `json:"total_target"`
	TotalDone     int        `json:"total_done"`
	Progress      float64    `json:"progress"`
	Status        int8       `json:"status"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// PlanListRequest 练习计划列表请求
type PlanListRequest struct {
	PageRequest
	Status *int8 `form:"status"`
}

// ExecutePlanRequest 执行计划请求
type ExecutePlanRequest struct {
	QuestionNum int `json:"question_num" binding:"omitempty,min=1,max=100"`
	Duration    int `json:"duration"`
}

// ==================== 每日练习 ====================

// DailyPracticeInfo 每日练习信息
type DailyPracticeInfo struct {
	ID            uint   `json:"id"`
	UserID        uint   `json:"user_id"`
	Date          string `json:"date"`
	TotalCount    int    `json:"total_count"`
	CorrectCount  int    `json:"correct_count"`
	Accuracy      float64 `json:"accuracy"`
	Duration      int    `json:"duration"`
	SessionCount  int    `json:"session_count"`
	WrongCount    int    `json:"wrong_count"`
	NewQuestion   int    `json:"new_question"`
	CategoryStats string `json:"category_stats"`
	IsCompleted   bool   `json:"is_completed"`
	Target        int    `json:"target"`
}

// CompleteDailyRequest 完成每日练习请求
type CompleteDailyRequest struct {
	QuestionCount int `json:"question_count" binding:"required,min=1"`
	CorrectCount  int `json:"correct_count" binding:"omitempty,min=0"`
	Duration      int `json:"duration"`
}

// ==================== 练习打卡 ====================

// CheckinRequest 打卡请求
type CheckinRequest struct {
	QuestionCount int `json:"question_count"`
	Duration      int `json:"duration"`
}

// CheckinInfo 打卡信息
type CheckinInfo struct {
	ID            uint   `json:"id"`
	UserID        uint   `json:"user_id"`
	Date          string `json:"date"`
	IsCheckin     bool   `json:"is_checkin"`
	QuestionCount int    `json:"question_count"`
	Duration      int    `json:"duration"`
	Streak        int    `json:"streak"`
	Reward        int    `json:"reward"`
}

// CheckinCalendarRequest 打卡日历请求
type CheckinCalendarRequest struct {
	Year  int `form:"year" binding:"required"`
	Month int `form:"month" binding:"required,min=1,max=12"`
}

// CheckinCalendarItem 打卡日历项
type CheckinCalendarItem struct {
	Date          string `json:"date"`
	IsCheckin     bool   `json:"is_checkin"`
	QuestionCount int    `json:"question_count"`
	Streak        int    `json:"streak"`
}

// ==================== 排行榜 ====================

// LeaderboardRequest 排行榜请求
type LeaderboardRequest struct {
	RankType int8 `form:"rank_type" binding:"required,oneof=1 2 3 4"`
	Limit    int  `form:"limit" binding:"omitempty,min=1,max=100"`
}

// LeaderboardItem 排行榜项
type LeaderboardItem struct {
	UserID     uint    `json:"user_id"`
	Nickname   string  `json:"nickname"`
	Avatar     string  `json:"avatar"`
	RankPos    int     `json:"rank_pos"`
	Score      float64 `json:"score"`
	Accuracy   float64 `json:"accuracy"`
	TotalCount int     `json:"total_count"`
	Duration   int     `json:"duration"`
}

// ==================== 闯关模式 ====================

// ChallengeLevelInfo 关卡信息
type ChallengeLevelInfo struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Level         int    `json:"level"`
	CategoryID    *uint  `json:"category_id"`
	QuestionCount int    `json:"question_count"`
	PassAccuracy  float64 `json:"pass_accuracy"`
	PassScore     int    `json:"pass_score"`
	TimeLimit     int    `json:"time_limit"`
	Difficulty    int8   `json:"difficulty"`
	Icon          string `json:"icon"`
	Badge         string `json:"badge"`
	PreLevel      int    `json:"pre_level"`
	Status        int8   `json:"status"`
	Sort          int    `json:"sort"`
	// 用户进度
	UserStatus    int8   `json:"user_status"`
	BestAccuracy  float64 `json:"best_accuracy"`
	Attempts      int    `json:"attempts"`
	IsLocked      bool   `json:"is_locked"`
}

// StartChallengeRequest 开始闯关请求
type StartChallengeRequest struct {
	LevelID uint `json:"level_id" binding:"required"`
}

// ChallengeSessionInfo 闯关会话信息
type ChallengeSessionInfo struct {
	ID            uint      `json:"id"`
	UserID        uint      `json:"user_id"`
	LevelID       uint      `json:"level_id"`
	LevelName     string    `json:"level_name"`
	TotalCount    int       `json:"total_count"`
	PassAccuracy  float64   `json:"pass_accuracy"`
	TimeLimit     int       `json:"time_limit"`
	Status        int8      `json:"status"`
	StartedAt     time.Time `json:"started_at"`
}

// SubmitChallengeRequest 提交闯关请求
type SubmitChallengeRequest struct {
	Answers  []ChallengeAnswer `json:"answers" binding:"required"`
	Duration int               `json:"duration"`
}

// ChallengeAnswer 闯关答题
type ChallengeAnswer struct {
	QuestionID uint   `json:"question_id" binding:"required"`
	Answer     string `json:"answer" binding:"required"`
	Duration   int    `json:"duration"`
}

// ChallengeResultInfo 闯关结果
type ChallengeResultInfo struct {
	LevelID      uint    `json:"level_id"`
	LevelName    string  `json:"level_name"`
	IsPassed     bool    `json:"is_passed"`
	Accuracy     float64 `json:"accuracy"`
	Duration     int     `json:"duration"`
	Score        int     `json:"score"`
	BestAccuracy float64 `json:"best_accuracy"`
	Attempts     int     `json:"attempts"`
}

// ChallengeProgressInfo 闯关进度
type ChallengeProgressInfo struct {
	TotalLevels    int                      `json:"total_levels"`
	PassedLevels   int                      `json:"passed_levels"`
	TotalScore     int                      `json:"total_score"`
	CurrentLevel   int                      `json:"current_level"`
	ProgressList   []ChallengeLevelProgress `json:"progress_list"`
}

// ChallengeLevelProgress 关卡进度
type ChallengeLevelProgress struct {
	LevelID      uint    `json:"level_id"`
	LevelName    string  `json:"level_name"`
	Status       int8    `json:"status"`
	BestAccuracy float64 `json:"best_accuracy"`
	Attempts     int     `json:"attempts"`
	Score        int     `json:"score"`
}
