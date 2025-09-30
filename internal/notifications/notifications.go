package notifications

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strconv"
	"text/template"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

type EmailNotifier struct {
	from      string
	dialer    *gomail.Dialer
	templates *template.Template
}

func NewEmailNotifier() *EmailNotifier {
	godotenv.Load()

	host := os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	from := os.Getenv("SMTP_FROM")
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")
	tmpl := os.Getenv("TMPL")

	port, _ := strconv.Atoi(portStr)

	notifier := NewEmailNotifierWithParams(host, port, from, user, pass)

	notifier.templates = template.Must(template.ParseFiles(tmpl))

	return notifier
}

func NewEmailNotifierWithParams(host string, port int, from, smtpUser, smtpPass string) *EmailNotifier {
	d := gomail.NewDialer(host, port, smtpUser, smtpPass)
	return &EmailNotifier{
		from:   from,
		dialer: d,
	}
}

func (n *EmailNotifier) SendVerificationCode(ctx context.Context, to, code string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", n.from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", "Код подтверждения регистрации")

	var body bytes.Buffer
	data := struct{ Code string }{Code: code}
	if err := n.templates.ExecuteTemplate(&body, "email_verification.html", data); err != nil {
		return fmt.Errorf("ошибка генерации шаблона: %v", err)
	}

	msg.SetBody("text/html", body.String())

	return n.dialer.DialAndSend(msg)
}

func (n *EmailNotifier) SendLoginNotification(ctx context.Context, to string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", n.from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", "Вход в аккаунт")
	msg.SetBody("text/plain", "Вы только что вошли в свой аккаунт. Если это были не вы — срочно смените пароль!")

	return n.dialer.DialAndSend(msg)
}
