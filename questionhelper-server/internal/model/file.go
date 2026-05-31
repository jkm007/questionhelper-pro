package model

import (
	"time"

	"gorm.io/gorm"
)

// File 文件表，存储上传文件的元数据信息
type File struct {
	ID                uint           `gorm:"primarykey" json:"id"`
	Name              string         `gorm:"size:255;not null" json:"name"`
	Original          string         `gorm:"size:255;not null" json:"original"`
	Path              string         `gorm:"size:500;not null" json:"path"`
	URL               string         `gorm:"size:500" json:"url"`
	Size              int64          `gorm:"not null" json:"size"`
	Type              string         `gorm:"size:50" json:"type"`
	Extension         string         `gorm:"size:20" json:"extension"`
	MD5               string         `gorm:"size:32;index;comment:文件MD5哈希，用于去重" json:"md5"`
	BusinessType      string         `gorm:"size:50;index;comment:业务类型(如question/exam/class)" json:"business_type"`
	BusinessID        uint           `gorm:"index;comment:业务记录ID" json:"business_id"`
	StorageType       string         `gorm:"size:20;default:'local';comment:存储类型:local/minio/oss" json:"storage_type"`
	Bucket            string         `gorm:"size:100;comment:存储桶名称" json:"bucket"`
	Width             int            `gorm:"comment:图片宽度(像素)" json:"width"`
	Height            int            `gorm:"comment:图片高度(像素)" json:"height"`
	Duration          int            `gorm:"comment:视频/音频时长(秒)" json:"duration"`
	ThumbnailIDs      string         `gorm:"type:text;comment:缩略图ID列表(JSON数组)" json:"thumbnail_ids"`
	ReferenceCount    int            `gorm:"default:0;comment:被引用次数" json:"reference_count"`
	VirusScanStatus   string         `gorm:"size:20;default:'pending';comment:病毒扫描状态:pending/scanning/clean/infected" json:"virus_scan_status"`
	ContentCheckStatus string        `gorm:"size:20;default:'pending';comment:内容审核状态:pending/checking/approved/rejected" json:"content_check_status"`
	WatermarkApplied  bool           `gorm:"default:false;comment:是否已加水印" json:"watermark_applied"`
	IsPublic          bool           `gorm:"default:true;comment:是否公开访问" json:"is_public"`
	Status            string         `gorm:"size:20;default:'active';comment:状态:active=正常,disabled=禁用,deleted=已删除" json:"status"`
	UploaderID        uint           `gorm:"index" json:"uploader_id"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

func (File) TableName() string {
	return "files"
}
