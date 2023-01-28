package helper

import (
	"cloud-disk/disk/internal/config"
	"crypto/tls"
	"net/smtp"

	"github.com/jordan-wright/email"
)

func SendEmailCode(mail, code string) error {
	e := email.NewEmail()
	e.From = "vczs <vczvip@163.com>"
	e.To = []string{mail}
	e.Subject = "vczs平台验证码"
	e.HTML = []byte("你的验证码为：<h1>" + code + "</h1>")
	err := e.SendWithTLS("smtp.163.com:465", smtp.PlainAuth("", "vczvip@163.com", config.MailPassword, "smtp.163.com"), &tls.Config{InsecureSkipVerify: true, ServerName: "smtp.163.com"})
	if err != nil {
		return err
	}
	return nil
}
