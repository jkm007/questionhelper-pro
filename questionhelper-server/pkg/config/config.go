package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Cfg 全局配置变量
var Cfg *Config

type Config struct {
	Server  ServerConfig  `yaml:"server"`
	MySQL   MySQLConfig   `yaml:"mysql"`
	Redis   RedisConfig   `yaml:"redis"`
	JWT     JWTConfig     `yaml:"jwt"`
	Auth    AuthConfig    `yaml:"auth"`
	Captcha CaptchaConfig `yaml:"captcha"`
	Limit   LimitConfig   `yaml:"limit"`
	Log     LogConfig     `yaml:"log"`
	OSS     OSSConfig     `yaml:"oss"`
	Email   EmailConfig   `yaml:"email"`
	SMS     SMSConfig     `yaml:"sms"`
	OAuth   OAuthConfig   `yaml:"oauth"`
}

// OAuthConfig 第三方登录配置
type OAuthConfig struct {
	GitHub OAuthProvider `yaml:"github"`
	Google OAuthProvider `yaml:"google"`
	WeChat OAuthProvider `yaml:"wechat"`
}

// OAuthProvider 第三方登录提供商配置
type OAuthProvider struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	RedirectURI  string `yaml:"redirect_uri"`
}

type ServerConfig struct {
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"`
}

type MySQLConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type JWTConfig struct {
	Secret           string `yaml:"secret"`
	Expire           int    `yaml:"expire"`             // Access Token 有效期（秒）
	RefreshExpire    int    `yaml:"refresh_expire"`     // Refresh Token 有效期（秒）
	RememberMeExpire int    `yaml:"remember_me_expire"` // 记住我 Refresh Token 有效期（秒），默认30天
}

type AuthConfig struct {
	BcryptCost       int    `yaml:"bcrypt_cost"`        // bcrypt 加密 cost，默认 12
	JWTIssuer        string `yaml:"jwt_issuer"`         // JWT 签发者
	ResetTokenExpire int    `yaml:"reset_token_expire"` // 密码重置 Token 有效期（秒），默认 900
}

type CaptchaConfig struct {
	Type              string `yaml:"type"`                // captcha 类型: digit / letter / math
	Length            int    `yaml:"length"`              // 验证码长度
	Width             int    `yaml:"width"`               // 图片宽度
	Height            int    `yaml:"height"`              // 图片高度
	Expire            int    `yaml:"expire"`              // 有效期（秒）
	MaxVerifyAttempts int    `yaml:"max_verify_attempts"` // 最大验证错误次数
}

type LimitConfig struct {
	LoginIPPerHour       int `yaml:"login_ip_per_hour"`        // 同一 IP 每小时最大登录次数
	SMSPerMinute         int `yaml:"sms_per_minute"`           // 短信发送频率（次/分钟/手机）
	SMSPerDay            int `yaml:"sms_per_day"`              // 每日短信发送上限
	SMSIPPerHour         int `yaml:"sms_ip_per_hour"`          // 同一 IP 每小时最大发送次数
	EmailPerMinute       int `yaml:"email_per_minute"`         // 邮箱发送频率（次/分钟/邮箱）
	EmailPerDay          int `yaml:"email_per_day"`            // 每日邮箱发送上限
	RegisterPerIPPerHour int `yaml:"register_per_ip_per_hour"` // 同一 IP 每小时最大注册次数
}

type LogConfig struct {
	Level      string `yaml:"level"`
	Format     string `yaml:"format"`
	Output     string `yaml:"output"`      // 日志文件路径，留空输出到 stdout
	MaxSize    int    `yaml:"max_size"`     // 单个日志文件最大尺寸 (MB)
	MaxBackups int    `yaml:"max_backups"`  // 保留旧日志文件最大数量
	MaxAge     int    `yaml:"max_age"`      // 保留旧日志文件最大天数
}

type OSSConfig struct {
	Type      string `yaml:"type"`
	Endpoint  string `yaml:"endpoint"`
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
	Bucket    string `yaml:"bucket"`
	CDN       string `yaml:"cdn"`
}

