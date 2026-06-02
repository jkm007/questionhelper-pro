package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/service/auth"
	apperr "questionhelper-server/pkg/errors"
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

// GetOAuthURL 获取第三方授权 URL
func (ac *AuthController) GetOAuthURL(c *gin.Context) {
	provider := c.Param("provider")

	url, err := auth.GetOAuthURL(provider)
	if err != nil {
		errMsg := err.Error()
		switch {
		case strings.Contains(errMsg, "不支持"):
			response.ErrorWithCode(c, http.StatusBadRequest, apperr.ErrParam.Code, errMsg)
		default:
			response.ErrorWithCode(c, http.StatusInternalServerError, apperr.ErrInternal.Code, errMsg)
		}
		return
	}

	response.Success(c, gin.H{"url": url})
}

// OAuthLogin 第三方登录（通过授权码）
func (ac *AuthController) OAuthLogin(c *gin.Context) {
	provider := c.Param("provider")

	var req struct {
		Code  string `json:"code" binding:"required"`
		State string `json:"state"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	result, err := auth.OAuthLogin(provider, req.Code, req.State)
	if err != nil {
		errMsg := err.Error()
		switch {
		case strings.Contains(errMsg, "不支持"):
			response.ErrorWithCode(c, http.StatusBadRequest, apperr.ErrParam.Code, errMsg)
		case strings.Contains(errMsg, "state"):
			response.ErrorWithCode(c, http.StatusBadRequest, apperr.ErrParam.Code, errMsg)
		case strings.Contains(errMsg, "禁用") || strings.Contains(errMsg, "注销中"):
			response.ErrorWithCode(c, http.StatusForbidden, apperr.ErrAccountDisabled.Code, errMsg)
		default:
			response.ErrorWithCode(c, http.StatusInternalServerError, apperr.ErrInternal.Code, errMsg)
		}
		return
	}

	response.Success(c, result)
}

// OAuthCallback 第三方登录回调（处理授权码重定向）
func (ac *AuthController) OAuthCallback(c *gin.Context) {
	provider := c.Param("provider")
	code := c.Query("code")
	state := c.Query("state")

	if code == "" {
		response.Error(c, http.StatusBadRequest, "缺少授权码")
		return
	}

	result, err := auth.OAuthLogin(provider, code, state)
	if err != nil {
		errMsg := err.Error()
		switch {
		case strings.Contains(errMsg, "不支持"):
			response.ErrorWithCode(c, http.StatusBadRequest, apperr.ErrParam.Code, errMsg)
		case strings.Contains(errMsg, "state"):
			response.ErrorWithCode(c, http.StatusBadRequest, apperr.ErrParam.Code, errMsg)
		case strings.Contains(errMsg, "禁用") || strings.Contains(errMsg, "注销中"):
			response.ErrorWithCode(c, http.StatusForbidden, apperr.ErrAccountDisabled.Code, errMsg)
		default:
			response.ErrorWithCode(c, http.StatusInternalServerError, apperr.ErrInternal.Code, errMsg)
		}
		return
	}

	response.Success(c, result)
}

// OAuthBind 绑定第三方账号（已认证用户）
func (ac *AuthController) OAuthBind(c *gin.Context) {
	userID := c.GetUint("user_id")
	provider := c.Param("provider")

	var req dto.OAuthBindRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	if err := auth.OAuthBind(userID, provider, req.Code, req.State); err != nil {
		errMsg := err.Error()
		switch {
		case strings.Contains(errMsg, "不支持"):
			response.ErrorWithCode(c, http.StatusBadRequest, apperr.ErrParam.Code, errMsg)
		case strings.Contains(errMsg, "已绑定"):
			response.ErrorWithCode(c, http.StatusConflict, apperr.ErrUserExists.Code, errMsg)
		case strings.Contains(errMsg, "已被其他用户"):
			response.ErrorWithCode(c, http.StatusConflict, apperr.ErrUserExists.Code, errMsg)
		default:
			response.ErrorWithCode(c, http.StatusInternalServerError, apperr.ErrInternal.Code, errMsg)
		}
		return
	}

	response.SuccessWithMessage(c, "绑定成功", nil)
}

// OAuthUnbind 解绑第三方账号（已认证用户）
func (ac *AuthController) OAuthUnbind(c *gin.Context) {
	userID := c.GetUint("user_id")
	provider := c.Param("provider")

	if err := auth.OAuthUnbind(userID, provider); err != nil {
		errMsg := err.Error()
		switch {
		case strings.Contains(errMsg, "未绑定"):
			response.ErrorWithCode(c, http.StatusNotFound, apperr.ErrNotFound.Code, errMsg)
		case strings.Contains(errMsg, "无法解绑"):
			response.ErrorWithCode(c, http.StatusBadRequest, apperr.ErrParam.Code, errMsg)
		default:
			response.ErrorWithCode(c, http.StatusInternalServerError, apperr.ErrInternal.Code, errMsg)
		}
		return
	}

	response.SuccessWithMessage(c, "解绑成功", nil)
}
