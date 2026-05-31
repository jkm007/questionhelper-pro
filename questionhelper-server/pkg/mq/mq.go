package mq

import "context"

// Message 消息结构
type Message struct {
	Topic   string
	Body    []byte
	Headers map[string]string
}

// Handler 消息处理器
type Handler func(ctx context.Context, msg *Message) error

// MQ 消息队列接口
type MQ interface {
	// Publish 发布消息
	Publish(ctx context.Context, topic string, body []byte) error
	// Subscribe 订阅消息
	Subscribe(ctx context.Context, topic string, handler Handler) error
	// Close 关闭连接
	Close() error
}
