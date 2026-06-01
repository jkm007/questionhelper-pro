package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"questionhelper-server/internal/service/auth"
	"questionhelper-server/pkg/response"
)

// GetCaptcha 获取验证码
func (ac *AuthController) GetCaptcha(c *gin.Context) {
	id, b64s, err := auth.GetCaptcha()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "生成验证码失败")
		return
	}

	response.Success(c, gin.H{
		"captchaId":    id,
		"captchaBase64": b64s,
	})
}

// VerifyCaptcha 验证验证码
func (ac *AuthController) VerifyCaptcha(c *gin.Context) {
	var req struct {
		CaptchaID  string `json:"captchaId" binding:"required"`
		CaptchaCode string `json:"captchaCode" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	if err := auth.VerifyCaptcha(req.CaptchaID, req.CaptchaCode); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessWithMessage(c, "验证码正确", nil)
}

// SendSmsCode 发送短信验证码（占位）
func (ac *AuthController) SendSmsCode(c *gin.Context) {
	var req struct {
		Phone     string `json:"phone" binding:"required"`
		CaptchaID string `json:"captchaId"`
		CaptchaCode string `json:"captchaCode"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// TODO: 实现短信发送
	response.SuccessWithMessage(c, "短信验证码已发送", nil)
}

// SendEmailCode 发送邮箱验证码（占位）
func (ac *AuthController) SendEmailCode(c *gin.Context) {
	var req struct {
		Email     string `json:"email" binding:"required,email"`
		CaptchaID string `json:"captchaId"`
		CaptchaCode string `json:"captchaCode"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// TODO: 实现邮箱发送
	response.SuccessWithMessage(c, "邮箱验证码已发送", nil)
}
