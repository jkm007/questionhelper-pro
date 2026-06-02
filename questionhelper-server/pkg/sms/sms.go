package sms

import (
	"fmt"

	"questionhelper-server/pkg/logger"
)

// SMSConfig 短信服务配置
type SMSConfig struct {
	Provider     string `yaml:"provider"`      // 短信服务商: aliyun / tencent / mock
	AccessKey    string `yaml:"access_key"`    // AccessKey ID
	AccessSecret string `yaml:"access_secret"` // AccessKey Secret
	SignName     string `yaml:"sign_name"`     // 短信签名
	TemplateCode string `yaml:"template_code"` // 短信模板 Code
}

var cfg *SMSConfig

// Init 初始化短信服务配置
func Init(config *SMSConfig) {
	cfg = config
}

// SendSMS 发送短信验证码
// 当前为占位实现，仅记录日志。实际服务商对接（阿里云/腾讯云）可后续补充。
func SendSMS(phone, code string) error {
	if cfg == nil {
		return fmt.Errorf("短信服务未初始化")
	}

	// 占位实现：记录日志，模拟发送成功
	logger.Infof("[SMS] 发送短信验证码: phone=%s, code=%s, provider=%s, sign=%s, template=%s",
		phone, code, cfg.Provider, cfg.SignName, cfg.TemplateCode)

	// TODO: 根据 cfg.Provider 调用对应的短信服务商 API
	// switch cfg.Provider {
	// case "aliyun":
	//     return sendViaAliyun(phone, code)
	// case "tencent":
	//     return sendViaTencent(phone, code)
	// }

	return nil
}
