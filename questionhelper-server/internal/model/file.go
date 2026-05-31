package model

import (
	"time"

	"gorm.io/gorm"
)

type File struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	Name       string         `gorm:"size:255;not null" json:"name"`
	Original   string         `gorm:"size:255;not null" json:"original"`
	Path       string         `gorm:"size:500;not null" json:"path"`
	URL        string         `gorm:"size:500" json:"url"`
	Size       int64          `gorm:"not null" json:"size"`
	Type       string         `gorm:"size:50" json:"type"`
	Extension  string         `gorm:"size:20" json:"extension"`
	UploaderID uint           `gorm:"index" json:"uploader_id"`
	CreatedAt  time.Time      `json:"created_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (File) TableName() string {
	return "files"
}
