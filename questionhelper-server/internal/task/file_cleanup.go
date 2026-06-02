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

	// 清理24小时前的临时文件（软删除）
	result := database.DB.Exec(`
		UPDATE files
		SET deleted_at = NOW(), status = 'deleted'
		WHERE deleted_at IS NULL
		AND created_at < ?
		AND path LIKE '%/temp/%'
	`, time.Now().Add(-24*time.Hour))

	if result.Error != nil {
		logger.Errorf("清理临时文件失败: %v", result.Error)
	} else {
		logger.Infof("软删除临时文件: %d 条", result.RowsAffected)
	}

	// 第一步：软删除无引用的孤立文件（reference_count = 0 且超过7天）
	result = database.DB.Exec(`
		UPDATE files
		SET deleted_at = NOW(), status = 'deleted'
		WHERE deleted_at IS NULL
		AND reference_count = 0
		AND status = 'active'
		AND created_at < ?
	`, time.Now().Add(-7*24*time.Hour))

	if result.Error != nil {
		logger.Errorf("软删除孤立文件失败: %v", result.Error)
	} else {
		logger.Infof("软删除孤立文件: %d 条", result.RowsAffected)
	}

	// 第二步：物理删除已软删除超过30天的孤立文件
	result = database.DB.Exec(`
		DELETE FROM files
		WHERE reference_count = 0
		AND deleted_at IS NOT NULL
		AND deleted_at < NOW() - INTERVAL 30 DAY
	`)

	if result.Error != nil {
		logger.Errorf("物理删除孤立文件失败: %v", result.Error)
	} else {
		logger.Infof("物理删除孤立文件: %d 条", result.RowsAffected)
	}
}
