package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"gorm.io/gorm"

	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/model"
	"questionhelper-server/internal/repository/user"
	"questionhelper-server/pkg/config"
	"questionhelper-server/pkg/database"
	"questionhelper-server/pkg/jwt"
	"questionhelper-server/pkg/logger"
)

// OAuth provider 常量
const (
	ProviderGitHub = "github"
	ProviderGoogle = "google"
	ProviderWeChat = "wechat"
)

// OAuthTokenResponse OAuth 令牌响应
type OAuthTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
	Error       string `json:"error"`
}

// GitHubUser GitHub 用户信息
type GitHubUser struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

// GoogleUser Google 用户信息
type GoogleUser struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Picture string `json:"picture"`
}

// WeChatTokenResponse 微信令牌响应
type WeChatTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	OpenID       string `json:"openid"`
	UnionID      string `json:"unionid"`
	Scope        string `json:"scope"`
	ErrCode      int    `json:"errcode"`
	ErrMsg       string `json:"errmsg"`
}

// WeChatUser 微信用户信息
type WeChatUser struct {
	OpenID     string `json:"openid"`
	Nickname   string `json:"nickname"`
	Sex        int    `json:"sex"`
	Province   string `json:"province"`
	City       string `json:"city"`
	Country    string `json:"country"`
	HeadImgURL string `json:"headimgurl"`
	UnionID    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

// GetOAuthURL 生成第三方授权 URL
func GetOAuthURL(provider string) (string, error) {
	cfg := config.Cfg
	if cfg == nil {
		return "", errors.New("配置未初始化")
	}

	state, err := generateState()
	if err != nil {
		return "", fmt.Errorf("生成 state 失败: %w", err)
	}

	// 将 state 存入 Redis 用于回调验证（5分钟有效）
	ctx := context.Background()
	stateKey := "oauth:state:" + state
	if err := database.RDB.Set(ctx, stateKey, provider, 5*time.Minute).Err(); err != nil {
		logger.Errorf("存储 OAuth state 失败: %v", err)
	}

	switch provider {
	case ProviderGitHub:
		return buildGitHubURL(cfg.OAuth.GitHub, state), nil
	case ProviderGoogle:
		return buildGoogleURL(cfg.OAuth.Google, state), nil
	case ProviderWeChat:
		return buildWeChatURL(cfg.OAuth.WeChat, state), nil
	default:
		return "", fmt.Errorf("不支持的第三方登录: %s", provider)
	}
}

// OAuthLogin 第三方登录（授权码换令牌 -> 获取用户信息 -> 查找或创建本地用户 -> 返回 JWT）
func OAuthLogin(provider, code, state string) (*dto.LoginResponse, error) {
	cfg := config.Cfg
	if cfg == nil {
		return nil, errors.New("配置未初始化")
	}

	// 验证 state
	if err := verifyState(state, provider); err != nil {
		return nil, fmt.Errorf("state 验证失败: %w", err)
	}

	// 根据 provider 交换令牌并获取用户信息
	var oauthUserID, oauthUsername, oauthAvatar, oauthEmail string
	var accessToken string

	switch provider {
	case ProviderGitHub:
		token, err := exchangeGitHubCode(cfg.OAuth.GitHub, code)
		if err != nil {
			return nil, fmt.Errorf("GitHub 授权失败: %w", err)
		}
		accessToken = token
		githubUser, err := getGitHubUser(token)
		if err != nil {
			return nil, fmt.Errorf("获取 GitHub 用户信息失败: %w", err)
		}
		oauthUserID = fmt.Sprintf("%d", githubUser.ID)
		oauthUsername = githubUser.Login
		oauthAvatar = githubUser.AvatarURL
		oauthEmail = githubUser.Email

	case ProviderGoogle:
		token, err := exchangeGoogleCode(cfg.OAuth.Google, code)
		if err != nil {
			return nil, fmt.Errorf("Google 授权失败: %w", err)
		}
		accessToken = token
		googleUser, err := getGoogleUser(token)
		if err != nil {
			return nil, fmt.Errorf("获取 Google 用户信息失败: %w", err)
		}
		oauthUserID = googleUser.ID
		oauthUsername = googleUser.Name
		oauthAvatar = googleUser.Picture
		oauthEmail = googleUser.Email

	case ProviderWeChat:
		tokenResp, err := exchangeWeChatCode(cfg.OAuth.WeChat, code)
		if err != nil {
			return nil, fmt.Errorf("微信授权失败: %w", err)
		}
		accessToken = tokenResp.AccessToken
		oauthUserID = tokenResp.OpenID
		wechatUser, err := getWeChatUser(tokenResp.AccessToken, tokenResp.OpenID)
		if err != nil {
			return nil, fmt.Errorf("获取微信用户信息失败: %w", err)
		}
		oauthUsername = wechatUser.Nickname
		oauthAvatar = wechatUser.HeadImgURL
		if tokenResp.UnionID != "" {
			oauthUserID = tokenResp.UnionID
		}

	default:
		return nil, fmt.Errorf("不支持的第三方登录: %s", provider)
	}

	// 查找已有的 OAuth 绑定记录
	oauthRecord, err := user.FindOAuthUser(provider, oauthUserID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("查询 OAuth 绑定失败: %w", err)
	}

	var localUser *model.User

	if oauthRecord != nil && oauthRecord.ID > 0 {
		// 已绑定，直接登录
		localUser, err = user.FindByID(oauthRecord.UserID)
		if err != nil {
			return nil, fmt.Errorf("查找本地用户失败: %w", err)
		}

		// 更新 OAuth 令牌
		oauthRecord.AccessToken = accessToken
		oauthRecord.ProviderUsername = oauthUsername
		oauthRecord.ProviderAvatar = oauthAvatar
		if err := user.UpdateOAuthUser(oauthRecord); err != nil {
			logger.Errorf("更新 OAuth 令牌失败: %v", err)
		}
	} else {
		// 未绑定，创建新用户
		localUser, err = createOAuthUser(provider, oauthUserID, oauthUsername, oauthAvatar, oauthEmail, accessToken)
		if err != nil {
			return nil, fmt.Errorf("创建用户失败: %w", err)
		}
	}

	// 检查用户状态
	if localUser.Status == 0 {
		return nil, errors.New("账号已被禁用")
	}
	if localUser.Status == 2 {
		return nil, errors.New("账号正在注销中")
	}

	// 更新最后登录信息
	now := time.Now()
	if err := user.UpdateByID(localUser.ID, map[string]interface{}{
		"last_login_at": &now,
	}); err != nil {
		logger.Errorf("更新最后登录信息失败: %v", err)
	}

	// 记录安全日志
	recordSecurityLog(localUser.ID, "login", fmt.Sprintf("通过 %s 登录", provider), "")

	// 生成 JWT
	jwtCfg := &cfg.JWT
	jti := jwt.GenerateJTI()

	roleIDs := make([]uint, 0, len(localUser.Roles))
	for _, role := range localUser.Roles {
		roleIDs = append(roleIDs, role.ID)
	}

	accessTokenJWT, err := jwt.GenerateTokenWithJTI(localUser.ID, localUser.Username, roleIDs, jti, jwtCfg.Expire, "access")
	if err != nil {
		return nil, fmt.Errorf("生成访问令牌失败: %w", err)
	}

	refreshJTI := jwt.GenerateJTI()
	refreshTokenJWT, err := jwt.GenerateTokenWithJTI(localUser.ID, localUser.Username, roleIDs, refreshJTI, jwtCfg.RefreshExpire, "refresh")
	if err != nil {
		return nil, fmt.Errorf("生成刷新令牌失败: %w", err)
	}

	return &dto.LoginResponse{
		AccessToken:  accessTokenJWT,
		RefreshToken: refreshTokenJWT,
		TokenType:    "Bearer",
		ExpiresIn:    jwtCfg.Expire,
		User: dto.UserInfo{
			ID:        localUser.ID,
			Username:  localUser.Username,
			Nickname:  localUser.Nickname,
			Avatar:    localUser.Avatar,
			Roles:     buildRoleInfos(localUser.Roles),
			CreatedAt: localUser.CreatedAt,
		},
	}, nil
}

// OAuthBind 绑定第三方账号到已有用户
func OAuthBind(userID uint, provider, code, state string) error {
	cfg := config.Cfg
	if cfg == nil {
		return errors.New("配置未初始化")
	}

	// 验证 state
	if err := verifyState(state, provider); err != nil {
		return fmt.Errorf("state 验证失败: %w", err)
	}

	// 根据 provider 交换令牌并获取用户信息
	var oauthUserID, oauthUsername, oauthAvatar, accessToken string

	switch provider {
	case ProviderGitHub:
		token, err := exchangeGitHubCode(cfg.OAuth.GitHub, code)
		if err != nil {
			return fmt.Errorf("GitHub 授权失败: %w", err)
		}
		accessToken = token
		githubUser, err := getGitHubUser(token)
		if err != nil {
			return fmt.Errorf("获取 GitHub 用户信息失败: %w", err)
		}
		oauthUserID = fmt.Sprintf("%d", githubUser.ID)
		oauthUsername = githubUser.Login
		oauthAvatar = githubUser.AvatarURL

	case ProviderGoogle:
		token, err := exchangeGoogleCode(cfg.OAuth.Google, code)
		if err != nil {
			return fmt.Errorf("Google 授权失败: %w", err)
		}
		accessToken = token
		googleUser, err := getGoogleUser(token)
		if err != nil {
			return fmt.Errorf("获取 Google 用户信息失败: %w", err)
		}
		oauthUserID = googleUser.ID
		oauthUsername = googleUser.Name
		oauthAvatar = googleUser.Picture

	case ProviderWeChat:
		tokenResp, err := exchangeWeChatCode(cfg.OAuth.WeChat, code)
		if err != nil {
			return fmt.Errorf("微信授权失败: %w", err)
		}
		accessToken = tokenResp.AccessToken
		oauthUserID = tokenResp.OpenID
		if tokenResp.UnionID != "" {
			oauthUserID = tokenResp.UnionID
		}
		wechatUser, err := getWeChatUser(tokenResp.AccessToken, tokenResp.OpenID)
		if err != nil {
			return fmt.Errorf("获取微信用户信息失败: %w", err)
		}
		oauthUsername = wechatUser.Nickname
		oauthAvatar = wechatUser.HeadImgURL

	default:
		return fmt.Errorf("不支持的第三方登录: %s", provider)
	}

	// 检查该第三方账号是否已被其他用户绑定
	existing, err := user.FindOAuthUser(provider, oauthUserID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("查询 OAuth 绑定失败: %w", err)
	}
	if existing != nil && existing.ID > 0 {
		if existing.UserID == userID {
			return errors.New("该第三方账号已绑定当前用户")
		}
		return errors.New("该第三方账号已被其他用户绑定")
	}

	// 检查当前用户是否已绑定该 provider
	oauths, err := user.FindOAuthUsersByUserID(userID)
	if err != nil {
		return fmt.Errorf("查询用户绑定失败: %w", err)
	}
	for _, o := range oauths {
		if o.Provider == provider {
			return errors.New("当前用户已绑定该第三方账号")
		}
	}

	// 创建绑定记录
	oauthUser := &model.OAuthUser{
		UserID:           userID,
		Provider:         provider,
		ProviderType:     "web",
		ProviderUserID:   oauthUserID,
		ProviderUsername: oauthUsername,
		ProviderAvatar:   oauthAvatar,
		AccessToken:      accessToken,
	}

	if err := user.CreateOAuthUser(oauthUser); err != nil {
		return fmt.Errorf("创建绑定失败: %w", err)
	}

	// 记录安全日志
	recordSecurityLog(userID, "bind", fmt.Sprintf("绑定 %s 账号", provider), "")

	logger.Infof("用户 %d 成功绑定 %s (provider_user_id: %s)", userID, provider, oauthUserID)
	return nil
}

// OAuthUnbind 解绑第三方账号
func OAuthUnbind(userID uint, provider string) error {
	// 检查绑定是否存在
	oauths, err := user.FindOAuthUsersByUserID(userID)
	if err != nil {
		return fmt.Errorf("查询绑定失败: %w", err)
	}

	found := false
	for _, o := range oauths {
		if o.Provider == provider {
			found = true
			break
		}
	}
	if !found {
		return errors.New("未绑定该第三方账号")
	}

	// 检查用户是否有真实密码（防止解绑后无法登录）
	u, err := user.FindByID(userID)
	if err != nil {
		return errors.New("用户不存在")
	}
	if isOAuthOnlyUser(u.Password) {
		// 检查是否还有其他绑定
		if len(oauths) <= 1 {
			return errors.New("当前账号无密码且仅绑定了一个第三方账号，无法解绑。请先设置密码")
		}
	}

	// 删除绑定
	if err := user.DeleteOAuthUser(userID, provider); err != nil {
		return fmt.Errorf("解绑失败: %w", err)
	}

	// 记录安全日志
	recordSecurityLog(userID, "unbind", fmt.Sprintf("解绑 %s 账号", provider), "")

	logger.Infof("用户 %d 成功解绑 %s", userID, provider)
	return nil
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
		case ProviderWeChat:
			status.Wechat = providerStatus
		case ProviderGitHub:
			status.Github = providerStatus
		case ProviderGoogle:
			status.Google = providerStatus
		}
	}

	return status, nil
}

