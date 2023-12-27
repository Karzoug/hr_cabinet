package email

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/mail"
	"net/smtp"
	"sync"
)

type Mail struct {
	smtpHost  string
	from      *mail.Address
	auth      smtp.Auth
	tlsconfig *tls.Config
}

func New(cfg Config) *Mail {
	var m Mail
	m.smtpHost = cfg.SMTPHost
	m.from = &mail.Address{Name: cfg.Name, Address: cfg.FromMail}
	m.auth = smtp.PlainAuth("", cfg.Login, cfg.AppPass, m.smtpHost)
	m.tlsconfig = &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         m.smtpHost,
	}
	return &m
}

func (m *Mail) SendSSLMail(subject, msg, recipient string) error {
	to := mail.Address{Name: "", Address: recipient}
	// Setup message
	message := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s", m.from.String(), to.String(), subject, msg)

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", m.smtpHost, 465), m.tlsconfig)
	if err != nil {
		log.Printf("conn error sending mail %v", err)
		return err
	}

	c, err := smtp.NewClient(conn, m.smtpHost)
	if err != nil {
		log.Printf("NewClient error sending mail %v", err)
		return err
	}
	// Auth
	if err = c.Auth(m.auth); err != nil {
		log.Printf("auth error sending mail %v", err)
		return err
	}
	// To && From
	if err = c.Mail(m.from.Address); err != nil {
		log.Printf("mail error sending mail %v", err)
		return err
	}

	if err = c.Rcpt(to.Address); err != nil {
		log.Printf("rcpt error sending mail %v", err)
		return err
	}
	// Data
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

// SendMails - concurrently sending mails to multiple recipients.
func (m *Mail) SendMails(subject, msg string, recipients []string) {
	var mailwg sync.WaitGroup
	mailwg.Add(len(recipients))
	for _, v := range recipients {
		go func(recipient string) {
			defer mailwg.Done()
			m.SendSSLMail(subject, msg, recipient)
		}(v)
	}
	mailwg.Wait()
}
