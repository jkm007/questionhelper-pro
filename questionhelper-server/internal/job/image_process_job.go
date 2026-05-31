package job

import (
	"context"
	"encoding/json"

	"questionhelper-server/pkg/logger"
	"questionhelper-server/pkg/mq"
)

// ImageProcessJob 异步图片处理任务
type ImageProcessJob struct {
	mq mq.MQ
}

func NewImageProcessJob(mq mq.MQ) *ImageProcessJob {
	return &ImageProcessJob{mq: mq}
}

// ImageProcessMessage 图片处理消息
type ImageProcessMessage struct {
	FileID uint   `json:"file_id"`
	Path   string `json:"path"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Quality int   `json:"quality"`
}

// Publish 发布图片处理任务
func (j *ImageProcessJob) Publish(ctx context.Context, msg ImageProcessMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return j.mq.Publish(ctx, "job:image_process", data)
}

// Handle 处理图片处理任务
func (j *ImageProcessJob) Handle(ctx context.Context, msg *mq.Message) error {
	var processMsg ImageProcessMessage
	if err := json.Unmarshal(msg.Body, &processMsg); err != nil {
		logger.Errorf("解析图片处理消息失败: %v", err)
		return err
	}

	logger.Infof("处理图片处理任务: file_id=%d, path=%s", processMsg.FileID, processMsg.Path)

	// TODO: 压缩图片
	// TODO: 生成缩略图
	// TODO: 更新文件记录

	return nil
}
