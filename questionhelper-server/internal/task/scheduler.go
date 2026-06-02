package task

import (
	"github.com/robfig/cron/v3"
	"questionhelper-server/pkg/logger"
)

// Scheduler 定时任务调度器
type Scheduler struct {
	cron *cron.Cron
}

// NewScheduler 创建调度器
func NewScheduler() *Scheduler {
	return &Scheduler{
		cron: cron.New(),
	}
}

// Start 启动调度器，注册所有定时任务
func (s *Scheduler) Start() {
	// Token 黑名单清理 — 每小时执行一次
	s.cron.AddJob("0 * * * *", &TokenCleanupTask{})

	// 日志清理 — 每天凌晨 3 点
	s.cron.AddJob("0 3 * * *", &LogCleanupTask{})

	// 考试提醒 — 每 5 分钟检查一次
	s.cron.AddJob("*/5 * * * *", &ExamRemindTask{})

	// 作业提醒 — 每 10 分钟检查一次
	s.cron.AddJob("*/10 * * * *", &HomeworkRemindTask{})

	// 文件清理 — 每天凌晨 4 点
	s.cron.AddJob("0 4 * * *", &FileCleanupTask{})

	// 通知过期清理 — 每天凌晨 2 点
	s.cron.AddJob("0 2 * * *", &NotificationExpireTask{})

	// 成绩预警 — 每小时执行一次
	s.cron.AddJob("0 * * * *", &ScoreWarningTask{})

	// 设备清理 — 每天凌晨 5 点
	s.cron.AddJob("0 5 * * *", &DeviceCleanupTask{})

	// 账号注销处理 — 每天凌晨 1 点
	s.cron.AddJob("0 1 * * *", &AccountDeactivateTask{})

	// 练习超时自动结束 — 每小时检查一次
	s.cron.AddJob("0 * * * *", &PracticeTimeoutTask{})

	s.cron.Start()
	logger.Info("定时任务调度器已启动")
}

// Stop 停止调度器
func (s *Scheduler) Stop() {
	s.cron.Stop()
	logger.Info("定时任务调度器已停止")
}
