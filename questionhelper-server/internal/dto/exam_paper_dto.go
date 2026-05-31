package dto

// PaperPreviewResponse 试卷预览响应
type PaperPreviewResponse struct {
	PaperInfo
	Sections []PaperSection `json:"sections"` // 按题型分组
}

// PaperSection 试卷分组
type PaperSection struct {
	Type     int8           `json:"type"`
	TypeName string         `json:"type_name"`
	Count    int            `json:"count"`
	Score    float64        `json:"score"`
	Questions []PaperQuestionInfo `json:"questions"`
}

// PaperQuestionInfo 试卷题目信息
type PaperQuestionInfo struct {
	ID         uint         `json:"id"`
	QuestionID uint         `json:"question_id"`
	Score      float64      `json:"score"`
	Sort       int          `json:"sort"`
	Question   *QuestionInfo `json:"question,omitempty"`
}

// CopyPaperRequest 复制试卷请求
type CopyPaperRequest struct {
	Title string `json:"title" binding:"omitempty,max=200"`
}

// PublishPaperRequest 发布试卷请求
type PublishPaperRequest struct {
	Status int8 `json:"status" binding:"required,oneof=1 2"` // 1=发布,2=归档
}

// SaveTemplateRequest 保存为模板请求
type SaveTemplateRequest struct {
	Category string `json:"category" binding:"omitempty,max=50"`
	Tags     string `json:"tags" binding:"omitempty,max=500"`
}

// TemplateListRequest 模板列表请求
type TemplateListRequest struct {
	PageRequest
	Category string `form:"category"`
	Keyword  string `form:"keyword"`
}

// ImportPaperRequest 导入试卷请求
type ImportPaperRequest struct {
	Title      string `json:"title" binding:"required,max=200"`
	CategoryID uint   `json:"category_id"`
	Visibility int8   `json:"visibility"`
}

// ExportPaperResponse 导出试卷响应
type ExportPaperResponse struct {
	Title       string             `json:"title"`
	Description string             `json:"description"`
	TotalScore  float64            `json:"total_score"`
	Questions   []ExportQuestion   `json:"questions"`
}

// ExportQuestion 导出题目
type ExportQuestion struct {
	Sort       int            `json:"sort"`
	Title      string         `json:"title"`
	Type       int8           `json:"type"`
	Difficulty int8           `json:"difficulty"`
	Score      float64        `json:"score"`
	Answer     string         `json:"answer"`
	Options    []OptionInfo   `json:"options,omitempty"`
}

// PaperStatsResponse 试卷统计响应
type PaperStatsResponse struct {
	UsageCount   int              `json:"usage_count"`
	AvgScore     float64          `json:"avg_score"`
	MaxScore     float64          `json:"max_score"`
	MinScore     float64          `json:"min_score"`
	PassRate     float64          `json:"pass_rate"`
	ScoreDistribution []ScoreDist `json:"score_distribution"`
}

// ScoreDist 分数分布
type ScoreDist struct {
	Range string `json:"range"` // 0-10, 10-20, ...
	Count int    `json:"count"`
}
