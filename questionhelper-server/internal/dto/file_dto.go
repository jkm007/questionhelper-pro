package dto

// ==================== File Upload ====================

// ImageUploadResponse 图片上传响应
type ImageUploadResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Original  string `json:"original"`
	URL       string `json:"url"`
	Thumbnail string `json:"thumbnail"`
	Size      int64  `json:"size"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
}

// BatchUploadResponse 批量上传响应
type BatchUploadResponse struct {
	Success []FileUploadResult `json:"success"`
	Failed  []FileUploadError  `json:"failed"`
	Total   int                `json:"total"`
}

// FileUploadResult 单个文件上传结果
type FileUploadResult struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Original string `json:"original"`
	URL      string `json:"url"`
	Size     int64  `json:"size"`
}

// FileUploadError 单个文件上传失败
type FileUploadError struct {
	Filename string `json:"filename"`
	Error    string `json:"error"`
}

// ==================== File Reference ====================

// AddReferenceRequest 添加引用请求
type AddReferenceRequest struct {
	BusinessType string `json:"business_type" binding:"required"`
	BusinessID   uint   `json:"business_id" binding:"required"`
	FieldName    string `json:"field_name"`
}

// FileReferenceInfo 文件引用信息
type FileReferenceInfo struct {
	ID           uint   `json:"id"`
	FileID       uint   `json:"file_id"`
	BusinessType string `json:"business_type"`
	BusinessID   uint   `json:"business_id"`
	FieldName    string `json:"field_name"`
	CreatorID    uint   `json:"creator_id"`
	CreatedAt    string `json:"created_at"`
}

// ==================== File Access Log ====================

// FileAccessLogInfo 文件访问日志信息
type FileAccessLogInfo struct {
	ID         uint   `json:"id"`
	FileID     uint   `json:"file_id"`
	UserID     uint   `json:"user_id"`
	AccessType string `json:"access_type"`
	IP         string `json:"ip"`
	UserAgent  string `json:"user_agent"`
	AccessedAt string `json:"accessed_at"`
}

// ==================== Hotlink Protection ====================

// HotlinkRuleRequest 防盗链规则请求
type HotlinkRuleRequest struct {
	Name           string   `json:"name" binding:"required"`
	AllowedDomains []string `json:"allowed_domains"`
	BlockedDomains []string `json:"blocked_domains"`
	IsActive       *bool    `json:"is_active"`
	DefaultAction  string   `json:"default_action"`
}

// HotlinkRuleInfo 防盗链规则信息
type HotlinkRuleInfo struct {
	ID             uint     `json:"id"`
	Name           string   `json:"name"`
	AllowedDomains []string `json:"allowed_domains"`
	BlockedDomains []string `json:"blocked_domains"`
	IsActive       bool     `json:"is_active"`
	DefaultAction  string   `json:"default_action"`
	CreatedAt      string   `json:"created_at"`
	UpdatedAt      string   `json:"updated_at"`
}

// ==================== Watermark Config ====================

// WatermarkConfigRequest 水印配置请求
type WatermarkConfigRequest struct {
	Name      string  `json:"name" binding:"required"`
	Type      string  `json:"type" binding:"required"`
	Position  string  `json:"position"`
	Opacity   float32 `json:"opacity"`
	Text      string  `json:"text"`
	ImagePath string  `json:"image_path"`
	FontSize  int     `json:"font_size"`
	FontColor string  `json:"font_color"`
	IsActive  *bool   `json:"is_active"`
}

// WatermarkConfigInfo 水印配置信息
type WatermarkConfigInfo struct {
	ID        uint    `json:"id"`
	Name      string  `json:"name"`
	Type      string  `json:"type"`
	Position  string  `json:"position"`
	Opacity   float32 `json:"opacity"`
	Text      string  `json:"text"`
	ImagePath string  `json:"image_path"`
	FontSize  int     `json:"font_size"`
	FontColor string  `json:"font_color"`
	IsActive  bool    `json:"is_active"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

// ==================== File Cleanup ====================

// CleanupOrphanResponse 孤立文件清理响应
type CleanupOrphanResponse struct {
	CleanedCount int    `json:"cleaned_count"`
	FreedSpace   int64  `json:"freed_space"`
	LogID        uint   `json:"log_id"`
}

// CleanupLogInfo 清理日志信息
type CleanupLogInfo struct {
	ID           uint   `json:"id"`
	CleanedCount int    `json:"cleaned_count"`
	FreedSpace   int64  `json:"freed_space"`
	FileIDs      string `json:"file_ids"`
	Status       string `json:"status"`
	StartedAt    string `json:"started_at"`
	CompletedAt  string `json:"completed_at"`
	CreatedAt    string `json:"created_at"`
}

// ==================== Storage Statistics ====================

// StorageStatistics 存储统计信息
type StorageStatistics struct {
	TotalFiles     int64            `json:"total_files"`
	TotalSize      int64            `json:"total_size"`
	ImageCount     int64            `json:"image_count"`
	ImageSize      int64            `json:"image_size"`
	DocumentCount  int64            `json:"document_count"`
	DocumentSize   int64            `json:"document_size"`
	VideoCount     int64            `json:"video_count"`
	VideoSize      int64            `json:"video_size"`
	AudioCount     int64            `json:"audio_count"`
	AudioSize      int64            `json:"audio_size"`
	OtherCount     int64            `json:"other_count"`
	OtherSize      int64            `json:"other_size"`
	OrphanCount    int64            `json:"orphan_count"`
	OrphanSize     int64            `json:"orphan_size"`
	DailyUploads   []DailyUploadStat `json:"daily_uploads"`
	TopUploaders   []UploaderStat    `json:"top_uploaders"`
}

// DailyUploadStat 每日上传统计
type DailyUploadStat struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
	Size  int64  `json:"size"`
}

// UploaderStat 上传者统计
type UploaderStat struct {
	UserID    uint   `json:"user_id"`
	Count     int64  `json:"count"`
	TotalSize int64  `json:"total_size"`
}
