package model

import (
	"time"

	"gorm.io/gorm"
)

type Exam struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	Title         string         `gorm:"size:200;not null" json:"title"`
	Description   string         `gorm:"type:text" json:"description"`
	PaperID       uint           `gorm:"index" json:"paper_id"`
	Paper         Paper          `json:"paper,omitempty"`
	ClassID       *uint          `gorm:"index" json:"class_id"`
	StartTime     time.Time      `gorm:"not null" json:"start_time"`
	EndTime       time.Time      `gorm:"not null" json:"end_time"`
	Duration      int            `gorm:"not null;comment:考试时长(分钟)" json:"duration"`
	TotalScore    float64        `gorm:"not null" json:"total_score"`
	PassScore     float64        `gorm:"not null" json:"pass_score"`
	MaxAttempts   int            `gorm:"default:1" json:"max_attempts"`
	Shuffle       bool           `gorm:"default:false;comment:是否打乱顺序" json:"shuffle"`
	ShowAnswer    int8           `gorm:"default:0;comment:是否显示答案:0=不显示,1=交卷后,2=考试后" json:"show_answer"`
	AntiCheat      int8           `gorm:"default:0;comment:防作弊级别:0=无,1=基础,2=严格" json:"anti_cheat"`
	Status         int8           `gorm:"default:0;comment:状态:0=未发布,1=进行中,2=已结束" json:"status"`
	ExamPassword   string         `gorm:"size:100;comment:考试密码" json:"exam_password"`
	StatusPause    bool           `gorm:"default:false;comment:是否暂停" json:"status_pause"`
	ShowRanking    bool           `gorm:"default:false;comment:是否显示排名" json:"show_ranking"`
	RankingType    int            `gorm:"default:1;comment:排名方式:1=按分数,2=按用时,3=按正确率" json:"ranking_type"`
	OriginalEndTime *time.Time    `json:"original_end_time"`                                // 原始结束时间
	ExtendReason   string         `gorm:"size:500;comment:延期原因" json:"extend_reason"`
	RemindSent     bool           `gorm:"default:false;comment:开考提醒是否已发送" json:"remind_sent"`
	EndRemindSent  bool           `gorm:"default:false;comment:结束提醒是否已发送" json:"end_remind_sent"`
	CreatorID      uint           `gorm:"index" json:"creator_id"`
	Creator       User           `json:"creator,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Exam) TableName() string {
	return "exams"
}

type Paper struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	Title        string         `gorm:"size:200;not null" json:"title"`
	Description  string         `gorm:"type:text" json:"description"`
	TotalScore   float64        `gorm:"not null" json:"total_score"`
	TotalCount   int            `gorm:"not null" json:"total_count"`
	Type         int8           `gorm:"default:1;comment:类型:1=手动组卷,2=智能组卷" json:"type"`
	Status       int8           `gorm:"default:0;comment:状态:0=草稿,1=已发布,2=已归档" json:"status"`
	Visibility   int8           `gorm:"default:1;comment:可见性:1=公开,2=私有,3=班级" json:"visibility"`
	IsTemplate   bool           `gorm:"default:false" json:"is_template"`
	Category     string         `gorm:"size:50;index" json:"category"`
	Tags         string         `gorm:"size:500" json:"tags"`
	UsageCount   int            `gorm:"default:0" json:"usage_count"`
	AvgScore     float64        `gorm:"default:0" json:"avg_score"`
	PassRate     float64        `gorm:"default:0" json:"pass_rate"`
	ClassID      *uint          `gorm:"index" json:"class_id,omitempty"`
	CreatorID    uint           `gorm:"index" json:"creator_id"`
	Creator      User           `json:"creator,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Paper) TableName() string {
	return "papers"
}

type PaperQuestion struct {
	ID         uint    `gorm:"primarykey" json:"id"`
	PaperID    uint    `gorm:"index;not null" json:"paper_id"`
	QuestionID uint    `gorm:"index;not null" json:"question_id"`
	Score      float64 `gorm:"not null" json:"score"`
	Sort       int     `gorm:"default:0" json:"sort"`
	Snapshot   string  `gorm:"type:text" json:"snapshot"` // 题目快照(JSON)
}

func (PaperQuestion) TableName() string {
	return "paper_questions"
}

type ExamRecord struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	ExamID       uint           `gorm:"index;not null" json:"exam_id"`
	UserID       uint           `gorm:"index;not null" json:"user_id"`
	Score        float64        `json:"score"`
	ObjScore     float64        `json:"obj_score"`                              // 客观题得分
	SubjScore    float64        `json:"subj_score"`                             // 主观题得分
	Status       int8           `gorm:"default:0;comment:状态:0=进行中,1=已交卷,2=已阅卷" json:"status"`
	SubmitType   int8           `gorm:"default:1;comment:交卷方式:1=手动,2=自动,3=强制" json:"submit_type"`
	StartTime    time.Time      `json:"start_time"`
	SubmitTime   *time.Time     `json:"submit_time"`
	DurationUsed int            `gorm:"default:0;comment:实际用时(秒)" json:"duration_used"`
	IP           string         `gorm:"size:50" json:"ip"`
	SwitchCount    int            `gorm:"default:0;comment:切屏次数" json:"switch_count"`
	IPChanges      int            `gorm:"default:0;comment:IP变更次数" json:"ip_changes"`
	OverallComment  string         `gorm:"type:text;comment:总体评语" json:"overall_comment"`
	ScoreWarningSent bool          `gorm:"default:false;comment:成绩预警是否已发送" json:"score_warning_sent"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ExamRecord) TableName() string {
	return "exam_records"
}

type AnswerRecord struct {
	ID           uint    `gorm:"primarykey" json:"id"`
	RecordID     uint    `gorm:"index;not null" json:"record_id"`
	QuestionID   uint    `gorm:"index;not null" json:"question_id"`
	Answer       string  `gorm:"type:text" json:"answer"`
	Score        float64 `json:"score"`
	MaxScore     float64 `json:"max_score"`                                // 该题满分
	IsCorrect    bool    `json:"is_correct"`
	IsMarked     bool    `gorm:"default:false" json:"is_marked"`           // 是否标记
	IsReviewed   bool    `gorm:"default:false" json:"is_reviewed"`         // 是否已阅卷
	ReviewNote   string  `gorm:"size:500" json:"review_note"`              // 阅卷备注
	Comment      string  `gorm:"type:text;comment:评语/批注" json:"comment"` // 评语/批注
	ReviewedBy   *uint   `json:"reviewed_by"`
	ReviewedAt   *time.Time `json:"reviewed_at"`
}

func (AnswerRecord) TableName() string {
	return "answer_records"
}