// ==================== GitHub OAuth ====================

func buildGitHubURL(cfg config.OAuthProvider, state string) string {
	params := url.Values{}
	params.Set("client_id", cfg.ClientID)
	params.Set("redirect_uri", cfg.RedirectURI)
	params.Set("scope", "user:email")
	params.Set("state", state)
	return "https://github.com/login/oauth/authorize?" + params.Encode()
}

func exchangeGitHubCode(cfg config.OAuthProvider, code string) (string, error) {
	data := url.Values{}
	data.Set("client_id", cfg.ClientID)
	data.Set("client_secret", cfg.ClientSecret)
	data.Set("code", code)

	req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token", strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求 GitHub 令牌失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	var tokenResp OAuthTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", fmt.Errorf("解析令牌响应失败: %w", err)
	}

	if tokenResp.Error != "" {
		return "", fmt.Errorf("GitHub 返回错误: %s", tokenResp.Error)
	}
	if tokenResp.AccessToken == "" {
		return "", errors.New("GitHub 未返回 access_token")
	}

	return tokenResp.AccessToken, nil
}

func getGitHubUser(accessToken string) (*GitHubUser, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求 GitHub 用户信息失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var githubUser GitHubUser
	if err := json.Unmarshal(body, &githubUser); err != nil {
		return nil, fmt.Errorf("解析用户信息失败: %w", err)
	}

	if githubUser.ID == 0 {
		return nil, errors.New("GitHub 返回无效用户信息")
	}

	// 如果没有邮箱，尝试通过 API 获取
	if githubUser.Email == "" {
		email, err := getGitHubEmail(accessToken)
		if err == nil {
			githubUser.Email = email
		}
	}

	return &githubUser, nil
}

