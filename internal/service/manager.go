package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

type NotificationRepository interface {
	//interface for db
	Save(context context.Context, notification *Notification) error
	MarkAsSent(ctx context.Context, id int64) error
	GetPending(ctx context.Context) ([]Notification, error)
}

type NotificationService struct {
	repo    NotificationRepository
	senders map[NotificationType]NotificationSender
	logger  *slog.Logger
}

func NewNotificationService(senders []NotificationSender, repo NotificationRepository,
	logger *slog.Logger) *NotificationService {
	senderMap := make(map[NotificationType]NotificationSender)
	for _, sender := range senders {
		senderMap[sender.Type()] = sender
	}
	return &NotificationService{
		senders: senderMap,
		repo:    repo,
		logger:  logger,
	}
}

func (service *NotificationService) Handle(ctx context.Context,
	notification *Notification) error {
	notification.CreatedAt = time.Now()

	notification.IsSended = false

	err := service.repo.Save(ctx, notification)
	if err != nil {
		service.logger.Error("failed to save notification", slog.Any("error", err))
		return fmt.Errorf("iled to save to repo: %w", err)
	}
	return nil
}

func (service *NotificationService) Send(ctx context.Context, notification *Notification) error {
	sender, exist := service.senders[notification.Type]
	if !exist {
		service.logger.Error("sender not found", "type", string(notification.Type))
		return fmt.Errorf("sender not found for type: %s", notification.Type)
	}

	err := sender.SendNotification(ctx, *notification)
	if err != nil {
		service.logger.Error("failed to send notification", slog.Any("error", err))
		return fmt.Errorf("sender.SendNotification")
	}
	service.logger.Info("notification processed successfully", slog.Int64("id", notification.ID))
	return nil
}
