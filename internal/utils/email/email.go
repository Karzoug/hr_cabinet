package email

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/mail"
	"net/smtp"
	"sync"
)

var (
	smtpHost  string
	from      *mail.Address
	auth      smtp.Auth
	tlsconfig *tls.Config
)

type Config struct {
	Name     string `env:"NAME" env-default:"Картотека сотрудника"`
	FromMail string `env:"FROM_MAIL"`
	Login    string `env:"LOGIN"`
	AppPass  string `env:"APP_PASS"`
	SMTPHost string `env:"SMTP_HOST"`
}

func EmailInit(cfg *Config) {
	smtpHost = cfg.SMTPHost
	from = &mail.Address{Name: cfg.Name, Address: cfg.FromMail}
	auth = smtp.PlainAuth("", cfg.Login, cfg.AppPass, smtpHost)
	tlsconfig = &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpHost,
	}
}

func SendSSLMail(subject, msg, recipient string) error {
	to := mail.Address{Name: "", Address: recipient}
	// Setup message
	message := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s", from.String(), to.String(), subject, msg)

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", smtpHost, 465), tlsconfig)
	if err != nil {
		log.Printf("conn error sending mail %v", err)
		return err
	}

	c, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		log.Printf("NewClient error sending mail %v", err)
		return err
	}
	// Auth
	if err = c.Auth(auth); err != nil {
		log.Printf("auth error sending mail %v", err)
		return err
	}
	// To && From
	if err = c.Mail(from.Address); err != nil {
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
func SendMails(subject, msg string, recipients []string) {
	var mailwg sync.WaitGroup
	mailwg.Add(len(recipients))
	for _, v := range recipients {
		go func(recipient string) {
			defer mailwg.Done()
			SendSSLMail(subject, msg, recipient)
		}(v)
	}
	mailwg.Wait()
}
