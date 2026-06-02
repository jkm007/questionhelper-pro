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

// @Summary      获取第三方绑定状态
// @Description  获取当前用户已绑定的第三方账号状态
// @Tags         用户认证
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "获取失败"
// @Router       /user/oauth/status [get]
// @Security     BearerAuth
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

// @Summary      获取第三方授权 URL
// @Description  获取指定第三方平台的 OAuth 授权跳转 URL
// @Tags         认证
// @Accept       json
// @Produce      json
// @Param        provider  path      string  true  "第三方平台(wechat/github/google)"
// @Success      200       {object}  response.Response  "成功"
// @Failure      400       {object}  response.Response  "不支持的平台"
// @Router       /auth/oauth/{provider}/url [get]
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

// @Summary      第三方登录
// @Description  通过第三方平台授权码登录，未注册用户将自动注册
// @Tags         认证
// @Accept       json
// @Produce      json
// @Param        provider  path      string  true   "第三方平台(wechat/github/google)"
// @Param        data      body      object  true   "授权码请求体"  Schema({"code":"string","state":"string"})
// @Success      200       {object}  response.Response  "登录成功"
// @Failure      400       {object}  response.Response  "参数错误"
// @Failure      403       {object}  response.Response  "账号被禁用"
// @Router       /auth/oauth/{provider} [post]
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

// @Summary      第三方登录回调
// @Description  处理第三方平台 OAuth 授权码回调重定向
// @Tags         认证
// @Accept       json
// @Produce      json
// @Param        provider  path      string  true  "第三方平台(wechat/github/google)"
// @Param        code      query     string  true  "授权码"
// @Param        state     query     string  false "状态参数"
// @Success      200       {object}  response.Response  "登录成功"
// @Failure      400       {object}  response.Response  "参数错误"
// @Failure      403       {object}  response.Response  "账号被禁用"
// @Router       /auth/oauth/{provider}/callback [get]
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

// @Summary      绑定第三方账号
// @Description  已认证用户绑定第三方平台账号
// @Tags         用户认证
// @Accept       json
// @Produce      json
// @Param        provider  path      string              true  "第三方平台(wechat/github/google)"
// @Param        data      body      dto.OAuthBindRequest  true  "绑定请求体"
// @Success      200       {object}  response.Response    "绑定成功"
// @Failure      400       {object}  response.Response    "参数错误"
// @Failure      409       {object}  response.Response    "已绑定或已被其他用户绑定"
// @Router       /user/oauth/bind/{provider} [post]
// @Security     BearerAuth
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

// @Summary      解绑第三方账号
// @Description  已认证用户解绑第三方平台账号
// @Tags         用户认证
// @Accept       json
// @Produce      json
// @Param        provider  path      string  true  "第三方平台(wechat/github/google)"
// @Success      200       {object}  response.Response  "解绑成功"
// @Failure      400       {object}  response.Response  "参数错误"
// @Failure      404       {object}  response.Response  "未绑定该平台"
// @Router       /user/oauth/unbind/{provider} [delete]
// @Security     BearerAuth
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
