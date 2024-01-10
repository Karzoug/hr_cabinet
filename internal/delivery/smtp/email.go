package smtp

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/mail"
	"net/smtp"
	"sync"
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
	m.auth = smtp.PlainAuth("", cfg.Login, cfg.AppPass, m.smtpHost)
	m.tlsconfig = &tls.Config{
		ServerName: m.smtpHost,
	}
	return &m
}

func (m *email) SendMessage(recipient, subject, msg string) error {
	to := mail.Address{Name: "", Address: recipient}

	message := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s", m.from.String(), to.String(), subject, msg)

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", m.smtpHost, m.smtpPort), m.tlsconfig)
	if err != nil {
		log.Printf("conn error sending mail %v", err)
		return err
	}

	c, err := smtp.NewClient(conn, m.smtpHost)
	if err != nil {
		log.Printf("NewClient error sending mail %v", err)
		return err
	}

	if err = c.Auth(m.auth); err != nil {
		log.Printf("auth error sending mail %v", err)
		return err
	}

	if err = c.Mail(m.from.Address); err != nil {
		log.Printf("mail error sending mail %v", err)
		return err
	}

	if err = c.Rcpt(to.Address); err != nil {
		log.Printf("rcpt error sending mail %v", err)
		return err
	}

	w, err := c.Data()
	if err != nil {
		log.Printf("data error sending mail %v", err)
		return err
	}

	if _, err = w.Write([]byte(message)); err != nil {
		log.Printf("write error sending mail %v", err)
		return err
	}

	if err = w.Close(); err != nil {
		log.Printf("close error sending mail %v", err)
		return err
	}

	if err = c.Quit(); err != nil {
		log.Printf("quit error sending mail %v", err)
		return err
	}
	return nil
}

// SendMessageToMany - concurrently sending mails to multiple recipients.
func (m *email) SendMessageToMany(subject, msg string, recipients []string) {
	var mailWG sync.WaitGroup
	mailWG.Add(len(recipients))
	for _, v := range recipients {
		go func(recipient string) {
			defer mailWG.Done()
			m.SendMessage(subject, msg, recipient)
		}(v)
	}
	mailWG.Wait()
}
