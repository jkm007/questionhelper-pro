package model

import (
	"time"

	"gorm.io/gorm"
)

// MockExamConfig 模拟考试配置表
type MockExamConfig struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	Name          string         `gorm:"size:200;not null" json:"name"`                // 配置名称
	Description   string         `gorm:"type:text" json:"description"`                 // 描述
	PaperID       *uint          `gorm:"index" json:"paper_id"`                        // 关联试卷ID(指定试卷模式)
	CategoryID    *uint          `gorm:"index" json:"category_id"`                     // 分类ID(随机抽题模式)
	QuestionCount int            `gorm:"default:0" json:"question_count"`              // 题目数量
	TotalScore    float64        `gorm:"default:0" json:"total_score"`                 // 总分
	PassScore     float64        `gorm:"default:0" json:"pass_score"`                  // 及格分
	Duration      int            `gorm:"not null;comment:考试时长(分钟)" json:"duration"` // 考试时长
	Shuffle       bool           `gorm:"default:false" json:"shuffle"`                 // 是否打乱题目顺序
	ShowAnswer    int8           `gorm:"default:1;comment:答案显示:0=不显示,1=交卷后,2=查看记录时" json:"show_answer"`
	DifficultyMix string         `gorm:"size:100;comment:难度配比(JSON,如 {1:30,2:50,3:20})" json:"difficulty_mix"`
	Mode          int8           `gorm:"default:1;comment:模式:1=指定试卷,2=随机抽题,3=错题重练,4=收藏题练习" json:"mode"`
	MaxAttempts   int            `gorm:"default:0;comment:最大练习次数,0=不限" json:"max_attempts"`
	Status        int8           `gorm:"default:1;comment:状态:0=禁用,1=启用" json:"status"`
	CreatorID     uint           `gorm:"index" json:"creator_id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (MockExamConfig) TableName() string {
	return "mock_exam_configs"
}

// PracticePlan 练习计划表
type PracticePlan struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	UserID        uint           `gorm:"index;not null" json:"user_id"`
	Name          string         `gorm:"size:200;not null" json:"name"`                  // 计划名称
	Description   string         `gorm:"type:text" json:"description"`                   // 计划描述
	PlanType      int8           `gorm:"default:1;comment:类型:1=每日计划,2=每周计划,3=自定义" json:"plan_type"`
	CategoryID    *uint          `gorm:"index" json:"category_id"`                       // 目标分类
	QuestionType  *int8          `gorm:"comment:题型:1=单选,2=多选,3=判断,4=填空,5=简答" json:"question_type"`
	Difficulty    *int8          `gorm:"comment:难度:1=简单,2=中等,3=困难" json:"difficulty"`
	DailyCount    int            `gorm:"default:0;comment:每日目标题目数" json:"daily_count"`
	DailyDuration int            `gorm:"default:0;comment:每日目标时长(分钟)" json:"daily_duration"`
	StartDate     time.Time      `gorm:"not null" json:"start_date"`
	EndDate       *time.Time     `json:"end_date"`
	TotalTarget   int            `gorm:"default:0;comment:总目标题目数" json:"total_target"`
	TotalDone     int            `gorm:"default:0;comment:已完成题目数" json:"total_done"`
	Progress      float64        `gorm:"default:0;comment:完成进度(百分比)" json:"progress"`
	Status        int8           `gorm:"default:1;comment:状态:0=暂停,1=进行中,2=已完成,3=已过期" json:"status"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (PracticePlan) TableName() string {
	return "practice_plans"
}

// PracticePlanLog 练习计划执行记录表
type PracticePlanLog struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	PlanID      uint      `gorm:"index;not null" json:"plan_id"`
	UserID      uint      `gorm:"index;not null" json:"user_id"`
	SessionID   *uint     `gorm:"index" json:"session_id"`                  // 关联练习会话
	QuestionNum int       `gorm:"default:0;comment:本次完成题目数" json:"question_num"`
	Duration    int       `gorm:"default:0;comment:本次练习时长(秒)" json:"duration"`
	Accuracy    float64   `gorm:"default:0;comment:本次正确率" json:"accuracy"`
	Date        string    `gorm:"size:10;index;not null;comment:记录日期(YYYY-MM-DD)" json:"date"`
	CreatedAt   time.Time `json:"created_at"`
}

