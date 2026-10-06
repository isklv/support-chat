package mailer

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"strings"

	"support-chat/internal/config"
)

type Mailer interface {
	Send(to, subject, bodyHTML string) error
}

type smtpMailer struct {
	cfg *config.Config
}

type noopMailer struct{}

func (n *noopMailer) Send(to, subject, bodyHTML string) error {
	return nil
}

func New(cfg *config.Config) Mailer {
	if cfg.SMTPHost == "" || cfg.SMTPFrom == "" {
		return &noopMailer{}
	}
	return &smtpMailer{cfg: cfg}
}

func (m *smtpMailer) Send(to, subject, bodyHTML string) error {
	addr := fmt.Sprintf("%s:%d", m.cfg.SMTPHost, m.cfg.SMTPPort)
	from := m.cfg.SMTPFrom

	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	var msg strings.Builder
	for k, v := range headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msg.WriteString("\r\n")
	msg.WriteString(bodyHTML)

	var auth smtp.Auth
	if m.cfg.SMTPUser != "" && m.cfg.SMTPPass != "" {
		auth = smtp.PlainAuth("", m.cfg.SMTPUser, m.cfg.SMTPPass, m.cfg.SMTPHost)
	}

	// For port 465 (SMTPS) or TLS
	if m.cfg.SMTPPort == 465 {
		tlsConfig := &tls.Config{
			ServerName: m.cfg.SMTPHost,
		}
		conn, err := tls.Dial("tcp", addr, tlsConfig)
		if err != nil {
			return fmt.Errorf("tls dial smtp: %w", err)
		}
		defer conn.Close()

		client, err := smtp.NewClient(conn, m.cfg.SMTPHost)
		if err != nil {
			return fmt.Errorf("smtp client: %w", err)
		}
		defer client.Close()

		if auth != nil {
			if err := client.Auth(auth); err != nil {
				return fmt.Errorf("smtp auth: %w", err)
			}
		}

		if err := client.Mail(from); err != nil {
			return err
		}
		if err := client.Rcpt(to); err != nil {
			return err
		}
		w, err := client.Data()
		if err != nil {
			return err
		}
		if _, err := w.Write([]byte(msg.String())); err != nil {
			return err
		}
		return w.Close()
	}

	err := smtp.SendMail(addr, auth, from, []string{to}, []byte(msg.String()))
	if err != nil {
		log.Printf("mailer: send to %s failed: %v", to, err)
		return err
	}
	return nil
}
