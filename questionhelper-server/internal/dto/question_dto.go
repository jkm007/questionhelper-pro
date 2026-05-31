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
	UserID     *uint  `form:"-" json:"-"` // 数据权限过滤用，不由前端直接传入
}

// ImportQuestionRequest 导入题目请求
type ImportQuestionRequest struct {
	CategoryID uint `json:"category_id"`
	Visibility int8 `json:"visibility" binding:"required,oneof=1 2 3"`
}

// CreateCategoryRequest 创建分类请求
type CreateCategoryRequest struct {
	ParentID *uint  `json:"parent_id"`
	Name     string `json:"name" binding:"required,max=100"`
	Sort     int    `json:"sort"`
}

// UpdateCategoryRequest 更新分类请求
type UpdateCategoryRequest struct {
	ParentID *uint  `json:"parent_id"`
	Name     string `json:"name" binding:"omitempty,max=100"`
	Sort     int    `json:"sort"`
}

// CreateKnowledgeRequest 创建知识点请求
type CreateKnowledgeRequest struct {
	CategoryID uint   `json:"category_id" binding:"required"`
	Name       string `json:"name" binding:"required,max=100"`
}

// UpdateKnowledgeRequest 更新知识点请求
type UpdateKnowledgeRequest struct {
	CategoryID uint   `json:"category_id"`
	Name       string `json:"name" binding:"omitempty,max=100"`
}

// CreateSensitiveWordRequest 创建敏感词请求
type CreateSensitiveWordRequest struct {
	Word string `json:"word" binding:"required"`
}

// ReviewInfo 审核信息
type ReviewInfo struct {
	ID         uint   `json:"id"`
	QuestionID uint   `json:"question_id"`
	ReviewerID uint   `json:"reviewer_id"`
	Status     int8   `json:"status"`
	Reason     string `json:"reason"`
	CreatedAt  string `json:"created_at"`
}
