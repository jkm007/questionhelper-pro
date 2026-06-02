package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"questionhelper-server/internal/router"
	"questionhelper-server/internal/service/file"
	"questionhelper-server/internal/task"
	"questionhelper-server/internal/ws"
	"questionhelper-server/pkg/config"
	"questionhelper-server/pkg/database"
	"questionhelper-server/pkg/email"
	"questionhelper-server/pkg/jwt"
	"questionhelper-server/pkg/logger"
	"questionhelper-server/pkg/mq"
	"questionhelper-server/pkg/mq/handler"
	"questionhelper-server/pkg/sms"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	logger.Init(cfg.Log)
	defer logger.Sync()

	database.InitMySQL(cfg.MySQL)
	database.InitRedis(cfg.Redis)

	jwt.Init(cfg.JWT.Secret, cfg.Auth.JWTIssuer)

	// 初始化文件存储
	if err := file.Init(cfg.OSS); err != nil {
		logger.Warnf("文件存储初始化失败，使用本地存储: %v", err)
	}

	// 初始化 WebSocket Hub
	hub := ws.NewHub()
	go hub.Run()
	logger.Info("WebSocket Hub 已启动")

	// 初始化消息队列并注册消费者
	mqInstance := mq.NewRedisStream()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	registerMQConsumers(ctx, mqInstance)

	// 启动定时任务调度器
	scheduler := task.NewScheduler()
	scheduler.Start()
	defer scheduler.Stop()

	// 初始化延迟队列处理器
	delayQueue := mq.NewRedisDelayQueue("qh:delay:queue")
	go func() {
		if err := delayQueue.Process(func(taskType string, payload []byte) error {
			logger.Infof("处理延迟任务: type=%s", taskType)
			return nil
		}); err != nil {
			logger.Errorf("延迟队列处理器退出: %v", err)
		}
	}()

	// 将 config.EmailConfig 转换为 email.EmailConfig
	email.Init(&email.EmailConfig{
		SMTPHost:    cfg.Email.SMTPHost,
		SMTPPort:    cfg.Email.SMTPPort,
		Username:    cfg.Email.Username,
		Password:    cfg.Email.Password,
		FromAddress: cfg.Email.FromAddress,
		FromName:    cfg.Email.FromName,
		UseTLS:      cfg.Email.UseTLS,
	})

	sms.Init(&sms.SMSConfig{
		Provider:     cfg.SMS.Provider,
		AccessKey:    cfg.SMS.AccessKey,
		AccessSecret: cfg.SMS.AccessSecret,
		SignName:     cfg.SMS.SignName,
		TemplateCode: cfg.SMS.TemplateCode,
	})

	r := router.Setup(cfg, hub)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	go func() {
		logger.Infof("服务启动成功: %s", addr)
		if err := r.Run(addr); err != nil {
			log.Fatalf("服务启动失败: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("服务正在关闭...")

	// 优雅关闭
	cancel()
	mqInstance.Close()
}

// registerMQConsumers 注册所有 MQ 消费者
func registerMQConsumers(ctx context.Context, mqInstance mq.MQ) {
	// 考试提交消费者
	go func() {
		if err := mqInstance.Subscribe(ctx, "job:exam_submit", handler.ExamSubmitHandler); err != nil {
			logger.Errorf("考试提交消费者退出: %v", err)
		}
	}()

	// 通知消费者
	go func() {
		if err := mqInstance.Subscribe(ctx, "job:notification", handler.NotificationHandler); err != nil {
			logger.Errorf("通知消费者退出: %v", err)
		}
	}()

	logger.Info("MQ 消费者已注册")
}
