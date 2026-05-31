package dto

import "time"

// PracticeSessionInfo 练习会话信息
type PracticeSessionInfo struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"user_id"`
	CategoryID   *uint     `json:"category_id"`
	TotalCount   int       `json:"total_count"`
	CorrectCount int       `json:"correct_count"`
	Accuracy     float64   `json:"accuracy"`
	Duration     int       `json:"duration"`
	Status       int8      `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

// PracticeRecordInfo 练习记录信息
type PracticeRecordInfo struct {
	ID         uint   `json:"id"`
	SessionID  uint   `json:"session_id"`
	QuestionID uint   `json:"question_id"`
	Answer     string `json:"answer"`
	IsCorrect  bool   `json:"is_correct"`
	Duration   int    `json:"duration"`
}

// StartPracticeRequest 开始练习请求
type StartPracticeRequest struct {
	CategoryID *uint `json:"category_id"`
	Type       *int8 `json:"type" binding:"omitempty,oneof=1 2 3 4 5"`
	Difficulty *int8 `json:"difficulty" binding:"omitempty,oneof=1 2 3"`
	Count      int   `json:"count" binding:"omitempty,min=1,max=100"`
}

// SubmitPracticeAnswerRequest 提交练习答案请求
type SubmitPracticeAnswerRequest struct {
	QuestionID uint   `json:"question_id" binding:"required"`
	Answer     string `json:"answer" binding:"required"`
	Duration   int    `json:"duration"`
}

// PracticeListRequest 练习历史列表请求
type PracticeListRequest struct {
	PageRequest
}
