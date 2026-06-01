package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"questionhelper-server/pkg/database"
	"questionhelper-server/pkg/logger"
)

// delayTask 延迟任务结构
type delayTask struct {
	TaskType string `json:"task_type"`
	Payload  []byte `json:"payload"`
}

// RedisDelayQueue 基于 Redis ZSet 的延迟队列
type RedisDelayQueue struct {
	client *redis.Client
	key    string
}

// NewRedisDelayQueue 创建延迟队列实例
func NewRedisDelayQueue(key string) *RedisDelayQueue {
	return &RedisDelayQueue{
		client: database.RDB,
		key:    key,
	}
}

// Add 添加延迟任务，score 为任务触发的 Unix 时间戳
func (q *RedisDelayQueue) Add(taskType string, payload []byte, delay time.Duration) error {
	task := &delayTask{
		TaskType: taskType,
		Payload:  payload,
	}

	data, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("序列化任务失败: %w", err)
	}

	score := float64(time.Now().Add(delay).Unix())
	return q.client.ZAdd(context.Background(), q.key, redis.Z{
		Score:  score,
		Member: string(data),
	}).Err()
}

// Size 返回队列中待处理任务数量
func (q *RedisDelayQueue) Size() (int64, error) {
	return q.client.ZCard(context.Background(), q.key).Result()
}

// Process 轮询处理到期任务，阻塞直到 context 取消
func (q *RedisDelayQueue) Process(handler func(taskType string, payload []byte) error) error {
	ctx := context.Background()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			now := float64(time.Now().Unix())

			// 取出一条到期任务（原子操作）
			results, err := q.client.ZPopMin(ctx, q.key, 1).Result()
			if err != nil {
				logger.Errorf("延迟队列取出任务失败: %v", err)
				time.Sleep(time.Second)
				continue
			}

			if len(results) == 0 {
				// 队列为空，等待后重试
				time.Sleep(500 * time.Millisecond)
				continue
			}

			// 检查任务是否到期
			if results[0].Score > now {
				// 未到期，放回队列
				q.client.ZAdd(ctx, q.key, redis.Z{
					Score:  results[0].Score,
					Member: results[0].Member,
				})
				time.Sleep(500 * time.Millisecond)
				continue
			}

			// 解析任务
			member, ok := results[0].Member.(string)
			if !ok {
				logger.Errorf("延迟队列任务格式错误")
				continue
			}

			var task delayTask
			if err := json.Unmarshal([]byte(member), &task); err != nil {
				logger.Errorf("延迟队列反序列化任务失败: %v", err)
				continue
			}

			// 执行任务
			if err := handler(task.TaskType, task.Payload); err != nil {
				logger.Errorf("延迟队列处理任务失败: %v", err)
				// 任务执行失败，重新放回队列，延迟 5 秒重试
				q.client.ZAdd(ctx, q.key, redis.Z{
					Score:  float64(time.Now().Add(5 * time.Second).Unix()),
					Member: member,
				})
			}
		}
	}
}