func (PracticePlanLog) TableName() string {
	return "practice_plan_logs"
}

// DailyPractice 每日练习记录表
type DailyPractice struct {
	ID            uint      `gorm:"primarykey" json:"id"`
	UserID        uint      `gorm:"uniqueIndex:idx_user_date;not null" json:"user_id"`
	Date          string    `gorm:"uniqueIndex:idx_user_date;size:10;not null;comment:日期(YYYY-MM-DD)" json:"date"`
	TotalCount    int       `gorm:"default:0;comment:总练习题数" json:"total_count"`
	CorrectCount  int       `gorm:"default:0;comment:正确题数" json:"correct_count"`
	Accuracy      float64   `gorm:"default:0;comment:正确率" json:"accuracy"`
	Duration      int       `gorm:"default:0;comment:总练习时长(秒)" json:"duration"`
	SessionCount  int       `gorm:"default:0;comment:练习会话数" json:"session_count"`
	WrongCount    int       `gorm:"default:0;comment:错题数" json:"wrong_count"`
	NewQuestion   int       `gorm:"default:0;comment:新题数" json:"new_question"`
	CategoryStats string    `gorm:"type:text;comment:分类练习统计(JSON)" json:"category_stats"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (DailyPractice) TableName() string {
	return "daily_practices"
}

// PracticeCheckin 练习打卡表
type PracticeCheckin struct {
	ID            uint      `gorm:"primarykey" json:"id"`
	UserID        uint      `gorm:"uniqueIndex:idx_user_checkin_date;not null" json:"user_id"`
	Date          string    `gorm:"uniqueIndex:idx_user_checkin_date;size:10;not null;comment:打卡日期(YYYY-MM-DD)" json:"date"`
	IsCheckin     bool      `gorm:"default:true" json:"is_checkin"`
	QuestionCount int       `gorm:"default:0;comment:打卡时练习题数" json:"question_count"`
	Duration      int       `gorm:"default:0;comment:打卡时练习时长(秒)" json:"duration"`
	Streak        int       `gorm:"default:0;comment:连续打卡天数" json:"streak"`
	Reward        int       `gorm:"default:0;comment:奖励积分" json:"reward"`
	CreatedAt     time.Time `json:"created_at"`
}

func (PracticeCheckin) TableName() string {
	return "practice_checkins"
}

// UserAbilityProfile 用户能力模型表
type UserAbilityProfile struct {
	ID                uint      `gorm:"primarykey" json:"id"`
	UserID            uint      `gorm:"uniqueIndex;not null" json:"user_id"`
	TotalQuestions    int       `gorm:"default:0;comment:总练习题数" json:"total_questions"`
	TotalCorrect      int       `gorm:"default:0;comment:总正确题数" json:"total_correct"`
	TotalDuration     int       `gorm:"default:0;comment:总练习时长(秒)" json:"total_duration"`
	OverallAccuracy   float64   `gorm:"default:0;comment:总正确率" json:"overall_accuracy"`
	ConsecutiveDays   int       `gorm:"default:0;comment:当前连续练习天数" json:"consecutive_days"`
	MaxConsecutive    int       `gorm:"default:0;comment:最大连续练习天数" json:"max_consecutive"`
	TotalCheckins     int       `gorm:"default:0;comment:总打卡天数" json:"total_checkins"`
	Level             int       `gorm:"default:1;comment:用户等级" json:"level"`
	Experience        int       `gorm:"default:0;comment:经验值" json:"experience"`
	CategoryAbilities string    `gorm:"type:text;comment:分类能力值(JSON)" json:"category_abilities"`
	WeakPoints        string    `gorm:"type:text;comment:薄弱知识点(JSON)" json:"weak_points"`
	StrongPoints      string    `gorm:"type:text;comment:擅长知识点(JSON)" json:"strong_points"`
	RankScore         float64   `gorm:"default:0;comment:排行分数" json:"rank_score"`
	LastPracticeAt    *time.Time `json:"last_practice_at"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

func (UserAbilityProfile) TableName() string {
	return "user_ability_profiles"
}

// PracticeLeaderboard 练习排行榜快照表
type PracticeLeaderboard struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	UserID     uint      `gorm:"index;not null" json:"user_id"`
	Nickname   string    `gorm:"size:50" json:"nickname"`
	Avatar     string    `gorm:"size:255" json:"avatar"`
	RankType   int8      `gorm:"index;not null;comment:排行类型:1=日榜,2=周榜,3=月榜,4=总榜" json:"rank_type"`
	RankDate   string    `gorm:"size:10;index;not null;comment:排行日期(YYYY-MM-DD)" json:"rank_date"`
	RankPos    int       `gorm:"default:0;comment:排名" json:"rank_pos"`
	Score      float64   `gorm:"default:0;comment:排行分数" json:"score"`
	Accuracy   float64   `gorm:"default:0;comment:正确率" json:"accuracy"`
	TotalCount int       `gorm:"default:0;comment:练习题数" json:"total_count"`
	Duration   int       `gorm:"default:0;comment:练习时长(秒)" json:"duration"`
	CreatedAt  time.Time `json:"created_at"`
}

