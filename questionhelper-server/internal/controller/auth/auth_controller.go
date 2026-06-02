package auth

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/service/auth"
	"questionhelper-server/pkg/config"
	apperr "questionhelper-server/pkg/errors"
	"questionhelper-server/pkg/response"
)

type AuthController struct {
	jwtCfg *config.JWTConfig
}

func NewAuthController(jwtCfg *config.JWTConfig) *AuthController {
	return &AuthController{jwtCfg: jwtCfg}
}

// @Summary      用户登录
// @Description  用户通过用户名和密码登录，支持图形验证码验证
// @Tags         认证
// @Accept       json
// @Produce      json
// @Param        data  body      dto.LoginRequest  true  "登录请求体"
// @Success      200   {object}  response.Response  "成功"
// @Failure      400   {object}  response.Response  "参数错误"
// @Failure      401   {object}  response.Response  "认证失败"
// @Failure      403   {object}  response.Response  "账号被禁用"
// @Failure      429   {object}  response.Response  "请求过于频繁"
// @Router       /auth/login [post]
// Login 用户登录
func (ac *AuthController) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 获取客户端信息
	req.DeviceInfo = c.ClientIP()
	req.UserAgent = c.GetHeader("User-Agent")

	result, err := auth.Login(&req, ac.jwtCfg)
	if err != nil {
		errMsg := err.Error()
		switch {
		case strings.Contains(errMsg, "频繁"):
			response.ErrorWithCode(c, http.StatusTooManyRequests, apperr.ErrTooManyRequests.Code, errMsg)
		case strings.Contains(errMsg, "验证码"):
			response.ErrorWithCode(c, http.StatusUnauthorized, apperr.ErrCaptchaWrong.Code, errMsg)
		case strings.Contains(errMsg, "禁用") || strings.Contains(errMsg, "注销中"):
			response.ErrorWithCode(c, http.StatusForbidden, apperr.ErrAccountDisabled.Code, errMsg)
		case strings.Contains(errMsg, "锁定"):
			response.ErrorWithCode(c, http.StatusForbidden, apperr.ErrAccountDisabled.Code, errMsg)
		default:
			response.ErrorWithCode(c, http.StatusUnauthorized, apperr.ErrUnauthorized.Code, errMsg)
		}
		return
	}

	// T15: set refresh token as httpOnly cookie for web clients
	c.SetCookie("refresh_token", result.RefreshToken, 7*24*60*60, "/", "", false, true)

	response.Success(c, result)
}

// @Summary      用户注册
// @Description  新用户注册，注册成功后自动登录并返回令牌
// @Tags         认证
// @Accept       json
// @Produce      json
// @Param        data  body      dto.RegisterRequest  true  "注册请求体"
// @Success      200   {object}  response.Response  "注册成功"
// @Failure      400   {object}  response.Response  "参数错误"
// @Failure      409   {object}  response.Response  "用户已存在"
// @Router       /auth/register [post]
// Register 用户注册（T07: 注册后自动登录）
func (ac *AuthController) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	result, err := auth.Register(&req, ac.jwtCfg)
	if err != nil {
		errMsg := err.Error()
		switch {
		case strings.Contains(errMsg, "已存在") || strings.Contains(errMsg, "已被注册"):
			response.ErrorWithCode(c, http.StatusConflict, apperr.ErrUserExists.Code, errMsg)
		case strings.Contains(errMsg, "验证码已过期") || strings.Contains(errMsg, "验证码已过期或不存在"):
			response.ErrorWithCode(c, http.StatusBadRequest, apperr.ErrCaptchaExpired.Code, errMsg)
		case strings.Contains(errMsg, "验证码"):
			response.ErrorWithCode(c, http.StatusBadRequest, apperr.ErrCaptchaWrong.Code, errMsg)
		default:
			response.ErrorWithCode(c, http.StatusBadRequest, apperr.ErrParam.Code, errMsg)
		}
		return
	}

	// T15: set refresh token as httpOnly cookie for web clients
	c.SetCookie("refresh_token", result.RefreshToken, 7*24*60*60, "/", "", false, true)

	response.SuccessWithMessage(c, "注册成功", result)
}

