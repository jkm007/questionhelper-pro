package dto

import "time"

// WrongQuestionInfo 错题信息
type WrongQuestionInfo struct {
	ID         uint         `json:"id"`
	UserID     uint         `json:"user_id"`
	QuestionID uint         `json:"question_id"`
	Question   QuestionInfo `json:"question"`
	Source     int8         `json:"source"`
	SourceID   uint         `json:"source_id"`
	WrongCount int          `json:"wrong_count"`
	LastAnswer string       `json:"last_answer"`
	Mastered   bool         `json:"mastered"`
	CreatedAt  time.Time    `json:"created_at"`
}

// WrongListRequest 错题列表请求
type WrongListRequest struct {
	PageRequest
	Mastered *bool `form:"mastered"`
	Source   *int8 `form:"source"`
}

// ReviewWrongRequest 复习错题请求
type ReviewWrongRequest struct {
	Answer string `json:"answer" binding:"required"`
}

// WrongAnalysisInfo 错题分析信息
type WrongAnalysisInfo struct {
	TotalCount    int                `json:"total_count"`
	MasteredCount int                `json:"mastered_count"`
	ByCategory    []CategoryWrongInfo `json:"by_category"`
	ByType        []TypeWrongInfo     `json:"by_type"`
	ByDifficulty  []DifficultyWrongInfo `json:"by_difficulty"`
}

// CategoryWrongInfo 分类错题统计
type CategoryWrongInfo struct {
	CategoryID   uint   `json:"category_id"`
	CategoryName string `json:"category_name"`
	Count        int    `json:"count"`
}

// TypeWrongInfo 题型错题统计
type TypeWrongInfo struct {
	Type  int8 `json:"type"`
	Count int  `json:"count"`
}

// DifficultyWrongInfo 难度错题统计
type DifficultyWrongInfo struct {
	Difficulty int8 `json:"difficulty"`
	Count      int  `json:"count"`
}
