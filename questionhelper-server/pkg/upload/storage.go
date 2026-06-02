package upload

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"questionhelper-server/pkg/config"
)

// Storage 文件存储接口
type Storage interface {
	// Save 保存文件，返回可访问的 URL
	Save(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error)
	// Delete 删除文件
	Delete(ctx context.Context, objectName string) error
	// GetURL 获取文件访问 URL
	GetURL(objectName string) string
}

// ==================== Local Storage ====================

// LocalStorage 本地文件存储
type LocalStorage struct {
	baseDir string // 本地存储根目录
	baseURL string // URL 前缀，如 /uploads
}

// NewLocalStorage 创建本地存储
func NewLocalStorage(baseDir, baseURL string) *LocalStorage {
	return &LocalStorage{
		baseDir: baseDir,
		baseURL: baseURL,
	}
}

func (s *LocalStorage) Save(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error) {
	// 确保目录存在
	destPath := filepath.Join(s.baseDir, objectName)
	dir := filepath.Dir(destPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("创建目录失败: %w", err)
	}

	dst, err := os.Create(destPath)
	if err != nil {
		return "", fmt.Errorf("创建文件失败: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, reader); err != nil {
		return "", fmt.Errorf("写入文件失败: %w", err)
	}

	return s.GetURL(objectName), nil
}

func (s *LocalStorage) Delete(ctx context.Context, objectName string) error {
	destPath := filepath.Join(s.baseDir, objectName)
	return os.Remove(destPath)
}

func (s *LocalStorage) GetURL(objectName string) string {
	return s.baseURL + "/" + objectName
}

// ==================== MinIO Storage ====================

// MinIOStorage MinIO 对象存储
type MinIOStorage struct {
	client *minio.Client
	bucket string
	cdn    string
}

// NewMinIOStorage 创建 MinIO 存储
func NewMinIOStorage(cfg config.OSSConfig) (*MinIOStorage, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, fmt.Errorf("创建 MinIO 客户端失败: %w", err)
	}

	// 确保 bucket 存在
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("检查 bucket 失败: %w", err)
	}
	if !exists {
		if err := client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("创建 bucket 失败: %w", err)
		}
	}

	return &MinIOStorage{
		client: client,
		bucket: cfg.Bucket,
		cdn:    cfg.CDN,
	}, nil
}

func (s *MinIOStorage) Save(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error) {
	_, err := s.client.PutObject(ctx, s.bucket, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("上传到 MinIO 失败: %w", err)
	}
	return s.GetURL(objectName), nil
}

func (s *MinIOStorage) Delete(ctx context.Context, objectName string) error {
	return s.client.RemoveObject(ctx, s.bucket, objectName, minio.RemoveObjectOptions{})
}

func (s *MinIOStorage) GetURL(objectName string) string {
	if s.cdn != "" {
		url := fmt.Sprintf("%s/%s", s.cdn, objectName)
		if strings.HasPrefix(url, "http") {
			return url
		}
		return "http://" + url
	}
	return fmt.Sprintf("http://%s/%s/%s", s.client.EndpointURL().Host, s.bucket, objectName)
}

// ==================== Storage Factory ====================

// NewStorage 根据配置创建存储实例
func NewStorage(cfg config.OSSConfig) (Storage, error) {
	switch strings.ToLower(cfg.Type) {
	case "minio":
		return NewMinIOStorage(cfg)
	case "local", "":
		return NewLocalStorage("./uploads", "/uploads"), nil
	default:
		return NewLocalStorage("./uploads", "/uploads"), nil
	}
}

// SaveFromMultipart 从 multipart.FileHeader 保存文件的便捷方法
func SaveFromMultipart(ctx context.Context, storage Storage, file *multipart.FileHeader, objectName string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("打开文件失败: %w", err)
	}
	defer src.Close()

	contentType := file.Header.Get("Content-Type")
	return storage.Save(ctx, objectName, src, file.Size, contentType)
}
