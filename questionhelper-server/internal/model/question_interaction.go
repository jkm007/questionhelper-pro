package model

import (
	"time"

	"gorm.io/gorm"
)

// QuestionNote 题目笔记表
type QuestionNote struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	UserID     uint           `gorm:"uniqueIndex:idx_user_question;not null" json:"user_id"`
	QuestionID uint           `gorm:"uniqueIndex:idx_user_question;not null" json:"question_id"`
	Content    string         `gorm:"type:text;not null" json:"content"`                          // 笔记内容
	IsPublic   bool           `gorm:"default:false;comment:是否公开" json:"is_public"`             // 是否对其他用户可见
	LikeCount  int            `gorm:"default:0" json:"like_count"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (QuestionNote) TableName() string {
	return "question_notes"
}

// QuestionDifficultyRating 题目难度评价表
type QuestionDifficultyRating struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	UserID     uint      `gorm:"uniqueIndex:idx_user_difficulty;not null" json:"user_id"`
	QuestionID uint      `gorm:"uniqueIndex:idx_user_difficulty;not null" json:"question_id"`
	Rating     int8      `gorm:"not null;comment:难度评分:1=很简单,2=简单,3=一般,4=困难,5=很困难" json:"rating"`
	IsCorrect  bool      `gorm:"comment:评价时是否答对" json:"is_correct"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (QuestionDifficultyRating) TableName() string {
	return "question_difficulty_ratings"
}

// QuestionQualityRating 题目质量评分表
type QuestionQualityRating struct {
	ID             uint      `gorm:"primarykey" json:"id"`
	UserID         uint      `gorm:"uniqueIndex:idx_user_quality;not null" json:"user_id"`
	QuestionID     uint      `gorm:"uniqueIndex:idx_user_quality;not null" json:"question_id"`
	Score          int8      `gorm:"not null;comment:质量评分:1-5星" json:"score"`
	ClarityScore   int8      `gorm:"default:0;comment:题目清晰度评分:1-5" json:"clarity_score"`
	RelevanceScore int8      `gorm:"default:0;comment:内容相关性评分:1-5" json:"relevance_score"`
	Comment        string    `gorm:"size:500;comment:评价内容" json:"comment"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (QuestionQualityRating) TableName() string {
	return "question_quality_ratings"
}

// QuestionCorrection 题目纠错表
type QuestionCorrection struct {
	ID           uint       `gorm:"primarykey" json:"id"`
	UserID       uint       `gorm:"index;not null" json:"user_id"`
	QuestionID   uint       `gorm:"index;not null" json:"question_id"`
	CorrectionType int8     `gorm:"not null;comment:纠错类型:1=答案错误,2=解析错误,3=选项错误,4=题目不完整,5=内容重复,6=其他" json:"correction_type"`
	Description  string     `gorm:"type:text;not null;comment:纠错描述" json:"description"`
	Screenshot   string     `gorm:"size:500;comment:截图URL" json:"screenshot"`
	Status       int8       `gorm:"default:0;comment:状态:0=待处理,1=已采纳,2=已驳回,3=已修复" json:"status"`
	AdminReply   string     `gorm:"size:500;comment:管理员回复" json:"admin_reply"`
	ProcessedBy  *uint      `json:"processed_by"`
	ProcessedAt  *time.Time `json:"processed_at"`
	Reward       int        `gorm:"default:0;comment:奖励积分" json:"reward"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (QuestionCorrection) TableName() string {
	return "question_corrections"
}
