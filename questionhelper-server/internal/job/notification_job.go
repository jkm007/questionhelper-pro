package job

import (
	"context"
	"encoding/json"

	"questionhelper-server/pkg/logger"
	"questionhelper-server/pkg/mq"
)

// NotificationJob 异步通知任务
type NotificationJob struct {
	mq mq.MQ
}

func NewNotificationJob(mq mq.MQ) *NotificationJob {
	return &NotificationJob{mq: mq}
}

// NotificationMessage 通知消息
type NotificationMessage struct {
	UserID  uint   `json:"user_id"`
	Type    int    `json:"type"`
	Title   string `json:"title"`
	Content string `json:"content"`
	TargetType string `json:"target_type"`
	TargetID   uint   `json:"target_id"`
}

// Publish 发布通知任务
func (j *NotificationJob) Publish(ctx context.Context, msg NotificationMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return j.mq.Publish(ctx, "job:notification", data)
}

// Handle 处理通知任务
func (j *NotificationJob) Handle(ctx context.Context, msg *mq.Message) error {
	var notifMsg NotificationMessage
	if err := json.Unmarshal(msg.Body, &notifMsg); err != nil {
		logger.Errorf("解析通知消息失败: %v", err)
		return err
	}

	logger.Infof("处理通知任务: user_id=%d, title=%s", notifMsg.UserID, notifMsg.Title)

	// TODO: 保存通知到数据库
	// TODO: 发送 WebSocket 推送
	// TODO: 发送短信通知（如果需要）
	// TODO: 发送邮件通知（如果需要）

	return nil
}
