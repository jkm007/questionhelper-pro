package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/repository/user"
	"questionhelper-server/pkg/cache/key"
	"questionhelper-server/pkg/config"
	"questionhelper-server/pkg/database"
	"questionhelper-server/pkg/jwt"
	"questionhelper-server/pkg/logger"
)

// RefreshToken 刷新令牌
func RefreshToken(req *dto.RefreshTokenRequest, cfg *config.JWTConfig) (*dto.LoginResponse, error) {
	// 解析刷新令牌
	claims, err := jwt.ParseToken(req.RefreshToken)
	if err != nil {
		return nil, errors.New("刷新令牌无效")
	}

	// 检查是否在黑名单中
	ctx := context.Background()
	blacklistKey := key.TokenBlacklistKey(claims.JTI)
	exists, err := database.RDB.Exists(ctx, blacklistKey).Result()
	if err != nil {
		return nil, fmt.Errorf("检查令牌状态失败: %w", err)
	}
	if exists > 0 {
		return nil, errors.New("刷新令牌已失效")
	}

	// 查找用户
	u, err := user.FindByID(claims.UserID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 检查用户状态
	if u.Status != 1 {
		return nil, errors.New("账号已被禁用")
	}

	// 获取角色ID列表
	roleIDs := make([]uint, 0, len(u.Roles))
	roleCodes := make([]string, 0, len(u.Roles))
	for _, role := range u.Roles {
		roleIDs = append(roleIDs, role.ID)
		roleCodes = append(roleCodes, role.Code)
	}

	// 生成新 Token
	jti := jwt.GenerateJTI()
	accessToken, err := jwt.GenerateTokenWithJTI(u.ID, u.Username, roleIDs, jti, cfg.Expire)
	if err != nil {
		return nil, fmt.Errorf("生成访问令牌失败: %w", err)
	}

	refreshJTI := jwt.GenerateJTI()
	refreshToken, err := jwt.GenerateTokenWithJTI(u.ID, u.Username, roleIDs, refreshJTI, cfg.RefreshExpire)
	if err != nil {
		return nil, fmt.Errorf("生成刷新令牌失败: %w", err)
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    cfg.Expire,
		User: dto.UserInfo{
			ID:        u.ID,
			Username:  u.Username,
			Nickname:  u.Nickname,
			Avatar:    u.Avatar,
			Roles:     buildRoleInfos(u.Roles),
			CreatedAt: u.CreatedAt,
		},
	}, nil
}

// LogoutAll 退出所有设备
func LogoutAll(userID uint) error {
	// 获取用户所有设备
	devices, err := user.FindDevicesByUserID(userID)
	if err != nil {
		return fmt.Errorf("获取设备列表失败: %w", err)
	}

	ctx := context.Background()
	for _, device := range devices {
		// 将每个设备的 Token JTI 加入黑名单
		blacklistKey := key.TokenBlacklistKey(device.TokenJTI)
		database.RDB.Set(ctx, blacklistKey, "1", 7*24*time.Hour)
	}

	// 删除所有设备记录
	if err := user.DeleteDevicesByUserID(userID); err != nil {
		return fmt.Errorf("删除设备记录失败: %w", err)
	}

	logger.Infof("用户 %d 已退出所有设备", userID)
	return nil
}

// KickDevice 踢出设备
func KickDevice(userID, deviceID uint) error {
	device, err := user.FindDeviceByID(deviceID)
	if err != nil {
		return errors.New("设备不存在")
	}

	if device.UserID != userID {
		return errors.New("无权操作")
	}

	// 将设备 Token JTI 加入黑名单
	ctx := context.Background()
	blacklistKey := key.TokenBlacklistKey(device.TokenJTI)
	database.RDB.Set(ctx, blacklistKey, "1", 7*24*time.Hour)

	// 删除设备记录
	return user.DeleteDevice(deviceID)
}
