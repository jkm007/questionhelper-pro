package model

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	TargetType   int8           `gorm:"not null;comment:目标类型:1=题目,2=考试,3=班级" json:"target_type"`
	TargetID     uint           `gorm:"index;not null" json:"target_id"`
	UserID       uint           `gorm:"index;not null" json:"user_id"`
	User         User           `json:"user,omitempty"`
	Content      string         `gorm:"type:text;not null" json:"content"`
	ParentID     *uint          `gorm:"index" json:"parent_id"`
	Images       string         `gorm:"type:text;comment:图片URL列表(JSON数组)" json:"images"`
	Mentions     string         `gorm:"type:text;comment:@的用户ID列表(JSON数组)" json:"mentions"`
	ReplyToUserID *uint         `gorm:"index;comment:回复的用户ID" json:"reply_to_user_id"`
	LikeCount    int            `gorm:"default:0" json:"like_count"`
	IsPinned     bool           `gorm:"default:false;comment:是否置顶" json:"is_pinned"`
	IsFeatured   bool           `gorm:"default:false;comment:是否精选" json:"is_featured"`
	IsOfficial   bool           `gorm:"default:false;comment:是否官方解答" json:"is_official"`
	Status       int8           `gorm:"default:1;comment:状态:0=隐藏,1=正常,2=举报" json:"status"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Comment) TableName() string {
	return "comments"
}

type CommentLike struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CommentID uint      `gorm:"uniqueIndex:idx_comment_user;not null" json:"comment_id"`
	UserID    uint      `gorm:"uniqueIndex:idx_comment_user;not null" json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (CommentLike) TableName() string {
	return "comment_likes"
}

type CommentReport struct {
	ID          uint       `gorm:"primarykey" json:"id"`
	CommentID   uint       `gorm:"index;not null" json:"comment_id"`
	UserID      uint       `gorm:"index;not null" json:"user_id"`
	Reason      string     `gorm:"size:200;not null" json:"reason"`
	ReasonType  string     `gorm:"size:30;comment:举报原因类型:spam=垃圾信息,abuse=辱骂,infringement=侵权,other=其他" json:"reason_type"`
	Status      int8       `gorm:"default:0;comment:状态:0=待处理,1=已处理,2=已驳回" json:"status"`
	HandlerID   *uint      `gorm:"index;comment:处理人ID" json:"handler_id"`
	HandleRemark string    `gorm:"size:500;comment:处理备注" json:"handle_remark"`
	HandledAt   *time.Time `gorm:"comment:处理时间" json:"handled_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (CommentReport) TableName() string {
	return "comment_reports"
}
