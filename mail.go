package main

import (
	"os"
	"time"
	"net/smtp"
	"strings"
)

func mail(content string) error {
	host := os.Getenv("MAIL_HOST")
	port := os.Getenv("MAIL_PORT")
	addr := host + ":" + port
	to := os.Getenv("MAIL_TO")
	from := os.Getenv("MAIL_FROM")
	pw := os.Getenv("MAIL_PASSWORD")
	auth := smtp.PlainAuth("", from, pw, host)

	subject := "RSS Feeds for " + time.Now().Format("Jan 02, 2006")

	msg := strings.Builder{}
	msg.WriteString("From: \"Feed Update\" <" + from + ">\n")
	msg.WriteString("To: " + to + "\n")
	msg.WriteString("Subject: " + subject + "\n")
	msg.WriteString("MIME-version: 1.0;\n")
	msg.WriteString("Content-Type: text/html;charset=\"UTF-8\";\n")
	msg.WriteString("\n")
	msg.WriteString(content)

	err := smtp.SendMail(addr, auth, from, []string{to}, []byte(msg.String()))
	if err != nil {
		return err
	}
	return nil
}
