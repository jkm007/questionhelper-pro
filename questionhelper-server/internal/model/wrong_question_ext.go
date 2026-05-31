package model

import (
	"time"

	"gorm.io/gorm"
)

// WrongQuestionTag 错题标签表
type WrongQuestionTag struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	UserID    uint           `gorm:"index;not null" json:"user_id"`          // 标签所属用户
	Name      string         `gorm:"size:50;not null" json:"name"`           // 标签名称
	Color     string         `gorm:"size:20" json:"color"`                   // 标签颜色
	Icon      string         `gorm:"size:50" json:"icon"`                    // 标签图标
	Sort      int            `gorm:"default:0" json:"sort"`                  // 排序
	Count     int            `gorm:"default:0" json:"count"`                 // 关联错题数量
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (WrongQuestionTag) TableName() string {
	return "wrong_question_tags"
}

// WrongQuestionTagRelation 错题标签关联表
type WrongQuestionTagRelation struct {
	ID               uint      `gorm:"primarykey" json:"id"`
	WrongQuestionID  uint      `gorm:"uniqueIndex:idx_wq_tag;not null" json:"wrong_question_id"`
	TagID            uint      `gorm:"uniqueIndex:idx_wq_tag;not null" json:"tag_id"`
	CreatedAt        time.Time `json:"created_at"`
}

func (WrongQuestionTagRelation) TableName() string {
	return "wrong_question_tag_relations"
}

// WrongQuestionNote 错题备注表
type WrongQuestionNote struct {
	ID              uint      `gorm:"primarykey" json:"id"`
	WrongQuestionID uint      `gorm:"index;not null" json:"wrong_question_id"`
	UserID          uint      `gorm:"index;not null" json:"user_id"`
	Content         string    `gorm:"type:text;not null" json:"content"` // 备注内容
	IsPinned        bool      `gorm:"default:false" json:"is_pinned"`    // 是否置顶
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (WrongQuestionNote) TableName() string {
	return "wrong_question_notes"
}

// WrongQuestionAttachment 错题附件表
type WrongQuestionAttachment struct {
	ID              uint      `gorm:"primarykey" json:"id"`
	WrongQuestionID uint      `gorm:"index;not null" json:"wrong_question_id"`
	UserID          uint      `gorm:"index;not null" json:"user_id"`
	FileID          uint      `gorm:"index" json:"file_id"`                  // 关联文件表
	FileName        string    `gorm:"size:200;not null" json:"file_name"`
	FileType        string    `gorm:"size:20;not null" json:"file_type"`     // image/audio/video/document
	FileURL         string    `gorm:"size:500;not null" json:"file_url"`
	FileSize        int64     `gorm:"default:0" json:"file_size"`
	Sort            int       `gorm:"default:0" json:"sort"`
	CreatedAt       time.Time `json:"created_at"`
}

func (WrongQuestionAttachment) TableName() string {
	return "wrong_question_attachments"
}

// WrongQuestionFavorite 错题收藏表
type WrongQuestionFavorite struct {
	ID              uint      `gorm:"primarykey" json:"id"`
	UserID          uint      `gorm:"uniqueIndex:idx_user_wq;not null" json:"user_id"`
	WrongQuestionID uint      `gorm:"uniqueIndex:idx_user_wq;not null" json:"wrong_question_id"`
	FolderID        uint      `gorm:"index;default:0" json:"folder_id"` // 收藏夹ID,0=默认
	Note            string    `gorm:"size:500" json:"note"`             // 收藏笔记
	CreatedAt       time.Time `json:"created_at"`
}

func (WrongQuestionFavorite) TableName() string {
	return "wrong_question_favorites"
}

// WrongQuestionReview 错题复习记录表
type WrongQuestionReview struct {
	ID              uint       `gorm:"primarykey" json:"id"`
	WrongQuestionID uint       `gorm:"index;not null" json:"wrong_question_id"`
	UserID          uint       `gorm:"index;not null" json:"user_id"`
	IsCorrect       bool       `json:"is_correct"`                                  // 本次是否回答正确
	Answer          string     `gorm:"type:text" json:"answer"`                     // 本次作答内容
	Duration        int        `gorm:"default:0;comment:用时(秒)" json:"duration"`    // 复习用时
	ReviewType      int8       `gorm:"default:1;comment:复习方式:1=手动,2=计划" json:"review_type"`
	NextReviewAt    *time.Time `json:"next_review_at"`                              // 下次复习时间
	CreatedAt       time.Time  `json:"created_at"`
}

func (WrongQuestionReview) TableName() string {
	return "wrong_question_reviews"
}

// WrongQuestionSync 错题同步记录表
type WrongQuestionSync struct {
	ID          uint       `gorm:"primarykey" json:"id"`
	UserID      uint       `gorm:"index;not null" json:"user_id"`
	DeviceID    string     `gorm:"size:100;not null" json:"device_id"`              // 设备标识
	DeviceType  string     `gorm:"size:20;not null" json:"device_type"`             // 设备类型:web/ios/android
	SyncType    int8       `gorm:"default:1;comment:同步类型:1=上传,2=下载" json:"sync_type"`
	Status      int8       `gorm:"default:0;comment:状态:0=进行中,1=成功,2=失败" json:"status"`
	ItemCount   int        `gorm:"default:0" json:"item_count"`                     // 同步条目数
	ErrorMsg    string     `gorm:"size:500" json:"error_msg"`                       // 错误信息
	StartedAt   time.Time  `json:"started_at"`
	FinishedAt  *time.Time `json:"finished_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

func (WrongQuestionSync) TableName() string {
	return "wrong_question_syncs"
}
