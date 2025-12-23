package email

import (
	"fmt"
	"net/smtp"
	"strings"
	"time"
)

type SMTPClient struct {
	host string
	port string
}

func NewSMTPClient(host, port string) *SMTPClient {
	return &SMTPClient{
		host: host,
		port: port,
	}
}

func (c *SMTPClient) SendEmail(email Email) (*EmailResponse, error) {
	if err := email.Validate(); err != nil {
		return nil, err
	}

	addr := fmt.Sprintf("%s:%s", c.host, c.port)

	headers := make([]string, 0)
	headers = append(headers, fmt.Sprintf("From: %s", email.From))
	headers = append(headers, fmt.Sprintf("To: %s", email.To))
	headers = append(headers, fmt.Sprintf("Subject: %s", email.Subject))
	headers = append(headers, "MIME-Version: 1.0")

	var body string
	if email.HtmlBody != "" {
		headers = append(headers, "Content-Type: text/html; charset=UTF-8")
		body = email.HtmlBody
	} else {
		headers = append(headers, "Content-Type: text/plain; charset=UTF-8")
		body = email.TextBody
	}

	msg := []byte(strings.Join(headers, "\r\n") + "\r\n\r\n" + body)

	err := smtp.SendMail(addr, nil, email.From, []string{email.To}, msg)
	if err != nil {
		return nil, fmt.Errorf("sending email via SMTP: %w", err)
	}

	return &EmailResponse{
		MessageID:   "smtp-dev-mode",
		SubmittedAt: time.Now(),
	}, nil
}
