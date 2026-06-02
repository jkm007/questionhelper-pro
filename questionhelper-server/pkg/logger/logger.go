package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"questionhelper-server/pkg/config"
)

var log *zap.Logger
var sugar *zap.SugaredLogger

func Init(cfg config.LogConfig) {
	var zapCfg zap.Config
	if cfg.Format == "json" {
		zapCfg = zap.NewProductionConfig()
	} else {
		zapCfg = zap.NewDevelopmentConfig()
	}

	level, _ := zap.ParseAtomicLevel(cfg.Level)
	zapCfg.Level = level

	// 配置日志输出目标
	if cfg.Output != "" {
		// 同时输出到文件和 stdout
		zapCfg.OutputPaths = []string{"stdout", cfg.Output}
		zapCfg.ErrorOutputPaths = []string{"stderr", cfg.Output}

		// 使用 lumberjack 实现日志轮转
		fileWriter := &lumberjack.Logger{
			Filename:   cfg.Output,
			MaxSize:    cfg.MaxSize,    // MB
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,     // 天
			Compress:   true,
		}

		// 创建自定义 core: 同时写文件(带轮转)和 stdout
		encoder := zapCfg.EncoderConfig
		var fileEncoder zapcore.Encoder
		if cfg.Format == "json" {
			fileEncoder = zapcore.NewJSONEncoder(encoder)
		} else {
			fileEncoder = zapcore.NewConsoleEncoder(encoder)
		}

		fileCore := zapcore.NewCore(
			fileEncoder,
			zapcore.AddSync(fileWriter),
			level,
		)

		stdoutCore := zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoder),
			zapcore.AddSync(os.Stdout),
			level,
		)

		log = zap.New(zapcore.NewTee(fileCore, stdoutCore),
			zap.AddCaller(),
			zap.AddStacktrace(zapcore.ErrorLevel),
		)
	} else {
		// 仅输出到 stdout
		var err error
		log, err = zapCfg.Build(zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
		if err != nil {
			panic(err)
		}
	}

	sugar = log.Sugar()
}

func Sync() {
	if log != nil {
		log.Sync()
	}
}

func Info(args ...interface{}) {
	sugar.Info(args...)
}

func Infof(template string, args ...interface{}) {
	sugar.Infof(template, args...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	sugar.Infow(msg, keysAndValues...)
}

func Error(args ...interface{}) {
	sugar.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	sugar.Errorf(template, args...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	sugar.Errorw(msg, keysAndValues...)
}

func Warn(args ...interface{}) {
	sugar.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	sugar.Warnf(template, args...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	sugar.Warnw(msg, keysAndValues...)
}

func Debug(args ...interface{}) {
	sugar.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	sugar.Debugf(template, args...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	sugar.Debugw(msg, keysAndValues...)
}

func Fatal(args ...interface{}) {
	sugar.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	sugar.Fatalf(template, args...)
}
