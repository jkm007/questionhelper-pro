package dto

// ==================== 基础统计 ====================

// UserStatistics 用户统计
type UserStatistics struct {
	TotalUsers    int64 `json:"total_users"`
	ActiveUsers   int64 `json:"active_users"`
	NewUsersToday int64 `json:"new_users_today"`
	NewUsersWeek  int64 `json:"new_users_week"`
}

// PracticeStatistics 练习统计
type PracticeStatistics struct {
	TotalSessions  int64   `json:"total_sessions"`
	TotalQuestions int64   `json:"total_questions"`
	AvgAccuracy    float64 `json:"avg_accuracy"`
	TotalDuration  int64   `json:"total_duration"`
}

// ExamStatistics 考试统计
type ExamStatistics struct {
	TotalExams    int64   `json:"total_exams"`
	TotalRecords  int64   `json:"total_records"`
	AvgScore      float64 `json:"avg_score"`
	PassRate      float64 `json:"pass_rate"`
}

// ClassStatistics 班级统计
type ClassStatistics struct {
	TotalClasses   int64   `json:"total_classes"`
	TotalMembers   int64   `json:"total_members"`
	AvgMemberCount float64 `json:"avg_member_count"`
}

// RankInfo 排行榜信息
type RankInfo struct {
	Rank     int     `json:"rank"`
	UserID   uint    `json:"user_id"`
	UserName string  `json:"user_name"`
	Avatar   string  `json:"avatar"`
	Score    float64 `json:"score"`
	Count    int     `json:"count"`
}

// RankingRequest 排行榜请求
type RankingRequest struct {
	PageRequest
	Type int `form:"type" binding:"required,oneof=1 2 3"` // 1:练习 2:考试 3:积分
}

// DashboardInfo 仪表盘信息
type DashboardInfo struct {
	UserStats     UserStatistics     `json:"user_stats"`
	PracticeStats PracticeStatistics `json:"practice_stats"`
	ExamStats     ExamStatistics     `json:"exam_stats"`
	ClassStats    ClassStatistics    `json:"class_stats"`
}

// ==================== 用户留存分析 ====================

// RetentionRequest 留存统计请求
type RetentionRequest struct {
	Period    string `form:"period" binding:"required,oneof=day week month"` // 统计周期
	StartDate string `form:"start_date" binding:"required"`                  // 开始日期
	EndDate   string `form:"end_date" binding:"required"`                    // 结束日期
}

// RetentionItem 留存统计项
type RetentionItem struct {
	Date          string  `json:"date"`
	NewUsers      int     `json:"new_users"`
	RetainedUsers int     `json:"retained_users"`
	RetentionRate float64 `json:"retention_rate"`
	Period        string  `json:"period"`
}

// ==================== 用户流失分析 ====================

// ChurnRequest 流失统计请求
type ChurnRequest struct {
	Period    string `form:"period" binding:"required,oneof=day week month"` // 统计周期
	StartDate string `form:"start_date" binding:"required"`                  // 开始日期
	EndDate   string `form:"end_date" binding:"required"`                    // 结束日期
}

// ChurnItem 流失统计项
type ChurnItem struct {
	Date         string  `json:"date"`
	ChurnedUsers int     `json:"churned_users"`
	ChurnRate    float64 `json:"churn_rate"`
	ChurnReasons string  `json:"churn_reasons"`
	Period       string  `json:"period"`
}

// ==================== 用户行为事件 ====================

// CreateEventRequest 创建行为事件请求
type CreateEventRequest struct {
	EventType  string `json:"event_type" binding:"required,max=50"`  // 事件类型
	EventName  string `json:"event_name" binding:"required,max=100"` // 事件名称
	Page       string `json:"page" binding:"max=200"`                // 页面
	Element    string `json:"element" binding:"max=200"`             // 元素
	ExtraData  string `json:"extra_data"`                            // 扩展数据(JSON)
	SessionID  string `json:"session_id" binding:"max=64"`           // 会话ID
	DeviceType string `json:"device_type" binding:"max=20"`          // 设备类型
}