func getGitHubEmail(accessToken string) (string, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}
	if err := json.Unmarshal(body, &emails); err != nil {
		return "", err
	}

	for _, e := range emails {
		if e.Primary && e.Verified {
			return e.Email, nil
		}
	}
	return "", errors.New("no primary email found")
}

// ==================== Google OAuth ====================

func buildGoogleURL(cfg config.OAuthProvider, state string) string {
	params := url.Values{}
	params.Set("client_id", cfg.ClientID)
	params.Set("redirect_uri", cfg.RedirectURI)
	params.Set("response_type", "code")
	params.Set("scope", "openid email profile")
	params.Set("state", state)
	return "https://accounts.google.com/o/oauth2/v2/auth?" + params.Encode()
}

func exchangeGoogleCode(cfg config.OAuthProvider, code string) (string, error) {
	data := url.Values{}
	data.Set("client_id", cfg.ClientID)
	data.Set("client_secret", cfg.ClientSecret)
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", cfg.RedirectURI)

	req, err := http.NewRequest("POST", "https://oauth2.googleapis.com/token", strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求 Google 令牌失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		Error       string `json:"error"`
		ErrorDesc   string `json:"error_description"`
	}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", fmt.Errorf("解析令牌响应失败: %w", err)
	}

	if tokenResp.Error != "" {
		return "", fmt.Errorf("Google 返回错误: %s - %s", tokenResp.Error, tokenResp.ErrorDesc)
	}
	if tokenResp.AccessToken == "" {
		return "", errors.New("Google 未返回 access_token")
	}

	return tokenResp.AccessToken, nil
}

