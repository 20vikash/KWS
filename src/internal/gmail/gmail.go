package gmail

import (
	"fmt"
	env "kws/kws/internal"
	"log"

	"gopkg.in/gomail.v2"
)

func SendMail(to string, token string) {
	m := gomail.NewMessage()

	m.SetHeader("From", env.GetGmail())
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Hello")

	url := fmt.Sprintf("http://localhost:8080/verify?token=%s", token)

	m.SetBody("text/html", fmt.Sprintf("<html>Click <a href='%s'>here</a> to activate your account. This link will expire in 1 day.</html>", url))

	d := gomail.NewDialer("smtp.gmail.com", 587, env.GetGmail(), env.GetGmailAppPassword())

	if err := d.DialAndSend(m); err != nil {
		log.Fatal("Cannot send mail")
	}
}
