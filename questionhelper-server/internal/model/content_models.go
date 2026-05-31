package model

import (
	"time"

	"gorm.io/gorm"
)

// ContentReview 内容审核表
type ContentReview struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	ContentType string         `gorm:"size:50;index;not null" json:"content_type"`           // 内容类型(question/exam/article等)
	ContentID   uint           `gorm:"index;not null" json:"content_id"`                     // 内容ID
	ReviewerID  *uint          `gorm:"index" json:"reviewer_id"`                             // 审核人ID
	Status      int8           `gorm:"default:0;comment:0=待审核,1=通过,2=驳回" json:"status"` // 审核状态
	Opinion     string         `gorm:"size:500" json:"opinion"`                              // 审核意见
	ReviewedAt  *time.Time     `json:"reviewed_at"`                                          // 审核时间
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ContentReview) TableName() string {
	return "content_reviews"
}

// CreatorLevel 创作者等级表
type CreatorLevel struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Level     int       `gorm:"uniqueIndex;not null" json:"level"`    // 等级值
	Name      string    `gorm:"size:50;not null" json:"name"`         // 等级名称
	MinPoints int       `gorm:"not null" json:"min_points"`           // 最低积分
	MaxPoints int       `gorm:"not null" json:"max_points"`           // 最高积分
	Benefits  string    `gorm:"type:text" json:"benefits"`            // 权益说明(JSON)
	Icon      string    `gorm:"size:200" json:"icon"`                 // 等级图标
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (CreatorLevel) TableName() string {
	return "creator_levels"
}

// CreatorPoint 创作者积分表
type CreatorPoint struct {
	ID              uint      `gorm:"primarykey" json:"id"`
	UserID          uint      `gorm:"uniqueIndex;not null" json:"user_id"`   // 用户ID
	TotalPoints     int       `gorm:"default:0" json:"total_points"`         // 累计积分
	AvailablePoints int       `gorm:"default:0" json:"available_points"`     // 可用积分
	Level           int       `gorm:"default:1" json:"level"`                // 当前等级
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (CreatorPoint) TableName() string {
	return "creator_points"
}

// CreatorPointLog 积分变动记录表
type CreatorPointLog struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	UserID     uint      `gorm:"index;not null" json:"user_id"`                               // 用户ID
	ChangeType string    `gorm:"size:20;not null;comment:earn/spend/expire" json:"change_type"` // 变动类型
	Points     int       `gorm:"not null" json:"points"`                                      // 变动积分
	Balance    int       `gorm:"not null" json:"balance"`                                     // 变动后余额
	Reason     string    `gorm:"size:200" json:"reason"`                                      // 变动原因
	RelatedID  uint      `gorm:"index" json:"related_id"`                                     // 关联业务ID
	CreatedAt  time.Time `json:"created_at"`
}

func (CreatorPointLog) TableName() string {
	return "creator_point_logs"
}

// CreatorAgreement 创作者协议表
type CreatorAgreement struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Title       string         `gorm:"size:200;not null" json:"title"`   // 协议标题
	Content     string         `gorm:"type:text;not null" json:"content"` // 协议内容
	Version     string         `gorm:"size:20;not null" json:"version"`  // 协议版本
	IsActive    bool           `gorm:"default:false" json:"is_active"`   // 是否当前生效
	EffectiveAt *time.Time     `json:"effective_at"`                      // 生效时间
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (CreatorAgreement) TableName() string {
	return "creator_agreements"
}

// CreatorAgreementSign 创作者协议签署记录表
type CreatorAgreementSign struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	AgreementID uint      `gorm:"uniqueIndex:idx_agreement_user;not null" json:"agreement_id"` // 协议ID
	UserID      uint      `gorm:"uniqueIndex:idx_agreement_user;not null" json:"user_id"`      // 用户ID
	SignedAt    time.Time `json:"signed_at"`                                                    // 签署时间
	CreatedAt   time.Time `json:"created_at"`
}

func (CreatorAgreementSign) TableName() string {
	return "creator_agreement_signs"
}