// EventAnalysisRequest 行为事件分析请求
type EventAnalysisRequest struct {
	EventType string `form:"event_type"`                           // 事件类型筛选
	StartDate string `form:"start_date" binding:"required"`        // 开始日期
	EndDate   string `form:"end_date" binding:"required"`          // 结束日期
	GroupBy   string `form:"group_by" binding:"oneof=date event_type page device_type"` // 分组维度
}

// EventAnalysisItem 行为事件分析项
type EventAnalysisItem struct {
	Dimension string `json:"dimension"` // 维度值
	Count     int64  `json:"count"`     // 事件数量
	Users     int64  `json:"users"`     // 独立用户数
}

// EventSummary 行为事件摘要
type EventSummary struct {
	TotalEvents int64               `json:"total_events"`
	TotalUsers  int64               `json:"total_users"`
	Items       []EventAnalysisItem `json:"items"`
}

// ==================== 用户分群 ====================

// CreateSegmentRequest 创建分群请求
type CreateSegmentRequest struct {
	Name        string `json:"name" binding:"required,max=100"`
	Description string `json:"description" binding:"max=500"`
	Rules       string `json:"rules" binding:"required"` // 分群规则(JSON)
}

// UpdateSegmentRequest 更新分群请求
type UpdateSegmentRequest struct {
	Name        string `json:"name" binding:"max=100"`
	Description string `json:"description" binding:"max=500"`
	Rules       string `json:"rules"`
	IsActive    *bool  `json:"is_active"`
}

// SegmentInfo 分群信息
type SegmentInfo struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Rules       string `json:"rules"`
	UserCount   int    `json:"user_count"`
	IsActive    bool   `json:"is_active"`
	CreatorID   uint   `json:"creator_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// SegmentListRequest 分群列表请求
type SegmentListRequest struct {
	PageRequest
	IsActive *bool  `form:"is_active"` // 是否启用
	Keyword  string `form:"keyword"`   // 关键字搜索
}

// SegmentMemberInfo 分群成员信息
type SegmentMemberInfo struct {
	ID       uint   `json:"id"`
	UserID   uint   `json:"user_id"`
	UserName string `json:"user_name"`
	Avatar   string `json:"avatar"`
	JoinedAt string `json:"joined_at"`
}

// SegmentDetail 分群详情
type SegmentDetail struct {
	SegmentInfo
	Members []SegmentMemberInfo `json:"members"`
}

// ==================== 用户路径分析 ====================

// PathAnalysisRequest 路径分析请求
type PathAnalysisRequest struct {
	StartDate  string `form:"start_date" binding:"required"` // 开始日期
	EndDate    string `form:"end_date" binding:"required"`   // 结束日期
	DeviceType string `form:"device_type"`                   // 设备类型
	Limit      int    `form:"limit"`                         // 返回条数
}

// PathItem 路径分析项
type PathItem struct {
	Page       string  `json:"page"`
	Count      int64   `json:"count"`
	Users      int64   `json:"users"`
	AvgTime    float64 `json:"avg_time"`    // 平均停留时长(秒)
	BounceRate float64 `json:"bounce_rate"` // 跳出率
}

// PathTransition 路径流转
type PathTransition struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Count int64  `json:"count"`
	Users int64  `json:"users"`
}

// PathAnalysisResult 路径分析结果
type PathAnalysisResult struct {
	Pages       []PathItem       `json:"pages"`
	Transitions []PathTransition `json:"transitions"`
}

// ==================== 转化漏斗 ====================

// CreateFunnelRequest 创建漏斗请求
type CreateFunnelRequest struct {
	Name        string `json:"name" binding:"required,max=100"`
	Description string `json:"description" binding:"max=500"`
	Steps       string `json:"steps" binding:"required"` // 漏斗步骤(JSON)
}

// UpdateFunnelRequest 更新漏斗请求
type UpdateFunnelRequest struct {
	Name        string `json:"name" binding:"max=100"`
	Description string `json:"description" binding:"max=500"`
	Steps       string `json:"steps"`
	IsActive    *bool  `json:"is_active"`
}

// FunnelInfo 漏斗信息
type FunnelInfo struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Steps       string `json:"steps"`
	IsActive    bool   `json:"is_active"`
	CreatorID   uint   `json:"creator_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// FunnelStatsRequest 漏斗统计请求
