package worker

import (
	"context"
	"log/slog"
	"time"

	"Notification-Service/internal/service"
)

type Worker struct {
	repo   service.NotificationRepository
	svc    *service.NotificationService
	logger *slog.Logger
}

func NewWorker(repo service.NotificationRepository, svc *service.NotificationService, logger *slog.Logger) *Worker {
	return &Worker{repo: repo, svc: svc, logger: logger}
}

func (w *Worker) Start(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.processPending(ctx)
		}
	}
}

func (w *Worker) processPending(ctx context.Context) {
	pending, err := w.repo.GetPending(ctx)
	if err != nil {
		w.logger.Error("failed to get pending notifications", "error", err)
		return
	}

	for _, n := range pending {
		err := w.svc.Send(ctx, &n)
		if err != nil {
			w.logger.Error("failed to send", "id", n.ID, "error", err)
			continue
		}

		w.repo.MarkAsSent(ctx, n.ID)
		w.logger.Info("notification sent by worker", "id", n.ID)
	}
}
