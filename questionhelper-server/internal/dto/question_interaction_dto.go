package dto

import "time"

// ==================== 题目笔记 ====================

// NoteInfo 笔记信息
type NoteInfo struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"user_id"`
	QuestionID uint      `json:"question_id"`
	Content    string    `json:"content"`
	IsPublic   bool      `json:"is_public"`
	LikeCount  int       `json:"like_count"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// CreateNoteRequest 创建笔记请求
type CreateNoteRequest struct {
	Content  string `json:"content" binding:"required,max=5000"`
	IsPublic bool   `json:"is_public"`
}

// UpdateNoteRequest 更新笔记请求
type UpdateNoteRequest struct {
	Content  string `json:"content" binding:"omitempty,max=5000"`
	IsPublic *bool  `json:"is_public"`
}

// NoteListRequest 笔记列表请求
type NoteListRequest struct {
	PageRequest
	IsPublic *bool `form:"is_public"`
}

// ==================== 题目评价 ====================

// DifficultyRatingRequest 难度评价请求
type DifficultyRatingRequest struct {
	Rating    int8 `json:"rating" binding:"required,min=1,max=5"`
	IsCorrect bool `json:"is_correct"`
}

// QualityRatingRequest 质量评价请求
type QualityRatingRequest struct {
	Score          int8   `json:"score" binding:"required,min=1,max=5"`
	ClarityScore   int8   `json:"clarity_score" binding:"omitempty,min=1,max=5"`
	RelevanceScore int8   `json:"relevance_score" binding:"omitempty,min=1,max=5"`
	Comment        string `json:"comment" binding:"omitempty,max=500"`
}

// RatingSummary 评价汇总
type RatingSummary struct {
	AvgDifficulty  float64 `json:"avg_difficulty"`
	AvgQuality     float64 `json:"avg_quality"`
	TotalRatings   int     `json:"total_ratings"`
	DifficultyDist map[int]int `json:"difficulty_dist"`
}

// ==================== 题目纠错 ====================

// CreateCorrectionRequest 创建纠错请求
type CreateCorrectionRequest struct {
	CorrectionType int8   `json:"correction_type" binding:"required,min=1,max=6"`
	Description    string `json:"description" binding:"required,max=1000"`
	Screenshot     string `json:"screenshot" binding:"omitempty,max=500"`
}

// CorrectionInfo 纠错信息
type CorrectionInfo struct {
	ID             uint       `json:"id"`
	UserID         uint       `json:"user_id"`
	QuestionID     uint       `json:"question_id"`
	CorrectionType int8       `json:"correction_type"`
	Description    string     `json:"description"`
	Screenshot     string     `json:"screenshot"`
	Status         int8       `json:"status"`
	AdminReply     string     `json:"admin_reply"`
	ProcessedBy    *uint      `json:"processed_by"`
	ProcessedAt    *time.Time `json:"processed_at"`
	Reward         int        `json:"reward"`
	CreatedAt      time.Time  `json:"created_at"`
}

// CorrectionListRequest 纠错列表请求
type CorrectionListRequest struct {
	PageRequest
	Status *int8 `form:"status"`
}
