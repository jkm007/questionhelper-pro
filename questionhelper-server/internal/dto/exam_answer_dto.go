package dto

// SaveAnswerRequest 保存答案请求
type SaveAnswerRequest struct {
	QuestionID uint   `json:"question_id" binding:"required"`
	Answer     string `json:"answer"`
	IsMarked   *bool  `json:"is_marked"` // 是否标记
}

// SaveAnswerBatchRequest 批量保存答案请求
type SaveAnswerBatchRequest struct {
	Answers []SaveAnswerRequest `json:"answers" binding:"required,min=1"`
}

// StandardAnswersResponse 标准答案响应
type StandardAnswersResponse struct {
	ExamID    uint                `json:"exam_id"`
	Questions []StandardAnswer   `json:"questions"`
}

// StandardAnswer 标准答案
type StandardAnswer struct {
	QuestionID uint         `json:"question_id"`
	Title      string       `json:"title"`
	Type       int8         `json:"type"`
	Score      float64      `json:"score"`
	Answer     string       `json:"answer"`
	Analysis   string       `json:"analysis"`
	Options    []OptionInfo `json:"options,omitempty"`
	UserAnswer string       `json:"user_answer,omitempty"`
	UserScore  float64      `json:"user_score"`
	IsCorrect  bool         `json:"is_correct"`
}

// MarkQuestionRequest 标记题目请求
type MarkQuestionRequest struct {
	QuestionID uint `json:"question_id" binding:"required"`
	IsMarked   bool `json:"is_marked"`
}

// ExamGuideResponse 答题指引响应
type ExamGuideResponse struct {
	ExamID      uint   `json:"exam_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Duration    int    `json:"duration"`
	TotalScore  float64 `json:"total_score"`
	PassScore   float64 `json:"pass_score"`
	Rules       []string `json:"rules"` // 考试规则
	Tips        []string `json:"tips"`  // 答题技巧
}

// FeedbackRequest 考后反馈请求
type FeedbackRequest struct {
	Rating  int8   `json:"rating" binding:"required,min=1,max=5"` // 评分1-5
	Content string `json:"content" binding:"omitempty,max=500"`
}
