package task

import (
	"context"

	"questionhelper-server/pkg/database"
	"questionhelper-server/pkg/logger"
)

// TokenCleanupTask Token 黑名单清理任务
type TokenCleanupTask struct{}

// Run 执行 Token 黑名单清理
func (t *TokenCleanupTask) Run() {
	ctx := context.Background()

	// 查找所有 qh:token:bl:* 的 key（Redis 会自动清理过期的 key，这里做统计）
	keys, err := database.RDB.Keys(ctx, "qh:token:bl:*").Result()
	if err != nil {
		logger.Errorf("查找 Token 黑名单失败: %v", err)
		return
	}

	logger.Infof("Token 黑名单清理完成，当前数量: %d", len(keys))
}
