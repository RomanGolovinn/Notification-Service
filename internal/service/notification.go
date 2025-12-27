package service

import "context"

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
}

type NotificationSender interface {
	SendNotification(context context.Context, notification Notification) error
	Type() NotificationType
}
