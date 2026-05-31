package model

import (
	"time"

	"gorm.io/gorm"
)

// CommentBlacklist 评论黑名单，被拉黑的用户无法在指定目标下发表评论
type CommentBlacklist struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	UserID     uint           `gorm:"uniqueIndex:idx_blacklist_user_target;not null;comment:被拉黑用户ID" json:"user_id"`
	User       User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	TargetType int8           `gorm:"uniqueIndex:idx_blacklist_user_target;not null;comment:目标类型:1=题目,2=考试,3=班级,0=全局" json:"target_type"`
	TargetID   uint           `gorm:"uniqueIndex:idx_blacklist_user_target;not null;comment:目标ID(全局拉黑时为0)" json:"target_id"`
	Reason     string         `gorm:"size:200;comment:拉黑原因" json:"reason"`
	OperatorID uint           `gorm:"comment:操作人ID" json:"operator_id"`
	ExpiresAt  *time.Time     `gorm:"comment:过期时间(空表示永久)" json:"expires_at"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (CommentBlacklist) TableName() string {
	return "comment_blacklists"
}

// Sticker 表情包/贴纸，用户可在评论中使用
type Sticker struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Name      string         `gorm:"size:50;not null;comment:表情名称" json:"name"`
	Category  string         `gorm:"size:30;index;comment:所属分类" json:"category"`
	ImageURL  string         `gorm:"size:500;not null;comment:表情图片URL" json:"image_url"`
	SortOrder int            `gorm:"default:0;comment:排序值(越大越靠前)" json:"sort_order"`
	Status    int8           `gorm:"default:1;comment:状态:0=禁用,1=启用" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Sticker) TableName() string {
	return "stickers"
}

// CommentAuditRule 评论自动审核规则，用于敏感内容过滤和自动处理
type CommentAuditRule struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `gorm:"size:100;not null;comment:规则名称" json:"name"`
	RuleType    string         `gorm:"size:30;not null;comment:规则类型:keyword=关键词,regex=正则,length=长度,repeat=重复" json:"rule_type"`
	Pattern     string         `gorm:"type:text;not null;comment:匹配模式(关键词/正则表达式/数值)" json:"pattern"`
	Action      int8           `gorm:"default:1;comment:处理动作:1=标记待审,2=自动隐藏,3=自动删除,4=拒绝发布" json:"action"`
	Priority    int            `gorm:"default:0;comment:优先级(数值越大优先级越高)" json:"priority"`
	Description string         `gorm:"size:200;comment:规则描述" json:"description"`
	Status      int8           `gorm:"default:1;comment:状态:0=禁用,1=启用" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (CommentAuditRule) TableName() string {
	return "comment_audit_rules"
}
