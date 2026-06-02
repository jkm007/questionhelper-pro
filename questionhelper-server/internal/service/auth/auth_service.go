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
	"questionhelper-server/internal/repository/application"
	classrepo "questionhelper-server/internal/repository/class"
	"questionhelper-server/internal/repository/exam"
	"questionhelper-server/internal/repository/user"
	"questionhelper-server/pkg/cache/key"
	"questionhelper-server/pkg/captcha"
	"questionhelper-server/pkg/config"
	"questionhelper-server/pkg/database"
	"questionhelper-server/pkg/email"
	"questionhelper-server/pkg/encrypt"
	"questionhelper-server/pkg/jwt"
	"questionhelper-server/pkg/logger"
	"questionhelper-server/pkg/sensitive"
	"questionhelper-server/pkg/sms"
)

const (
	maxLoginFailCount    = 5                // 最大登录失败次数
	lockDuration         = 15 * time.Minute // 锁定时长
	loginRateLimit       = 10               // IP 登录尝试次数上限
	loginRateLimitWindow = 5 * time.Minute  // IP 登录限制时间窗口
)

// FindByIdentifier 通过用户名/邮箱/手机号查找用户（T01: 支持多种登录方式）
func FindByIdentifier(identifier string) (*model.User, error) {
	// 优先按用户名查找
	u, err := user.FindByUsername(identifier)
	if err == nil {
		return u, nil
	}
	// 按邮箱查找
	u, err = user.FindByEmail(identifier)
	if err == nil {
		return u, nil
	}
	// 按手机号查找
	return user.FindByPhone(identifier)
}