// @Summary      用户退出
// @Description  退出当前设备登录，清除刷新令牌
// @Tags         用户认证
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response  "退出成功"
// @Failure      500  {object}  response.Response  "退出失败"
// @Router       /user/logout [post]
// @Security     BearerAuth
// Logout 用户退出
func (ac *AuthController) Logout(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if len(token) > 7 {
		token = token[7:] // Remove "Bearer " prefix
	}

	if err := auth.Logout(token); err != nil {
		response.Error(c, http.StatusInternalServerError, "退出失败")
		return
	}

	// T15: clear refresh token cookie on logout
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	response.SuccessWithMessage(c, "退出成功", nil)
}

// @Summary      退出所有设备
// @Description  退出当前用户的所有已登录设备
// @Tags         用户认证
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response  "已退出所有设备"
// @Failure      500  {object}  response.Response  "退出失败"
// @Router       /user/logout-all [post]
// @Security     BearerAuth
// LogoutAll 退出所有设备
func (ac *AuthController) LogoutAll(c *gin.Context) {
	userID := c.GetUint("user_id")

	if err := auth.LogoutAll(userID); err != nil {
		response.Error(c, http.StatusInternalServerError, "退出失败")
		return
	}

	// T15: clear refresh token cookie on logout all
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	response.SuccessWithMessage(c, "已退出所有设备", nil)
}

// @Summary      请求重置密码
// @Description  通过用户名/手机号/邮箱申请重置密码，需提供图形验证码
// @Tags         认证
// @Accept       json
// @Produce      json
// @Param        data  body      dto.RequestPasswordResetRequest  true  "重置密码请求体"
// @Success      200   {object}  response.Response  "重置密码请求已提交"
// @Failure      400   {object}  response.Response  "参数错误"
// @Router       /auth/password/reset-request [post]
// RequestPasswordReset 请求重置密码
func (ac *AuthController) RequestPasswordReset(c *gin.Context) {
	var req dto.RequestPasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	if err := auth.RequestPasswordReset(&req); err != nil {
		errMsg := err.Error()
		switch {
		case strings.Contains(errMsg, "验证码"):
			response.ErrorWithCode(c, http.StatusBadRequest, apperr.ErrCaptchaWrong.Code, errMsg)
		default:
			response.Error(c, http.StatusInternalServerError, errMsg)
		}
		return
	}

	response.SuccessWithMessage(c, "重置密码请求已提交，请查看手机或邮箱", nil)
}

// @Summary      重置密码
// @Description  通过验证码重置密码
// @Tags         认证
// @Accept       json
// @Produce      json
// @Param        data  body      dto.ResetPasswordRequest  true  "重置密码请求体"
// @Success      200   {object}  response.Response  "密码重置成功"
// @Failure      400   {object}  response.Response  "参数错误"
// @Router       /auth/password/reset [post]
// ResetPassword 重置密码
func (ac *AuthController) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	if err := auth.ResetPassword(req.Account, req.Code, req.NewPassword); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessWithMessage(c, "密码重置成功", nil)
}

// @Summary      申请注销账号
// @Description  已认证用户申请注销账号，需要验证密码
// @Tags         用户认证
// @Accept       json
// @Produce      json
// @Param        data  body      dto.DeactivateRequest  true  "注销请求体"
// @Success      200   {object}  response.Response  "验证码已发送"
// @Failure      400   {object}  response.Response  "参数错误"
// @Failure      404   {object}  response.Response  "用户不存在"
// @Failure      403   {object}  response.Response  "冷却期未过"
// @Router       /user/account/deactivate [post]
// @Security     BearerAuth
// DeactivateAccount 申请注销账号
func (ac *AuthController) DeactivateAccount(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.DeactivateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	if err := auth.DeactivateAccount(userID, req.Password); err != nil {
		errMsg := err.Error()
		switch {
		case strings.Contains(errMsg, "用户不存在"):
			response.ErrorWithCode(c, http.StatusNotFound, apperr.ErrUserNotFound.Code, errMsg)
		case strings.Contains(errMsg, "密码错误"):
			response.ErrorWithCode(c, http.StatusBadRequest, apperr.ErrPasswordWrong.Code, errMsg)
		case strings.Contains(errMsg, "30天内不可再次申请注销"):
			response.ErrorWithCode(c, http.StatusForbidden, apperr.ErrDeactivateCooldown.Code, errMsg)
		default:
			response.ErrorWithCode(c, http.StatusBadRequest, apperr.ErrParam.Code, errMsg)
		}
		return
	}

	response.Success(c, gin.H{
		"verifyType": "sms",
		"expireIn":   300,
	})
}

