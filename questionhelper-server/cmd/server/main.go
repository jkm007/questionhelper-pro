package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"questionhelper-server/internal/router"
	"questionhelper-server/pkg/config"
	"questionhelper-server/pkg/database"
	"questionhelper-server/pkg/email"
	"questionhelper-server/pkg/jwt"
	"questionhelper-server/pkg/logger"
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

	jwt.Init(cfg.JWT.Secret)

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

	r := router.Setup(cfg)

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
}
