package consts

// 题目状态
const (
	QuestionStatusDraft    = 0 // 草稿
	QuestionStatusNormal   = 1 // 正常
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
