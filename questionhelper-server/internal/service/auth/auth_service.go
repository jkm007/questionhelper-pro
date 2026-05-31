package auth

import (
	"context"
	crand "crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"gorm.io/gorm"

	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	"questionhelper-server/internal/repository/user"
	"questionhelper-server/pkg/captcha"
	"questionhelper-server/pkg/cache/key"
	"questionhelper-server/pkg/config"
	"questionhelper-server/pkg/database"
	"questionhelper-server/pkg/encrypt"
	"questionhelper-server/pkg/jwt"
	"questionhelper-server/pkg/logger"
)

const (
	maxLoginFailCount = 5           // 最大登录失败次数
	lockDuration      = 15 * time.Minute // 锁定时长
)

// Login 用户登录
func Login(req *dto.LoginRequest, cfg *config.JWTConfig) (*dto.LoginResponse, error) {
	// 查找用户
	u, err := user.FindByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户名或密码错误")
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	// 检查用户状态
	if u.Status == 0 {
		return nil, errors.New("账号已被禁用")
	}
	if u.Status == 2 {
		return nil, errors.New("账号正在注销中")
	}

	// 检查是否被锁定
	if u.LockUntil != nil && u.LockUntil.After(time.Now()) {
		remaining := time.Until(*u.LockUntil).Minutes()
		return nil, fmt.Errorf("账号已被锁定，请 %.0f 分钟后再试", remaining)
	}

	// 验证验证码（密码错误3次后需要验证码）
	if u.LoginFailCount >= 3 {
		if req.CaptchaID == "" || req.Captcha == "" {
			return nil, errors.New("请输入验证码")
		}
		if !captcha.VerifyCaptcha(req.CaptchaID, req.Captcha) {
			return nil, errors.New("验证码错误")
		}
	}

	// 验证密码
	if !encrypt.CheckPassword(req.Password, u.Password) {
		// 使用原子操作增加失败次数
		if err := user.IncrementLoginFailCount(u.ID); err != nil {
			logger.Errorf("更新登录失败次数失败: %v", err)
		}
		// 重新获取用户信息以获取最新的失败次数
		u, _ = user.FindByUsername(req.Username)
		newFailCount := u.LoginFailCount
		if newFailCount >= maxLoginFailCount {
			lockUntil := time.Now().Add(lockDuration)
			user.UpdateByID(u.ID, map[string]interface{}{"lock_until": lockUntil})
			logger.Warnf("用户 %s 登录失败 %d 次，账号已锁定 %v", req.Username, newFailCount, lockDuration)
		}

		// 记录登录日志
		recordLoginLog(u, req, false, "密码错误")

		return nil, fmt.Errorf("用户名或密码错误，还可尝试 %d 次", maxLoginFailCount-newFailCount)
	}

	// 密码正确，重置失败次数
	if u.LoginFailCount > 0 || u.LockUntil != nil {
		if err := user.UpdateByID(u.ID, map[string]interface{}{
			"login_fail_count": 0,
			"lock_until":       nil,
		}); err != nil {
			logger.Errorf("重置登录失败次数失败: %v", err)
		}
	}

	// 更新最后登录信息
	now := time.Now()
	if err := user.UpdateByID(u.ID, map[string]interface{}{
		"last_login_at": &now,
		"last_login_ip": req.DeviceInfo,
	}); err != nil {
		logger.Errorf("更新最后登录信息失败: %v", err)
	}

	// 记录登录日志
	recordLoginLog(u, req, true, "")

	// 记录安全日志
	recordSecurityLog(u.ID, "login", "用户登录成功", req.DeviceInfo)

	// 获取角色ID列表
	roleIDs := make([]uint, 0, len(u.Roles))
	roleCodes := make([]string, 0, len(u.Roles))
	for _, role := range u.Roles {
		roleIDs = append(roleIDs, role.ID)
		roleCodes = append(roleCodes, role.Code)
	}

	// 生成 JTI
	jti := jwt.GenerateJTI()

	// 计算 Token 有效期
	accessTokenExpire := cfg.Expire
	refreshTokenExpire := cfg.RefreshExpire
	if req.RememberMe {
		refreshTokenExpire = cfg.RememberMeExpire // 30天
	}

	// 生成 Access Token
	accessToken, err := jwt.GenerateTokenWithJTI(u.ID, u.Username, roleIDs, jti, accessTokenExpire)
	if err != nil {
		return nil, fmt.Errorf("生成访问令牌失败: %w", err)
	}

	// 生成 Refresh Token
	refreshJTI := jwt.GenerateJTI()
	refreshToken, err := jwt.GenerateTokenWithJTI(u.ID, u.Username, roleIDs, refreshJTI, refreshTokenExpire)
	if err != nil {
		return nil, fmt.Errorf("生成刷新令牌失败: %w", err)
	}

	// 记录登录设备
	recordDevice(u.ID, req, jti)

	logger.Infof("用户 %s 登录成功", req.Username)

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    accessTokenExpire,
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

// Register 用户注册
func Register(req *dto.RegisterRequest) error {
	// 检查用户名是否已存在
	exists, err := user.ExistsByUsername(req.Username)
	if err != nil {
		return fmt.Errorf("检查用户名失败: %w", err)
	}
	if exists {
		return errors.New("用户名已存在")
	}

	// 检查手机号
	if req.Phone != "" {
		exists, err = user.ExistsByPhone(req.Phone)
		if err != nil {
			return fmt.Errorf("检查手机号失败: %w", err)
		}
		if exists {
			return errors.New("手机号已被注册")
		}
	}

	// 检查邮箱
	if req.Email != "" {
		exists, err = user.ExistsByEmail(req.Email)
		if err != nil {
			return fmt.Errorf("检查邮箱失败: %w", err)
		}
		if exists {
			return errors.New("邮箱已被注册")
		}
	}

	// 验证验证码
	if req.CaptchaID != "" && req.CaptchaCode != "" {
		if !captcha.VerifyCaptcha(req.CaptchaID, req.CaptchaCode) {
			return errors.New("验证码错误")
		}
	}

	// 加密密码
	hashedPassword, err := encrypt.HashPassword(req.Password)
	if err != nil {
		return fmt.Errorf("加密密码失败: %w", err)
	}

	// 设置昵称
	nickname := req.Nickname
	if nickname == "" {
		nickname = req.Username
	}

	// 设置注册来源
	source := req.Source
	if source == "" {
		source = "web"
	}

	// 创建用户
	u := &model.User{
		Username:       req.Username,
		Password:       hashedPassword,
		Nickname:       nickname,
		Phone:          req.Phone,
		Email:          req.Email,
		Status:         1,
		RegisterSource: source,
	}

	if err := user.Create(u); err != nil {
		return fmt.Errorf("创建用户失败: %w", err)
	}

	// 创建默认隐私设置
	privacy := &model.UserPrivacy{UserID: u.ID}
	if err := user.CreatePrivacy(privacy); err != nil {
		logger.Errorf("创建隐私设置失败: %v", err)
	}

	// 保存密码历史
	if err := savePasswordHistory(u.ID, hashedPassword); err != nil {
		logger.Errorf("保存密码历史失败: %v", err)
	}

	logger.Infof("用户注册成功: %s", req.Username)
	return nil
}

// Logout 用户退出
func Logout(tokenString string) error {
	claims, err := jwt.ParseToken(tokenString)
	if err != nil {
		return nil // token 无效直接返回成功
	}

	// 使用 JTI 作为黑名单 key
	ctx := context.Background()
	blacklistKey := key.TokenBlacklistKey(claims.JTI)
	ttl := time.Until(claims.ExpiresAt.Time)
	if ttl <= 0 {
		return nil
	}

	return database.RDB.Set(ctx, blacklistKey, "1", ttl).Err()
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

// GetCaptcha 获取验证码
func GetCaptcha() (string, string, error) {
	return captcha.GenerateCaptcha()
}

// RequestPasswordReset 请求重置密码
func RequestPasswordReset(account string) error {
	u, err := user.FindByUsername(account)
	if err != nil {
		// 不暴露用户是否存在
		return nil
	}

	// 生成重置令牌
	resetToken := jwt.GenerateJTI()
	ctx := context.Background()
	resetKey := key.ResetTokenKey(u.ID)

	// 存储到 Redis，15分钟有效期
	if err := database.RDB.Set(ctx, resetKey, resetToken, 15*time.Minute).Err(); err != nil {
		return fmt.Errorf("生成重置令牌失败: %w", err)
	}

	// TODO: 发送短信/邮箱通知
	logger.Infof("密码重置令牌已生成: user_id=%d", u.ID)
	return nil
}

// ResetPassword 重置密码
func ResetPassword(account, code, newPassword string) error {
	u, err := user.FindByUsername(account)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 验证重置令牌
	ctx := context.Background()
	resetKey := key.ResetTokenKey(u.ID)
	storedToken, err := database.RDB.Get(ctx, resetKey).Result()
	if err != nil || storedToken != code {
		return errors.New("重置令牌无效或已过期")
	}

	// 加密新密码
	hashedPassword, err := encrypt.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("加密密码失败: %w", err)
	}

	// 更新密码
	now := time.Now()
	if err := user.UpdateByID(u.ID, map[string]interface{}{
		"password":            hashedPassword,
		"password_changed_at": &now,
	}); err != nil {
		return fmt.Errorf("更新密码失败: %w", err)
	}

	// 保存密码历史
	if err := savePasswordHistory(u.ID, hashedPassword); err != nil {
		logger.Errorf("保存密码历史失败: %v", err)
	}

	// 删除已使用的令牌
	database.RDB.Del(ctx, resetKey)

	// 强制退出所有设备
	LogoutAll(u.ID)

	logger.Infof("用户 %s 密码重置成功", account)
	return nil
}

// DeactivateAccount 申请注销账号
func DeactivateAccount(userID uint, password string) error {
	u, err := user.FindByID(userID)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 验证密码
	if !encrypt.CheckPassword(password, u.Password) {
		return errors.New("密码错误")
	}

	// 生成确认验证码
	code := generateVerifyCode()
	ctx := context.Background()
	deactivateKey := key.DeactivateCodeKey(userID)

	// 存储到 Redis，5分钟有效期
	if err := database.RDB.Set(ctx, deactivateKey, code, 5*time.Minute).Err(); err != nil {
		return fmt.Errorf("生成确认码失败: %w", err)
	}

	// TODO: 发送短信/邮箱通知
	logger.Infof("注销确认码已生成: user_id=%d", userID)
	return nil
}

// ConfirmDeactivate 确认注销
func ConfirmDeactivate(userID uint, code string) error {
	// 验证确认码
	ctx := context.Background()
	deactivateKey := key.DeactivateCodeKey(userID)
	storedCode, err := database.RDB.Get(ctx, deactivateKey).Result()
	if err != nil || storedCode != code {
		return errors.New("确认码无效或已过期")
	}

	// 更新状态为注销中
	now := time.Now()
	if err := user.UpdateByID(userID, map[string]interface{}{
		"status":    2,
		"logout_at": &now,
	}); err != nil {
		return fmt.Errorf("更新用户状态失败: %w", err)
	}

	// 删除确认码
	database.RDB.Del(ctx, deactivateKey)

	logger.Infof("用户 %d 注销申请已确认，进入30天冷静期", userID)
	return nil
}

// CancelDeactivate 取消注销
func CancelDeactivate(userID uint) error {
	u, err := user.FindByID(userID)
	if err != nil {
		return errors.New("用户不存在")
	}

	if u.Status != 2 {
		return errors.New("账号未处于注销状态")
	}

	if err := user.UpdateByID(userID, map[string]interface{}{
		"status":    1,
		"logout_at": nil,
	}); err != nil {
		return fmt.Errorf("取消注销失败: %w", err)
	}

	logger.Infof("用户 %d 已取消注销", userID)
	return nil
}

// ExportUserData 导出用户数据
func ExportUserData(userID uint) (*dto.UserDataExport, error) {
	u, err := user.FindByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	export := &dto.UserDataExport{
		ExportAt: time.Now(),
		User: map[string]interface{}{
			"id":       u.ID,
			"username": u.Username,
			"nickname": u.Nickname,
			"email":    u.Email,
			"phone":    u.Phone,
		},
		// TODO: 获取考试、练习、错题、收藏、评论数据
		Exams:          []interface{}{},
		Practices:      []interface{}{},
		WrongQuestions: []interface{}{},
		Favorites:      []interface{}{},
		Comments:       []interface{}{},
	}

	return export, nil
}

// recordLoginLog 记录登录日志
// 注意: IP 字段当前使用 req.DeviceInfo（语义为设备信息/User-Agent），
// 理想情况下应由客户端单独传递 IP 或从请求头 X-Forwarded-For 中提取。
// 此处暂用 DeviceInfo 兼容，后续应添加专门的 IP 字段到 LoginRequest。
func recordLoginLog(u *model.User, req *dto.LoginRequest, success bool, message string) {
	status := int8(1)
	if !success {
		status = 0
	}

	log := &model.LoginLog{
		UserID:    &u.ID,
		Username:  u.Username,
		IP:        req.DeviceInfo, // TODO: 应使用独立的 IP 字段，DeviceInfo 语义为设备信息
		Status:    status,
		Msg:       message,
		LoginType: "password",
	}

	if err := user.CreateLoginLog(log); err != nil {
		logger.Errorf("记录登录日志失败: %v", err)
	}
}

// recordSecurityLog 记录安全日志
func recordSecurityLog(userID uint, eventType, detail, ip string) {
	securityLog := &model.SecurityLog{
		UserID:      userID,
		EventType:   eventType,
		EventDetail: detail,
		IP:          ip,
		Status:      1,
	}

	if err := user.CreateSecurityLog(securityLog); err != nil {
		logger.Errorf("记录安全日志失败: %v", err)
	}
}

// recordDevice 记录登录设备
func recordDevice(userID uint, req *dto.LoginRequest, jti string) {
	// 根据 DeviceID 推断设备类型
	deviceType := "web"
	id := strings.ToLower(req.DeviceID)
	if strings.Contains(id, "mobile") || strings.Contains(id, "app") {
		deviceType = "mobile"
	}

	device := &model.LoginDevice{
		UserID:       userID,
		DeviceID:     req.DeviceID,
		DeviceType:   deviceType,
		IP:           req.DeviceInfo, // TODO: 应使用独立的 IP 字段
		TokenJTI:     jti,
		LastActiveAt: time.Now(),
		IsCurrent:    true,
	}

	if err := user.CreateDevice(device); err != nil {
		logger.Errorf("记录登录设备失败: %v", err)
	}
}

// savePasswordHistory 保存密码历史
func savePasswordHistory(userID uint, password string) error {
	history := &model.PasswordHistory{
		UserID:   userID,
		Password: password,
	}
	return user.CreatePasswordHistory(history)
}

// generateVerifyCode 生成6位验证码
func generateVerifyCode() string {
	n, _ := crand.Int(crand.Reader, big.NewInt(1000000))
	return fmt.Sprintf("%06d", n.Int64())
}

// buildRoleInfos 构建角色信息列表
func buildRoleInfos(roles []model.Role) []dto.RoleInfo {
	infos := make([]dto.RoleInfo, 0, len(roles))
	for _, role := range roles {
		infos = append(infos, dto.RoleInfo{
			ID:          role.ID,
			Name:        role.Name,
			Code:        role.Code,
			Description: role.Description,
		})
	}
	return infos
}

// GetUserProfile 获取用户个人信息
func GetUserProfile(userID uint) (*dto.UserInfo, error) {
	u, err := user.FindByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	return &dto.UserInfo{
		ID:        u.ID,
		Username:  u.Username,
		Nickname:  u.Nickname,
		Email:     u.Email,
		Phone:     u.Phone,
		Avatar:    u.Avatar,
		Gender:    u.Gender,
		Birthday:  u.Birthday,
		Bio:       u.Bio,
		Status:    u.Status,
		IsReal:    u.IsReal,
		Roles:     buildRoleInfos(u.Roles),
		CreatedAt: u.CreatedAt,
	}, nil
}

// GetUserDevices 获取用户登录设备
func GetUserDevices(userID uint) ([]dto.DeviceResponse, error) {
	devices, err := user.FindDevicesByUserID(userID)
	if err != nil {
		return nil, err
	}

	result := make([]dto.DeviceResponse, 0, len(devices))
	for _, d := range devices {
		result = append(result, dto.DeviceResponse{
			ID:           d.ID,
			DeviceType:   d.DeviceType,
			DeviceName:   d.DeviceName,
			Browser:      d.Browser,
			OS:           d.OS,
			IP:           d.IP,
			Location:     d.Location,
			LastActiveAt: d.LastActiveAt.Format("2006-01-02 15:04:05"),
			IsCurrent:    d.IsCurrent,
			CreatedAt:    d.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return result, nil
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

// GetSecurityLogs 获取安全日志
func GetSecurityLogs(userID uint, page, pageSize int) ([]dto.SecurityLogResponse, int64, error) {
	logs, total, err := user.ListSecurityLogs(userID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	result := make([]dto.SecurityLogResponse, 0, len(logs))
	for _, l := range logs {
		result = append(result, dto.SecurityLogResponse{
			ID:          l.ID,
			EventType:   l.EventType,
			EventDetail: l.EventDetail,
			IP:          l.IP,
			Status:      l.Status,
			CreatedAt:   l.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return result, total, nil
}

// GetOAuthStatus 获取第三方绑定状态
func GetOAuthStatus(userID uint) (*dto.OAuthStatusResponse, error) {
	oauths, err := user.FindOAuthUsersByUserID(userID)
	if err != nil {
		return nil, err
	}

	status := &dto.OAuthStatusResponse{}
	for _, oauth := range oauths {
		providerStatus := &dto.OAuthProviderStatus{
			Bound:    true,
			Nickname: oauth.ProviderUsername,
			BoundAt:  oauth.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		switch oauth.Provider {
		case "wechat":
			status.Wechat = providerStatus
		case "github":
			status.Github = providerStatus
		case "google":
			status.Google = providerStatus
		}
	}

	return status, nil
}

// VerifyCaptcha 验证验证码
func VerifyCaptcha(captchaID, captchaCode string) error {
	if !captcha.VerifyCaptcha(captchaID, captchaCode) {
		return errors.New("验证码错误")
	}
	return nil
}

// ChangePassword 修改密码
func ChangePassword(userID uint, oldPassword, newPassword string) error {
	u, err := user.FindByID(userID)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 验证旧密码
	if !encrypt.CheckPassword(oldPassword, u.Password) {
		return errors.New("旧密码错误")
	}

	// 检查新密码是否与旧密码相同
	if oldPassword == newPassword {
		return errors.New("新密码不能与旧密码相同")
	}

	// 检查新密码是否与最近 3 次密码相同
	recentPasswords, err := user.GetRecentPasswords(userID, 3)
	if err == nil {
		for _, history := range recentPasswords {
			if encrypt.CheckPassword(newPassword, history.Password) {
				return errors.New("新密码不能与最近 3 次密码相同")
			}
		}
	}

	// 加密新密码
	hashedPassword, err := encrypt.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("加密密码失败: %w", err)
	}

	// 更新密码
	now := time.Now()
	if err := user.UpdateByID(userID, map[string]interface{}{
		"password":            hashedPassword,
		"password_changed_at": &now,
	}); err != nil {
		return fmt.Errorf("更新密码失败: %w", err)
	}

	// 保存密码历史
	if err := savePasswordHistory(userID, hashedPassword); err != nil {
		logger.Errorf("保存密码历史失败: %v", err)
	}

	// 记录安全日志
	securityLog := &model.SecurityLog{
		UserID:      userID,
		EventType:   "password_change",
		EventDetail: "用户修改密码",
		Status:      1,
	}
	if err := user.CreateSecurityLog(securityLog); err != nil {
		logger.Errorf("记录安全日志失败: %v", err)
	}

	logger.Infof("用户 %d 密码修改成功", userID)
	return nil
}