// CreatorPortfolio 创作者作品集表
type CreatorPortfolio struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	UserID      uint           `gorm:"index;not null" json:"user_id"`    // 用户ID
	Name        string         `gorm:"size:100;not null" json:"name"`    // 作品集名称
	Description string         `gorm:"size:500" json:"description"`      // 作品集描述
	CoverImage  string         `gorm:"size:255" json:"cover_image"`      // 封面图
	IsPublic    bool           `gorm:"default:true" json:"is_public"`    // 是否公开
	ViewCount   int            `gorm:"default:0" json:"view_count"`      // 浏览次数
	LikeCount   int            `gorm:"default:0" json:"like_count"`      // 点赞次数
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (CreatorPortfolio) TableName() string {
	return "creator_portfolios"
}

// CreatorPortfolioItem 作品集内容关联表
type CreatorPortfolioItem struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	PortfolioID uint      `gorm:"uniqueIndex:idx_portfolio_content;not null" json:"portfolio_id"` // 作品集ID
	ContentType string    `gorm:"uniqueIndex:idx_portfolio_content;size:50;not null" json:"content_type"` // 内容类型
	ContentID   uint      `gorm:"uniqueIndex:idx_portfolio_content;not null" json:"content_id"`   // 内容ID
	SortOrder   int       `gorm:"default:0" json:"sort_order"`                                    // 排序
	CreatedAt   time.Time `json:"created_at"`
}

func (CreatorPortfolioItem) TableName() string {
	return "creator_portfolio_items"
}

// ContentVersion 内容版本表
type ContentVersion struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	ContentType string    `gorm:"size:50;index;not null" json:"content_type"` // 内容类型
	ContentID   uint      `gorm:"index;not null" json:"content_id"`           // 内容ID
	Version     int       `gorm:"not null" json:"version"`                    // 版本号
	Data        string    `gorm:"type:text" json:"data"`                      // 版本数据(JSON)
	CreatorID   uint      `gorm:"index" json:"creator_id"`                    // 创建者ID
	Remark      string    `gorm:"size:200" json:"remark"`                     // 版本备注
	CreatedAt   time.Time `json:"created_at"`
}

func (ContentVersion) TableName() string {
	return "content_versions"
}

// ContentTag 内容标签表
type ContentTag struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `gorm:"uniqueIndex;size:50;not null" json:"name"` // 标签名称
	Slug        string         `gorm:"uniqueIndex;size:50;not null" json:"slug"` // 标签别名
	Description string         `gorm:"size:200" json:"description"`              // 标签描述
	Color       string         `gorm:"size:20" json:"color"`                     // 标签颜色
	UsageCount  int            `gorm:"default:0" json:"usage_count"`             // 使用次数
	IsActive    bool           `gorm:"default:true" json:"is_active"`            // 是否启用
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ContentTag) TableName() string {
	return "content_tags"
}

// ContentTagRelation 内容标签关联表
type ContentTagRelation struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	TagID       uint      `gorm:"uniqueIndex:idx_tag_content;not null" json:"tag_id"`        // 标签ID
	ContentType string    `gorm:"uniqueIndex:idx_tag_content;size:50;not null" json:"content_type"` // 内容类型
	ContentID   uint      `gorm:"uniqueIndex:idx_tag_content;not null" json:"content_id"`    // 内容ID
	CreatedAt   time.Time `json:"created_at"`
}

func (ContentTagRelation) TableName() string {
	return "content_tag_relations"
}

// ContentFavorite 内容收藏表
type ContentFavorite struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	UserID      uint      `gorm:"uniqueIndex:idx_user_content_folder;not null" json:"user_id"`        // 用户ID
	ContentType string    `gorm:"uniqueIndex:idx_user_content_folder;size:50;not null" json:"content_type"` // 内容类型
	ContentID   uint      `gorm:"uniqueIndex:idx_user_content_folder;not null" json:"content_id"`      // 内容ID
	FolderName  string    `gorm:"uniqueIndex:idx_user_content_folder;size:100;default:default" json:"folder_name"` // 收藏夹名称
	CreatedAt   time.Time `json:"created_at"`
}

func (ContentFavorite) TableName() string {
	return "content_favorites"
}

