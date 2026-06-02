package database

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"questionhelper-server/pkg/config"
)

var DB *gorm.DB

func InitMySQL(cfg config.MySQLConfig) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	// 根据项目日志级别设置 GORM logger 级别
	logLevel := gormlogger.Info
	if config.Cfg != nil {
		switch strings.ToLower(config.Cfg.Log.Level) {
		case "debug":
			logLevel = gormlogger.Info
		case "warn":
			logLevel = gormlogger.Warn
		case "error":
			logLevel = gormlogger.Error
		default:
			logLevel = gormlogger.Silent
		}
	}

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(logLevel),
	})
	if err != nil {
		panic(fmt.Sprintf("连接 MySQL 失败: %v", err))
	}

	sqlDB, err := DB.DB()
	if err != nil {
		panic(fmt.Sprintf("获取 DB 实例失败: %v", err))
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
}
