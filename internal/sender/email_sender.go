package sender

import (
	"context"
	"fmt"
	"net/smtp"

	"Notification-Service/internal/service"
)

type EmailSender struct {
	config EmailConfig
	//implement NotificationSenders
}

type EmailConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	From     string `json:"from"`
}

func NewEmailSender(config EmailConfig) *EmailSender {
	return &EmailSender{config: config}
}

func (s *EmailSender) SendNotification(ctx context.Context, notification service.Notification) error {
	subject := "Уведомление"
	mime := "Mime-Version: 1.0;\nContent-Type: text; charset=\"UTF-8\";\n\n"
	msg := []byte(subject + mime + notification.Message)

	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	err := smtp.SendMail(addr, auth, s.config.From, []string{notification.Address}, msg)

	if err != nil {
		return fmt.Errorf("failed to send email to %s: %w", notification.Address, err)
	}
	return nil
}

func (s *EmailSender) Type() service.NotificationType {
	return service.Email
}
