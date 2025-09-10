package notifications

import (
	"context"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

type EmailNotifier struct {
	from   string
	dialer *gomail.Dialer
}

func NewEmailNotifier() *EmailNotifier {
	godotenv.Load()

	host := os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	from := os.Getenv("SMTP_FROM")
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")

	port, _ := strconv.Atoi(portStr)

	return NewEmailNotifierWithParams(host, port, from, user, pass)
}

func NewEmailNotifierWithParams(host string, port int, from, smtpUser, smtpPass string) *EmailNotifier {
	d := gomail.NewDialer(host, port, smtpUser, smtpPass)
	return &EmailNotifier{from: from, dialer: d}
}

func (n *EmailNotifier) SendRegistrationSuccess(ctx context.Context, to string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", n.from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", "Добро пожаловать!")
	msg.SetBody("text/plain", "Ваша почта успешно зарегистрирована.")

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
