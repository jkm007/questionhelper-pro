package job

import (
	"context"
	"encoding/json"

	"questionhelper-server/pkg/logger"
	"questionhelper-server/pkg/mq"
)

// ExportJob 异步导出任务
type ExportJob struct {
	mq mq.MQ
}

func NewExportJob(mq mq.MQ) *ExportJob {
	return &ExportJob{mq: mq}
}

// ExportMessage 导出消息
type ExportMessage struct {
	UserID uint   `json:"user_id"`
	Type   string `json:"type"` // question, exam, score
	Params map[string]interface{} `json:"params"`
}

// Publish 发布导出任务
func (j *ExportJob) Publish(ctx context.Context, msg ExportMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return j.mq.Publish(ctx, "job:export", data)
}

// Handle 处理导出任务
func (j *ExportJob) Handle(ctx context.Context, msg *mq.Message) error {
	var exportMsg ExportMessage
	if err := json.Unmarshal(msg.Body, &exportMsg); err != nil {
		logger.Errorf("解析导出消息失败: %v", err)
		return err
	}

	logger.Infof("处理导出任务: user_id=%d, type=%s", exportMsg.UserID, exportMsg.Type)

	// TODO: 查询数据并生成 Excel
	// TODO: 上传文件到 OSS
	// TODO: 发送导出完成通知

	return nil
}
