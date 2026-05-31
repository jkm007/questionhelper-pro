package dto

import "time"

// ==================== 试卷共享/收藏 ====================

// SharePaperRequest 试卷共享请求
type SharePaperRequest struct {
	TargetID   uint  `json:"target_id" binding:"required"`
	TargetType int8  `json:"target_type" binding:"required,oneof=1 2"` // 1=用户,2=班级
	Permission int8  `json:"permission" binding:"required,oneof=1 2"`  // 1=只读,2=可编辑
}

// FavoritePaperRequest 收藏试卷请求
type FavoritePaperRequest struct {
	Note string `json:"note" binding:"omitempty,max=500"`
}

// ImportPaperFromFileRequest 从文件导入试卷请求
type ImportPaperFromFileRequest struct {
	Title       string `json:"title" binding:"required,max=200"`
	Description string `json:"description"`
	CategoryID  uint   `json:"category_id"`
	Visibility  int8   `json:"visibility" binding:"omitempty,oneof=1 2 3"`
}

// ==================== 考试操作 ====================

// ExtendExamRequest 延长考试请求
type ExtendExamRequest struct {
	Minutes int    `json:"minutes" binding:"required,min=1"`
	Reason  string `json:"reason" binding:"required,max=500"`
}

// PauseExamRequest 暂停考试请求
type PauseExamRequest struct {
	Reason string `json:"reason" binding:"omitempty,max=500"`
}

// ResumeExamRequest 恢复考试请求
type ResumeExamRequest struct {
	Reason string `json:"reason" binding:"omitempty,max=500"`
}

// SwitchScreenRequest 切屏上报请求
type SwitchScreenRequest struct {
	Detail    string `json:"detail"`
	UserAgent string `json:"user_agent"`
}

// ==================== 成绩复核 ====================

// ScoreReviewRequest 申请成绩复核请求
type ScoreReviewRequest struct {
	Reason string `json:"reason" binding:"required,max=1000"`
}

// ScoreReviewListRequest 复核列表请求
type ScoreReviewListRequest struct {
	PageRequest
	Status  *int8  `form:"status"`
	ExamID  *uint  `form:"exam_id"`
	Keyword string `form:"keyword"`
}

// HandleScoreReviewRequest 处理复核请求
type HandleScoreReviewRequest struct {
	Status     int8   `json:"status" binding:"required,oneof=1 2"` // 1=已复核,2=已驳回
	NewScore   float64 `json:"new_score"`
	ReviewNote string `json:"review_note" binding:"required,max=1000"`
}

// ScoreReviewInfo 复核信息
type ScoreReviewInfo struct {
	ID         uint       `json:"id"`
	RecordID   uint       `json:"record_id"`
	UserID     uint       `json:"user_id"`
	Username   string     `json:"username"`
	Nickname   string     `json:"nickname"`
	ExamID     uint       `json:"exam_id"`
	ExamTitle  string     `json:"exam_title"`
	Reason     string     `json:"reason"`
	OldScore   float64    `json:"old_score"`
	NewScore   float64    `json:"new_score"`
	ReviewerID *uint      `json:"reviewer_id"`
	ReviewNote string     `json:"review_note"`
	Status     int8       `json:"status"`
	ReviewedAt *time.Time `json:"reviewed_at"`
	CreatedAt  time.Time  `json:"created_at"`
}

// ==================== 考试公告 ====================

// CreateExamNoticeRequest 创建考试公告请求
type CreateExamNoticeRequest struct {
	Title    string `json:"title" binding:"required,max=200"`
	Content  string `json:"content" binding:"required"`
	Priority int8   `json:"priority" binding:"omitempty,oneof=0 1 2"`
	IsPinned bool   `json:"is_pinned"`
}

// ExamNoticeInfo 公告信息
type ExamNoticeInfo struct {
	ID        uint      `json:"id"`
	ExamID    uint      `json:"exam_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Priority  int8      `json:"priority"`
	IsPinned  bool      `json:"is_pinned"`
	CreatorID uint      `json:"creator_id"`
	Status    int8      `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// ==================== 考试查询增强 ====================

// ExamRankingInfo 排名信息
type ExamRankingInfo struct {
	ID           uint    `json:"id"`
	ExamID       uint    `json:"exam_id"`
	UserID       uint    `json:"user_id"`
	Username     string  `json:"username"`
	Nickname     string  `json:"nickname"`
	Score        float64 `json:"score"`
	ObjScore     float64 `json:"obj_score"`
	SubjScore    float64 `json:"subj_score"`
	RankPos      int     `json:"rank_pos"`
	DurationUsed int     `json:"duration_used"`
	Accuracy     float64 `json:"accuracy"`
	SubmitTime   string  `json:"submit_time"`
}

// ExamUpcomingInfo 即将开始的考试信息
type ExamUpcomingInfo struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Duration    int       `json:"duration"`
	TotalScore  float64   `json:"total_score"`
	PassScore   float64   `json:"pass_score"`
	Status      int8      `json:"status"`
	TimeLeft    int       `json:"time_left"` // 距离开始的秒数
}

// ExamStatisticsResponse 考试统计响应
type ExamStatisticsResponse struct {
	BasicInfo     ExamBasicInfo    `json:"basic_info"`
	ScoreStats    ScoreStatistics  `json:"score_stats"`
	QuestionStats []QuestionStat   `json:"question_stats"`
	TopStudents   []StudentScore   `json:"top_students"`
	WarningCount  int              `json:"warning_count"`
	OnlineCount   int              `json:"online_count"`
}
