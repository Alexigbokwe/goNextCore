package mail

import (
	"fmt"
	"github.com/Alexigbokwe/gonext-framework/core/config"
	"github.com/Alexigbokwe/gonext-framework/core/logger"
	"net/smtp"
	"strings"

	"go.uber.org/zap"
)

// Mailer defines the interface for sending emails
type Mailer interface {
	Send(to []string, subject string, body string) error
	SendHTML(to []string, subject string, body string) error
}

// SMTPMailer implementation
type SMTPMailer struct {
	cfg *config.Config
}

func NewMailer(cfg *config.Config) Mailer {
	return &SMTPMailer{cfg: cfg}
}

func (m *SMTPMailer) Send(to []string, subject string, body string) error {
	return m.sendMail(to, subject, body, "text/plain")
}

func (m *SMTPMailer) SendHTML(to []string, subject string, body string) error {
	return m.sendMail(to, subject, body, "text/html")
}

func (m *SMTPMailer) sendMail(to []string, subject string, body string, contentType string) error {
	mailCfg := m.cfg.Mail

	// If no config, log warning (dev mode)
	if mailCfg.Host == "" {
		logger.Log.Warn("SMTP Host not configured, skipping email",
			zap.Strings("to", to),
			zap.String("subject", subject),
		)
		return nil
	}

	auth := smtp.PlainAuth("", mailCfg.Username, mailCfg.Password, mailCfg.Host)
	addr := fmt.Sprintf("%s:%d", mailCfg.Host, mailCfg.Port)

	mime := "MIME-version: 1.0;\nContent-Type: " + contentType + "; charset=\"UTF-8\";\n\n"
	msg := []byte("To: " + strings.Join(to, ",") + "\r\n" +
		"Subject: " + subject + "\r\n" +
		mime + "\r\n" +
		body)

	err := smtp.SendMail(addr, auth, mailCfg.From, to, msg)
	if err != nil {
		logger.Log.Error("Failed to send email", zap.Error(err))
		return err
	}

	logger.Log.Info("Email sent successfully", zap.Strings("to", to))
	return nil
}

// MockMailer for testing
type MockMailer struct {
	SentEmails []string
}

func (m *MockMailer) Send(to []string, subject string, body string) error {
	m.SentEmails = append(m.SentEmails, fmt.Sprintf("To: %v, Sub: %s", to, subject))
	return nil
}

func (m *MockMailer) SendHTML(to []string, subject string, body string) error {
	return m.Send(to, subject, body)
}
