package consts

// 考试状态
const (
	ExamStatusDraft    = 0 // 未发布
	ExamStatusOngoing  = 1 // 进行中
	ExamStatusFinished = 2 // 已结束
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
