package dto

import "time"

// WrongQuestionInfo 错题信息
type WrongQuestionInfo struct {
	ID            uint              `json:"id"`
	UserID        uint              `json:"user_id"`
	QuestionID    uint              `json:"question_id"`
	Question      QuestionInfo      `json:"question"`
	Source        int8              `json:"source"`
	SourceID      uint              `json:"source_id"`
	WrongCount    int               `json:"wrong_count"`
	LastAnswer    string            `json:"last_answer"`
	Mastered      bool              `json:"mastered"`
	CorrectStreak int               `json:"correct_streak"`
	LastReviewAt  *time.Time        `json:"last_review_at"`
	NextReviewAt  *time.Time        `json:"next_review_at"`
	ReviewCount   int               `json:"review_count"`
	IsFavorite    bool              `json:"is_favorite"`
	Tags          []WrongTagInfo    `json:"tags,omitempty"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}

// WrongListRequest 错题列表请求
type WrongListRequest struct {
	PageRequest
	Mastered *bool `form:"mastered"`
	Source   *int8 `form:"source"`
}

// WrongSearchRequest 错题搜索请求
type WrongSearchRequest struct {
	PageRequest
	Keyword    string `form:"keyword"`
	CategoryID *uint  `form:"category_id"`
	Difficulty *int8  `form:"difficulty"`
	Type       *int8  `form:"type"`
	Source     *int8  `form:"source"`
	Mastered   *bool  `form:"mastered"`
	IsFavorite *bool  `form:"is_favorite"`
	TagID      *uint  `form:"tag_id"`
	StartDate  string `form:"start_date"`
	EndDate    string `form:"end_date"`
	SortBy     string `form:"sort_by"`    // created_at, wrong_count, difficulty
	SortOrder  string `form:"sort_order"` // asc, desc
}

// ReviewWrongRequest 复习错题请求
type ReviewWrongRequest struct {
	Answer string `json:"answer" binding:"required"`
}

// ReviewWrongAdvancedRequest 高级复习请求
type ReviewWrongAdvancedRequest struct {
	IsCorrect  bool   `json:"is_correct"`
	Answer     string `json:"answer"`
	Duration   int    `json:"duration"`    // 用时(秒)
	ReviewType int8   `json:"review_type"` // 1=手动,2=计划
}

// WrongAnalysisInfo 错题分析信息
type WrongAnalysisInfo struct {
	TotalCount    int                  `json:"total_count"`
	MasteredCount int                  `json:"mastered_count"`
	ByCategory    []CategoryWrongInfo  `json:"by_category"`
	ByType        []TypeWrongInfo      `json:"by_type"`
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

// ==================== 错题标签 ====================

// WrongTagInfo 错题标签信息
type WrongTagInfo struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	Icon      string    `json:"icon"`
	Sort      int       `json:"sort"`
	Count     int       `json:"count"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateWrongTagRequest 创建错题标签请求
type CreateWrongTagRequest struct {
	Name  string `json:"name" binding:"required,max=50"`
	Color string `json:"color" binding:"omitempty,max=20"`
	Icon  string `json:"icon" binding:"omitempty,max=50"`
	Sort  int    `json:"sort"`
}

// UpdateWrongTagRequest 更新错题标签请求
type UpdateWrongTagRequest struct {
	Name  string `json:"name" binding:"omitempty,max=50"`
	Color string `json:"color" binding:"omitempty,max=20"`
	Icon  string `json:"icon" binding:"omitempty,max=50"`
	Sort  int    `json:"sort"`
}

// WrongTagIDsRequest 错题标签ID列表请求
type WrongTagIDsRequest struct {
	TagIDs []uint `json:"tag_ids" binding:"required,min=1"`
}

// ==================== 错题备注 ====================

// WrongNoteInfo 错题备注信息
type WrongNoteInfo struct {
	ID              uint      `json:"id"`
	WrongQuestionID uint      `json:"wrong_question_id"`
	UserID          uint      `json:"user_id"`
	Content         string    `json:"content"`
	IsPinned        bool      `json:"is_pinned"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// CreateWrongNoteRequest 创建错题备注请求
type CreateWrongNoteRequest struct {
	Content string `json:"content" binding:"required"`
}

// UpdateWrongNoteRequest 更新错题备注请求
type UpdateWrongNoteRequest struct {
	Content  string `json:"content" binding:"omitempty"`
	IsPinned *bool  `json:"is_pinned"`
}

// ==================== 错题附件 ====================

// WrongAttachmentInfo 错题附件信息
type WrongAttachmentInfo struct {
	ID              uint      `json:"id"`
	WrongQuestionID uint      `json:"wrong_question_id"`
	UserID          uint      `json:"user_id"`
	FileID          uint      `json:"file_id"`
	FileName        string    `json:"file_name"`
	FileType        string    `json:"file_type"`
	FileURL         string    `json:"file_url"`
	FileSize        int64     `json:"file_size"`
	Sort            int       `json:"sort"`
	CreatedAt       time.Time `json:"created_at"`
}

// ==================== 错题收藏 ====================

// WrongFavoriteInfo 错题收藏信息
type WrongFavoriteInfo struct {
	ID              uint               `json:"id"`
	UserID          uint               `json:"user_id"`
	WrongQuestionID uint               `json:"wrong_question_id"`
	FolderID        uint               `json:"folder_id"`
	Note            string             `json:"note"`
	WrongQuestion   *WrongQuestionInfo `json:"wrong_question,omitempty"`
	CreatedAt       time.Time          `json:"created_at"`
}

// FavoriteWrongRequest 收藏错题请求
type FavoriteWrongRequest struct {
	FolderID uint   `json:"folder_id"`
	Note     string `json:"note" binding:"omitempty,max=500"`
}

// WrongFavoriteListRequest 错题收藏列表请求
type WrongFavoriteListRequest struct {
	PageRequest
	FolderID *uint `form:"folder_id"`
}

// ==================== 错题复习 ====================

// WrongReviewInfo 错题复习记录信息
type WrongReviewInfo struct {
	ID              uint       `json:"id"`
	WrongQuestionID uint       `json:"wrong_question_id"`
	UserID          uint       `json:"user_id"`
	IsCorrect       bool       `json:"is_correct"`
	Answer          string     `json:"answer"`
	Duration        int        `json:"duration"`
	ReviewType      int8       `json:"review_type"`
	NextReviewAt    *time.Time `json:"next_review_at"`
	CreatedAt       time.Time  `json:"created_at"`
}

// WrongReviewListRequest 错题复习历史请求
type WrongReviewListRequest struct {
	PageRequest
	WrongQuestionID *uint `form:"wrong_question_id"`
	IsCorrect       *bool `form:"is_correct"`
}

// ==================== 错题导出 ====================

// WrongExportRequest 错题导出请求
type WrongExportRequest struct {
	IDs       []uint `json:"ids"`
	Format    string `json:"format" binding:"required,oneof=csv excel"` // csv, excel
	Keyword   string `json:"keyword"`
	CategoryID *uint `json:"category_id"`
	Difficulty *int8 `json:"difficulty"`
	Mastered   *bool `json:"mastered"`
}

// ==================== 错题分析 ====================

// WrongTrendRequest 错题趋势请求
type WrongTrendRequest struct {
	StartDate string `form:"start_date" binding:"required"`
	EndDate   string `form:"end_date" binding:"required"`
	Interval  string `form:"interval"` // day, week, month
}

// WrongTrendInfo 错题趋势信息
type WrongTrendInfo struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

// WrongCategoryAnalysisInfo 错题分类分析信息
type WrongCategoryAnalysisInfo struct {
	CategoryID   uint    `json:"category_id"`
	CategoryName string  `json:"category_name"`
	TotalCount   int     `json:"total_count"`
	MasteredCount int    `json:"mastered_count"`
	AvgWrongCount float64 `json:"avg_wrong_count"`
}

// WrongAccuracyInfo 正确率分析信息
type WrongAccuracyInfo struct {
	TotalReviews   int     `json:"total_reviews"`
	CorrectReviews int     `json:"correct_reviews"`
	AccuracyRate   float64 `json:"accuracy_rate"`
	ByDifficulty   []DifficultyAccuracyInfo `json:"by_difficulty"`
	ByType         []TypeAccuracyInfo       `json:"by_type"`
}

// DifficultyAccuracyInfo 按难度正确率
type DifficultyAccuracyInfo struct {
	Difficulty   int8    `json:"difficulty"`
	TotalCount   int     `json:"total_count"`
	CorrectCount int     `json:"correct_count"`
	AccuracyRate float64 `json:"accuracy_rate"`
}

// TypeAccuracyInfo 按题型正确率
type TypeAccuracyInfo struct {
	Type         int8    `json:"type"`
	TotalCount   int     `json:"total_count"`
	CorrectCount int     `json:"correct_count"`
	AccuracyRate float64 `json:"accuracy_rate"`
}

// WrongBatchDeleteRequest 错题批量删除请求
type WrongBatchDeleteRequest struct {
	IDs []uint `json:"ids" binding:"required,min=1"`
}
