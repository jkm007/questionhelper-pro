package job

import (
	"context"
	"encoding/json"

	"questionhelper-server/pkg/logger"
	"questionhelper-server/pkg/mq"
)

// ExamSubmitJob 异步交卷任务
type ExamSubmitJob struct {
	mq mq.MQ
}

func NewExamSubmitJob(mq mq.MQ) *ExamSubmitJob {
	return &ExamSubmitJob{mq: mq}
}

// ExamSubmitMessage 考试提交消息
type ExamSubmitMessage struct {
	RecordID uint `json:"record_id"`
	ExamID   uint `json:"exam_id"`
	UserID   uint `json:"user_id"`
	AutoSubmit bool `json:"auto_submit"` // 是否超时自动提交
}

// Publish 发布交卷任务
func (j *ExamSubmitJob) Publish(ctx context.Context, msg ExamSubmitMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return j.mq.Publish(ctx, "job:exam_submit", data)
}

// Handle 处理交卷任务
func (j *ExamSubmitJob) Handle(ctx context.Context, msg *mq.Message) error {
	var submitMsg ExamSubmitMessage
	if err := json.Unmarshal(msg.Body, &submitMsg); err != nil {
		logger.Errorf("解析交卷消息失败: %v", err)
		return err
	}

	logger.Infof("处理交卷任务: record_id=%d, auto_submit=%v", submitMsg.RecordID, submitMsg.AutoSubmit)

	// TODO: 自动评分（客观题）
	// TODO: 更新考试记录状态
	// TODO: 生成成绩报告
	// TODO: 发送交卷通知

	return nil
}
