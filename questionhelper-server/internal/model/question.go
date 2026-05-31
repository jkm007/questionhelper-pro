package model

import (
	"time"

	"gorm.io/gorm"
)

type Question struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	Title         string         `gorm:"type:text;not null" json:"title"`
	Content       string         `gorm:"type:text" json:"content"`
	ContentHash   string         `gorm:"size:64;index" json:"-"`                          // 内容MD5哈希(查重)
	Type          int8           `gorm:"not null;comment:题型:1=单选,2=多选,3=判断,4=填空,5=简答" json:"type"`
	IsSubjective  bool           `gorm:"default:false" json:"is_subjective"`              // 是否主观题
	Difficulty    int8           `gorm:"default:2;comment:难度:1=简单,2=中等,3=困难" json:"difficulty"`
	Answer        string         `gorm:"type:text" json:"answer"`
	Analysis      string         `gorm:"type:text" json:"analysis"`
	CategoryID    uint           `gorm:"index" json:"category_id"`
	Category      Category       `json:"category,omitempty"`
	KnowledgeIDs  string         `gorm:"size:500" json:"knowledge_ids"`
	Options       []Option       `gorm:"foreignKey:QuestionID" json:"options,omitempty"`
	Visibility    int8           `gorm:"default:1;comment:可见性:1=公开,2=私有,3=班级" json:"visibility"`
	ClassID       *uint          `gorm:"index" json:"class_id,omitempty"`                 // 班级ID(班级可见时)
	CreatorID     uint           `gorm:"index" json:"creator_id"`
	Creator       User           `json:"creator,omitempty"`
	Version       int            `gorm:"default:1" json:"version"`                        // 版本号
	Status        int8           `gorm:"default:0;comment:状态:0=草稿,1=已发布,2=已归档,3=待审核,4=需修改,5=审核超时" json:"status"`
	ViewCount     int            `gorm:"default:0" json:"view_count"`
	LikeCount     int            `gorm:"default:0" json:"like_count"`
	FavoriteCount int            `gorm:"default:0" json:"favorite_count"`
	AnswerCount   int            `gorm:"default:0" json:"answer_count"`                   // 答题次数
	CorrectRate   float64        `gorm:"default:0" json:"correct_rate"`                   // 正确率
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Question) TableName() string {
	return "questions"
}

type Option struct {
	ID         uint   `gorm:"primarykey" json:"id"`
	QuestionID uint   `gorm:"index;not null" json:"question_id"`
	Label      string `gorm:"size:10;not null" json:"label"` // A/B/C/D
	Content    string `gorm:"type:text;not null" json:"content"`
	IsCorrect  bool   `gorm:"default:false" json:"is_correct"`
	Sort       int    `gorm:"default:0" json:"sort"`
}

func (Option) TableName() string {
	return "options"
}

type Category struct {
	ID          uint       `gorm:"primarykey" json:"id"`
	ParentID    *uint      `gorm:"index" json:"parent_id"`
	Name        string     `gorm:"size:100;not null" json:"name"`
	Code        string     `gorm:"size:50;index" json:"code"`              // 分类编码
	Path        string     `gorm:"size:500;index" json:"path"`             // 路径(如 /1/2/3)
	Icon        string     `gorm:"size:50" json:"icon"`
	Color       string     `gorm:"size:20" json:"color"`                   // 颜色标识
	Description string     `gorm:"size:200" json:"description"`
	Sort        int        `gorm:"default:0" json:"sort"`
	Status      int8       `gorm:"default:1" json:"status"`
	QuestionCount int      `gorm:"default:0" json:"question_count"`        // 题目数量
	Children    []Category `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

func (Category) TableName() string {
	return "categories"
}

type Knowledge struct {
	ID            uint    `gorm:"primarykey" json:"id"`
	CategoryID    uint    `gorm:"index" json:"category_id"`
	ParentID      *uint   `gorm:"index" json:"parent_id"`                  // 父知识点(支持树形)
	Name          string  `gorm:"size:100;not null" json:"name"`
	Code          string  `gorm:"size:50;index" json:"code"`               // 知识点编码
	Description   string  `gorm:"size:200" json:"description"`
	Weight        int     `gorm:"default:1" json:"weight"`                 // 权重(重要程度)
	Sort          int     `gorm:"default:0" json:"sort"`
	Status        int8    `gorm:"default:1" json:"status"`
	QuestionCount int     `gorm:"default:0" json:"question_count"`         // 关联题目数
	UsageCount    int     `gorm:"default:0" json:"usage_count"`            // 使用次数
	CorrectRate   float64 `gorm:"default:0" json:"correct_rate"`           // 平均正确率
}

func (Knowledge) TableName() string {
	return "knowledge_points"
}
