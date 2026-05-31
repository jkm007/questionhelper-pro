package model

import (
	"time"

	"gorm.io/gorm"
)

// PaperShare 试卷共享表
type PaperShare struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	PaperID   uint           `gorm:"index;not null" json:"paper_id"`           // 试卷ID
	SharerID  uint           `gorm:"index;not null" json:"sharer_id"`          // 共享人ID
	TargetID  uint           `gorm:"index;not null" json:"target_id"`          // 目标用户/班级ID
	TargetType int8          `gorm:"default:1;comment:目标类型:1=用户,2=班级" json:"target_type"`
	Permission int8         `gorm:"default:1;comment:权限:1=只读,2=可编辑" json:"permission"`
	Status    int8           `gorm:"default:1;comment:状态:0=已撤销,1=有效" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (PaperShare) TableName() string {
	return "paper_shares"
}

// PaperFavorite 试卷收藏表
type PaperFavorite struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"uniqueIndex:idx_user_paper;not null" json:"user_id"`   // 用户ID
	PaperID   uint      `gorm:"uniqueIndex:idx_user_paper;not null" json:"paper_id"`  // 试卷ID
	Note      string    `gorm:"size:500" json:"note"`                                  // 收藏备注
	CreatedAt time.Time `json:"created_at"`
}

func (PaperFavorite) TableName() string {
	return "paper_favorites"
}

// ExamReminder 考试提醒表
type ExamReminder struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	ExamID     uint           `gorm:"index;not null" json:"exam_id"`                    // 考试ID
	UserID     uint           `gorm:"index;not null" json:"user_id"`                    // 用户ID
	RemindType int8           `gorm:"default:1;comment:提醒类型:1=开始提醒,2=结束提醒,3=自定义" json:"remind_type"`
	RemindTime time.Time      `gorm:"not null" json:"remind_time"`                      // 提醒时间
	Content    string         `gorm:"size:500" json:"content"`                           // 提醒内容
	IsSent     bool           `gorm:"default:false" json:"is_sent"`                      // 是否已发送
	SentAt     *time.Time     `json:"sent_at"`                                           // 发送时间
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ExamReminder) TableName() string {
	return "exam_reminders"
}

// ExamExtension 考试延期记录表
type ExamExtension struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	ExamID      uint           `gorm:"index;not null" json:"exam_id"`                  // 考试ID
	OperatorID  uint           `gorm:"index;not null" json:"operator_id"`              // 操作人ID
	OldEndTime  time.Time      `gorm:"not null" json:"old_end_time"`                   // 原结束时间
	NewEndTime  time.Time      `gorm:"not null" json:"new_end_time"`                   // 新结束时间
	ExtendMinutes int          `gorm:"not null;comment:延期时长(分钟)" json:"extend_minutes"`
	Reason      string         `gorm:"size:500;not null" json:"reason"`                // 延期原因
	Status      int8           `gorm:"default:1;comment:状态:0=已撤销,1=有效" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ExamExtension) TableName() string {
	return "exam_extensions"
}

// ExamPause 考试暂停记录表
type ExamPause struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	ExamID     uint           `gorm:"index;not null" json:"exam_id"`                   // 考试ID
	OperatorID uint           `gorm:"index;not null" json:"operator_id"`               // 操作人ID
	Action     int8           `gorm:"not null;comment:操作:1=暂停,2=恢复" json:"action"`
	Reason     string         `gorm:"size:500" json:"reason"`                          // 操作原因
	PausedAt   *time.Time     `json:"paused_at"`                                       // 暂停时间
	ResumedAt  *time.Time     `json:"resumed_at"`                                      // 恢复时间
	Duration   int            `gorm:"default:0;comment:暂停时长(秒)" json:"duration"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ExamPause) TableName() string {
	return "exam_pauses"
}

