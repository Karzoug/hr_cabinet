package smtp

import (
	"fmt"
	"net/mail"
	"net/smtp"
)

type mock struct {
	smtpAddr string
	from     *mail.Address
}

func NewMock(cfg Config) *mock {
	return &mock{
		smtpAddr: fmt.Sprintf("%s:%d", cfg.SMTPHost, cfg.SMTPPort),
		from:     &mail.Address{Name: cfg.Name, Address: cfg.From},
	}
}

func (m *mock) SendMessage(recipient, subject, message string) error {
	to := mail.Address{Name: "", Address: recipient}

	msg := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s", m.from.String(), to.String(), subject, message)

	err := smtp.SendMail(m.smtpAddr, nil, m.from.String(), []string{recipient}, []byte(msg))
	return err
}
