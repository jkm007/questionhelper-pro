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

// SwitchRole 切换角色并重新签发令牌
func (ac *AuthController) SwitchRole(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.SwitchRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 从 Authorization 头获取当前 token
	currentToken := c.GetHeader("Authorization")
	if len(currentToken) > 7 {
		currentToken = currentToken[7:] // Remove "Bearer " prefix
	}

	result, err := auth.SwitchRole(userID, req.RoleID, ac.jwtCfg, currentToken)
	if err != nil {
		errMsg := err.Error()
		switch {
		case strings.Contains(errMsg, "用户不存在"):
			response.ErrorWithCode(c, http.StatusNotFound, apperr.ErrUserNotFound.Code, errMsg)
		case strings.Contains(errMsg, "禁用"):
			response.ErrorWithCode(c, http.StatusForbidden, apperr.ErrAccountDisabled.Code, errMsg)
		case strings.Contains(errMsg, "无此角色"):
			response.ErrorWithCode(c, http.StatusForbidden, apperr.ErrForbidden.Code, errMsg)
		default:
			response.ErrorWithCode(c, http.StatusBadRequest, apperr.ErrParam.Code, errMsg)
		}
		return
	}

	// T15: set refresh token as httpOnly cookie for web clients
	c.SetCookie("refresh_token", result.RefreshToken, 7*24*60*60, "/", "", false, true)

	response.Success(c, result)
}

// RefreshToken 刷新令牌
func (ac *AuthController) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	result, err := auth.RefreshToken(&req, ac.jwtCfg)
	if err != nil {
		errMsg := err.Error()
		switch {
		case strings.Contains(errMsg, "令牌"):
			response.ErrorWithCode(c, http.StatusUnauthorized, apperr.ErrRefreshTokenInvalid.Code, errMsg)
		case strings.Contains(errMsg, "用户不存在"):
			response.ErrorWithCode(c, http.StatusUnauthorized, apperr.ErrUserNotFound.Code, errMsg)
		case strings.Contains(errMsg, "禁用"):
			response.ErrorWithCode(c, http.StatusForbidden, apperr.ErrAccountDisabled.Code, errMsg)
		default:
			response.ErrorWithCode(c, http.StatusUnauthorized, apperr.ErrUnauthorized.Code, errMsg)
		}
		return
	}

	// T15: set new refresh token as httpOnly cookie for web clients
	c.SetCookie("refresh_token", result.RefreshToken, 7*24*60*60, "/", "", false, true)

	response.Success(c, result)
}
