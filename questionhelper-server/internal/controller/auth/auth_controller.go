package auth

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/service/auth"
	"questionhelper-server/pkg/config"
	"questionhelper-server/pkg/response"
)

type AuthController struct {
	jwtCfg *config.JWTConfig
}

func NewAuthController(jwtCfg *config.JWTConfig) *AuthController {
	return &AuthController{jwtCfg: jwtCfg}
}

// Login 用户登录
func (ac *AuthController) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 获取客户端信息
	req.DeviceInfo = c.ClientIP()

	result, err := auth.Login(&req, ac.jwtCfg)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	response.Success(c, result)
}

// Register 用户注册
func (ac *AuthController) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	if err := auth.Register(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessWithMessage(c, "注册成功", nil)
}

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

	response.SuccessWithMessage(c, "退出成功", nil)
}

// LogoutAll 退出所有设备
func (ac *AuthController) LogoutAll(c *gin.Context) {
	userID := c.GetUint("user_id")

	if err := auth.LogoutAll(userID); err != nil {
		response.Error(c, http.StatusInternalServerError, "退出失败")
		return
	}

	response.SuccessWithMessage(c, "已退出所有设备", nil)
}

// RequestPasswordReset 请求重置密码
func (ac *AuthController) RequestPasswordReset(c *gin.Context) {
	var req struct {
		Account string `json:"account" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	if err := auth.RequestPasswordReset(req.Account); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithMessage(c, "重置密码请求已提交，请查看手机或邮箱", nil)
}

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

// DeactivateAccount 申请注销账号
func (ac *AuthController) DeactivateAccount(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.DeactivateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	if err := auth.DeactivateAccount(userID, req.Password); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, gin.H{
		"verifyType": "sms",
		"expireIn":   300,
	})
}

// ConfirmDeactivate 确认注销
func (ac *AuthController) ConfirmDeactivate(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.DeactivateConfirmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	if err := auth.ConfirmDeactivate(userID, req.Code); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessWithMessage(c, "注销申请已提交，30天后将永久删除账号", nil)
}

// CancelDeactivate 取消注销
func (ac *AuthController) CancelDeactivate(c *gin.Context) {
	userID := c.GetUint("user_id")

	if err := auth.CancelDeactivate(userID); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessWithMessage(c, "已取消注销", nil)
}

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

// KickAllDevices 踢出所有设备
func (ac *AuthController) KickAllDevices(c *gin.Context) {
	userID := c.GetUint("user_id")

	if err := auth.LogoutAll(userID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithMessage(c, "已退出所有设备", nil)
}

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

// ChangePassword 修改密码
func (ac *AuthController) ChangePassword(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	if err := auth.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessWithMessage(c, "密码修改成功", nil)
}

// UpdateSecuritySettings 更新安全设置（占位）
func (ac *AuthController) UpdateSecuritySettings(c *gin.Context) {
	// TODO: 实现安全设置更新
	response.SuccessWithMessage(c, "安全设置已更新", nil)
}
