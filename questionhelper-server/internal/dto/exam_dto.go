package dto

import "time"

// ExamInfo 考试信息
type ExamInfo struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	PaperID     uint      `json:"paper_id"`
	ClassID     *uint     `json:"class_id"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Duration    int       `json:"duration"`
	TotalScore  float64   `json:"total_score"`
	PassScore   float64   `json:"pass_score"`
	MaxAttempts int       `json:"max_attempts"`
	Shuffle     bool      `json:"shuffle"`
	ShowAnswer  int8      `json:"show_answer"`
	AntiCheat   int8      `json:"anti_cheat"`
	Status      int8      `json:"status"`
	CreatorID   uint      `json:"creator_id"`
}

// PaperInfo 试卷信息
type PaperInfo struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	TotalScore  float64 `json:"total_score"`
	TotalCount  int    `json:"total_count"`
	Type        int8   `json:"type"`
}

// ExamRecordInfo 考试记录信息
type ExamRecordInfo struct {
	ID         uint       `json:"id"`
	ExamID     uint       `json:"exam_id"`
	UserID     uint       `json:"user_id"`
	Score      float64    `json:"score"`
	Status     int8       `json:"status"`
	StartTime  time.Time  `json:"start_time"`
	SubmitTime *time.Time `json:"submit_time"`
}

// CreateExamRequest 创建考试请求
type CreateExamRequest struct {
	Title       string    `json:"title" binding:"required,max=200"`
	Description string    `json:"description"`
	PaperID     uint      `json:"paper_id" binding:"required"`
	ClassID     *uint     `json:"class_id"`
	StartTime   time.Time `json:"start_time" binding:"required"`
	EndTime     time.Time `json:"end_time" binding:"required"`
	Duration    int       `json:"duration" binding:"required,min=1"`
	TotalScore  float64   `json:"total_score" binding:"required"`
	PassScore   float64   `json:"pass_score" binding:"required"`
	MaxAttempts int       `json:"max_attempts" binding:"omitempty,min=1"`
	Shuffle     bool      `json:"shuffle"`
	ShowAnswer  int8      `json:"show_answer" binding:"omitempty,oneof=0 1 2"`
	AntiCheat   int8      `json:"anti_cheat" binding:"omitempty,oneof=0 1 2"`
}

// UpdateExamRequest 更新考试请求
type UpdateExamRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Duration    int       `json:"duration"`
	TotalScore  float64   `json:"total_score"`
	PassScore   float64   `json:"pass_score"`
	MaxAttempts int       `json:"max_attempts"`
	Shuffle     *bool     `json:"shuffle"`
	ShowAnswer  int8      `json:"show_answer"`
	AntiCheat   int8      `json:"anti_cheat"`
}

// SubmitExamRequest 提交考试请求
type SubmitExamRequest struct {
	Answers []AnswerRequest `json:"answers" binding:"required"`
}

// AnswerRequest 答题请求
type AnswerRequest struct {
	QuestionID uint   `json:"question_id" binding:"required"`
	Answer     string `json:"answer" binding:"required"`
}

// ExamListRequest 考试列表请求
type ExamListRequest struct {
	PageRequest
	Status  *int8  `form:"status"`
	ClassID *uint  `form:"class_id"`
	Keyword string `form:"keyword"`
}

// CreatePaperRequest 创建试卷请求
type CreatePaperRequest struct {
	Title       string `json:"title" binding:"required,max=200"`
	Description string `json:"description"`
	Type        int8   `json:"type" binding:"required,oneof=1 2"`
}