type FunnelStatsRequest struct {
	StartDate string `form:"start_date" binding:"required"` // 开始日期
	EndDate   string `form:"end_date" binding:"required"`   // 结束日期
}

// FunnelStepStat 漏斗步骤统计
type FunnelStepStat struct {
	StepIndex      int     `json:"step_index"`
	StepName       string  `json:"step_name"`
	UserCount      int     `json:"user_count"`
	ConversionRate float64 `json:"conversion_rate"` // 相对转化率
	TotalRate      float64 `json:"total_rate"`      // 总体转化率
}

// FunnelStatsResult 漏斗统计结果
type FunnelStatsResult struct {
	FunnelID   uint              `json:"funnel_id"`
	FunnelName string            `json:"funnel_name"`
	StartDate  string            `json:"start_date"`
	EndDate    string            `json:"end_date"`
	Steps      []FunnelStepStat  `json:"steps"`
}

// ==================== 数据预警 ====================

// CreateAlertRuleRequest 创建预警规则请求
type CreateAlertRuleRequest struct {
	Name       string  `json:"name" binding:"required,max=100"`
	MetricName string  `json:"metric_name" binding:"required,max=100"` // 指标名
	Condition  string  `json:"condition" binding:"required,oneof=gt lt eq gte lte"` // 条件
	Threshold  float64 `json:"threshold" binding:"required"`            // 阈值
	Duration   int     `json:"duration"`                                // 持续时间(分钟)
	NotifyType string  `json:"notify_type" binding:"required,oneof=system email sms"` // 通知方式
}