func getGoogleUser(accessToken string) (*GoogleUser, error) {
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求 Google 用户信息失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var googleUser GoogleUser
	if err := json.Unmarshal(body, &googleUser); err != nil {
		return nil, fmt.Errorf("解析用户信息失败: %w", err)
	}

	if googleUser.ID == "" {
		return nil, errors.New("Google 返回无效用户信息")
	}

	return &googleUser, nil
}

// ==================== WeChat OAuth ====================

func buildWeChatURL(cfg config.OAuthProvider, state string) string {
	params := url.Values{}
	params.Set("appid", cfg.ClientID)
	params.Set("redirect_uri", cfg.RedirectURI)
	params.Set("response_type", "code")
	params.Set("scope", "snsapi_login")
	params.Set("state", state)
	return "https://open.weixin.qq.com/connect/qrconnect?" + params.Encode() + "#wechat_redirect"
}

func exchangeWeChatCode(cfg config.OAuthProvider, code string) (*WeChatTokenResponse, error) {
	params := url.Values{}
	params.Set("appid", cfg.ClientID)
	params.Set("secret", cfg.ClientSecret)
	params.Set("code", code)
	params.Set("grant_type", "authorization_code")

	apiURL := "https://api.weixin.qq.com/sns/oauth2/access_token?" + params.Encode()

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("请求微信令牌失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var tokenResp WeChatTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("解析令牌响应失败: %w", err)
	}

	if tokenResp.ErrCode != 0 {
		return nil, fmt.Errorf("微信返回错误: %d - %s", tokenResp.ErrCode, tokenResp.ErrMsg)
	}
	if tokenResp.AccessToken == "" {
		return nil, errors.New("微信未返回 access_token")
	}

	return &tokenResp, nil
}

