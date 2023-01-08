package utils

import (
	"crypto/tls"
	"net/smtp"

	"github.com/jordan-wright/email"
)

/**
 * @Author jiang
 * @Description 发送邮件验证码
 * @Date 12:00 2023/1/8
 **/
func SendEmailCode(toUserEmail, code string) error {
	e := email.NewEmail()

	mailUserName := "jiang2381385276@163.com" //邮箱账号
	mailPassword := "OSJXVUTKLANNJZIP"        //邮箱授权码
	Subject := "验证码"                          //发送的主题

	e.From = "心灵树洞 <jiang2381385276@163.com>"
	e.To = []string{toUserEmail}
	e.Subject = Subject
	e.HTML = []byte("你好，感谢您使用心灵树洞，您的账号正在使用邮箱验证，本次请求的验证码为：" + code + ", 5分钟内有效")
	err := e.SendWithTLS("smtp.163.com:465", smtp.PlainAuth("", mailUserName, mailPassword, "smtp.163.com"),
		&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.163.com"})
	return err
}
