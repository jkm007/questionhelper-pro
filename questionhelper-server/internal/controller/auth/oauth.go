package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"questionhelper-server/internal/service/auth"
	"questionhelper-server/pkg/response"
)

// GetOAuthStatus 获取第三方绑定状态
func (ac *AuthController) GetOAuthStatus(c *gin.Context) {
	userID := c.GetUint("user_id")

	status, err := auth.GetOAuthStatus(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, status)
}

// GetOAuthURL 获取第三方授权 URL（占位）
func (ac *AuthController) GetOAuthURL(c *gin.Context) {
	provider := c.Param("provider")

	// TODO: 实现第三方 OAuth
	response.Success(c, gin.H{
		"url": "https://example.com/oauth/" + provider,
	})
}

// OAuthLogin 第三方登录（占位）
func (ac *AuthController) OAuthLogin(c *gin.Context) {
	provider := c.Param("provider")

	// TODO: 实现第三方 OAuth
	response.Error(c, http.StatusNotImplemented, provider+" 登录暂未实现")
}

// OAuthBind 绑定第三方账号（占位）
func (ac *AuthController) OAuthBind(c *gin.Context) {
	provider := c.Param("provider")

	// TODO: 实现第三方绑定
	response.Error(c, http.StatusNotImplemented, provider+" 绑定暂未实现")
}

// OAuthUnbind 解绑第三方账号（占位）
func (ac *AuthController) OAuthUnbind(c *gin.Context) {
	provider := c.Param("provider")

	// TODO: 实现第三方解绑
	response.Error(c, http.StatusNotImplemented, provider+" 解绑暂未实现")
}