// @Summary      确认注销
// @Description  通过验证码确认注销账号，30天后将永久删除
// @Tags         用户认证
// @Accept       json
// @Produce      json
// @Param        data  body      dto.DeactivateConfirmRequest  true  "确认注销请求体"
// @Success      200   {object}  response.Response  "注销申请已提交"
// @Failure      400   {object}  response.Response  "验证码错误或已过期"
// @Router       /user/account/deactivate/confirm [post]
// @Security     BearerAuth
// ConfirmDeactivate 确认注销
func (ac *AuthController) ConfirmDeactivate(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.DeactivateConfirmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	if err := auth.ConfirmDeactivate(userID, req.Code); err != nil {
		response.ErrorWithCode(c, http.StatusBadRequest, apperr.ErrCaptchaExpired.Code, err.Error())
		return
	}

	// T15: clear refresh token cookie on account deactivation
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	response.SuccessWithMessage(c, "注销申请已提交，30天后将永久删除账号", nil)
}

// @Summary      取消注销
// @Description  取消已提交的账号注销申请
// @Tags         用户认证
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response  "已取消注销"
// @Failure      400  {object}  response.Response  "参数错误"
// @Failure      404  {object}  response.Response  "用户不存在"
// @Router       /user/account/cancel-deactivate [post]
// @Security     BearerAuth
// CancelDeactivate 取消注销
func (ac *AuthController) CancelDeactivate(c *gin.Context) {
	userID := c.GetUint("user_id")

	if err := auth.CancelDeactivate(userID); err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "用户不存在") {
			response.ErrorWithCode(c, http.StatusNotFound, apperr.ErrUserNotFound.Code, errMsg)
		} else {
			response.ErrorWithCode(c, http.StatusBadRequest, apperr.ErrParam.Code, errMsg)
		}
		return
	}

	response.SuccessWithMessage(c, "已取消注销", nil)
}

// @Summary      导出个人数据
// @Description  导出当前用户的个人数据（考试、练习、错题、收藏、评论等）
// @Tags         用户认证
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "导出失败"
// @Router       /user/account/export [get]
// @Security     BearerAuth
// ExportUserData 导出个人数据
func (ac *AuthController) ExportUserData(c *gin.Context) {
	userID := c.GetUint("user_id")

	data, err := auth.ExportUserData(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, data)
}

// @Summary      获取个人信息
// @Description  获取当前登录用户的个人资料
// @Tags         用户认证
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "获取失败"
// @Router       /user/profile [get]
// @Security     BearerAuth
// GetProfile 获取个人信息
func (ac *AuthController) GetProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	user, err := auth.GetUserProfile(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, user)
}

// @Summary      获取登录设备列表
// @Description  获取当前用户的所有已登录设备信息
// @Tags         用户认证
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "获取失败"
// @Router       /user/devices [get]
// @Security     BearerAuth
// GetUserDevices 获取登录设备列表
func (ac *AuthController) GetUserDevices(c *gin.Context) {
	userID := c.GetUint("user_id")

	devices, err := auth.GetUserDevices(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{"devices": devices})
}

