package handler

import (
	"context"
	"encoding/json"

	"questionhelper-server/pkg/logger"
	"questionhelper-server/pkg/mq"
)

// NotificationMessage 通知消息
type NotificationMessage struct {
	UserID  uint   `json:"user_id"`
	Type    int    `json:"type"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// NotificationHandler 通知消息处理器
func NotificationHandler(ctx context.Context, msg *mq.Message) error {
	var notification NotificationMessage
	if err := json.Unmarshal(msg.Body, &notification); err != nil {
		logger.Errorf("解析通知消息失败: %v", err)
		return err
	}

	logger.Infof("处理通知消息: user_id=%d, title=%s", notification.UserID, notification.Title)

	// TODO: 保存通知到数据库
	// TODO: 发送 WebSocket 推送

	return nil
}
