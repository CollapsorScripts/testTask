package notifications

import (
	"auth/pkg/logger"
	"crypto/tls"
	"errors"
	"gopkg.in/mail.v2"
	"strings"
)

func SendEmail(message string, subject, to string) error {
	if len(message) == 0 {
		return errors.New("сообщение не может быть пустым")
	}

	if !strings.Contains(to, "@") {
		return errors.New("неверная почта")
	}

	from := "securety@notification.ru"
	m := mail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", message)

	d := mail.NewDialer("mail.notification.ru", 465, from, "")

	d.SSL = true
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	logger.Info("Успешно отправлено уведомление \"%s\" на почту: %s", message, to)
	return nil
}