// ContentSearchLog 内容搜索历史表
type ContentSearchLog struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	UserID      uint      `gorm:"index;not null" json:"user_id"`    // 用户ID
	Keyword     string    `gorm:"size:200;not null" json:"keyword"` // 搜索关键词
	ResultCount int       `gorm:"default:0" json:"result_count"`    // 搜索结果数
	SearchedAt  time.Time `json:"searched_at"`                      // 搜索时间
	CreatedAt   time.Time `json:"created_at"`
}

func (ContentSearchLog) TableName() string {
	return "content_search_logs"
}

// ReviewWorkflow 审核流程模板表
type ReviewWorkflow struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`        // 流程名称
	Description string         `gorm:"size:500" json:"description"`          // 流程描述
	ContentType string         `gorm:"size:50;index;not null" json:"content_type"` // 适用内容类型
	Steps       string         `gorm:"type:text" json:"steps"`               // 流程步骤(JSON)
	IsActive    bool           `gorm:"default:true" json:"is_active"`        // 是否启用
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ReviewWorkflow) TableName() string {
	return "review_workflows"
}

// ReviewInstance 审核流程实例表
type ReviewInstance struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	WorkflowID  uint           `gorm:"index;not null" json:"workflow_id"`                                       // 流程模板ID
	ContentType string         `gorm:"size:50;index;not null" json:"content_type"`                             // 内容类型
	ContentID   uint           `gorm:"index;not null" json:"content_id"`                                       // 内容ID
	CurrentStep int            `gorm:"default:0" json:"current_step"`                                          // 当前步骤
	Status      int8           `gorm:"default:0;comment:0=待审核,1=通过,2=驳回,3=已取消" json:"status"`          // 审核状态
	CreatorID   uint           `gorm:"index" json:"creator_id"`                                                // 提交者ID
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ReviewInstance) TableName() string {
	return "review_instances"
}

// ReviewStepLog 审核步骤记录表
type ReviewStepLog struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	InstanceID uint      `gorm:"index;not null" json:"instance_id"`                              // 审核实例ID
	StepIndex  int       `gorm:"not null" json:"step_index"`                                     // 步骤序号
	StepName   string    `gorm:"size:100" json:"step_name"`                                      // 步骤名称
	ReviewerID uint      `gorm:"index" json:"reviewer_id"`                                       // 审核人ID
	Action     string    `gorm:"size:20;not null;comment:pass/reject/return" json:"action"`      // 操作动作
	Opinion    string    `gorm:"size:500" json:"opinion"`                                        // 审核意见
	OperatedAt time.Time `json:"operated_at"`                                                    // 操作时间
	CreatedAt  time.Time `json:"created_at"`
}

func (ReviewStepLog) TableName() string {
	return "review_step_logs"
}

// ReviewNotification 审核通知表
type ReviewNotification struct {
	ID         uint       `gorm:"primarykey" json:"id"`
	InstanceID uint       `gorm:"index;not null" json:"instance_id"` // 审核实例ID
	UserID     uint       `gorm:"index;not null" json:"user_id"`     // 接收用户ID
	Type       string     `gorm:"size:50;not null" json:"type"`       // 通知类型
	Message    string     `gorm:"size:500" json:"message"`            // 通知内容
	IsRead     bool       `gorm:"default:false" json:"is_read"`      // 是否已读
	ReadAt     *time.Time `json:"read_at"`                            // 阅读时间
	CreatedAt  time.Time  `json:"created_at"`
}

func (ReviewNotification) TableName() string {
	return "review_notifications"
}

// ReviewReply 审核意见回复表
type ReviewReply struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	InstanceID uint      `gorm:"index;not null" json:"instance_id"` // 审核实例ID
	StepLogID  uint      `gorm:"index;not null" json:"step_log_id"` // 审核步骤记录ID
	UserID     uint      `gorm:"index;not null" json:"user_id"`     // 回复用户ID
	Content    string    `gorm:"type:text;not null" json:"content"`  // 回复内容
	ReplyAt    time.Time `json:"reply_at"`                           // 回复时间
	CreatedAt  time.Time `json:"created_at"`
}

func (ReviewReply) TableName() string {
	return "review_replies"
}
