package model

import (
	"time"

	"gorm.io/gorm"
)

type Class struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Description string         `gorm:"size:500" json:"description"`
	Cover       string         `gorm:"size:255" json:"cover"`
	Code        string         `gorm:"uniqueIndex;size:20" json:"code"` // 加入码
	CreatorID   uint           `gorm:"index" json:"creator_id"`
	Creator     User           `json:"creator,omitempty"`
	MemberCount int            `gorm:"default:0" json:"member_count"`
	Status      int8           `gorm:"default:1" json:"status"`
	ExpiresAt   *time.Time     `gorm:"index;comment:班级有效期" json:"expires_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Class) TableName() string {
	return "classes"
}

type ClassMember struct {
	ID       uint      `gorm:"primarykey" json:"id"`
	ClassID  uint      `gorm:"index;not null" json:"class_id"`
	UserID   uint      `gorm:"index;not null" json:"user_id"`
	Role     int8      `gorm:"default:1;comment:角色:1=学生,2=教师,3=管理员" json:"role"`
	IsPinned bool      `gorm:"default:false;comment:是否置顶" json:"is_pinned"`
	JoinedAt time.Time `json:"joined_at"`
}

func (ClassMember) TableName() string {
	return "class_members"
}

type Homework struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	ClassID     uint           `gorm:"index;not null" json:"class_id"`
	Title       string         `gorm:"size:200;not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	Deadline    time.Time      `gorm:"not null" json:"deadline"`
	CreatorID   uint           `gorm:"index" json:"creator_id"`
	RemindSent  bool           `gorm:"default:false" json:"remind_sent"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Homework) TableName() string {
	return "homeworks"
}

type ClassNotice struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	ClassID   uint           `gorm:"index;not null" json:"class_id"`
	Title     string         `gorm:"size:200;not null" json:"title"`
	Content   string         `gorm:"type:text" json:"content"`
	CreatorID uint           `gorm:"index" json:"creator_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ClassNotice) TableName() string {
	return "class_notices"
}
