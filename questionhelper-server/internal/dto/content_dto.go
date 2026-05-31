package dto

// ==================== 创作者申请 ====================

// ApplyCreatorRequest 申请成为创作者
type ApplyCreatorRequest struct {
	Reason      string `json:"reason" binding:"required,max=500"`
	RealName    string `json:"real_name" binding:"required,max=50"`
	IDCard      string `json:"id_card" binding:"required"`
	Phone       string `json:"phone" binding:"required"`
	Email       string `json:"email" binding:"omitempty,email"`
	Description string `json:"description" binding:"max=1000"`
}

// ==================== 创作者信息 ====================

// CreatorProfileInfo 创作者资料
type CreatorProfileInfo struct {
	UserID      uint   `json:"user_id"`
	Username    string `json:"username"`
	Nickname    string `json:"nickname"`
	Avatar      string `json:"avatar"`
	Bio         string `json:"bio"`
	Level       int    `json:"level"`
	LevelName   string `json:"level_name"`
	TotalPoints int    `json:"total_points"`
	IsCreator   bool   `json:"is_creator"`
	Agreement   bool   `json:"agreement"`
	CreatedAt   string `json:"created_at"`
}

// CreatorLevelInfo 创作者等级信息
type CreatorLevelInfo struct {
	ID        uint   `json:"id"`
	Level     int    `json:"level"`
	Name      string `json:"name"`
	MinPoints int    `json:"min_points"`
	MaxPoints int    `json:"max_points"`
	Benefits  string `json:"benefits"`
	Icon      string `json:"icon"`
	IsCurrent bool   `json:"is_current"`
}

// CreatorPointsInfo 积分信息
type CreatorPointsInfo struct {
	UserID          uint   `json:"user_id"`
	TotalPoints     int    `json:"total_points"`
	AvailablePoints int    `json:"available_points"`
	Level           int    `json:"level"`
	LevelName       string `json:"level_name"`
	NextLevel       int    `json:"next_level"`
	NextLevelName   string `json:"next_level_name"`
	PointsNeeded    int    `json:"points_needed"`
}

// CreatorPointLogInfo 积分记录
type CreatorPointLogInfo struct {
	ID         uint   `json:"id"`
	ChangeType string `json:"change_type"`
	Points     int    `json:"points"`
	Balance    int    `json:"balance"`
	Reason     string `json:"reason"`
	CreatedAt  string `json:"created_at"`
}

// CreatorPointLogsRequest 积分记录查询
type CreatorPointLogsRequest struct {
	PageRequest
	ChangeType string `form:"change_type"`
}

// ==================== 创作者协议 ====================

// CreatorAgreementInfo 协议信息
type CreatorAgreementInfo struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Version   string `json:"version"`
	IsActive  bool   `json:"is_active"`
	Signed    bool   `json:"signed"`
	SignedAt  string `json:"signed_at,omitempty"`
	CreatedAt string `json:"created_at"`
}

// SignAgreementRequest 签署协议
type SignAgreementRequest struct {
	AgreementID uint `json:"agreement_id" binding:"required"`
}

// ==================== 创作者作品集 ====================

// PortfolioInfo 作品集信息
type PortfolioInfo struct {
	ID          uint   `json:"id"`
	UserID      uint   `json:"user_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CoverImage  string `json:"cover_image"`
	IsPublic    bool   `json:"is_public"`
	ViewCount   int    `json:"view_count"`
	LikeCount   int    `json:"like_count"`
	ItemCount   int    `json:"item_count"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// CreatePortfolioRequest 创建作品集
type CreatePortfolioRequest struct {
	Name        string `json:"name" binding:"required,max=100"`
	Description string `json:"description" binding:"max=500"`
	CoverImage  string `json:"cover_image"`
	IsPublic    *bool  `json:"is_public" binding:"required"`
}

// UpdatePortfolioRequest 更新作品集
type UpdatePortfolioRequest struct {
	Name        string `json:"name" binding:"omitempty,max=100"`
	Description string `json:"description" binding:"max=500"`
	CoverImage  string `json:"cover_image"`
	IsPublic    *bool  `json:"is_public"`
}

// PortfolioListRequest 作品集列表请求
type PortfolioListRequest struct {
	PageRequest
	UserID *uint `form:"user_id"`
}

// ==================== 内容版本管理 ====================

// ContentVersionInfo 版本信息
type ContentVersionInfo struct {
	ID          uint   `json:"id"`
	ContentType string `json:"content_type"`
	ContentID   uint   `json:"content_id"`
	Version     int    `json:"version"`
	Data        string `json:"data"`
	CreatorID   uint   `json:"creator_id"`
	CreatorName string `json:"creator_name"`
	Remark      string `json:"remark"`
	CreatedAt   string `json:"created_at"`
}

// RollbackVersionRequest 回滚版本
type RollbackVersionRequest struct {
	VersionID uint `json:"version_id" binding:"required"`
}

// ==================== 内容标签 ====================

// ContentTagInfo 内容标签信息
type ContentTagInfo struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Color       string `json:"color"`
	UsageCount  int    `json:"usage_count"`
}

