package task

import (
	"time"

	"go.uber.org/zap"
	"questionhelper-server/internal/repository/user"
)

// DeviceCleanupTask 不活跃设备清理任务
type DeviceCleanupTask struct{}

// Run 执行清理
func (t *DeviceCleanupTask) Run() {
	// 清理 30 天不活跃的设备
	threshold := time.Now().AddDate(0, 0, -30)
	count, err := user.DeleteInactiveDevices(threshold)
	if err != nil {
		zap.L().Error("清理不活跃设备失败", zap.Error(err))
		return
	}

	if count > 0 {
		zap.L().Info("清理不活跃设备完成", zap.Int64("count", count))
	}
}
