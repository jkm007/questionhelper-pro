package task

import (
	"time"

	"questionhelper-server/pkg/database"
	"questionhelper-server/pkg/logger"
)

// FileCleanupTask 文件清理任务
type FileCleanupTask struct{}

// Run 执行文件清理
func (t *FileCleanupTask) Run() {
	logger.Info("执行文件清理任务")

	// 清理24小时前的临时文件
	result := database.DB.Exec(`
		DELETE FROM files
		WHERE created_at < ?
		AND path LIKE '%/temp/%'
	`, time.Now().Add(-24*time.Hour))

	if result.Error != nil {
		logger.Errorf("清理临时文件失败: %v", result.Error)
	} else {
		logger.Infof("清理临时文件: %d 条", result.RowsAffected)
	}

	// 清理孤立文件记录（未关联到任何实体的文件）
	result = database.DB.Exec(`
		DELETE FROM files
		WHERE created_at < ?
		AND uploader_id IS NULL
	`, time.Now().Add(-7*24*time.Hour))

	if result.Error != nil {
		logger.Errorf("清理孤立文件失败: %v", result.Error)
	} else {
		logger.Infof("清理孤立文件: %d 条", result.RowsAffected)
	}
}
