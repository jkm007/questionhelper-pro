package dto

// QuestionInfo 题目信息
type QuestionInfo struct {
	ID           uint         `json:"id"`
	Title        string       `json:"title"`
	Content      string       `json:"content"`
	Type         int8         `json:"type"`
	Difficulty   int8         `json:"difficulty"`
	Answer       string       `json:"answer"`
	Analysis     string       `json:"analysis"`
	CategoryID   uint         `json:"category_id"`
	Category     CategoryInfo `json:"category"`
	Options      []OptionInfo `json:"options"`
	Visibility   int8         `json:"visibility"`
	CreatorID    uint         `json:"creator_id"`
	CreatorName  string       `json:"creator_name"`
	Status       int8         `json:"status"`
	ViewCount    int          `json:"view_count"`
	LikeCount    int          `json:"like_count"`
}

// OptionInfo 选项信息
type OptionInfo struct {
	ID        uint   `json:"id"`
	Label     string `json:"label"`
	Content   string `json:"content"`
	IsCorrect bool   `json:"is_correct"`
}

// CategoryInfo 分类信息
type CategoryInfo struct {
	ID       uint           `json:"id"`
	ParentID *uint          `json:"parent_id"`
	Name     string         `json:"name"`
	Children []CategoryInfo `json:"children,omitempty"`
}

// CreateQuestionRequest 创建题目请求
type CreateQuestionRequest struct {
	Title        string              `json:"title" binding:"required"`
	Content      string              `json:"content"`
	Type         int8                `json:"type" binding:"required,oneof=1 2 3 4 5"`
	Difficulty   int8                `json:"difficulty" binding:"required,oneof=1 2 3"`
	Answer       string              `json:"answer"`
	Analysis     string              `json:"analysis"`
	CategoryID   uint                `json:"category_id"`
	Options      []CreateOptionRequest `json:"options"`
	Visibility   int8                `json:"visibility" binding:"required,oneof=1 2 3"`
	KnowledgeIDs []uint              `json:"knowledge_ids"`
}

// CreateOptionRequest 创建选项请求
type CreateOptionRequest struct {
	Label     string `json:"label" binding:"required"`
	Content   string `json:"content" binding:"required"`
	IsCorrect bool   `json:"is_correct"`
}

// UpdateQuestionRequest 更新题目请求
type UpdateQuestionRequest struct {
	Title        string              `json:"title"`
	Content      string              `json:"content"`
	Type         int8                `json:"type" binding:"omitempty,oneof=1 2 3 4 5"`
	Difficulty   int8                `json:"difficulty" binding:"omitempty,oneof=1 2 3"`
	Answer       string              `json:"answer"`
	Analysis     string              `json:"analysis"`
	CategoryID   *uint               `json:"category_id"`
	Options      []CreateOptionRequest `json:"options"`
	Visibility   int8                `json:"visibility" binding:"omitempty,oneof=1 2 3"`
	KnowledgeIDs []uint              `json:"knowledge_ids"`
}

// QuestionListRequest 题目列表请求
type QuestionListRequest struct {
	PageRequest
	Keyword    string `form:"keyword"`
	CategoryID *uint  `form:"category_id"`
	Type       *int8  `form:"type"`
	Difficulty *int8  `form:"difficulty"`
	Visibility *int8  `form:"visibility"`
	Status     *int8  `form:"status"`
}

// ImportQuestionRequest 导入题目请求
type ImportQuestionRequest struct {
	CategoryID uint `json:"category_id"`
	Visibility int8 `json:"visibility" binding:"required,oneof=1 2 3"`
}
