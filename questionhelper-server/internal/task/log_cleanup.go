package task

import (
	"time"

	"questionhelper-server/pkg/database"
	"questionhelper-server/pkg/logger"
)

// LogCleanupTask 日志清理任务
type LogCleanupTask struct{}

// Run 执行日志清理
func (t *LogCleanupTask) Run() {
	logger.Info("执行日志清理任务")

	// 清理30天前的操作日志
	result := database.DB.Exec(`
		DELETE FROM operation_logs
		WHERE created_at < ?
	`, time.Now().AddDate(0, 0, -30))

	if result.Error != nil {
		logger.Errorf("清理操作日志失败: %v", result.Error)
	} else {
		logger.Infof("清理操作日志: %d 条", result.RowsAffected)
	}

	// 清理90天前的登录日志
	result = database.DB.Exec(`
		DELETE FROM login_logs
		WHERE created_at < ?
	`, time.Now().AddDate(0, 0, -90))

	if result.Error != nil {
		logger.Errorf("清理登录日志失败: %v", result.Error)
	} else {
		logger.Infof("清理登录日志: %d 条", result.RowsAffected)
	}
}
