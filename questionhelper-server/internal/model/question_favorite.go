package model

import (
	"time"
)

// FavoriteFolder 收藏夹表
type FavoriteFolder struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	UserID      uint      `gorm:"index;not null" json:"user_id"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	Description string    `gorm:"size:200" json:"description"`
	Icon        string    `gorm:"size:50" json:"icon"`
	Sort        int       `gorm:"default:0" json:"sort"`
	Count       int       `gorm:"default:0" json:"count"` // 收藏数量
	IsDefault   bool      `gorm:"default:false" json:"is_default"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (FavoriteFolder) TableName() string {
	return "favorite_folders"
}

// QuestionFavorite 题目收藏表
type QuestionFavorite struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	UserID     uint      `gorm:"uniqueIndex:idx_user_question;not null" json:"user_id"`
	QuestionID uint      `gorm:"uniqueIndex:idx_user_question;not null" json:"question_id"`
	FolderID   uint      `gorm:"index;default:0" json:"folder_id"` // 收藏夹ID,0=默认
	Note       string    `gorm:"size:500" json:"note"`             // 收藏笔记
	CreatedAt  time.Time `json:"created_at"`
}

func (QuestionFavorite) TableName() string {
	return "question_favorites"
}

// QuestionLike 题目点赞表（用于防重复点赞）
type QuestionLike struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	UserID     uint      `gorm:"uniqueIndex:idx_user_like_question;not null" json:"user_id"`
	QuestionID uint      `gorm:"uniqueIndex:idx_user_like_question;not null" json:"question_id"`
	CreatedAt  time.Time `json:"created_at"`
}

func (QuestionLike) TableName() string {
	return "question_likes"
}