func getWeChatUser(accessToken, openID string) (*WeChatUser, error) {
	params := url.Values{}
	params.Set("access_token", accessToken)
	params.Set("openid", openID)
	params.Set("lang", "zh_CN")

	apiURL := "https://api.weixin.qq.com/sns/userinfo?" + params.Encode()

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("请求微信用户信息失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var wechatUser WeChatUser
	if err := json.Unmarshal(body, &wechatUser); err != nil {
		return nil, fmt.Errorf("解析用户信息失败: %w", err)
	}

	if wechatUser.ErrCode != 0 {
		return nil, fmt.Errorf("微信返回错误: %d - %s", wechatUser.ErrCode, wechatUser.ErrMsg)
	}

	return &wechatUser, nil
}

// ==================== 辅助函数 ====================

// generateState 生成随机 state 字符串
func generateState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// verifyState 验证 OAuth state（防 CSRF）
func verifyState(state, provider string) error {
	if state == "" {
		return errors.New("state 参数缺失")
	}

	ctx := context.Background()
	stateKey := "oauth:state:" + state

	storedProvider, err := database.RDB.Get(ctx, stateKey).Result()
	if err != nil {
		return errors.New("state 已过期或无效")
	}

	// 验证后删除 state（一次性使用）
	database.RDB.Del(ctx, stateKey)

	if storedProvider != provider {
		return errors.New("state 与 provider 不匹配")
	}

	return nil
}

// createOAuthUser 创建 OAuth 新用户
func createOAuthUser(provider, oauthUserID, oauthUsername, oauthAvatar, oauthEmail, accessToken string) (*model.User, error) {
	// 生成唯一用户名
	username := generateOAuthUsername(provider, oauthUsername)

	// 检查用户名是否已存在，如果存在则添加后缀
	exists, _ := user.ExistsByUsername(username)
	if exists {
		suffix, _ := rand.Read(make([]byte, 3))
		username = fmt.Sprintf("%s_%x", username, suffix)
	}

	// 创建用户（OAuth 用户使用随机密码，用户后续可设置真实密码）
	randomPassword, _ := rand.Read(make([]byte, 16))
	u := &model.User{
		Username:       username,
		Password:       fmt.Sprintf("oauth_%x", randomPassword), // OAuth 占位密码，不可用于登录
		Nickname:       oauthUsername,
		Avatar:         oauthAvatar,
		Email:          oauthEmail,
		Status:         1,
		RegisterSource: "oauth",
	}

	if err := user.Create(u); err != nil {
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

	// 创建默认隐私设置
	privacy := &model.UserPrivacy{UserID: u.ID}
	if err := user.CreatePrivacy(privacy); err != nil {
		logger.Errorf("创建隐私设置失败: %v", err)
	}

	// 创建 OAuth 绑定记录
	oauthUser := &model.OAuthUser{
		UserID:           u.ID,
		Provider:         provider,
		ProviderType:     "web",
		ProviderUserID:   oauthUserID,
		ProviderUsername: oauthUsername,
		ProviderAvatar:   oauthAvatar,
		AccessToken:      accessToken,
	}
	if err := user.CreateOAuthUser(oauthUser); err != nil {
		logger.Errorf("创建 OAuth 绑定失败: %v", err)
	}

	logger.Infof("通过 %s 创建新用户: %s (id: %d)", provider, username, u.ID)
	return u, nil
}

// generateOAuthUsername 根据 provider 和用户名生成本地用户名
func generateOAuthUsername(provider, oauthUsername string) string {
	prefix := provider
	if oauthUsername != "" {
		// 清理用户名中的特殊字符
		cleaned := strings.Map(func(r rune) rune {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '-' {
				return r
			}
			return '_'
		}, oauthUsername)
		if cleaned != "" {
			return prefix + "_" + cleaned
		}
	}
	return fmt.Sprintf("%s_user_%d", prefix, time.Now().UnixMilli()%1000000)
}

// isOAuthOnlyUser 判断用户是否仅通过 OAuth 登录（无真实密码）
func isOAuthOnlyUser(password string) bool {
	return strings.HasPrefix(password, "oauth_")
}