// Login 用户登录
func Login(req *dto.LoginRequest, cfg *config.JWTConfig) (*dto.LoginResponse, error) {
	// T20: IP 登录频率限制
	clientIP := req.DeviceInfo // 控制器已将 ClientIP 赋值给 DeviceInfo
	limitKey := key.LoginLimitKey(clientIP)
	ctx := context.Background()

	count, err := database.RDB.Incr(ctx, limitKey).Result()
	if err != nil {
		logger.Errorf("检查登录频率限制失败: %v", err)
	} else {
		if count == 1 {
			// 首次尝试，设置过期时间
			database.RDB.Expire(ctx, limitKey, loginRateLimitWindow)
		}
		if count > int64(loginRateLimit) {
			return nil, errors.New("登录尝试过于频繁，请稍后再试")
		}
	}

	// 查找用户（支持用户名/邮箱/手机号）
	u, err := FindByIdentifier(req.Username)
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
		// T30: 验证码尝试次数跟踪
		ok, tooMany := captcha.VerifyCaptcha(req.CaptchaID, req.Captcha, req.CaptchaID)
		if tooMany {
			return nil, errors.New("验证码错误次数过多，请重新获取验证码")
		}
		if !ok {
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

		return nil, errors.New("用户名或密码错误")
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

	// T16: 登录异常检测（新设备、IP 变更）
	detectLoginAnomaly(u.ID, req)

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

	// 生成 Access Token（T28: 添加 type 字段区分 token 类型）
	accessToken, err := jwt.GenerateTokenWithJTI(u.ID, u.Username, roleIDs, jti, accessTokenExpire, "access")
	if err != nil {
		return nil, fmt.Errorf("生成访问令牌失败: %w", err)
	}

	// 生成 Refresh Token（T28: 添加 type 字段区分 token 类型）
	refreshJTI := jwt.GenerateJTI()
	refreshToken, err := jwt.GenerateTokenWithJTI(u.ID, u.Username, roleIDs, refreshJTI, refreshTokenExpire, "refresh")
	if err != nil {
		return nil, fmt.Errorf("生成刷新令牌失败: %w", err)
	}

	// 记录登录设备
	recordDevice(u.ID, req, jti)

	// T20: 登录成功，清除 IP 频率限制计数器
	database.RDB.Del(ctx, limitKey)

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

// Register 用户注册（T07: 注册后自动登录，T08: 邮箱必填，T09: 密码强度校验）
func Register(req *dto.RegisterRequest, cfg *config.JWTConfig) (*dto.LoginResponse, error) {
	// 检查用户名是否已存在
	exists, err := user.ExistsByUsername(req.Username)
	if err != nil {
		return nil, fmt.Errorf("检查用户名失败: %w", err)
	}
	if exists {
		return nil, errors.New("用户名已存在")
	}

	// 检查手机号
	if req.Phone != "" {
		exists, err = user.ExistsByPhone(req.Phone)
		if err != nil {
			return nil, fmt.Errorf("检查手机号失败: %w", err)
		}
		if exists {
			return nil, errors.New("手机号已被注册")
		}
	}

	// 检查邮箱
	exists, err = user.ExistsByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("检查邮箱失败: %w", err)
	}
	if exists {
		return nil, errors.New("邮箱已被注册")
	}

	// 验证邮箱验证码（T08）
	if req.EmailCode != "" {
		if err := verifyEmailCode(req.Email, req.EmailCode); err != nil {
			return nil, err
		}
	}

	// 验证图形验证码（T30: 验证码尝试次数跟踪）
	ok, tooMany := captcha.VerifyCaptcha(req.CaptchaID, req.CaptchaCode, req.CaptchaID)
	if tooMany {
		return nil, errors.New("验证码错误次数过多，请重新获取验证码")
	}
	if !ok {
		return nil, errors.New("验证码错误")
	}

	// 加密密码
	hashedPassword, err := encrypt.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("加密密码失败: %w", err)
	}

	// 设置昵称
	nickname := req.Nickname
	if nickname == "" {
		nickname = req.Username
	}

	// T26: 敏感词过滤（注册时检查昵称）
	filter := sensitive.NewFilter()
	if filter.HasSensitive(nickname) {
		return nil, errors.New("昵称包含敏感词，请修改")
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
		return nil, fmt.Errorf("创建用户失败: %w", err)
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

	// 分配默认角色
	var roleIDs []uint
	var roleCodes []string
	defaultRoles, err := findDefaultRoles()
	if err != nil {
		logger.Errorf("查询默认角色失败: %v", err)
	} else if len(defaultRoles) > 0 {
		for _, role := range defaultRoles {
			if err := assignUserRole(u.ID, role.ID); err != nil {
				logger.Errorf("分配默认角色失败: user_id=%d, role_id=%d, err=%v", u.ID, role.ID, err)
			} else {
				roleIDs = append(roleIDs, role.ID)
				roleCodes = append(roleCodes, role.Code)
			}
		}
	}

	logger.Infof("用户注册成功: %s", req.Username)

	// 记录安全日志和设备信息
	recordSecurityLog(u.ID, "register", "用户注册成功", req.DeviceInfo)
	recordDevice(u.ID, &dto.LoginRequest{
		DeviceID:   req.DeviceID,
		DeviceInfo: req.DeviceInfo,
		UserAgent:  req.UserAgent,
	}, "")

	// T07: 注册后自动登录，返回 Token（T28: 添加 type 字段区分 token 类型）
	jti := jwt.GenerateJTI()
	accessToken, err := jwt.GenerateTokenWithJTI(u.ID, u.Username, roleIDs, jti, cfg.Expire, "access")
	if err != nil {
		return nil, fmt.Errorf("生成访问令牌失败: %w", err)
	}

	refreshJTI := jwt.GenerateJTI()
	refreshToken, err := jwt.GenerateTokenWithJTI(u.ID, u.Username, roleIDs, refreshJTI, cfg.RefreshExpire, "refresh")
	if err != nil {
		return nil, fmt.Errorf("生成刷新令牌失败: %w", err)
	}

	// 构建角色信息
	roleInfos := make([]dto.RoleInfo, 0, len(defaultRoles))
	for _, role := range defaultRoles {
		roleInfos = append(roleInfos, dto.RoleInfo{
			ID:   role.ID,
			Name: role.Name,
			Code: role.Code,
		})
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    cfg.Expire,
		User: dto.UserInfo{
			ID:        u.ID,
			Username:  u.Username,
			Nickname:  nickname,
			Roles:     roleInfos,
			CreatedAt: u.CreatedAt,
		},
	}, nil
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

// SendEmailCode 发送邮箱验证码
func SendEmailCode(emailAddr string, clientIP string) error {
	ctx := context.Background()

	// 检查每日发送上限
	dailyKey := key.EmailDailyKey(emailAddr)
	dailyCount, err := database.RDB.Incr(ctx, dailyKey).Result()
	if err != nil {
		logger.Errorf("检查每日发送次数失败: %v", err)
	} else {
		if dailyCount == 1 {
			database.RDB.Expire(ctx, dailyKey, key.DailyLimitExpire())
		}
		if dailyCount > int64(key.MaxEmailDaily) {
			return errors.New("今日发送次数已达上限")
		}
	}

	// 检查 IP 每小时发送上限
	if clientIP != "" {
		ipKey := key.EmailIPKey(clientIP)
		ipCount, err := database.RDB.Incr(ctx, ipKey).Result()
		if err != nil {
			logger.Errorf("检查 IP 发送次数失败: %v", err)
		} else {
			if ipCount == 1 {
				database.RDB.Expire(ctx, ipKey, key.EmailIPLimitExpire)
			}
			if ipCount > int64(key.MaxEmailIPHourly) {
				return errors.New("当前 IP 发送过于频繁，请稍后再试")
			}
		}
	}

	// 检查发送频率限制（1分钟内不能重复发送）
	limitKey := key.EmailLimitKey(emailAddr)
	exists, err := database.RDB.Exists(ctx, limitKey).Result()
	if err != nil {
		return fmt.Errorf("检查发送频率失败: %w", err)
	}
	if exists > 0 {
		return errors.New("发送过于频繁，请稍后再试")
	}

	// 生成验证码
	code := generateVerifyCode()

	// 存储验证码到 Redis，5分钟有效期
	codeKey := key.EmailKey(emailAddr)
	if err := database.RDB.Set(ctx, codeKey, code, key.CodeExpire*time.Second).Err(); err != nil {
		return fmt.Errorf("存储验证码失败: %w", err)
	}

	// 设置发送频率限制，1分钟内不能重复发送
	if err := database.RDB.Set(ctx, limitKey, "1", key.CodeLimitExpire*time.Second).Err(); err != nil {
		logger.Errorf("设置发送频率限制失败: %v", err)
	}

	// 发送邮件
	if err := email.SendVerificationCode(emailAddr, code); err != nil {
		logger.Errorf("发送邮箱验证码失败: %v", err)
		return errors.New("发送验证码失败，请稍后再试")
	}

	logger.Infof("邮箱验证码已发送: %s", emailAddr)
	return nil
}

// SendSmsCode 发送短信验证码
func SendSmsCode(phone string) error {
	ctx := context.Background()

	// 检查每日发送上限
	dailyKey := key.SmsDailyKey(phone)
	dailyCount, err := database.RDB.Incr(ctx, dailyKey).Result()
	if err != nil {
		logger.Errorf("检查每日发送次数失败: %v", err)
	} else {
		if dailyCount == 1 {
			database.RDB.Expire(ctx, dailyKey, key.DailyLimitExpire())
		}
		if dailyCount > int64(key.MaxSmsDaily) {
			return errors.New("今日发送次数已达上限")
		}
	}

	// 检查发送频率限制（60秒内不能重复发送）
	limitKey := key.SmsLimitKey(phone)
	exists, err := database.RDB.Exists(ctx, limitKey).Result()
	if err != nil {
		return fmt.Errorf("检查发送频率失败: %w", err)
	}
	if exists > 0 {
		return errors.New("发送过于频繁，请稍后再试")
	}

	// 生成6位验证码
	code := generateVerifyCode()

	// 存储验证码到 Redis，5分钟有效期
	codeKey := key.SmsKey(phone)
	if err := database.RDB.Set(ctx, codeKey, code, key.CodeExpire*time.Second).Err(); err != nil {
		return fmt.Errorf("存储验证码失败: %w", err)
	}

	// 设置发送频率限制，60秒内不能重复发送
	if err := database.RDB.Set(ctx, limitKey, "1", key.CodeLimitExpire*time.Second).Err(); err != nil {
		logger.Errorf("设置发送频率限制失败: %v", err)
	}

	// 发送短信
	if err := sms.SendSMS(phone, code); err != nil {
		logger.Errorf("发送短信验证码失败: %v", err)
		return errors.New("发送验证码失败，请稍后再试")
	}

	logger.Infof("短信验证码已发送: %s", phone)
	return nil
}

// GetCaptcha 获取验证码（T29: 支持 digit/letter/math 三种类型）
func GetCaptcha(captchaType captcha.CaptchaType) (string, string, error) {
	return captcha.GenerateCaptcha(captchaType)
}

// RequestPasswordReset 请求重置密码（T22: 发送邮件，T23: 验证码校验）
func RequestPasswordReset(req *dto.RequestPasswordResetRequest) error {
	// T23: 验证图形验证码（T30: 验证码尝试次数跟踪）
	ok, tooMany := captcha.VerifyCaptcha(req.CaptchaID, req.CaptchaCode, req.CaptchaID)
	if tooMany {
		return errors.New("验证码错误次数过多，请重新获取验证码")
	}
	if !ok {
		return errors.New("验证码错误")
	}

	// T21: 不暴露用户是否存在，统一返回成功
	u, err := FindByIdentifier(req.Account)
	if err != nil {
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

	// T22: 发送重置邮件（如果邮箱存在）
	if u.Email != "" {
		// TODO: 调用邮件服务发送重置链接
		logger.Infof("密码重置邮件应发送至: %s", u.Email)
	}

	logger.Infof("密码重置令牌已生成: user_id=%d", u.ID)
	return nil
}

// verifyEmailCode 验证邮箱验证码（T06/T08 辅助函数）
func verifyEmailCode(email, code string) error {
	ctx := context.Background()
	codeKey := key.EmailKey(email)
	storedCode, err := database.RDB.Get(ctx, codeKey).Result()
	if err != nil {
		return errors.New("验证码已过期或不存在")
	}

	// 检查验证尝试次数
	attemptsKey := key.EmailAttemptsKey(email, code)
	attempts, err := database.RDB.Incr(ctx, attemptsKey).Result()
	if err != nil {
		logger.Errorf("检查验证尝试次数失败: %v", err)
	} else {
		if attempts == 1 {
			database.RDB.Expire(ctx, attemptsKey, key.CodeExpire*time.Second)
		}
		if attempts > int64(key.MaxEmailCodeAttempts) {
			// 超过尝试次数，删除验证码
			database.RDB.Del(ctx, codeKey)
			database.RDB.Del(ctx, attemptsKey)
			return errors.New("验证码错误次数过多，请重新获取")
		}
	}

	if storedCode != code {
		return errors.New("验证码错误")
	}
	// 验证通过后删除验证码和尝试次数（单次有效）
	database.RDB.Del(ctx, codeKey)
	database.RDB.Del(ctx, attemptsKey)
	return nil
}

// ResetPassword 重置密码（T02: 检查历史密码，T21: 防枚举，T24: 记录安全日志）
func ResetPassword(account, code, newPassword string) error {
	// T21: 统一错误消息，不暴露用户是否存在
	u, err := user.FindByUsername(account)
	if err != nil {
		// 尝试邮箱查找
		u, err = user.FindByEmail(account)
		if err != nil {
			// 尝试手机号查找
			u, err = user.FindByPhone(account)
			if err != nil {
				return errors.New("重置令牌无效或已过期")
			}
		}
	}

	// 验证重置令牌
	ctx := context.Background()
	resetKey := key.ResetTokenKey(u.ID)
	storedToken, err := database.RDB.Get(ctx, resetKey).Result()
	if err != nil || storedToken != code {
		return errors.New("重置令牌无效或已过期")
	}

	// T02: 检查新密码是否与最近 3 次密码相同
	recentPasswords, err := user.GetRecentPasswords(u.ID, 3)
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

	// T24: 记录安全日志
	recordSecurityLog(u.ID, "password_reset", "用户重置密码", "")

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

	// T31: 检查注销取消冷却期（30天内不可再次申请注销）
	ctx := context.Background()
	cooldownKey := key.DeactivateCooldownKey(userID)
	exists, err := database.RDB.Exists(ctx, cooldownKey).Result()
	if err != nil {
		logger.Errorf("检查注销冷却期失败: %v", err)
	}
	if exists > 0 {
		return errors.New("注销取消后30天内不可再次申请注销")
	}

	// T25: 注销前前置条件检查
	// 检查是否有进行中的考试
	ongoingCount, err := exam.CountOngoingByUserID(userID)
	if err != nil {
		return fmt.Errorf("检查进行中考试失败: %w", err)
	}
	if ongoingCount > 0 {
		return errors.New("您有进行中的考试，请等待考试结束后再申请注销")
	}

	// 检查是否有待审核的角色申请
	pendingCount, err := application.CountPendingByUserID(userID)
	if err != nil {
		return fmt.Errorf("检查待审核角色申请失败: %w", err)
	}
	if pendingCount > 0 {
		return errors.New("您有待审核的角色申请，请等待审核结束后再申请注销")
	}

	// 检查是否有创建的班级
	classCount, err := classrepo.CountByCreatorID(userID)
	if err != nil {
		return fmt.Errorf("检查创建班级失败: %w", err)
	}
	if classCount > 0 {
		return errors.New("您创建的班级尚未转让，请先转让或解散班级")
	}

	// 生成确认验证码
	code := generateVerifyCode()
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

	// T31: 记录取消注销时间，设置30天冷却期
	ctx := context.Background()
	cooldownKey := key.DeactivateCooldownKey(userID)
	if err := database.RDB.Set(ctx, cooldownKey, time.Now().Unix(), key.DeactivateCooldownDuration).Err(); err != nil {
		logger.Errorf("设置注销冷却期失败: %v", err)
	}

	// 记录安全日志
	recordSecurityLog(userID, "cancel_deactivate", "用户取消注销账号", "")

	logger.Infof("用户 %d 已取消注销，30天冷却期已设置", userID)
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

// detectLoginAnomaly T16: 登录异常检测（新设备、IP 变更）
// 查询用户最近一次登录设备，对比当前 IP 和设备标识，
// 若存在异常则记录到 security_logs 表。
func detectLoginAnomaly(userID uint, req *dto.LoginRequest) {
	currentIP := req.DeviceInfo // 控制器已将 ClientIP 赋值给 DeviceInfo
	currentDeviceID := req.DeviceID

	// 查询用户最近一次登录设备
	lastDevice, err := user.FindLastDeviceByUserID(userID)
	if err != nil {
		// 无历史记录说明是首次登录，不做异常判断
		return
	}

	// 检测 IP 变更
	ipChanged := lastDevice.IP != "" && lastDevice.IP != currentIP

	// 检测新设备：在已有设备列表中查找当前 DeviceID
	isNewDevice := false
	if currentDeviceID != "" {
		_, err := user.FindDeviceByUserIDAndDeviceID(userID, currentDeviceID)
		isNewDevice = errors.Is(err, gorm.ErrRecordNotFound)
	}

	if !ipChanged && !isNewDevice {
		return // 无异常
	}

	// 构造异常详情
	var anomalyTypes []string
	if isNewDevice {
		anomalyTypes = append(anomalyTypes, "new_device")
	}
	if ipChanged {
		anomalyTypes = append(anomalyTypes, "ip_change")
	}
	anomalyType := strings.Join(anomalyTypes, ",")

	detail := fmt.Sprintf("登录异常: %s", anomalyType)

	securityLog := &model.SecurityLog{
		UserID:      userID,
		EventType:   "login_anomaly",
		EventDetail: detail,
		IP:          currentIP,
		UserAgent:   req.UserAgent,
		Status:      1,
	}

	if err := user.CreateSecurityLog(securityLog); err != nil {
		logger.Errorf("记录登录异常日志失败: %v", err)
	}

	logger.Warnf("用户 %d 登录异常: %s, IP=%s", userID, anomalyType, currentIP)
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
		return
	}

	// 将同一用户其他设备的 is_current 设为 false
	if device.ID > 0 {
		if err := database.DB.Model(&model.LoginDevice{}).
			Where("user_id = ? AND id != ?", userID, device.ID).
			Update("is_current", false).Error; err != nil {
			logger.Errorf("更新其他设备状态失败: %v", err)
		}
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

// findDefaultRoles 查询默认角色
func findDefaultRoles() ([]model.Role, error) {
	var roles []model.Role
	err := database.DB.Where("is_default = ? AND deleted_at IS NULL", true).Find(&roles).Error
	return roles, err
}

// assignUserRole 分配用户角色
func assignUserRole(userID, roleID uint) error {
	return database.DB.Exec("INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)", userID, roleID).Error
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

// VerifyCaptcha 验证验证码（T30: 支持尝试次数跟踪）
func VerifyCaptcha(captchaID, captchaCode string) error {
	ok, tooMany := captcha.VerifyCaptcha(captchaID, captchaCode, captchaID)
	if tooMany {
		return errors.New("验证码错误次数过多，请重新获取验证码")
	}
	if !ok {
		return errors.New("验证码错误")
	}
	return nil
}

// GetSecuritySettings 获取安全设置
func GetSecuritySettings(userID uint) (*dto.SecuritySettingsResponse, error) {
	privacy, err := user.FindPrivacyByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 不存在则创建默认设置
			privacy = &model.UserPrivacy{UserID: userID}
			if createErr := user.CreatePrivacy(privacy); createErr != nil {
				logger.Errorf("创建默认安全设置失败: %v", createErr)
			}
		} else {
			return nil, fmt.Errorf("获取安全设置失败: %w", err)
		}
	}

	return &dto.SecuritySettingsResponse{
		LoginNotification:    privacy.LoginNotification,
		PasswordChangeNotify: privacy.PasswordChangeNotify,
		DeviceManageNotify:   privacy.DeviceManageNotify,
	}, nil
}

// UpdateSecuritySettings 更新安全设置
func UpdateSecuritySettings(userID uint, req *dto.SecuritySettingsRequest) error {
	privacy, err := user.FindPrivacyByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 不存在则创建
			privacy = &model.UserPrivacy{UserID: userID}
			if err := user.CreatePrivacy(privacy); err != nil {
				return fmt.Errorf("创建安全设置失败: %w", err)
			}
		} else {
			return fmt.Errorf("获取安全设置失败: %w", err)
		}
	}

	// 仅更新传入的字段
	if req.LoginNotification != nil {
		privacy.LoginNotification = *req.LoginNotification
	}
	if req.PasswordChangeNotify != nil {
		privacy.PasswordChangeNotify = *req.PasswordChangeNotify
	}
	if req.DeviceManageNotify != nil {
		privacy.DeviceManageNotify = *req.DeviceManageNotify
	}

	if err := user.UpdatePrivacy(privacy); err != nil {
		return fmt.Errorf("更新安全设置失败: %w", err)
	}

	// 记录安全日志
	recordSecurityLog(userID, "security_settings_update", "用户更新安全设置", "")

	logger.Infof("用户 %d 安全设置已更新", userID)
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
