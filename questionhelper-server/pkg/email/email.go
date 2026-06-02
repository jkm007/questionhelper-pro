package email

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"
)

// EmailConfig 邮件服务配置
type EmailConfig struct {
	SMTPHost    string `yaml:"smtp_host"`
	SMTPPort    int    `yaml:"smtp_port"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	FromAddress string `yaml:"from_address"`
	FromName    string `yaml:"from_name"`
	UseTLS      bool   `yaml:"use_tls"`
}

var cfg *EmailConfig

// Init 初始化邮件服务配置
func Init(config *EmailConfig) {
	cfg = config
}

// SendEmail 发送邮件（支持 TLS 和 STARTTLS）
func SendEmail(to, subject, body string) error {
	if cfg == nil {
		return fmt.Errorf("邮件服务未初始化")
	}

	from := mail.Address{Name: cfg.FromName, Address: cfg.FromAddress}

	// 构造邮件头
	header := make(map[string]string)
	header["From"] = from.String()
	header["To"] = to
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=UTF-8"

	// 组装邮件内容
	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	addr := fmt.Sprintf("%s:%d", cfg.SMTPHost, cfg.SMTPPort)
	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.SMTPHost)

	if cfg.UseTLS {
		return sendWithTLS(addr, auth, from.Address, to, []byte(message))
	}
	return sendWithSTARTTLS(addr, auth, from.Address, to, []byte(message))
}

// sendWithTLS 通过隐式 TLS（端口 465）发送邮件
func sendWithTLS(addr string, auth smtp.Auth, from, to string, msg []byte) error {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return fmt.Errorf("解析地址失败: %w", err)
	}

	tlsConfig := &tls.Config{
		ServerName: host,
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("TLS 连接失败: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return fmt.Errorf("创建 SMTP 客户端失败: %w", err)
	}
	defer client.Close()

	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP 认证失败: %w", err)
	}
	if err = client.Mail(from); err != nil {
		return fmt.Errorf("设置发件人失败: %w", err)
	}
	if err = client.Rcpt(to); err != nil {
		return fmt.Errorf("设置收件人失败: %w", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("打开数据写入失败: %w", err)
	}
	if _, err = w.Write(msg); err != nil {
		return fmt.Errorf("写入邮件内容失败: %w", err)
	}
	if err = w.Close(); err != nil {
		return fmt.Errorf("关闭数据写入失败: %w", err)
	}

	return client.Quit()
}

// sendWithSTARTTLS 通过 STARTTLS（端口 587）发送邮件
func sendWithSTARTTLS(addr string, auth smtp.Auth, from, to string, msg []byte) error {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return fmt.Errorf("解析地址失败: %w", err)
	}

	tlsConfig := &tls.Config{
		ServerName: host,
	}

	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("连接 SMTP 服务器失败: %w", err)
	}
	defer client.Close()

	if ok, _ := client.Extension("STARTTLS"); ok {
		if err = client.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("STARTTLS 升级失败: %w", err)
		}
	}

	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP 认证失败: %w", err)
	}
	if err = client.Mail(from); err != nil {
		return fmt.Errorf("设置发件人失败: %w", err)
	}
	if err = client.Rcpt(to); err != nil {
		return fmt.Errorf("设置收件人失败: %w", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("打开数据写入失败: %w", err)
	}
	if _, err = w.Write(msg); err != nil {
		return fmt.Errorf("写入邮件内容失败: %w", err)
	}
	if err = w.Close(); err != nil {
		return fmt.Errorf("关闭数据写入失败: %w", err)
	}

	return client.Quit()
}

// SendVerificationCode 发送邮箱验证码
func SendVerificationCode(to, code string) error {
	subject := "题小助 - 邮箱验证码"
	body := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"></head>
<body style="margin:0;padding:0;background-color:#f4f4f7;font-family:'Helvetica Neue',Arial,sans-serif;">
  <table width="100%%" cellpadding="0" cellspacing="0" style="background-color:#f4f4f7;padding:40px 0;">
    <tr><td align="center">
      <table width="520" cellpadding="0" cellspacing="0" style="background-color:#ffffff;border-radius:8px;box-shadow:0 2px 8px rgba(0,0,0,0.08);">
        <tr><td style="padding:32px 40px 0;text-align:center;">
          <h2 style="margin:0;color:#333333;font-size:22px;">题小助</h2>
        </td></tr>
        <tr><td style="padding:24px 40px;">
          <p style="color:#555555;font-size:15px;line-height:1.6;">您好，您正在进行邮箱验证。请使用以下验证码完成操作：</p>
          <div style="margin:24px 0;text-align:center;">
            <span style="display:inline-block;background-color:#f0f5ff;border:1px solid #d6e4ff;border-radius:6px;padding:14px 32px;font-size:28px;font-weight:bold;letter-spacing:6px;color:#1677ff;">%s</span>
          </div>
          <p style="color:#999999;font-size:13px;line-height:1.6;">验证码 5 分钟内有效，请勿泄露给他人。如非本人操作，请忽略此邮件。</p>
        </td></tr>
        <tr><td style="padding:0 40px 32px;">
          <hr style="border:none;border-top:1px solid #eeeeee;margin:0 0 16px;">
          <p style="color:#bbbbbb;font-size:12px;text-align:center;margin:0;">此邮件由系统自动发送，请勿直接回复</p>
        </td></tr>
      </table>
    </td></tr>
  </table>
</body>
</html>`, code)

	return SendEmail(to, subject, body)
}