func (PracticeLeaderboard) TableName() string {
	return "practice_leaderboards"
}

// ChallengeLevel 闯关配置表
type ChallengeLevel struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	Name          string         `gorm:"size:200;not null" json:"name"`                // 关卡名称
	Description   string         `gorm:"type:text" json:"description"`                 // 关卡描述
	Level         int            `gorm:"index;not null;comment:关卡序号" json:"level"`   // 关卡序号
	CategoryID    *uint          `gorm:"index" json:"category_id"`                     // 关联分类
	QuestionCount int            `gorm:"not null;comment:关卡题目数" json:"question_count"` // 题目数量
	PassAccuracy  float64        `gorm:"not null;comment:通关正确率" json:"pass_accuracy"` // 通关所需正确率
	PassScore     int            `gorm:"default:0;comment:通关奖励积分" json:"pass_score"`
	TimeLimit     int            `gorm:"default:0;comment:时间限制(秒),0=不限" json:"time_limit"`
	Difficulty    int8           `gorm:"default:1;comment:难度:1=简单,2=中等,3=困难" json:"difficulty"`
	Icon          string         `gorm:"size:255" json:"icon"`
	Badge         string         `gorm:"size:255;comment:通关徽章图" json:"badge"`
	PreLevel      int            `gorm:"default:0;comment:前置关卡序号,0=无" json:"pre_level"`
	Status        int8           `gorm:"default:1;comment:状态:0=禁用,1=启用" json:"status"`
	Sort          int            `gorm:"default:0" json:"sort"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ChallengeLevel) TableName() string {
	return "challenge_levels"
}

// UserChallengeProgress 用户闯关进度表
type UserChallengeProgress struct {
	ID             uint      `gorm:"primarykey" json:"id"`
	UserID         uint      `gorm:"uniqueIndex:idx_user_challenge;not null" json:"user_id"`
	ChallengeLevel uint      `gorm:"uniqueIndex:idx_user_challenge;not null;comment:关卡ID" json:"challenge_level"`
	Status         int8      `gorm:"default:0;comment:状态:0=未开始,1=进行中,2=已通关,3=失败" json:"status"`
	BestAccuracy   float64   `gorm:"default:0;comment:最佳正确率" json:"best_accuracy"`
	BestDuration   int       `gorm:"default:0;comment:最佳用时(秒)" json:"best_duration"`
	Attempts       int       `gorm:"default:0;comment:尝试次数" json:"attempts"`
	PassedAt       *time.Time `json:"passed_at"`
	Score          int       `gorm:"default:0;comment:获得积分" json:"score"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

func (UserChallengeProgress) TableName() string {
	return "user_challenge_progress"
}
