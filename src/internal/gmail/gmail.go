package gmail

import (
	"errors"
	"fmt"
	"kws/kws/consts/config"
	env "kws/kws/internal"

	"gopkg.in/gomail.v2"
)

func SendMail(to string, token string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", env.GetGmail())
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Hello")

	url := fmt.Sprintf("https://%s/verify?token=%s", config.DOMAIN(), token)

	m.SetBody("text/html", fmt.Sprintf("<html>Click <a href='%s'>here</a> to activate your account. This link will expire in 1 day.</html>", url))

	d := gomail.NewDialer(env.GetSMTPHost(), env.GetSMTPPort(), env.GetGmail(), env.GetGmailAppPassword())

	if err := d.DialAndSend(m); err != nil {
		return errors.New("cannot send Email: " + err.Error())
	}

	return nil
}
