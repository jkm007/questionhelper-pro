package job

import (
	"context"
	"encoding/json"

	"questionhelper-server/pkg/logger"
	"questionhelper-server/pkg/mq"
)

// StatisticsJob 异步统计任务
type StatisticsJob struct {
	mq mq.MQ
}

func NewStatisticsJob(mq mq.MQ) *StatisticsJob {
	return &StatisticsJob{mq: mq}
}

// StatisticsMessage 统计消息
type StatisticsMessage struct {
	Type string `json:"type"` // daily, weekly, monthly
	Date string `json:"date"`
}

// Publish 发布统计任务
func (j *StatisticsJob) Publish(ctx context.Context, msg StatisticsMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return j.mq.Publish(ctx, "job:statistics", data)
}

// Handle 处理统计任务
func (j *StatisticsJob) Handle(ctx context.Context, msg *mq.Message) error {
	var statsMsg StatisticsMessage
	if err := json.Unmarshal(msg.Body, &statsMsg); err != nil {
		logger.Errorf("解析统计消息失败: %v", err)
		return err
	}

	logger.Infof("处理统计任务: type=%s, date=%s", statsMsg.Type, statsMsg.Date)

	// TODO: 统计用户数据
	// TODO: 统计题目数据
	// TODO: 统计考试数据
	// TODO: 统计班级数据

	return nil
}
