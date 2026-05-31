package upload

import (
	"context"
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"questionhelper-server/pkg/config"
)

type OSSClient struct {
	client *minio.Client
	bucket string
	cdn    string
}

func NewOSSClient(cfg config.OSSConfig) (*OSSClient, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, fmt.Errorf("创建 OSS 客户端失败: %w", err)
	}

	return &OSSClient{
		client: client,
		bucket: cfg.Bucket,
		cdn:    cfg.CDN,
	}, nil
}

func (o *OSSClient) Upload(ctx context.Context, file *multipart.FileHeader, objectName string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("打开文件失败: %w", err)
	}
	defer src.Close()

	contentType := file.Header.Get("Content-Type")
	_, err = o.client.PutObject(ctx, o.bucket, objectName, src, file.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("上传文件失败: %w", err)
	}

	url := fmt.Sprintf("%s/%s", o.cdn, objectName)
	if !strings.HasPrefix(url, "http") {
		url = fmt.Sprintf("http://%s/%s/%s", o.client.EndpointURL().Host, o.bucket, objectName)
	}

	return url, nil
}

func (o *OSSClient) Delete(ctx context.Context, objectName string) error {
	return o.client.RemoveObject(ctx, o.bucket, objectName, minio.RemoveObjectOptions{})
}
