package task

import (
	"time"

	"questionhelper-server/internal/model"
	"questionhelper-server/pkg/database"
	"questionhelper-server/pkg/logger"
)

// PracticeTimeoutTask 练习超时自动结束任务
type PracticeTimeoutTask struct{}

// Run 执行练习超时检查，自动结束超过24小时的进行中练习
func (t *PracticeTimeoutTask) Run() {
	result := database.DB.Model(&model.PracticeSession{}).
		Where("status = 0 AND created_at < ?", time.Now().Add(-24*time.Hour)).
		Update("status", 5) // 5=已超时

	if result.Error != nil {
		logger.Errorf("练习超时任务执行失败: %v", result.Error)
		return
	}

	if result.RowsAffected > 0 {
		logger.Infof("练习超时任务完成，自动结束 %d 个超时练习", result.RowsAffected)
	}
}
