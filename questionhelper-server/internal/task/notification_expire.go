package task

import (
	"time"

	"questionhelper-server/pkg/database"
	"questionhelper-server/pkg/logger"
)

// NotificationExpireTask 通知过期清理任务
type NotificationExpireTask struct{}

// Run 执行通知过期清理
func (t *NotificationExpireTask) Run() {
	logger.Info("执行通知过期清理任务")

	// 清理90天前的已读通知
	result := database.DB.Exec(`
		DELETE FROM notifications
		WHERE created_at < ?
		AND is_read = true
	`, time.Now().AddDate(0, 0, -90))

	if result.Error != nil {
		logger.Errorf("清理已读通知失败: %v", result.Error)
	} else {
		logger.Infof("清理已读通知: %d 条", result.RowsAffected)
	}

	// 清理180天前的所有通知
	result = database.DB.Exec(`
		DELETE FROM notifications
		WHERE created_at < ?
	`, time.Now().AddDate(0, 0, -180))

	if result.Error != nil {
		logger.Errorf("清理过期通知失败: %v", result.Error)
	} else {
		logger.Infof("清理过期通知: %d 条", result.RowsAffected)
	}
}
