package task

import (
	"context"

	"go.uber.org/zap"
	"questionhelper-server/pkg/database"
)

// TokenCleanupTask Token 黑名单清理任务
type TokenCleanupTask struct{}

// Run 执行 Token 黑名单清理
func (t *TokenCleanupTask) Run() {
	ctx := context.Background()

	// 查找所有 qh:token:bl:* 的 key（Redis 会自动清理过期的 key，这里做统计）
	keys, err := database.RDB.Keys(ctx, "qh:token:bl:*").Result()
	if err != nil {
		zap.L().Error("查找 Token 黑名单失败", zap.Error(err))
		return
	}

	zap.L().Info("Token 黑名单清理完成", zap.Int("count", len(keys)))
}