// ScoreReview 成绩复核表
type ScoreReview struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	RecordID    uint           `gorm:"index;not null" json:"record_id"`                // 考试记录ID
	UserID      uint           `gorm:"index;not null" json:"user_id"`                  // 申请人ID
	Reason      string         `gorm:"size:1000;not null" json:"reason"`               // 复核原因
	OldScore    float64        `json:"old_score"`                                       // 原始分数
	NewScore    float64        `json:"new_score"`                                       // 复核后分数
	ReviewerID  *uint          `gorm:"index" json:"reviewer_id"`                       // 复核人ID
	ReviewNote  string         `gorm:"size:1000" json:"review_note"`                   // 复核备注
	Status      int8           `gorm:"default:0;comment:状态:0=待复核,1=已复核,2=已驳回" json:"status"`
	ReviewedAt  *time.Time     `json:"reviewed_at"`                                     // 复核时间
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ScoreReview) TableName() string {
	return "score_reviews"
}

// ExamNotice 考试公告表
type ExamNotice struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	ExamID    uint           `gorm:"index;not null" json:"exam_id"`                   // 考试ID
	Title     string         `gorm:"size:200;not null" json:"title"`                  // 公告标题
	Content   string         `gorm:"type:text;not null" json:"content"`               // 公告内容
	Priority  int8           `gorm:"default:0;comment:优先级:0=普通,1=重要,2=紧急" json:"priority"`
	IsPinned  bool           `gorm:"default:false" json:"is_pinned"`                  // 是否置顶
	CreatorID uint           `gorm:"index;not null" json:"creator_id"`                // 创建人ID
	Status    int8           `gorm:"default:1;comment:状态:0=草稿,1=已发布" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ExamNotice) TableName() string {
	return "exam_notices"
}

// ExamRanking 成绩排名表
type ExamRanking struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	ExamID       uint      `gorm:"uniqueIndex:idx_exam_user;not null" json:"exam_id"`   // 考试ID
	UserID       uint      `gorm:"uniqueIndex:idx_exam_user;not null" json:"user_id"`   // 用户ID
	Score        float64   `gorm:"not null" json:"score"`                                // 总分
	ObjScore     float64   `json:"obj_score"`                                            // 客观题得分
	SubjScore    float64   `json:"subj_score"`                                           // 主观题得分
	RankPos      int       `gorm:"not null;comment:排名" json:"rank_pos"`
	DurationUsed int       `gorm:"default:0;comment:用时(秒)" json:"duration_used"`
	Accuracy     float64   `gorm:"default:0;comment:正确率" json:"accuracy"`
	SubmitTime   time.Time `json:"submit_time"`                                           // 交卷时间
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (ExamRanking) TableName() string {
	return "exam_rankings"
}

// ExamFeedback 考试反馈表
type ExamFeedback struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	ExamID      uint           `gorm:"index;not null" json:"exam_id"`                   // 考试ID
	UserID      uint           `gorm:"index;not null" json:"user_id"`                   // 用户ID
	FeedbackType int8          `gorm:"default:1;comment:反馈类型:1=题目问题,2=系统问题,3=其他" json:"feedback_type"`
	Content     string         `gorm:"type:text;not null" json:"content"`               // 反馈内容
	Images      string         `gorm:"type:text" json:"images"`                          // 图片附件(JSON数组)
	Reply       string         `gorm:"type:text" json:"reply"`                           // 回复内容
	ReplyerID   *uint          `gorm:"index" json:"replyer_id"`                          // 回复人ID
	ReplyAt     *time.Time     `json:"reply_at"`                                         // 回复时间
	Status      int8           `gorm:"default:0;comment:状态:0=待处理,1=已回复,2=已关闭" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ExamFeedback) TableName() string {
	return "exam_feedbacks"
}

// CommentTemplate 评语模板表
type CommentTemplate struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Name      string         `gorm:"size:100;not null" json:"name"`                    // 模板名称
	Content   string         `gorm:"type:text;not null" json:"content"`                // 评语内容
	Category  string         `gorm:"size:50;index" json:"category"`                    // 分类
	Sort      int            `gorm:"default:0" json:"sort"`                            // 排序
	UsageCount int           `gorm:"default:0" json:"usage_count"`                    // 使用次数
	CreatorID uint           `gorm:"index;not null" json:"creator_id"`                 // 创建人ID
	Status    int8           `gorm:"default:1;comment:状态:0=禁用,1=启用" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (CommentTemplate) TableName() string {
	return "comment_templates"
}
