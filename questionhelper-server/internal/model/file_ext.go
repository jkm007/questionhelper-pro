package model

import (
	"time"

	"gorm.io/gorm"
)

// FileReference 文件引用表，记录文件与业务数据之间的关联关系
type FileReference struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	FileID       uint           `gorm:"index;not null;comment:文件ID" json:"file_id"`
	File         File           `gorm:"foreignKey:FileID" json:"file,omitempty"`
	BusinessType string         `gorm:"size:50;not null;index;comment:业务类型(如question/exam/class)" json:"business_type"`
	BusinessID   uint           `gorm:"not null;index;comment:业务记录ID" json:"business_id"`
	FieldName    string         `gorm:"size:50;comment:业务字段名" json:"field_name"`
	CreatorID    uint           `gorm:"index;comment:创建人ID" json:"creator_id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (FileReference) TableName() string {
	return "file_references"
}

// FileAccessLog 文件访问日志表，记录文件的查看和下载行为
type FileAccessLog struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	FileID      uint      `gorm:"index;not null;comment:文件ID" json:"file_id"`
	File        File      `gorm:"foreignKey:FileID" json:"file,omitempty"`
	UserID      uint      `gorm:"index;comment:访问用户ID(未登录为0)" json:"user_id"`
	AccessType  string    `gorm:"size:20;not null;comment:访问类型:view=预览,download=下载" json:"access_type"`
	IP          string    `gorm:"size:45;comment:访问者IP" json:"ip"`
	UserAgent   string    `gorm:"size:500;comment:浏览器UserAgent" json:"user_agent"`
	AccessedAt  time.Time `gorm:"index;not null;comment:访问时间" json:"accessed_at"`
	CreatedAt   time.Time `json:"created_at"`
}

func (FileAccessLog) TableName() string {
	return "file_access_logs"
}

// HotlinkProtectionRule 防盗链配置表，控制文件访问来源限制
type HotlinkProtectionRule struct {
	ID             uint           `gorm:"primarykey" json:"id"`
	Name           string         `gorm:"size:100;not null;comment:规则名称" json:"name"`
	AllowedDomains string         `gorm:"type:text;comment:允许的域名列表(JSON数组)" json:"allowed_domains"`
	BlockedDomains string         `gorm:"type:text;comment:禁止的域名列表(JSON数组)" json:"blocked_domains"`
	IsActive       bool           `gorm:"default:true;comment:是否启用" json:"is_active"`
	DefaultAction  string         `gorm:"size:20;default:'allow';comment:默认动作:allow=放行,block=拦截" json:"default_action"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (HotlinkProtectionRule) TableName() string {
	return "hotlink_protection_rules"
}

// WatermarkConfig 水印配置表，定义文件水印的样式和行为
type WatermarkConfig struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	Name       string         `gorm:"size:100;not null;comment:配置名称" json:"name"`
	Type       string         `gorm:"size:20;not null;comment:水印类型:text=文字,image=图片" json:"type"`
	Position   string         `gorm:"size:20;default:'bottom-right';comment:水印位置:top-left/top-right/bottom-left/bottom-right/center" json:"position"`
	Opacity    float32        `gorm:"default:0.5;comment:透明度(0~1)" json:"opacity"`
	Text       string         `gorm:"size:200;comment:文字内容(文字水印时使用)" json:"text"`
	ImagePath  string         `gorm:"size:500;comment:水印图片路径(图片水印时使用)" json:"image_path"`
	FontSize   int            `gorm:"default:14;comment:字体大小(像素)" json:"font_size"`
	FontColor  string         `gorm:"size:20;default:'#000000';comment:字体颜色(十六进制)" json:"font_color"`
	IsActive   bool           `gorm:"default:true;comment:是否启用" json:"is_active"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (WatermarkConfig) TableName() string {
	return "watermark_configs"
}

// FileCleanupLog 孤立文件清理记录表，记录定期清理无引用文件的操作日志
type FileCleanupLog struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	CleanedCount int            `gorm:"default:0;comment:清理文件数量" json:"cleaned_count"`
	FreedSpace   int64          `gorm:"default:0;comment:释放空间(字节)" json:"freed_space"`
	FileIDs      string         `gorm:"type:text;comment:已清理文件ID列表(JSON数组)" json:"file_ids"`
	Status       string         `gorm:"size:20;default:'pending';comment:状态:pending=待执行,running=执行中,completed=已完成,failed=失败" json:"status"`
	StartedAt    *time.Time     `gorm:"comment:开始时间" json:"started_at"`
	CompletedAt  *time.Time     `gorm:"comment:完成时间" json:"completed_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (FileCleanupLog) TableName() string {
	return "file_cleanup_logs"
}
