package job

import (
	"context"
	"encoding/json"

	"questionhelper-server/pkg/logger"
	"questionhelper-server/pkg/mq"
)

// ImportJob 异步导入任务
type ImportJob struct {
	mq mq.MQ
}

func NewImportJob(mq mq.MQ) *ImportJob {
	return &ImportJob{mq: mq}
}

// ImportMessage 导入消息
type ImportMessage struct {
	UserID     uint   `json:"user_id"`
	FileID     uint   `json:"file_id"`
	CategoryID uint   `json:"category_id"`
	Visibility int8   `json:"visibility"`
}

// Publish 发布导入任务
func (j *ImportJob) Publish(ctx context.Context, msg ImportMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return j.mq.Publish(ctx, "job:import", data)
}

// Handle 处理导入任务
func (j *ImportJob) Handle(ctx context.Context, msg *mq.Message) error {
	var importMsg ImportMessage
	if err := json.Unmarshal(msg.Body, &importMsg); err != nil {
		logger.Errorf("解析导入消息失败: %v", err)
		return err
	}

	logger.Infof("处理导入任务: user_id=%d, file_id=%d", importMsg.UserID, importMsg.FileID)

	// TODO: 读取文件并解析题目
	// TODO: 批量创建题目
	// TODO: 发送导入完成通知

	return nil
}