// AddContentTagRequest 添加内容标签
type AddContentTagRequest struct {
	TagName string `json:"tag_name" binding:"required,max=50"`
}

// HotTagInfo 热门标签信息
type HotTagInfo struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Slug       string `json:"slug"`
	UsageCount int    `json:"usage_count"`
	Rank       int    `json:"rank"`
}

// ==================== 内容收藏（通用）====================

// ContentFavoriteInfo 内容收藏信息
type ContentFavoriteInfo struct {
	ID          uint   `json:"id"`
	UserID      uint   `json:"user_id"`
	ContentType string `json:"content_type"`
	ContentID   uint   `json:"content_id"`
	FolderName  string `json:"folder_name"`
	CreatedAt   string `json:"created_at"`
}

// AddContentFavoriteRequest 添加收藏
type AddContentFavoriteRequest struct {
	FolderName string `json:"folder_name" binding:"omitempty,max=100"`
}

// FavoriteListRequest 收藏列表请求
type FavoriteListRequest struct {
	PageRequest
	ContentType string `form:"content_type"`
	FolderName  string `form:"folder_name"`
}

// FavoriteListInfo 收藏列表信息
type FavoriteListInfo struct {
	ID          uint        `json:"id"`
	ContentType string      `json:"content_type"`
	ContentID   uint        `json:"content_id"`
	FolderName  string      `json:"folder_name"`
	Content     interface{} `json:"content,omitempty"`
	CreatedAt   string      `json:"created_at"`
}

// ==================== 内容预览 ====================

// ContentPreviewInfo 内容预览信息
type ContentPreviewInfo struct {
	ContentType string      `json:"content_type"`
	ContentID   uint        `json:"content_id"`
	Title       string      `json:"title"`
	Summary     string      `json:"summary"`
	Content     interface{} `json:"content"`
	CreatorID   uint        `json:"creator_id"`
	CreatorName string      `json:"creator_name"`
	CreatedAt   string      `json:"created_at"`
}

// ==================== 综合搜索 ====================

// SearchRequest 搜索请求
type SearchRequest struct {
	PageRequest
	Keyword     string `form:"keyword" binding:"required"`
	ContentType string `form:"content_type"`
	Sort        string `form:"sort"`
}

// SearchResultInfo 搜索结果信息
type SearchResultInfo struct {
	ContentType string      `json:"content_type"`
	ContentID   uint        `json:"content_id"`
	Title       string      `json:"title"`
	Summary     string      `json:"summary"`
	Score       float64     `json:"score"`
	Highlight   string      `json:"highlight"`
	Content     interface{} `json:"content,omitempty"`
}

// SearchSuggestionInfo 搜索建议信息
type SearchSuggestionInfo struct {
	Keyword string `json:"keyword"`
	Type    string `json:"type"`
}

// HotSearchInfo 热门搜索信息
type HotSearchInfo struct {
	Keyword string `json:"keyword"`
	Count   int    `json:"count"`
	Rank    int    `json:"rank"`
}

// SearchHistoryInfo 搜索历史信息
type SearchHistoryInfo struct {
	ID        uint   `json:"id"`
	Keyword   string `json:"keyword"`
	CreatedAt string `json:"created_at"`
}

// ==================== 审核流程 ====================

// ReviewListRequest 审核列表请求
type ReviewListRequest struct {
	PageRequest
	Status      int8   `form:"status"`
	ContentType string `form:"content_type"`
}

// ReviewInstanceInfo 审核实例信息
type ReviewInstanceInfo struct {
	ID            uint             `json:"id"`
	WorkflowID    uint             `json:"workflow_id"`
	ContentType   string           `json:"content_type"`
	ContentID     uint             `json:"content_id"`
	CurrentStep   int              `json:"current_step"`
	Status        int8             `json:"status"`
	CreatorID     uint             `json:"creator_id"`
	CreatorName   string           `json:"creator_name"`
	ContentTitle  string           `json:"content_title"`
	Steps         []ReviewStepInfo `json:"steps"`
	CreatedAt     string           `json:"created_at"`
	UpdatedAt     string           `json:"updated_at"`
}

// ReviewStepInfo 审核步骤信息
type ReviewStepInfo struct {
	ID         uint   `json:"id"`
	StepIndex  int    `json:"step_index"`
	StepName   string `json:"step_name"`
	ReviewerID uint   `json:"reviewer_id"`
	Action     string `json:"action"`
	Opinion    string `json:"opinion"`
	OperatedAt string `json:"operated_at"`
}

// ApproveReviewRequest 通过审核
type ApproveReviewRequest struct {
	Opinion string `json:"opinion" binding:"max=500"`
}

// RejectReviewRequest 拒绝审核
type RejectReviewRequest struct {
	Opinion string `json:"opinion" binding:"required,max=500"`
}

// AddReviewCommentRequest 添加审核意见
type AddReviewCommentRequest struct {
	Content string `json:"content" binding:"required,max=1000"`
}
