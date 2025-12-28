package service

import (
	"context"
	"time"
)

type NotificationType string

const (
	Email NotificationType = "email"
	SMS   NotificationType = "sms"
	HTTP  NotificationType = "http"
	// Notifications can be different types (can be expanded)
)

type Notification struct {
	ID       int64
	Type     NotificationType
	Address  string
	IsSended bool
	Message  string

	CreatedAt   time.Time
	ScheduledAt time.Time
	SentAt      *time.Time // pointer for correct nil value conserving
}

type NotificationSender interface {
	SendNotification(ctx context.Context, notification Notification) error
	Type() NotificationType
	// EmailSender, SMSSender, HTTPSender must implement
}
