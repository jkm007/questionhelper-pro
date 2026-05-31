package task

import (
	"time"

	"go.uber.org/zap"
	"questionhelper-server/internal/model"
	"questionhelper-server/internal/repository/user"
	"questionhelper-server/pkg/database"
)

// AccountDeactivateTask 账号注销处理任务
type AccountDeactivateTask struct{}

// Run 执行注销处理
func (t *AccountDeactivateTask) Run() {
	// 查找 30 天前申请注销的用户（status=2 且 logout_at 在 30 天前）
	threshold := time.Now().AddDate(0, 0, -30)

	var users []model.User
	err := database.DB.Where("status = ? AND logout_at IS NOT NULL AND logout_at < ?", 2, threshold).
		Find(&users).Error
	if err != nil {
		zap.L().Error("查找注销用户失败", zap.Error(err))
		return
	}

	for _, u := range users {
		// 执行数据删除
		if err := deleteUserData(u.ID); err != nil {
			zap.L().Error("删除用户数据失败", zap.Uint("user_id", u.ID), zap.Error(err))
			continue
		}

		// 标记为已删除（软删除）
		if err := user.DeleteByID(u.ID); err != nil {
			zap.L().Error("标记用户删除失败", zap.Uint("user_id", u.ID), zap.Error(err))
			continue
		}

		zap.L().Info("用户数据已永久删除", zap.Uint("user_id", u.ID))
	}

	if len(users) > 0 {
		zap.L().Info("账号注销处理完成", zap.Int("count", len(users)))
	}
}

// deleteUserData 删除用户相关数据
func deleteUserData(userID uint) error {
	// 删除用户设备
	if err := user.DeleteDevicesByUserID(userID); err != nil {
		return err
	}

	// 删除隐私设置
	database.DB.Where("user_id = ?", userID).Delete(&model.UserPrivacy{})

	// 删除第三方登录绑定
	database.DB.Where("user_id = ?", userID).Delete(&model.OAuthUser{})

	// 删除安全日志
	database.DB.Where("user_id = ?", userID).Delete(&model.SecurityLog{})

	// 删除密码历史
	database.DB.Where("user_id = ?", userID).Delete(&model.PasswordHistory{})

	// 删除登录日志
	database.DB.Where("user_id = ?", userID).Delete(&model.LoginLog{})

	return nil
}
