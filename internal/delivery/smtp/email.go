package smtp

import (
	"crypto/tls"
	"fmt"
	"net/mail"
	"net/smtp"
)

type email struct {
	smtpHost  string
	smtpPort  int
	from      *mail.Address
	auth      smtp.Auth
	tlsconfig *tls.Config
}

func New(cfg Config) *email {
	var m email
	m.smtpHost = cfg.SMTPHost
	m.smtpPort = cfg.SMTPPort
	m.from = &mail.Address{Name: cfg.Name, Address: cfg.From}
	m.auth = smtp.PlainAuth("", cfg.Login, cfg.Password, m.smtpHost)
	m.tlsconfig = &tls.Config{
		ServerName: m.smtpHost,
	}
	return &m
}

func (m *email) SendMessage(recipient, subject, msg string) error {
	const op = "email: send message"

	to := mail.Address{Name: "", Address: recipient}

	message := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s", m.from.String(), to.String(), subject, msg)

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", m.smtpHost, m.smtpPort), m.tlsconfig)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	c, err := smtp.NewClient(conn, m.smtpHost)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err = c.Auth(m.auth); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err = c.Mail(m.from.Address); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err = c.Rcpt(to.Address); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if _, err = w.Write([]byte(message)); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err = w.Close(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err = c.Quit(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
