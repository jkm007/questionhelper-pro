package mq

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"questionhelper-server/pkg/database"
	"questionhelper-server/pkg/logger"
)

type RedisStream struct {
	client *redis.Client
}

func NewRedisStream() *RedisStream {
	return &RedisStream{
		client: database.RDB,
	}
}

func (s *RedisStream) Publish(ctx context.Context, topic string, body []byte) error {
	return s.client.XAdd(ctx, &redis.XAddArgs{
		Stream: topic,
		Values: map[string]interface{}{
			"body": string(body),
		},
	}).Err()
}

func (s *RedisStream) Subscribe(ctx context.Context, topic string, handler Handler) error {
	group := fmt.Sprintf("group:%s", topic)
	consumer := fmt.Sprintf("consumer:%s", topic)

	// 创建消费者组
	s.client.XGroupCreateMkStream(ctx, topic, group, "0")

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			streams, err := s.client.XReadGroup(ctx, &redis.XReadGroupArgs{
				Group:    group,
				Consumer: consumer,
				Streams:  []string{topic, ">"},
				Block:    0,
				Count:    1,
			}).Result()

			if err != nil {
				logger.Errorf("读取消息失败: %v", err)
				continue
			}

			for _, stream := range streams {
				for _, msg := range stream.Messages {
					m := &Message{
						Topic: topic,
						Body:  []byte(msg.Values["body"].(string)),
					}

					if err := handler(ctx, m); err != nil {
						logger.Errorf("处理消息失败: %v", err)
						continue
					}

					s.client.XAck(ctx, topic, group, msg.ID)
				}
			}
		}
	}
}

func (s *RedisStream) Close() error {
	return nil
}
