package upload

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofrs/uuid"
)

// UploadResult 上传结果
type UploadResult struct {
	URL      string `json:"url"`
	Path     string `json:"path"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
}

// AllowedExts 允许的文件扩展名
var AllowedExts = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
	".pdf":  true,
	".doc":  true,
	".docx": true,
	".xls":  true,
	".xlsx": true,
	".ppt":  true,
	".pptx": true,
	".zip":  true,
	".rar":  true,
}

// ValidateFile 验证文件
func ValidateFile(file *multipart.FileHeader, maxSize int64) error {
	// 检查文件大小
	if file.Size > maxSize {
		return fmt.Errorf("文件大小超过限制: %dMB", maxSize/1024/1024)
	}

	// 检查文件扩展名
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !AllowedExts[ext] {
		return fmt.Errorf("不支持的文件类型: %s", ext)
	}

	return nil
}

// GenerateFilename 生成唯一文件名
func GenerateFilename(ext string) string {
	id, _ := uuid.NewV4()
	return fmt.Sprintf("%s/%s%s", time.Now().Format("2006/01/02"), id.String(), ext)
}