type EmailConfig struct {
	SMTPHost    string `yaml:"smtp_host"`
	SMTPPort    int    `yaml:"smtp_port"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	FromAddress string `yaml:"from_address"`
	FromName    string `yaml:"from_name"`
	UseTLS      bool   `yaml:"use_tls"`
}

type SMSConfig struct {
	Provider     string `yaml:"provider"`      // 短信服务商: aliyun / tencent / mock
	AccessKey    string `yaml:"access_key"`    // AccessKey ID
	AccessSecret string `yaml:"access_secret"` // AccessKey Secret
	SignName     string `yaml:"sign_name"`     // 短信签名
	TemplateCode string `yaml:"template_code"` // 短信模板 Code
}

func Load() (*Config, error) {
	configPath := "config/config.yaml"
	if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
		configPath = envPath
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 设置默认值
	setDefaults(cfg)

	// 环境变量覆盖（容器化部署场景）
	applyEnvOverrides(cfg)

	// 设置全局配置
	Cfg = cfg

	return cfg, nil
}

func setDefaults(cfg *Config) {
	// Auth 默认值
	if cfg.Auth.BcryptCost == 0 {
		cfg.Auth.BcryptCost = 12
	}
	if cfg.Auth.JWTIssuer == "" {
		cfg.Auth.JWTIssuer = "questionhelper"
	}
	if cfg.Auth.ResetTokenExpire == 0 {
		cfg.Auth.ResetTokenExpire = 900 // 15 分钟
	}

	// JWT 默认值
	if cfg.JWT.Expire == 0 {
		cfg.JWT.Expire = 7200 // 2 小时
	}
	if cfg.JWT.RefreshExpire == 0 {
		cfg.JWT.RefreshExpire = 604800 // 7 天
	}
	if cfg.JWT.RememberMeExpire == 0 {
		cfg.JWT.RememberMeExpire = 2592000 // 30 天
	}

	// Captcha 默认值
	if cfg.Captcha.Type == "" {
		cfg.Captcha.Type = "math"
	}
	if cfg.Captcha.Length == 0 {
		cfg.Captcha.Length = 4
	}
	if cfg.Captcha.Width == 0 {
		cfg.Captcha.Width = 120
	}
	if cfg.Captcha.Height == 0 {
		cfg.Captcha.Height = 40
	}
	if cfg.Captcha.Expire == 0 {
		cfg.Captcha.Expire = 300 // 5 分钟
	}
	if cfg.Captcha.MaxVerifyAttempts == 0 {
		cfg.Captcha.MaxVerifyAttempts = 5
	}

	// Limit 默认值
	if cfg.Limit.LoginIPPerHour == 0 {
		cfg.Limit.LoginIPPerHour = 50
	}
	if cfg.Limit.SMSPerMinute == 0 {
		cfg.Limit.SMSPerMinute = 1
	}
	if cfg.Limit.SMSPerDay == 0 {
		cfg.Limit.SMSPerDay = 10
	}
	if cfg.Limit.SMSIPPerHour == 0 {
		cfg.Limit.SMSIPPerHour = 50
	}
	if cfg.Limit.EmailPerMinute == 0 {
		cfg.Limit.EmailPerMinute = 1
	}
	if cfg.Limit.EmailPerDay == 0 {
		cfg.Limit.EmailPerDay = 10
	}
	if cfg.Limit.RegisterPerIPPerHour == 0 {
		cfg.Limit.RegisterPerIPPerHour = 10
	}

	// Email 默认值
	if cfg.Email.SMTPHost == "" {
		cfg.Email.SMTPHost = "smtp.qq.com"
	}
	if cfg.Email.SMTPPort == 0 {
		cfg.Email.SMTPPort = 465
	}
	if cfg.Email.FromName == "" {
		cfg.Email.FromName = "题小助"
	}

	// SMS 默认值
	if cfg.SMS.Provider == "" {
		cfg.SMS.Provider = "mock"
	}

	// Log 默认值
	if cfg.Log.MaxSize == 0 {
		cfg.Log.MaxSize = 100 // 100MB
	}
	if cfg.Log.MaxBackups == 0 {
		cfg.Log.MaxBackups = 5
	}
	if cfg.Log.MaxAge == 0 {
		cfg.Log.MaxAge = 30 // 30天
	}
}

// applyEnvOverrides 环境变量覆盖关键配置（容器化部署场景）
func applyEnvOverrides(cfg *Config) {
	if v := os.Getenv("MYSQL_HOST"); v != "" {
		cfg.MySQL.Host = v
	}
	if v := os.Getenv("MYSQL_PORT"); v != "" {
		fmt.Sscanf(v, "%d", &cfg.MySQL.Port)
	}
	if v := os.Getenv("MYSQL_USER"); v != "" {
		cfg.MySQL.User = v
	}
	if v := os.Getenv("MYSQL_PASSWORD"); v != "" {
		cfg.MySQL.Password = v
	}
	if v := os.Getenv("MYSQL_DBNAME"); v != "" {
		cfg.MySQL.DBName = v
	}
	if v := os.Getenv("REDIS_HOST"); v != "" {
		cfg.Redis.Host = v
	}
	if v := os.Getenv("REDIS_PORT"); v != "" {
		fmt.Sscanf(v, "%d", &cfg.Redis.Port)
	}
	if v := os.Getenv("REDIS_PASSWORD"); v != "" {
		cfg.Redis.Password = v
	}
	if v := os.Getenv("JWT_SECRET"); v != "" {
		cfg.JWT.Secret = v
	}
	if v := os.Getenv("SERVER_PORT"); v != "" {
		fmt.Sscanf(v, "%d", &cfg.Server.Port)
	}
	if v := os.Getenv("SERVER_MODE"); v != "" {
		cfg.Server.Mode = v
	}
	if v := os.Getenv("LOG_LEVEL"); v != "" {
		cfg.Log.Level = v
	}
	if v := os.Getenv("OSS_TYPE"); v != "" {
		cfg.OSS.Type = v
	}
}
