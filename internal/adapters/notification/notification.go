package notification

import (
	"fmt"

	gomail "gopkg.in/mail.v2"
)

type Config struct {
	EmailFrom     string `env:"NOTIFY_EMAIL_FROM"`
	EmailSMTPHost string `env:"NOTIFY_EMAIL_SMTP_HOST"`
	EmailSMTPPort int    `env:"NOTIFY_EMAIL_SMTP_PORT"`
	EmailPassword string `env:"NOTIFY_EMAIL_PASS"`
}

type Notify struct {
	cfg Config
}

func New(cfg Config) (*Notify, error) {
	n := &Notify{
		cfg: cfg,
	}

	return n, nil
}

func (n *Notify) SendEmail(to string, subject string, body []byte) error {
	// Create a new message
	message := gomail.NewMessage()

	// Set email headers
	message.SetHeader("From", n.cfg.EmailFrom)
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)

	// Set email body
	message.SetBody("text/html", string(body))

	// Set up the SMTP dialer
	dialer := gomail.NewDialer(n.cfg.EmailSMTPHost, n.cfg.EmailSMTPPort, n.cfg.EmailFrom, n.cfg.EmailPassword)

	// Send the email
	if err := dialer.DialAndSend(message); err != nil {
		return fmt.Errorf("failed send email: %w", err)
	}

	return nil
}