// @Summary      踢出设备
// @Description  将指定设备踢出登录
// @Tags         用户认证
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "设备ID"
// @Success      200  {object}  response.Response  "设备已退出"
// @Failure      400  {object}  response.Response  "设备ID无效或操作失败"
// @Router       /user/devices/{id} [delete]
// @Security     BearerAuth
// KickDevice 踢出设备
func (ac *AuthController) KickDevice(c *gin.Context) {
	userID := c.GetUint("user_id")
	deviceIDStr := c.Param("id")
	deviceID, err := strconv.ParseUint(deviceIDStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "设备ID无效")
		return
	}

	if err := auth.KickDevice(userID, uint(deviceID)); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessWithMessage(c, "设备已退出", nil)
}

// @Summary      踢出所有设备
// @Description  将当前用户的所有设备踢出登录
// @Tags         用户认证
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response  "已退出所有设备"
// @Failure      500  {object}  response.Response  "操作失败"
// @Router       /user/devices [delete]
// @Security     BearerAuth
// KickAllDevices 踢出所有设备
func (ac *AuthController) KickAllDevices(c *gin.Context) {
	userID := c.GetUint("user_id")

	if err := auth.LogoutAll(userID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	// T15: clear refresh token cookie when kicking all devices
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	response.SuccessWithMessage(c, "已退出所有设备", nil)
}

// @Summary      获取安全日志
// @Description  分页获取当前用户的安全操作日志
// @Tags         用户认证
// @Accept       json
// @Produce      json
// @Param        page      query  int  false  "页码"       default(1)
// @Param        pageSize  query  int  false  "每页数量"   default(10)
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "获取失败"
// @Router       /user/security/logs [get]
// @Security     BearerAuth
// GetSecurityLogs 获取安全日志
func (ac *AuthController) GetSecurityLogs(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	logs, total, err := auth.GetSecurityLogs(userID, page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":     logs,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// @Summary      修改密码
// @Description  已认证用户修改登录密码，需要验证旧密码
// @Tags         用户认证
// @Accept       json
// @Produce      json
// @Param        data  body      dto.ChangePasswordRequest  true  "修改密码请求体"
// @Success      200   {object}  response.Response  "密码修改成功"
// @Failure      400   {object}  response.Response  "参数错误"
// @Failure      404   {object}  response.Response  "用户不存在"
// @Router       /user/password [put]
// @Security     BearerAuth
// ChangePassword 修改密码
func (ac *AuthController) ChangePassword(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	if err := auth.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		errMsg := err.Error()
		switch {
		case strings.Contains(errMsg, "用户不存在"):
			response.ErrorWithCode(c, http.StatusNotFound, apperr.ErrUserNotFound.Code, errMsg)
		case strings.Contains(errMsg, "密码错误"):
			response.ErrorWithCode(c, http.StatusBadRequest, apperr.ErrPasswordWrong.Code, errMsg)
		default:
			response.ErrorWithCode(c, http.StatusBadRequest, apperr.ErrParam.Code, errMsg)
		}
		return
	}

	response.SuccessWithMessage(c, "密码修改成功", nil)
}

// @Summary      获取安全设置
// @Description  获取当前用户的安全通知设置
// @Tags         用户认证
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response  "成功"
// @Failure      500  {object}  response.Response  "获取失败"
// @Router       /user/security/settings [get]
// @Security     BearerAuth
// GetSecuritySettings 获取安全设置
func (ac *AuthController) GetSecuritySettings(c *gin.Context) {
	userID := c.GetUint("user_id")

	settings, err := auth.GetSecuritySettings(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, settings)
}

// @Summary      更新安全设置
// @Description  更新当前用户的安全通知设置
// @Tags         用户认证
// @Accept       json
// @Produce      json
// @Param        data  body      dto.SecuritySettingsRequest  true  "安全设置请求体"
// @Success      200   {object}  response.Response  "安全设置已更新"
// @Failure      400   {object}  response.Response  "参数错误"
// @Router       /user/security/settings [put]
// @Security     BearerAuth
// UpdateSecuritySettings 更新安全设置
func (ac *AuthController) UpdateSecuritySettings(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.SecuritySettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	if err := auth.UpdateSecuritySettings(userID, &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithMessage(c, "安全设置已更新", nil)
}
