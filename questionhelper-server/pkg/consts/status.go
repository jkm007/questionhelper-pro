package consts

// 通用状态
const (
	StatusDisabled = 0 // 禁用
	StatusEnabled  = 1 // 正常
)

// 题目状态
const (
	QuestionStatusDraft   = 0 // 草稿
	QuestionStatusNormal  = 1 // 正常
	QuestionStatusDisabled = 2 // 禁用
)

// 题目类型
const (
	QuestionTypeSingle   = 1 // 单选题
	QuestionTypeMultiple = 2 // 多选题
	QuestionTypeJudge    = 3 // 判断题
	QuestionTypeFill     = 4 // 填空题
	QuestionTypeShort    = 5 // 简答题
)

// 难度等级
const (
	DifficultyEasy   = 1 // 简单
	DifficultyMedium = 2 // 中等
	DifficultyHard   = 3 // 困难
)

// 可见性
const (
	VisibilityPublic  = 1 // 公开
	VisibilityPrivate = 2 // 私有
	VisibilityClass   = 3 // 班级
)

// 考试状态
const (
	ExamStatusDraft     = 0 // 未发布
	ExamStatusOngoing   = 1 // 进行中
	ExamStatusFinished  = 2 // 已结束
)

// 答题记录状态
const (
	RecordStatusOngoing  = 0 // 进行中
	RecordStatusSubmit   = 1 // 已交卷
	RecordStatusReviewed = 2 // 已阅卷
)

// 防作弊级别
const (
	AntiCheatNone   = 0 // 无
	AntiCheatBasic  = 1 // 基础
	AntiCheatStrict = 2 // 严格
)

// 显示答案
const (
	ShowAnswerNever      = 0 // 不显示
	ShowAnswerAfterSubmit = 1 // 交卷后
	ShowAnswerAfterExam  = 2 // 考试后
)