// AlertRuleInfo 预警规则信息
type AlertRuleInfo struct {
	ID         uint    `json:"id"`
	Name       string  `json:"name"`
	MetricName string  `json:"metric_name"`
	Condition  string  `json:"condition"`
	Threshold  float64 `json:"threshold"`
	Duration   int     `json:"duration"`
	NotifyType string  `json:"notify_type"`
	IsActive   bool    `json:"is_active"`
	CreatorID  uint    `json:"creator_id"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}

// AlertRuleListRequest 预警规则列表请求
type AlertRuleListRequest struct {
	PageRequest
	IsActive *bool  `form:"is_active"`
	Keyword  string `form:"keyword"`
}

// AlertRecordInfo 预警记录信息
type AlertRecordInfo struct {
	ID          uint    `json:"id"`
	RuleID      uint    `json:"rule_id"`
	RuleName    string  `json:"rule_name"`
	MetricValue float64 `json:"metric_value"`
	Threshold   float64 `json:"threshold"`
	Level       string  `json:"level"`
	Message     string  `json:"message"`
	HandledAt   string  `json:"handled_at,omitempty"`
	HandlerID   *uint   `json:"handler_id,omitempty"`
	CreatedAt   string  `json:"created_at"`
}

// AlertRecordListRequest 预警记录列表请求
type AlertRecordListRequest struct {
	PageRequest
	Level  string `form:"level" binding:"omitempty,oneof=info warning critical"` // 预警级别
	RuleID uint   `form:"rule_id"`                                               // 规则ID
}

// ==================== 数据订阅 ====================

// CreateSubscriptionRequest 创建订阅请求
type CreateSubscriptionRequest struct {
	ReportType string `json:"report_type" binding:"required,max=50"` // 报表类型
	Frequency  string `json:"frequency" binding:"required,oneof=daily weekly monthly"` // 订阅频率
	Channels   string `json:"channels" binding:"required"` // 通知渠道(JSON)
}

// SubscriptionInfo 订阅信息
type SubscriptionInfo struct {
	ID         uint   `json:"id"`
	UserID     uint   `json:"user_id"`
	ReportType string `json:"report_type"`
	Frequency  string `json:"frequency"`
	Channels   string `json:"channels"`
	IsActive   bool   `json:"is_active"`
	LastSentAt string `json:"last_sent_at,omitempty"`
	CreatedAt  string `json:"created_at"`
}

// ==================== 数据导出 ====================

// ExportRequest 数据导出请求
type ExportRequest struct {
	ExportType string `json:"export_type" binding:"required,max=50"` // 导出类型
	FileFormat string `json:"file_format" binding:"required,oneof=xlsx csv pdf"` // 文件格式
	Filters    string `json:"filters"` // 筛选条件(JSON)
}

// ExportInfo 导出信息
type ExportInfo struct {
	ID         uint   `json:"id"`
	UserID     uint   `json:"user_id"`
	ExportType string `json:"export_type"`
	FileFormat string `json:"file_format"`
	FilePath   string `json:"file_path"`
	FileSize   int64  `json:"file_size"`
	Status     string `json:"status"`
	CreatedAt  string `json:"created_at"`
}

// ==================== 数据对比 ====================

// CompareRequest 数据对比请求
type CompareRequest struct {
	SnapshotType string `form:"snapshot_type" binding:"required"` // 快照类型
	Period1Start string `form:"period1_start" binding:"required"` // 第一段开始
	Period1End   string `form:"period1_end" binding:"required"`   // 第一段结束
	Period2Start string `form:"period2_start" binding:"required"` // 第二段开始
	Period2End   string `form:"period2_end" binding:"required"`   // 第二段结束
}

// CompareResult 数据对比结果
type CompareResult struct {
	Period1 ComparePeriod `json:"period1"`
	Period2 ComparePeriod `json:"period2"`
	Diff    CompareDiff   `json:"diff"`
}

// ComparePeriod 对比时段数据
type ComparePeriod struct {
	StartDate string                 `json:"start_date"`
	EndDate   string                 `json:"end_date"`
	Data      map[string]interface{} `json:"data"`
}

// CompareDiff 对比差异
type CompareDiff struct {
	Values map[string]float64 `json:"values"` // 各指标变化值
	Rates  map[string]float64 `json:"rates"`  // 各指标变化率(%)
}

// ==================== 题目分析 ====================

// QuestionDifficultyRequest 题目难度分析请求
type QuestionDifficultyRequest struct {
	CategoryID uint `form:"category_id"` // 分类筛选
}

// QuestionDifficultyItem 题目难度分析项
type QuestionDifficultyItem struct {
	QuestionID   uint    `json:"question_id"`
	Title        string  `json:"title"`
	Type         int8    `json:"type"`
	Difficulty    int8    `json:"difficulty"`
	CorrectRate  float64 `json:"correct_rate"`
	AnswerCount  int     `json:"answer_count"`
	AvgDuration  float64 `json:"avg_duration"`
}

// QuestionDiscriminationRequest 题目区分度分析请求
type QuestionDiscriminationRequest struct {
	CategoryID uint `form:"category_id"` // 分类筛选
}

// QuestionDiscriminationItem 题目区分度分析项
type QuestionDiscriminationItem struct {
	QuestionID     uint    `json:"question_id"`
	Title          string  `json:"title"`
	Type           int8    `json:"type"`
	Difficulty     int8    `json:"difficulty"`
	Discrimination float64 `json:"discrimination"` // 区分度系数
	HighGroupRate  float64 `json:"high_group_rate"` // 高分组正确率
	LowGroupRate   float64 `json:"low_group_rate"`  // 低分组正确率
}

// ==================== 成绩预测与预警 ====================

// ScorePredictionRequest 成绩预测请求
type ScorePredictionRequest struct {
	UserID uint `form:"user_id"` // 用户ID(管理员可指定)
	ExamID uint `form:"exam_id"` // 考试ID
}

// ScorePrediction 成绩预测结果
type ScorePrediction struct {
	UserID          uint    `json:"user_id"`
	UserName        string  `json:"user_name"`
	ExamID          uint    `json:"exam_id"`
	ExamTitle       string  `json:"exam_title"`
	PredictedScore  float64 `json:"predicted_score"`
	ConfidenceLevel float64 `json:"confidence_level"` // 置信度
	HistoryAvg      float64 `json:"history_avg"`      // 历史平均分
	HistoryBest     float64 `json:"history_best"`     // 历史最高分
	Trend           string  `json:"trend"`            // 趋势:up/down/stable
}

// ScoreAlertRequest 成绩预警请求
type ScoreAlertRequest struct {
	UserID uint `form:"user_id"` // 用户ID
	ExamID uint `form:"exam_id"` // 考试ID
}

// ScoreAlertItem 成绩预警项
type ScoreAlertItem struct {
	UserID        uint    `json:"user_id"`
	UserName      string  `json:"user_name"`
	ExamID        uint    `json:"exam_id"`
	ExamTitle     string  `json:"exam_title"`
	LatestScore   float64 `json:"latest_score"`
	AvgScore      float64 `json:"avg_score"`
	Trend         string  `json:"trend"`
	AlertLevel    string  `json:"alert_level"` // warning/critical
	AlertMessage  string  `json:"alert_message"`
}

// ==================== 班级统计(教师视角) ====================

// ClassOverview 班级概览
type ClassOverview struct {
	ClassID        uint    `json:"class_id"`
	ClassName      string  `json:"class_name"`
	TotalStudents  int     `json:"total_students"`
	ActiveStudents int     `json:"active_students"`
	AvgScore       float64 `json:"avg_score"`
	PassRate       float64 `json:"pass_rate"`
	PracticeCount  int64   `json:"practice_count"`
	ExamCount      int64   `json:"exam_count"`
}

// ClassStudentItem 班级学生成绩项
type ClassStudentItem struct {
	UserID       uint    `json:"user_id"`
	UserName     string  `json:"user_name"`
	Avatar       string  `json:"avatar"`
	AvgScore     float64 `json:"avg_score"`
	BestScore    float64 `json:"best_score"`
	PracticeCount int    `json:"practice_count"`
	ExamCount    int     `json:"exam_count"`
	LastActiveAt string  `json:"last_active_at"`
}

// ClassStudentListRequest 班级学生成绩列表请求
type ClassStudentListRequest struct {
	PageRequest
	SortBy string `form:"sort_by" binding:"omitempty,oneof=avg_score practice_count exam_count last_active"` // 排序字段
}

// ClassPracticeStats 班级练习统计
type ClassPracticeStats struct {
	TotalSessions int64            `json:"total_sessions"`
	TotalQuestions int64           `json:"total_questions"`
	AvgAccuracy   float64          `json:"avg_accuracy"`
	TotalDuration int64            `json:"total_duration"`
	DailyStats    []PracticeDayStat `json:"daily_stats"`
}

// PracticeDayStat 练习日统计
type PracticeDayStat struct {
	Date      string  `json:"date"`
	Sessions  int     `json:"sessions"`
	Questions int     `json:"questions"`
	Accuracy  float64 `json:"accuracy"`
}

// ClassPracticeRequest 班级练习统计请求
type ClassPracticeRequest struct {
	StartDate string `form:"start_date"` // 开始日期
	EndDate   string `form:"end_date"`   // 结束日期
}

// ClassExamStats 班级考试统计
type ClassExamStats struct {
	TotalExams   int64          `json:"total_exams"`
	TotalRecords int64          `json:"total_records"`
	AvgScore     float64        `json:"avg_score"`
	PassRate     float64        `json:"pass_rate"`
	ExamList     []ClassExamItem `json:"exam_list"`
}

// ClassExamItem 班级考试项
type ClassExamItem struct {
	ExamID    uint    `json:"exam_id"`
	Title     string  `json:"title"`
	StartTime string  `json:"start_time"`
	AvgScore  float64 `json:"avg_score"`
	PassRate  float64 `json:"pass_rate"`
	SubCount  int     `json:"sub_count"` // 提交人数
}

// ClassExamRequest 班级考试统计请求
type ClassExamRequest struct {
	StartDate string `form:"start_date"` // 开始日期
	EndDate   string `form:"end_date"`   // 结束日期
}

// ClassQuestionStats 班级题目统计
type ClassQuestionStats struct {
	TotalQuestions int64                 `json:"total_questions"`
	CorrectRate    float64               `json:"correct_rate"`
	TypeStats      []QuestionTypeStat    `json:"type_stats"`
	DifficultyStats []QuestionDiffStat   `json:"difficulty_stats"`
}

// QuestionTypeStat 题目类型统计
type QuestionTypeStat struct {
	Type     int8    `json:"type"`
	Count    int64   `json:"count"`
	AvgRate  float64 `json:"avg_rate"`
}

// QuestionDiffStat 题目难度统计
type QuestionDiffStat struct {
	Difficulty int8    `json:"difficulty"`
	Count      int64   `json:"count"`
	AvgRate    float64 `json:"avg_rate"`
}

// ==================== 移动端统计 ====================

// MobileOverview 移动端个人概览
type MobileOverview struct {
	TotalPractice  int     `json:"total_practice"`
	TotalQuestions int     `json:"total_questions"`
	AvgAccuracy    float64 `json:"avg_accuracy"`
	TotalExams     int     `json:"total_exams"`
	AvgScore       float64 `json:"avg_score"`
	StudyDays      int     `json:"study_days"`
	StudyDuration  int64   `json:"study_duration"`
	Rank           int     `json:"rank"`
}

// MobilePracticeStats 移动端练习统计
type MobilePracticeStats struct {
	TodayCount    int              `json:"today_count"`
	WeekCount     int              `json:"week_count"`
	MonthCount    int              `json:"month_count"`
	TodayAccuracy float64          `json:"today_accuracy"`
	WeekAccuracy  float64          `json:"week_accuracy"`
	DailyStats    []PracticeDayStat `json:"daily_stats"`
}

// MobilePracticeRequest 移动端练习统计请求
type MobilePracticeRequest struct {
	Days int `form:"days"` // 最近N天
}

// MobileWrongStats 移动端错题统计
type MobileWrongStats struct {
	TotalWrong    int64               `json:"total_wrong"`
	UndoWrong     int64               `json:"undo_wrong"`
	DoneWrong     int64               `json:"done_wrong"`
	TypeStats     []WrongTypeStat     `json:"type_stats"`
	DailyStats    []WrongDayStat      `json:"daily_stats"`
}

// WrongTypeStat 错题类型统计
type WrongTypeStat struct {
	Type  int8  `json:"type"`
	Count int64 `json:"count"`
}

// WrongDayStat 错题日统计
type WrongDayStat struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

// MobileTrendRequest 移动端学习趋势请求
type MobileTrendRequest struct {
	Days int `form:"days"` // 最近N天
}

// MobileTrendItem 移动端趋势项
type MobileTrendItem struct {
	Date          string  `json:"date"`
	PracticeCount int     `json:"practice_count"`
	QuestionCount int     `json:"question_count"`
	Accuracy      float64 `json:"accuracy"`
	Duration      int     `json:"duration"`
}

// DateRangeRequest 通用日期范围请求
type DateRangeRequest struct {
	StartDate string `form:"start_date" binding:"required"`
	EndDate   string `form:"end_date" binding:"required"`
}
