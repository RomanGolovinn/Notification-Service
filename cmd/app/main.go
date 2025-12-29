package main

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"Notification-Service/internal/config"
	"Notification-Service/internal/handler"
	"Notification-Service/internal/repository"
	"Notification-Service/internal/sender"
	"Notification-Service/internal/service"
	"Notification-Service/internal/worker"

	_ "github.com/lib/pq"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		logger.Error("failed to load config", slog.Any("error", err))
		os.Exit(1)
	}
	logger.Info("config loaded successfully")

	db, err := sql.Open("postgres", cfg.DBDSN)
	if err != nil {
		logger.Error("failed to connect to db", slog.Any("error", err))
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		logger.Error("db is unreachable", slog.Any("error", err))
		os.Exit(1)
	}

	repo := repository.NewPostgresRepository(db)

	emailSender := sender.NewEmailSender(cfg.Email)

	senders := []service.NotificationSender{
		emailSender,
	}

	svc := service.NewNotificationService(senders, repo, logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wrk := worker.NewWorker(repo, svc, logger)
	go wrk.Start(ctx)
	logger.Info("background worker started")

	h := handler.NewNotificationHandler(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("/notification", h.Create)

	server := &http.Server{
		Addr:    cfg.ServerPort,
		Handler: mux,
	}

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		logger.Info("shutting down server...")
		cancel()
		server.Shutdown(context.Background())
	}()

	logger.Info("HTTP server starting", slog.String("port", cfg.ServerPort))
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("server failed", slog.Any("error", err))
		os.Exit(1)
	}
}
