package dto

// BatchQuestionRequest 批量操作题目请求
type BatchQuestionRequest struct {
	IDs []uint `json:"ids" binding:"required,min=1"`
}

// BatchMoveRequest 批量移动分类请求
type BatchMoveRequest struct {
	IDs        []uint `json:"ids" binding:"required,min=1"`
	CategoryID uint   `json:"category_id" binding:"required"`
}

// BatchResult 批量操作结果
type BatchResult struct {
	Total   int          `json:"total"`
	Success int          `json:"success"`
	Failed  int          `json:"failed"`
	Errors  []BatchError `json:"errors,omitempty"`
}

// BatchError 批量操作错误
type BatchError struct {
	ID     uint   `json:"id"`
	Reason string `json:"reason"`
}

// DuplicateResult 查重结果
type DuplicateResult struct {
	IsDuplicate bool             `json:"is_duplicate"`
	Count       int              `json:"count"`
	Questions   []QuestionInfo   `json:"questions,omitempty"`
}

// QuestionPreview 题目预览
type QuestionPreview struct {
	QuestionInfo
	RenderedContent string       `json:"rendered_content"` // 渲染后的HTML内容
	Options         []OptionInfo `json:"options"`
}
