package handler

import (
	"context"
	"encoding/json"

	"questionhelper-server/pkg/logger"
	"questionhelper-server/pkg/mq"
)

// ExamSubmitMessage 考试提交消息
type ExamSubmitMessage struct {
	RecordID uint `json:"record_id"`
	ExamID   uint `json:"exam_id"`
	UserID   uint `json:"user_id"`
}

// ExamSubmitHandler 考试提交处理器（异步评分）
func ExamSubmitHandler(ctx context.Context, msg *mq.Message) error {
	var submit ExamSubmitMessage
	if err := json.Unmarshal(msg.Body, &submit); err != nil {
		logger.Errorf("解析考试提交消息失败: %v", err)
		return err
	}

	logger.Infof("处理考试提交: record_id=%d, exam_id=%d", submit.RecordID, submit.ExamID)

	// TODO: 自动评分（客观题）
	// TODO: 生成成绩报告
	// TODO: 发送成绩通知

	return nil
}
