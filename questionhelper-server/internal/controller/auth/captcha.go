package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"questionhelper-server/internal/service/auth"
	"questionhelper-server/pkg/captcha"
	apperr "questionhelper-server/pkg/errors"
	"questionhelper-server/pkg/response"
)

// GetCaptcha 获取验证码（T29: 支持 digit/letter/math 三种类型，通过 query 参数 type 指定）
func (ac *AuthController) GetCaptcha(c *gin.Context) {
	captchaType := captcha.CaptchaType(c.DefaultQuery("type", string(captcha.CaptchaTypeDigit)))
	// 校验类型合法性
	switch captchaType {
	case captcha.CaptchaTypeDigit, captcha.CaptchaTypeLetter, captcha.CaptchaTypeMath:
	default:
		captchaType = captcha.CaptchaTypeDigit
	}

	id, b64s, err := auth.GetCaptcha(captchaType)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "生成验证码失败")
		return
	}

	response.Success(c, gin.H{
		"captchaId":     id,
		"captchaBase64": b64s,
		"captchaType":   string(captchaType),
	})
}

// VerifyCaptcha 验证验证码
func (ac *AuthController) VerifyCaptcha(c *gin.Context) {
	var req struct {
		CaptchaID   string `json:"captchaId" binding:"required"`
		CaptchaCode string `json:"captchaCode" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	if err := auth.VerifyCaptcha(req.CaptchaID, req.CaptchaCode); err != nil {
		response.ErrorWithCode(c, http.StatusBadRequest, apperr.ErrCaptchaWrong.Code, err.Error())
		return
	}

	response.SuccessWithMessage(c, "验证码正确", nil)
}

// SendSmsCode 发送短信验证码
func (ac *AuthController) SendSmsCode(c *gin.Context) {
	var req struct {
		Phone       string `json:"phone" binding:"required"`
		CaptchaID   string `json:"captchaId"`
		CaptchaCode string `json:"captchaCode"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 验证图形验证码（如果提供了的话）
	if req.CaptchaID != "" && req.CaptchaCode != "" {
		if err := auth.VerifyCaptcha(req.CaptchaID, req.CaptchaCode); err != nil {
			response.ErrorWithCode(c, http.StatusBadRequest, apperr.ErrCaptchaWrong.Code, "验证码错误")
			return
		}
	}

	if err := auth.SendSmsCode(req.Phone); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithMessage(c, "短信验证码已发送", nil)
}

// SendEmailCode 发送邮箱验证码
func (ac *AuthController) SendEmailCode(c *gin.Context) {
	var req struct {
		Email       string `json:"email" binding:"required,email"`
		CaptchaID   string `json:"captchaId"`
		CaptchaCode string `json:"captchaCode"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 验证图形验证码（如果提供了的话）
	if req.CaptchaID != "" && req.CaptchaCode != "" {
		if err := auth.VerifyCaptcha(req.CaptchaID, req.CaptchaCode); err != nil {
			response.ErrorWithCode(c, http.StatusBadRequest, apperr.ErrCaptchaWrong.Code, "验证码错误")
			return
		}
	}

	if err := auth.SendEmailCode(req.Email); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithMessage(c, "邮箱验证码已发送", nil)
}
