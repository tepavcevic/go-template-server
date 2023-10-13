package models

import "github.com/go-mail/mail/v2"

const (
	DefaultSender = "support@myservice.com"
)

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

func NewEmailService(config SMTPConfig) *EmailService {
	es := EmailService{
		dialer: mail.NewDialer(
			config.Host,
			config.Port,
			config.Username,
			config.Password,
		),
	}
	return &es
}

type EmailService struct {
	// DefaultSender is used as default sender when one isn't provided
	DefaultSender string

	// unexported fields
	dialer *mail.Dialer
}
